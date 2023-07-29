package wallet

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
	//priv, _ := service.bip44Key.ECPrivKey()
	//pub, _ := service.bip44Key.ECPubKey()
	//t.Log(priv.Key.String())
	//t.Log(pub.Y(), pub.Y(), err)
}

/*
var s = `{"meta":{"v":1,"labels":{"1":{"1":"Account label 1","2":"Account label 2"},"2":{"1":"addr label 1","2":"addr label 2"}}},"addresses":["m/44'/60'/0'/0/0'","m/44'/60'/0'/0/1'","m/44'/60'/0'/0/2'","m/44'/60'/0'/0/3'","m/44'/60'/0'/0/4'","m/44'/60'/0'/0/5'","m/44'/60'/0'/0/6'","m/44'/60'/0'/0/7'","m/44'/60'/0'/0/8'","m/44'/60'/0'/0/9'","m/44'/195'/0'/0/10'","m/44'/195'/0'/0/0'","m/44'/195'/0'/0/1'","m/44'/195'/0'/0/2'","m/44'/195'/0'/0/3'","m/44'/195'/0'/0/4'","m/44'/195'/0'/0/5'","m/44'/195'/0'/0/6'","m/44'/195'/0'/0/7'","m/44'/195'/0'/0/8'","m/44'/195'/0'/0/9'","m/44'/195'/0'/0/10'"]}`

func TestWalletService_MarshalJSON(t *testing.T) {
	service := &Wallet{}
	err := service.Init(&dto.InitWalletDTO{
		Mnemonic:   testMnemonic,
		Passphrase: ``,
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		dPath, err := types.CreatePath(types.Ethereum, 0, types.ChargeExternal, types.AddressIndex{InternalIndex: uint32(i), IsHardened: true})

		if err != nil {
			t.Fatal(err)
		}

		_, err = service.addAddress(dPath)

		if err != nil {
			t.Fatal(err)
		}

		//t.Log(addr.addr())
	}

	for i := 0; i < 10; i++ {
		dPath, err := types.CreatePath(types.Tron, 0, types.ChargeExternal, types.AddressIndex{InternalIndex: uint32(i), IsHardened: true})

		if err != nil {
			t.Fatal(err)
		}

		_, err = service.addAddress(dPath)

		if err != nil {
			t.Fatal(err)
		}

		//t.Log(addr.addr())
	}

	_, err = service.AddLabel(&dto.AddLabelDTO{
		LabelType: uint8(meta.AccountLabel),
		Title:     "Account label 1",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = service.AddLabel(&dto.AddLabelDTO{
		LabelType: uint8(meta.AccountLabel),
		Title:     "Account label 2",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = service.AddLabel(&dto.AddLabelDTO{
		LabelType: uint8(meta.AddressLabel),
		Title:     "addr label 1",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = service.AddLabel(&dto.AddLabelDTO{
		LabelType: uint8(meta.AddressLabel),
		Title:     "addr label 2",
	})

	if err != nil {
		t.Fatal(err)
	}

	str, err := json.Marshal(service)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, s, string(str))

	//fmt.Println(string(str), err)
}
*/
