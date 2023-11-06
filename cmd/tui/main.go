// Copyright 2023 The soikawallet Authors
// This file is part of soikawallet.
//
// soikawallet is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// soikawallet is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with  soikawallet. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"github.com/censync/soikawallet/service/tui"
	"github.com/censync/soikawallet/service/tui/config"
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

	provider := tui.NewTUIServiceProvider(cfg, &wg)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		defer wg.Done()
		<-signalChannel
		log.Println("INTERRUPT")
		provider.Web3Connection().Stop()
		provider.UI().Stop()
		log.Println("Stopped")
	}()

	provider.Web3Connection().Start()
	provider.UI().Start()

	//wg.Wait()
	//tui.Init().Start(*optionVerbose)
}
