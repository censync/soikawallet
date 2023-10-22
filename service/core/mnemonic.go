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
	"github.com/tyler-smith/go-bip39"
	"strings"
)

const (
	Entropy128 = 128
	Entropy160 = 160
	Entropy192 = 192
	Entropy224 = 224
	Entropy256 = 256
)

var (
	errMnemonicGenAttemptsReached = errors.New("reached attempts to generate mnemonic")

	mnemonicCount = map[int]int{
		12: Entropy128,
		15: Entropy160,
		18: Entropy192,
		21: Entropy224,
		24: Entropy256,
	}
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
		err = mnemonicCheck(mnemonic)
		if err != nil {
			// log.Printf("Generated bad mnemonic: %s\n", err)
			continue
		}
		return mnemonic, nil
	}
	return "", errMnemonicGenAttemptsReached
}

// 128 - 12, 160 - 15, 192 - 18, 224 - 21, 256 - 24
func mnemonicCheck(mnemonic string) error {
	entropy := entropyByMnemonic(mnemonic)

	if entropy == 0 {
		return errors.New("undefined entropy")
	}

	if !checkDuplicates(mnemonic) {
		return errors.New("mnemonic has duplicates prefix")
	}
	return nil
}

func entropyByMnemonic(mnemonic string) int {
	wordsCount := len(strings.Fields(mnemonic))
	return mnemonicCount[wordsCount]
}

func checkDuplicates(str string) bool {
	var index = map[string]bool{}
	arr := strings.Fields(str)
	for idx := range arr {
		prefix := substr(arr[idx], 4)
		if _, ok := index[prefix]; !ok {
			index[prefix] = true
		} else {
			return false
		}
	}
	return true
}

func substr(input string, length int) string {
	asRunes := []rune(input)

	if 0 >= len(asRunes) {
		return ""
	}

	if length > len(asRunes) {
		length = len(asRunes) - 0
	}

	return string(asRunes[0 : 0+length])
}
