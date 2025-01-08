package unisat

import (
	"encoding/json"
	"os"
)

type Withdraw struct {
	Type              string  `json:"type"`
	Valid             bool    `json:"valid"`
	TxID              string  `json:"txid"`
	Idx               int     `json:"idx"`
	Vout              int     `json:"vout"`
	Offset            int     `json:"offset"`
	InscriptionNumber int     `json:"inscriptionNumber"`
	InscriptionID     string  `json:"inscriptionId"`
	ContentType       string  `json:"contentType"`
	ContentBody       string  `json:"contentBody"`
	OldSatPoint       string  `json:"oldSatPoint"`
	NewSatPoint       string  `json:"newSatPoint"`
	From              string  `json:"from"`
	To                string  `json:"to"`
	Satoshi           float64 `json:"satoshi"`
	Data              struct {
		Tick   string `json:"tick"`
		Amount string `json:"amount"`
	} `json:"data"`
	Height    int    `json:"height"`
	TxIdx     int    `json:"txidx"`
	Blockhash string `json:"blockhash"`
	BlockTime int    `json:"blocktime"`

	BlockTimeHumanReadable string `json:"blocktimeHumanReadable"`
}

type WithdrawResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Height int        `json:"height"`
		Total  int        `json:"total"`
		Cursor int        `json:"cursor"`
		Detail []Withdraw `json:"detail"`
	} `json:"data"`
}

func parseWithdrawResponse() []Withdraw {
	bs, err := os.ReadFile("../external/unisat/withdraw.json")
	if err != nil {
		panic(err)
	}
	var resp WithdrawResponse
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		panic(err)
	}
	return resp.Data.Detail
}
