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
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	airgap "github.com/censync/go-airgap"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/core/internal/config/version"
	"github.com/censync/soikawallet/service/core/internal/network"
	types2 "github.com/censync/soikawallet/service/core/internal/types"
	"github.com/censync/soikawallet/service/core/internal/types/currencies"
	"github.com/censync/soikawallet/service/core/meta"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

const (
	fiatTitle  = "USD"
	fiatSymbol = "$"
)

var (
	errWalletAlreadyInitialized      = errors.New("core already initialized")
	errWalletKeyRootCannotInitialize = errors.New("cannot initialize root key")
)

type Wallet struct {
	// instanceId compressed public key for root key, used for identify device instance
	instanceId     []byte
	rootKey        *hdkeychain.ExtendedKey
	meta           *meta.Meta
	currenciesFiat *currencies.FiatCurrencies
}

func (s *Wallet) getNetworkProvider(ctx *types2.RPCContext) (types2.NetworkAdapter, error) {
	return network.WithContext(ctx)
}

func (s *Wallet) getRPCProvider(chainKey mhda.ChainKey) types2.RPCAdapter {
	return network.Get(chainKey)
}

// Init initializes static instance of wallet with mnemonic and optional passphrase.
// If result is successful, will be returned base58 encoded compressed root public key.
func (s *Wallet) Init(dto *dto.InitWalletDTO) (string, error) {
	var err error
	dto.Mnemonic = strings.TrimSpace(dto.Mnemonic)
	dto.Passphrase = strings.TrimSpace(dto.Passphrase)

	// Check for singleton
	if s.instanceId != nil {
		return "", errWalletAlreadyInitialized
	}

	// SkipMnemonicCheck flag used only for testing vectors
	if !dto.SkipMnemonicCheck {
		err = mnemonicCheck(dto.Mnemonic)
	}

	if err != nil {
		return "", err
	}

	rootSeed := pbkdf2.Key([]byte(dto.Mnemonic), []byte("mnemonic"+dto.Passphrase), 2048, 64, sha512.New)

	masterKey, err := generateKeyFromSeed(&rootSeed)

	if err != nil {
		return "", errWalletKeyRootCannotInitialize
	}

	masterPubKey, err := masterKey.ECPubKey()

	// ROOT pub key
	if err != nil {
		return "", errAddrKeyCannotCreate
	}

	*s = Wallet{
		instanceId:     masterPubKey.SerializeCompressed(),
		rootKey:        masterKey,
		meta:           meta.InitMeta(),
		currenciesFiat: currencies.NewFiatCurrencies(fiatTitle, fiatSymbol),
	}
	return s.getInstanceId(), nil
}

func (s *Wallet) getInstanceId() string {
	return base58.Encode(s.instanceId)
}

func (s *Wallet) GetAccountsByNetwork(dto *dto.GetAccountsByNetworkDTO) []*resp.AccountResponse {
	accountsIndex := map[mhda.AccountIndex]bool{}

	for _, addr := range s.meta.Addresses() {
		if addr.MHDA().Chain().Key() == dto.ChainKey {
			accountsIndex[addr.MHDA().DerivationPath().Account()] = true
		}
	}

	accounts := make([]*resp.AccountResponse, 0)

	for accountIndex := range accountsIndex {
		// Deprecated
		// 	accountPath, err := types.CreateAccountPath(types.CoinType(dto.NetworkType), accountIndex)
		// 	if err != nil {
		// 		continue
		//	}

		// Fix account path
		accounts = append(accounts, &resp.AccountResponse{
			// Path:        accountPath.String(),
			ChainKey: dto.ChainKey,
			Account:  accountIndex,
			Label:    s.meta.GetAccountLabel(dto.ChainKey, accountIndex),
		})
	}

	return accounts
}

func (s *Wallet) FlushKeys(dto *dto.FlushKeysDTO) {
	s.rootKey = nil
	/*for key := range s.addresses {
		if dto.Force || !s.addresses[key].isKeyDelivered {
			s.addresses[key].key.Free()
			s.addresses[key].key = nil
		}
	}*/
}

func (s *Wallet) Version() string {
	return version.VERSION
}

func (s *Wallet) ExportMeta() (*resp.AirGapMessage, error) {
	data, err := s.meta.MarshalJSON()
	if err != nil {
		return nil, err
	}
	airgapMsg := airgap.NewAirGap(airgap.VersionDefault, s.instanceId).
		CreateMessage().
		AddOperation(types2.OpMetaWallet, data)
	chunks, err := airgapMsg.MarshalB64Chunks()
	if err != nil {
		return nil, err
	}
	return &resp.AirGapMessage{
		Chunks: chunks,
	}, nil
}

func (s *Wallet) ExportMetaDebug() ([]byte, error) {
	data, err := s.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Wallet) MarshalJSON() ([]byte, error) {
	var strPaths []string
	// TODO: Add internal
	addresses := s.getAllAddresses()
	for index := range addresses {
		strPaths = append(strPaths, addresses[index].Path)
	}
	return json.Marshal(&struct {
		Meta      *meta.Meta `json:"meta"`
		Addresses []string   `json:"addresses"`
	}{
		Meta:      s.meta,
		Addresses: strPaths,
	})
}
