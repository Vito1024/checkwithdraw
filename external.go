package withdraw

import (
	"context"
	"errors"
)

const (
	OKLinkBRC20ActionWithdraw = "swapWithdraw"
)

type OKLinkService interface {
	GetFractalBitcoinBRC20TransactionDetail(ctx context.Context, txId string, opts ...func(*Option)) (OKLinkBRC20TransactionDetail, error)
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
	GetWithdrawTransactionsFromFile(ctx context.Context) []string
	FollowWithdrawTransactions(ctx context.Context) <-chan UnisatWithdrawTransaction
}

type UnisatWithdrawTransaction struct {
	Type              string `json:"type"`
	Valid             bool   `json:"valid"`
	TxId              string `json:"txid"`
	Idx               int    `json:"idx"`
	Vout              int    `json:"vout"`
	Offset            int    `json:"offset"`
	InscriptionNumber int    `json:"inscriptionNumber"`
	InscriptionId     string `json:"inscriptionId"`
	ContentType       string `json:"contentType"`
	ContentBody       string `json:"contentBody"`
	OldSatPoint       string `json:"oldSatPoint"`
	NewSatPoint       string `json:"newSatPoint"`
	From              string `json:"from"`
	To                string `json:"to"`
	Satoshi           uint64 `json:"satoshi"`
	Data              struct {
		Tick   string `json:"tick"`
		Amount string `json:"amount"`
	} `json:"data"`
	Height    int    `json:"height"`
	TxIdx     int    `json:"txidx"`
	BlockHash string `json:"blockhash"`
	BlockTime int    `json:"blocktime"`
}
