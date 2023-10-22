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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/censync/soikawallet/service/core/internal/types/protected_key"
	"github.com/censync/soikawallet/service/core/meta"
)

const (
	hardenedKeyStart = uint32(0x80000000) // 2^31
)

var (
	errWalletRootKeyNotSet             = errors.New("root key is not set")
	errAddrKeyCannotCreate             = errors.New("cannot create addr key")
	errAddrDerivationPathUndefinedType = errors.New("undefined derivation type")
	errAddrW3AlreadyPermitted          = errors.New("address already permitted for web3")
	errAddrW3NotPermitted              = errors.New("address not permitted for web3")
	errCannotCalculateDerivedKey       = errors.New("cannot calculate derived key")
)

func (s *Wallet) chargeDeriveKey(path *mhda.DerivationPath) (*ecdsa.PrivateKey, error) {
	if s.rootKey == nil {
		return nil, errWalletRootKeyNotSet
	}

	switch path.DerivationType() {
	case mhda.ROOT:
		return s.derivationKeyRoot()
	case mhda.BIP32:
		return s.derivationKeyBip32(path)
	case mhda.BIP44:
		return s.derivationKeyBip44(path)
	case mhda.BIP84:
		return s.derivationKeyBip84(path)
	case mhda.CIP11:
		return s.derivationKeyCip11(path)
	case mhda.ZIP32:
		return s.derivationKeyZip32(path)
	default:
		return nil, errAddrDerivationPathUndefinedType
	}
}
func (s *Wallet) derivationKeyRoot() (*ecdsa.PrivateKey, error) {
	if s.rootKey == nil {
		return nil, errWalletRootKeyNotSet
	}

	ecAddrKey, err := s.rootKey.ECPrivKey()

	if err != nil {
		return nil, errAddrKeyCannotCreate
	}

	return ecAddrKey.ToECDSA(), nil
}

func (s *Wallet) derivationKeyBip32(path *mhda.DerivationPath) (*ecdsa.PrivateKey, error) {
	if s.rootKey == nil {
		return nil, errWalletRootKeyNotSet
	}

	// BIP-32 level
	bip32Key, err := s.rootKey.Derive(hardenedKeyStart + 32) // ??
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	// m / account ' / charge / Address
	// BIP-32 account level
	accountKey, err := bip32Key.Derive(hardenedKeyStart + uint32(path.Account()))
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	// BIP-32 address level
	chargeKey, err := accountKey.Derive(uint32(path.Charge()))
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	var (
		key *hdkeychain.ExtendedKey
	)

	if path.IsHardenedAddress() {
		key, err = chargeKey.Derive(hardenedKeyStart + path.AddressIndex().Index)
	} else {
		key, err = chargeKey.Derive(path.AddressIndex().Index)
	}

	if err != nil {
		return nil, errAddrKeyCannotCreate
	}

	ecAddrKey, err := key.ECPrivKey()

	if err != nil {
		return nil, errAddrKeyCannotCreate
	}

	return ecAddrKey.ToECDSA(), nil
}

func (s *Wallet) derivationKeyBip44(path *mhda.DerivationPath) (*ecdsa.PrivateKey, error) {
	if s.rootKey == nil {
		return nil, errWalletRootKeyNotSet
	}

	// BIP-44 level
	bip44Key, err := s.rootKey.Derive(hardenedKeyStart + 44)
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	// m/44'/60'
	// BIP-44 network (coin) level
	networkKey, err := bip44Key.Derive(hardenedKeyStart + uint32(path.Coin()))
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	// m/44'/60'/0'
	// BIP-44 account level
	accountKey, err := networkKey.Derive(hardenedKeyStart + uint32(path.Account()))
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	// m/44'/60'/0'/0
	// BIP-44 charge level
	chargeKey, err := accountKey.Derive(uint32(path.Charge()))

	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}
	var (
		key *hdkeychain.ExtendedKey
	)

	if path.IsHardenedAddress() {
		key, err = chargeKey.Derive(hardenedKeyStart + path.AddressIndex().Index)
	} else {
		key, err = chargeKey.Derive(path.AddressIndex().Index)
	}

	// BIP-44 address level
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	ecAddrKey, err := key.ECPrivKey()

	if err != nil {
		return nil, errAddrKeyCannotCreate
	}

	return ecAddrKey.ToECDSA(), nil
}

