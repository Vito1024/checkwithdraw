package main

import (
	"context"
	"encoding/json"
	"os"
	"withdraw"
	"withdraw/config"
	"withdraw/external/oklink"
	"withdraw/external/unisat"
	checkallwithdraw "withdraw/service/check_all_withdraw"
)

type dep struct {
	config *config.Config

	unisatSvc withdraw.UnisatService
	oklinkSvc withdraw.OKLinkService
}

func main() {
	var dep dep
	dep.initConfig()
	dep.initExternal()

	withdrawSvc := checkallwithdraw.New(dep.oklinkSvc)
	notInOKLinkTdIDs := withdrawSvc.FilterNotBRC20WithdrawByOKLink(context.Background(), dep.unisatSvc.GetWithdrawTransactions(context.Background())...)

	bs, err := json.Marshal(notInOKLinkTdIDs)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("not_as_withdraw_by_oklink.json", bs, 0644)
	if err != nil {
		panic(err)
	}
}

func (d *dep) initConfig() {
	d.config = config.New()
}

func (d *dep) initExternal() {
	d.unisatSvc = unisat.New()
	d.oklinkSvc = oklink.New(d.config.OkLinkConfig)
}
