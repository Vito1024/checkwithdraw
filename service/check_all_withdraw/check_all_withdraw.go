package checkallwithdraw

import (
	"context"
	"withdraw"
)

type Service struct {
	oklinkSvc withdraw.OKLinkService
}

func New(oklinkSvc withdraw.OKLinkService) *Service {
	return &Service{
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
