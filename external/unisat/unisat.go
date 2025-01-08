package unisat

import (
	"context"
	"fmt"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (svc *Service) GetWithdrawTransactions(ctx context.Context) []string {
	withdraws := parseWithdrawResponse()
	fmt.Printf("withdraw transactions from unisat: %d\n", len(withdraws))

	txIds := make([]string, len(withdraws))
	for idx, withdraw := range withdraws {
		txIds[len(withdraws)-1-idx] = withdraw.TxID
	}

	return txIds
}