func (s *Wallet) derivationKeyBip84(path *mhda.DerivationPath) (*ecdsa.PrivateKey, error) {
	if s.rootKey == nil {
		return nil, errWalletRootKeyNotSet
	}
	// BIP-32 level
	bip84Key, err := s.rootKey.Derive(hardenedKeyStart + 84)
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	// m / 84 ' / 0 ' / account ' / charge / Address
	// BIP-84 network (coin) level (const BTC)
	networkKey, err := bip84Key.Derive(hardenedKeyStart + uint32(0)) // Coin, BTC=0
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	// BIP-84 account level
	accountKey, err := networkKey.Derive(hardenedKeyStart + uint32(path.Account()))
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	// BIP-84 charge level
	chargeKey, err := accountKey.Derive(uint32(path.Charge()))

	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	var (
		key *hdkeychain.ExtendedKey
	)

	if path.IsHardenedAddress() {
		key, err = chargeKey.Derive(hardenedKeyStart + path.AddressIndex().Index)
	} else {
		key, err = chargeKey.Derive(path.AddressIndex().Index)
	}

	// BIP-84 address level
	if err != nil {
		return nil, errCannotCalculateDerivedKey
	}

	ecAddrKey, err := key.ECPrivKey()

	if err != nil {
		return nil, errAddrKeyCannotCreate
	}

	return ecAddrKey.ToECDSA(), nil
}

func (s *Wallet) derivationKeyCip11(path *mhda.DerivationPath) (*ecdsa.PrivateKey, error) {
	return nil, nil
}

func (s *Wallet) derivationKeyZip32(path *mhda.DerivationPath) (*ecdsa.PrivateKey, error) {
	return nil, nil
}

// Deprecated
/*
func (s *Wallet) isAccountExists(networkType types.CoinType, accountIndex types.AccountIndex) bool {
	for _, addr := range s.addresses {
		if addr.NetworkType() == networkType && addr.Account() == accountIndex {
			return true
		}
	}
	return false
}*/

func (s *Wallet) addAddress(path mhda.MHDA) (addr *meta.Address, err error) {
	if s.rootKey == nil {
		return nil, errWalletRootKeyNotSet
	}

	// TODO: Updated preconfigured
	if !types.IsNetworkExists(path.Chain().Key()) {
		return nil, errChainKeyNotSupported
	}
	if s.meta.IsAddressExist(path.NSS()) {
		return nil, errors.New("addr already exists")
	}

	// Create addr

	ecAddrKey, err := s.chargeDeriveKey(path.DerivationPath())

	if err != nil {
		return nil, err
	}

	pubKey := ecAddrKey.Public().(*ecdsa.PublicKey)

	ctx := types.NewRPCContext(path.Chain().Key(), 0)
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	addr = meta.NewAddress(
		path,
		protected_key.NewProtectedKey(ecAddrKey),
		pubKey,
		provider.Address(pubKey), // TODO: Move addr marshaller from provider
	)

	s.meta.SetAddress(path.NSS(), addr)

	return addr, nil
}

func (s *Wallet) AddAddresses(dto *dto.AddAddressesDTO) (addresses []*resp.AddressResponse, err error) {
	if len(dto.MhdaPaths) == 0 {
		return nil, errors.New("derivation paths is not set")
	}
	for i := range dto.MhdaPaths {
		dPath, err := mhda.ParseNSS(dto.MhdaPaths[i])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot parse derivation path: %s", err))
		}
		addr, err := s.addAddress(dPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot create account: %s", err))
		}
		addresses = append(addresses, &resp.AddressResponse{
			Address:      addr.Address(),
			Path:         addr.MHDA().NSS(),
			IsExternal:   addr.IsExternal(),
			AddressIndex: addr.AddressIndex(),
			ChainKey:     addr.MHDA().Chain().Key(),
			Account:      addr.Account(),
			Label:        s.meta.GetAddressLabel(addr.MHDA().NSS()),
			IsW3:         addr.IsW3(),
		})
	}
	return addresses, nil
}

