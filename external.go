package withdraw

import (
	"context"
	"errors"
)

const (
	OKLinkBRC20ActionWithdraw = "swapWithdraw"
)

type OKLinkService interface {
	GetFractalBitcoinBRC20TransactionDetail(ctx context.Context, txId string) (OKLinkBRC20TransactionDetail, error)
	GetFractalBitcoinBRC20TransactionDetailBatch(ctx context.Context, txIds []string, opts ...func(*Option)) <-chan OKLinkBRC20TransactionDetailBatch
}

var (
	ErrTransactionNotBRC20Withdraw = errors.New("transaction is not a BRC20 withdraw")
)

type Option struct {
	RateLimit bool
	Progress  bool
}

func RequestOkLinkWithRateLimit() func(*Option) {
	return func(o *Option) {
		o.RateLimit = true
	}
}

func WithProgress() func(*Option) {
	return func(o *Option) {
		o.Progress = true
	}
}

type OKLinkBRC20TransactionDetail struct {
	Symbol             string `json:"symbol"`
	Action             string `json:"action"`
	ProtocolType       string `json:"protocolType"`
	State              string `json:"state"`
	Amount             string `json:"amount"`
	From               string `json:"from"`
	To                 string `json:"to"`
	InscriptionId      string `json:"inscriptionId"`
	InscriptionNumber  string `json:"inscriptionNumber"`
	OutputIndex        string `json:"outputIndex"`
	TokenInscriptionId string `json:"tokenInscriptionId"`
	TxId               string `json:"txId"`
	TransactionTime    string `json:"transactionTime"`
	Height             string `json:"height"`
	BlockHash          string `json:"blockHash"`
}

type OKLinkBRC20TransactionDetailBatch struct {
	OKLinkBRC20TransactionDetail
	Err error
}

type UnisatService interface {
	GetWithdrawTransactions(ctx context.Context) []string
}
