package wallet

import (
	"errors"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/util/seed"
	"github.com/tyler-smith/go-bip39"
)

var (
	errMnemonicGenAttemptsReached = errors.New("reached attempts to generate mnemonic")
)

// Separated for vectors testing
func generateKeyFromSeed(seed *[]byte) (*hdkeychain.ExtendedKey, error) {
	key, err := hdkeychain.NewMaster(*seed, &chaincfg.MainNetParams)

	if err != nil {
		return nil, errors.New("cannot initialize master key")
	}
	return key, nil
}

func (s *Wallet) GenerateMnemonic(dto *dto.GenerateMnemonicDTO) (string, error) {
	for attempts := 200; attempts > 0; attempts-- {
		entropy, err := bip39.NewEntropy(dto.BitSize)
		if err != nil {
			//log.Printf("Cannot initialize entropy for mnemonic: %s\n", err)
			continue
		}
		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			// log.Printf("Cannot generate mnemonic: %s\n", err)
			continue
		}
		err = seed.Check(mnemonic)
		if err != nil {
			// log.Printf("Generated bad mnemonic: %s\n", err)
			continue
		}
		return mnemonic, nil
	}
	return "", errMnemonicGenAttemptsReached
}
