package main

import (
	"flag"
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/service"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	optionVerbose = flag.Bool("v", false, "Show additional log messages")
)

func init() {
	flag.Parse()
}

func main() {
	cfg := &config.Config{
		Verbose: *optionVerbose,
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	provider := service.NewServiceProvider(cfg, &wg)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		defer wg.Done()
		<-signalChannel
		log.Println("INTERRUPT")
		provider.Web3Connection().Stop()
		provider.UI().Stop()
	}()

	provider.Web3Connection().Start()
	provider.UI().Start()

	wg.Wait()
	//tui.Init().Start(*optionVerbose)
}
