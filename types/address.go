package types

import mhda "github.com/censync/go-mhda"

/*
type CoinType uint32

type AccountIndex uint32

type ChargeType uint8 // 0 or 1

type AddressIndex struct {
	Index      uint32
	IsHardened bool
}

func (i *AddressIndex) MarshalJSON() ([]byte, error) {
	result := ""
	if i.IsHardened {
		result = fmt.Sprintf(
			`"%d'"`,
			i.Index,
		)
	} else {
		result = fmt.Sprintf(
			`"%d"`,
			i.Index,
		)
	}

	return []byte(result), nil
}
*/

type NodeIndex struct {
	mhda.ChainKey
	Index uint32
}

// meta

type TokenIndex struct {
	mhda.ChainKey
	Contract string
}