// TODO: Add error for response
func (s *Wallet) GetAddressesByAccount(dto *dto.GetAddressesByAccountDTO) []*resp.AddressResponse {
	var addresses []*resp.AddressResponse

	for _, addr := range s.meta.Addresses() {
		if (addr.MHDA().Chain().Key() == dto.ChainKey) &&
			addr.DerivationPath().Account() == mhda.AccountIndex(dto.AccountIndex) {
			addresses = append(addresses, &resp.AddressResponse{
				Address:      addr.Address(),
				Path:         addr.MHDA().NSS(),
				IsExternal:   addr.IsExternal(),
				AddressIndex: addr.AddressIndex(),
				ChainKey:     addr.MHDA().Chain().Key(),
				Account:      addr.Account(),
				Label:        s.meta.GetAddressLabel(addr.MHDA().NSS()),
				IsW3:         addr.IsW3(),
			})
		}
	}

	return addresses
}

func (s *Wallet) getAllAddresses() []*resp.AddressResponse {
	var addresses []*resp.AddressResponse
	for _, addr := range s.meta.Addresses() {
		addresses = append(addresses, &resp.AddressResponse{
			Address:      addr.Address(),
			Path:         addr.MHDA().NSS(),
			IsExternal:   addr.IsExternal(),
			AddressIndex: addr.AddressIndex(),
			ChainKey:     addr.MHDA().Chain().Key(),
			Account:      addr.Account(),
			Label:        s.meta.GetAddressLabel(addr.MHDA().NSS()),
			IsW3:         addr.IsW3(),
		})
	}
	return addresses
}

func (s *Wallet) GetTokensBalancesByAddress(dto *dto.GetAddressTokensByPathDTO) (tokens map[string]float64, err error) {
	result := map[string]float64{}

	addrPath, err := mhda.ParseNSS(dto.MhdaPath)
	if err != nil {
		return nil, err
	}

	addr := s.meta.GetAddress(addrPath.NSS())

	if addr == nil {
		return nil, err
	}

	ctx := types.NewRPCContext(addr.MHDA().Chain().Key(), addr.NodeIndex(), addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	balance, err := provider.GetBalance(ctx)

	if err != nil {
		return nil, err
	}

	result[provider.Currency()] = balance

	addressLinkedTokenConfigs, err := s.meta.GetAddressTokens(addr.Index())

	if len(addressLinkedTokenConfigs) > 0 {
		for _, tokenConfig := range addressLinkedTokenConfigs {
			humanBalance, err := provider.GetTokenBalance(ctx, tokenConfig.Contract(), tokenConfig.Decimals())
			if err != nil {
				return nil, err
			}
			floatBalance, _ := humanBalance.Float64()

			result[tokenConfig.Symbol()] = floatBalance
			// Show only non-zero balances
			/* if floatBalance != 0 {
				result[tokenConfig.Symbol()] = floatBalance
			}*/
		}
	}

	return result, nil
}

// SetAddressW3 Mark address as available for web3 iterations with WebExtension
func (s *Wallet) SetAddressW3(dto *dto.SetAddressW3DTO) error {
	addrPath, err := mhda.ParseNSS(dto.MhdaPath)
	if err != nil {
		return err
	}

	addr := s.meta.GetAddress(addrPath.NSS())

	if addr.IsW3() {
		return errAddrW3AlreadyPermitted
	}

	addr.SetW3()
	// TODO: Add save to meta

	return nil
}

// UnsetAddressW3 Unmark address is not available for web3 iterations.
// If address already have iterations and was delivered to web extension,
// returns error.
func (s *Wallet) UnsetAddressW3(dto *dto.SetAddressW3DTO) error {
	addrPath, err := mhda.ParseNSS(dto.MhdaPath)
	if err != nil {
		return err
	}

	addr := s.meta.GetAddress(addrPath.NSS())

	if !addr.IsW3() {
		return errAddrW3NotPermitted
	}

	addr.UnsetW3()
	// TODO: Check previous used as web3

	return nil
}
