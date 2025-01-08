package oklink

import (
	"context"
	"encoding/json"
	"testing"
	"withdraw"
	"withdraw/config"
)

func newSvc() *Service {
	return New(config.New().OkLinkConfig)
}

func TestGetFractalBitcoinBRC20TransactionDetail(t *testing.T) {
	txHash := "bd331d228d959479da5a100f26917cdd11dc99a69c65f77228092d359696381e"
	detail, err := newSvc().GetFractalBitcoinBRC20TransactionDetail(context.Background(), txHash)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := json.Marshal(detail)
	if err != nil {
		panic(err)
	}

	t.Logf("%s", bs)
}

func TestGetFractalBitcoinBRC20TransactionDetailBatch(t *testing.T) {
	txHashes := []string{
		"bd331d228d959479da5a100f26917cdd11dc99a69c65f77228092d359696381e",
	}
	results := newSvc().GetFractalBitcoinBRC20TransactionDetailBatch(
		context.Background(),
		txHashes,
		withdraw.RequestOkLinkWithRateLimit(),
	)

	t.Logf("%+v", results)
}

func TestGetFractalBitcoinAddressTokenTransactionList(t *testing.T) {
	inscriptionID := "01228e2219333f2188f234297576438b634bf16f80df041705d7785c206e9a92i0"
	address := "bc1pdlm2cgre95xq22pnzal5xzafljfe8gqgrdcdce26wwm298z6grtq65kpz6"
	transactionList, err := newSvc().GetFractalBitcoinAddressTokenTransactionList(context.Background(), address, "brc20", inscriptionID)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := json.Marshal(transactionList)
	if err != nil {
		panic(err)
	}
	t.Logf("%s", bs)
}

func TestGetFractalBitcoinAddressBrc20InscriptionList(t *testing.T) {
	address := "bc1pdlm2cgre95xq22pnzal5xzafljfe8gqgrdcdce26wwm298z6grtq65kpz6"
	resp, err := newSvc().GetFractalBitcoinAddressBrc20InscriptionList(
		context.Background(),
		address,
		RequestAddressBrc20InscriptionListWithSymbol("GLIZZY"),
	)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.Marshal(resp)
	t.Log(string(bs))
}
