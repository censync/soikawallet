package core

import (
	"github.com/censync/soikawallet/api/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testMnemonic = `ceiling code ginger fabric transfer gallery sort deputy engine blur believe sunny divert nephew brain tired result husband upper clock auction ritual correct inhale`

func TestWalletService_Init(t *testing.T) {
	service := &Wallet{}
	_, err := service.Init(&dto.InitWalletDTO{
		Mnemonic:   testMnemonic,
		Passphrase: ``,
	})
	assert.Nil(t, err)
	//priv, _ := service.rootKey.ECPrivKey()
	//pub, _ := service.rootKey.ECPubKey()
	//t.Log(priv.key.String())
	//t.Log(pub.Y(), pub.Y(), err)
}
