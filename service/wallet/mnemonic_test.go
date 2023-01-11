package wallet

import (
	"crypto/sha512"
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/types"
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

		err := service.Init(&dto.InitWalletDTO{
			Mnemonic:        vector.Mnemonic,
			Passphrase:      testPassphrase,
			SkipPrefixCheck: true,
		})

		assert.Nil(t, err)

		for _, testAddr := range vector.Addresses {
			path, err := types.ParsePath(testAddr.Path)

			assert.Nil(t, err)

			addr, err := service.addAddress(path)

			assert.Nil(t, err)

			assert.Equal(t, testAddr.Address, addr.Address())

			// WRONG: t.Logf("%s vs %s", testAddr.PublicKey, hexutil.Encode(crypto.FromECDSAPub(addr.pub)))

			// t.Logf("%s vs %s", testAddr.PrivateKey, hexutil.Encode(crypto.FromECDSA(addr.key.Get())))

			/*if addr != nil {
				if testAddr.Address != addr.Address() {
					t.Logf("%s: %s vs %s", vector.Entropy, testAddr.Address, addr.Address())
				}
			}*/
		}
	}
}
