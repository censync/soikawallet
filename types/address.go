package types

type CoinType uint32

type AccountIndex uint32

type ChargeType uint8 // 0 or 1

type AddressIndex struct {
	Index      uint32
	IsHardened bool
}

type NodeIndex struct {
	CoinType
	Index uint32
}

// meta

type TokenIndex struct {
	CoinType
	Index uint32
}
