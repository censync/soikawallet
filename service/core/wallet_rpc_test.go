package core

import (
	"os"
	"testing"
)

var (
	testTestnetMnemonic = os.Getenv("TESTNET_MNEMONIC")
	testnetPassphrase   = os.Getenv("TESTNET_PASSPHRASE")
)

func TestWalletService_ethereum(t *testing.T) {

}
