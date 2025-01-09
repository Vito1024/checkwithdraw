package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

var configPath = flag.String("config", "", "config file path")

func main() {
	withdraw.ParseEnv()
	flag.Parse()

	var dep dep
	dep.initConfig()
	dep.initExternal()

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	go func() {
		s := <-sig
		fmt.Printf("signal received: %v\n", s)
		cancel()
	}()

	withdrawSvc := checkallwithdraw.New(dep.config.CheckWithdraw, dep.unisatSvc, dep.oklinkSvc)
	withdrawSvc.FollowWithdrawTransactions(ctx)
}

func (d *dep) initConfig() {
	d.config = config.New(*configPath)
}

func (d *dep) initExternal() {
	d.unisatSvc = unisat.New(d.config.UnisatConfig)
	d.oklinkSvc = oklink.New(d.config.OkLinkConfig)
}
