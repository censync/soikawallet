package types

/*
const (
	ChargeExternal = 0
	ChargeInternal = 1
)

var (
	//rxAccountPath = regexp.MustCompile(`^m/44[Hh']/([0-9]+)[Hh']/([0-9]+)[Hh']?$`)
	rxAddressPath = regexp.MustCompile(`^m/44[Hh']/([0-9]+)[Hh']/([0-9]+)[Hh']/(0|1)/([0-9]+)([Hh'])?$`)
)


type MhdaPath struct {
	coinType CoinType
	account  AccountIndex
	charge   ChargeType
	index    AddressIndex
}

func CreatePath(
	coinType CoinType,
	account AccountIndex,
	charge ChargeType,
	index AddressIndex,
) (*MhdaPath, error) {
	if !IsNetworkExists(coinType) {
		return nil, errors.New("coinType is not exists in SLIP-44 list")
	}
	if charge > 1 {
		return nil, errors.New("charge can be 0 or 1")
	}
	return &MhdaPath{
		coinType: coinType,
		account:  account,
		charge:   charge,
		index:    index,
	}, nil
}

func (p *MhdaPath) CoinType() CoinType {
	return p.coinType
}

func (p *MhdaPath) Account() AccountIndex {
	return p.account
}

func (p *MhdaPath) Charge() ChargeType {
	return p.charge
}

func (p *MhdaPath) AddressIndex() AddressIndex {
	return p.index
}

func (p *MhdaPath) IsHardenedAddress() bool {
	return p.index.IsHardened
}

func (p *MhdaPath) String() string {
	var format = "m/44'/%d'/%d'/%d/%d"
	if p.index.IsHardened {
		format += `'`
	}
	return fmt.Sprintf(format, p.coinType, p.account, p.charge, p.index.Index)
}

func ParsePath(path string) (*MhdaPath, error) {
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

type AccountDerivationPath struct {
	network CoinType
	account AccountIndex
}

func CreateAccountPath(
	coinType CoinType,
	account AccountIndex,
) (*AccountDerivationPath, error) {
	if !IsNetworkExists(coinType) {
		return nil, errors.New("coinType is not exists in SLIP-44 list")
	}
	return &AccountDerivationPath{
		network: coinType,
		account: account,
	}, nil
}

func (p *AccountDerivationPath) CoinType() CoinType {
	return p.network
}

func (p *AccountDerivationPath) Account() AccountIndex {
	return p.account
}

func (p *AccountDerivationPath) String() string {
	return fmt.Sprintf("m/44'/%d'/%d'", p.network, p.account)
}
*/
