// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

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
