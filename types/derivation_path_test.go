package types

import (
	"testing"
)

const (
	pathEthereum         = `m/44'/60'/0'/0/0`
	pathTron             = `m/44'/195'/0'/0/0`
	pathEthereumHardened = `m/44'/60'/0'/0/0'`
	pathTronHardened     = `m/44'/195'/0'/0/0'`
)

func TestDerivationPath_Parse(t *testing.T) {
	path, err := ParsePath(pathEthereum)
	t.Log(path, err)
	path, err = ParsePath(pathTron)
	t.Log(path, err)
	path, err = ParsePath(pathEthereumHardened)
	t.Log(path, err)
	path, err = ParsePath(pathTronHardened)
	t.Log(path, err)
}

func TestDerivationPath_Validate(t *testing.T) {

}
