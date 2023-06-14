package types

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	ChargeExternal = 0
	ChargeInternal = 1
)

var (
	rxAccountPath = regexp.MustCompile(`^m/44[Hh']/([0-9]+)[Hh']/([0-9]+)[Hh']?$`)
	rxAddressPath = regexp.MustCompile(`^m/44[Hh']/([0-9]+)[Hh']/([0-9]+)[Hh']/(0|1)/([0-9]+)([Hh'])?$`)
)

type DerivationPath struct {
	coin    CoinType
	account AccountIndex
	charge  ChargeType
	index   AddressIndex
}

func CreatePath(
	coin CoinType,
	account AccountIndex,
	charge ChargeType,
	index AddressIndex,
) (*DerivationPath, error) {
	if !IsCoinExists(coin) {
		return nil, errors.New("coin is not exists in SLIP-44 list")
	}
	if charge > 1 {
		return nil, errors.New("charge can be 0 or 1")
	}
	return &DerivationPath{
		coin:    coin,
		account: account,
		charge:  charge,
		index:   index,
	}, nil
}

func (p *DerivationPath) Coin() CoinType {
	return p.coin
}

func (p *DerivationPath) Account() AccountIndex {
	return p.account
}

func (p *DerivationPath) Charge() ChargeType {
	return p.charge
}

func (p *DerivationPath) AddressIndex() AddressIndex {
	return p.index
}

func (p *DerivationPath) IsHardenedAddress() bool {
	return p.index.IsHardened
}

func (p *DerivationPath) String() string {
	var format = "m/44'/%d'/%d'/%d/%d"
	if p.index.IsHardened {
		format += `'`
	}
	return fmt.Sprintf(format, p.coin, p.account, p.charge, p.index.Index)
}

func ParsePath(path string) (*DerivationPath, error) {
	var isAddressHardened = false
	matches := rxAddressPath.FindStringSubmatch(path)
	if len(matches) < 5 {
		return nil, errors.New(fmt.Sprintf("cannot parse path: %s", path))
	}
	coinType, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, err
	}
	accountIndex, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return nil, err
	}
	chargeType, err := strconv.ParseUint(matches[3], 10, 32)
	if err != nil {
		return nil, err
	}
	addressIndex, err := strconv.ParseUint(matches[4], 10, 32)
	if err != nil {
		return nil, err
	}
	if len(matches) == 6 && matches[5] != "" {
		isAddressHardened = true
	}

	return CreatePath(
		CoinType(coinType),
		AccountIndex(accountIndex),
		ChargeType(chargeType),
		AddressIndex{
			Index:      uint32(addressIndex),
			IsHardened: isAddressHardened,
		},
	)
}

func Validate(path string) bool {
	_, err := ParsePath(path)
	return err == nil
}

type AccountDerivationPath struct {
	coin    CoinType
	account AccountIndex
}

func CreateAccountPath(
	coin CoinType,
	account AccountIndex,
) (*AccountDerivationPath, error) {
	if !IsCoinExists(coin) {
		return nil, errors.New("coin is not exists in SLIP-44 list")
	}
	return &AccountDerivationPath{
		coin:    coin,
		account: account,
	}, nil
}

func (p *AccountDerivationPath) Coin() CoinType {
	return p.coin
}

func (p *AccountDerivationPath) Account() AccountIndex {
	return p.account
}

func (p *AccountDerivationPath) String() string {
	return fmt.Sprintf("m/44'/%d'/%d'", p.coin, p.account)
}

func ParseAccountPath(path string) (*AccountDerivationPath, error) {
	matches := rxAddressPath.FindStringSubmatch(path)
	if len(matches) < 5 {
		return nil, errors.New(fmt.Sprintf("cannot parse path: %s", path))
	}
	coinType, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return nil, err
	}
	accountIndex, err := strconv.ParseUint(matches[2], 10, 32)

	return CreateAccountPath(
		CoinType(coinType),
		AccountIndex(accountIndex),
	)
}
