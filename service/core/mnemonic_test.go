package core

import (
	"crypto/sha512"
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/config/chain"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/pbkdf2"
	"testing"
)

func TestWallet_generateKeyFromSeedPositive(t *testing.T) {
	for _, vector := range testVectors {

		seed := pbkdf2.Key([]byte(vector.Mnemonic), []byte("mnemonic"+testPassphrase), 2048, 64, sha512.New)

		assert.Equal(t, vector.Seed, fmt.Sprintf("%x", seed))

		key, err := generateKeyFromSeed(&seed)

		assert.Nil(t, err)

		assert.Equal(t, vector.Key, key.String())
	}
}

func TestWallet_AddAddressesPositive(t *testing.T) {
	for _, vector := range testVectors {
		service := &Wallet{}

		_, err := service.Init(&dto.InitWalletDTO{
			Mnemonic:          vector.Mnemonic,
			Passphrase:        testPassphrase,
			SkipMnemonicCheck: true,
		})

		assert.Nil(t, err)

		for _, testAddr := range vector.Addresses {

			path, err := mhda.ParseDerivationPath(mhda.BIP44, testAddr.Path)

			assert.Nil(t, err)

			addrKey := mhda.NewAddress(chain.EthereumChain, path)

			// Temporary
			if addrKey.DerivationPath().Coin() != 60 {
				continue
			}

			addr, err := service.addAddress(addrKey)

			assert.Nil(t, err)

			assert.Equal(t, testAddr.Address, addr.Address())
		}
	}
}
