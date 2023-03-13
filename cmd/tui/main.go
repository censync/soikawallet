package main

import (
	"flag"
	"github.com/censync/soikawallet/service/ui"
)

var (
	optionVerbose = flag.Bool("v", false, "Show additional log messages")
)

func init() {
	flag.Parse()
}

func main() {
	ui.Init().Run(*optionVerbose)
}
