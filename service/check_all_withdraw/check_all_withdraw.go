package checkallwithdraw

import (
	"context"
	"fmt"
	"withdraw"
	"withdraw/config"
)

type Service struct {
	config config.CheckWithdraw

	unisatSvc withdraw.UnisatService
	oklinkSvc withdraw.OKLinkService
}

func New(config config.CheckWithdraw, unisatSvc withdraw.UnisatService, oklinkSvc withdraw.OKLinkService) *Service {
	return &Service{
		config:    config,
		unisatSvc: unisatSvc,
		oklinkSvc: oklinkSvc,
	}
}

func (svc *Service) FilterNotBRC20WithdrawByOKLink(ctx context.Context, withdrawTransactions ...string) []string {
	var notBRC20Withdraws []string

	out := svc.oklinkSvc.GetFractalBitcoinBRC20TransactionDetailBatch(
		ctx, withdrawTransactions,
		withdraw.RequestOkLinkWithRateLimit(),
		withdraw.WithProgress(),
	)
	for result := range out {
		if result.Err == withdraw.ErrTransactionNotBRC20Withdraw {
			notBRC20Withdraws = append(notBRC20Withdraws, result.TxId)
		}
	}

	return notBRC20Withdraws
}

func (svc *Service) FollowWithdrawTransactions(ctx context.Context) {
	for unisatWithdrawTx := range svc.unisatSvc.FollowWithdrawTransactions(ctx) {
		if unisatWithdrawTx.TxId == "" || len(unisatWithdrawTx.TxId) != 64 {
			continue
		}
		if inStrSlice(unisatWithdrawTx.From, svc.config.ExcludedAddresses) || inStrSlice(unisatWithdrawTx.To, svc.config.ExcludedAddresses) {
			continue
		}

		_, err := svc.oklinkSvc.GetFractalBitcoinBRC20TransactionDetail(ctx, unisatWithdrawTx.TxId, withdraw.RequestOkLinkWithRateLimit())
		if err == withdraw.ErrTransactionNotBRC20Withdraw {
			fmt.Printf("found mismatch withdraw tx: %s height: %d from: %s to: %s\n", unisatWithdrawTx.TxId, unisatWithdrawTx.Height, unisatWithdrawTx.From, unisatWithdrawTx.To)
			continue
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("match tx: %s\n", unisatWithdrawTx.TxId)
	}
}

func inStrSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
