package wallet

import (
	"crypto/sha512"
	"errors"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/wallet/meta"
	"github.com/censync/soikawallet/util/seed"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

func (s *Wallet) Init(dto *dto.InitWalletDTO) error {
	var err error
	dto.Mnemonic = strings.TrimSpace(dto.Mnemonic)
	dto.Passphrase = strings.TrimSpace(dto.Passphrase)

	if !dto.SkipPrefixCheck {
		err = seed.Check(dto.Mnemonic)
	}

	if err != nil {
		return err
	}

	rootSeed := pbkdf2.Key([]byte(dto.Mnemonic), []byte("mnemonic"+dto.Passphrase), 2048, 64, sha512.New)

	masterKey, err := generateKeyFromSeed(&rootSeed)

	bip44Key, err := masterKey.Derive(hardenedKeyStart + 44)
	if err != nil {
		return errors.New("cannot initialize BIP-44 key")
	}

	*s = Wallet{
		instanceId: uuid.New(),
		bip44Key:   bip44Key,
		addresses:  map[string]*address{},
		meta:       meta.InitMeta(),
	}
	return nil
}

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
	return "", errors.New("reached attempts to generate mnemonic")
}
