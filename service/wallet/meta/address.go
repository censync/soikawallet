package meta

import (
	"crypto/ecdsa"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/wallet/protected_key"
)

const (
	flagDisabled uint8 = 1 << iota
	flagDerived
	flagW3Enabled
	flagW3Derived
)

var subIndex = aIndex(1)

// aIndex is address index for internal indexing
// max 4294967295  addresses
type aIndex uint32

type Address struct {
	path mhda.MHDA
	key  *protected_key.ProtectedKey
	pub  *ecdsa.PublicKey
	addr string

	// TODO: Move to meta
	nodeIndex uint32

	isKeyDelivered bool

	flags uint8

	// subIndex is internal numeric serial index for binding values
	subIndex aIndex

	// level provides nesting level for bip derivation path
	level uint8 // 0 - root, 1 - account, 2 - charge, 3 - index
}

func NewAddress(path mhda.MHDA, key *protected_key.ProtectedKey, pub *ecdsa.PublicKey, addr string) *Address {
	subIndex++
	return &Address{
		path:           path,
		key:            key,
		pub:            pub,
		addr:           addr,
		nodeIndex:      0,
		isKeyDelivered: false,
		flags:          0,
		subIndex:       subIndex,
		level:          0,
	}
}

func (a *Address) Address() string {
	return a.addr
}

func (a *Address) MHDA() mhda.MHDA {
	return a.path
}

func (a *Address) Index() aIndex {
	return a.subIndex
}

func (a *Address) NodeIndex() uint32 {
	return a.nodeIndex
}

func (a *Address) Key() *protected_key.ProtectedKey {
	return a.key
}

func (a *Address) DerivationPath() *mhda.DerivationPath {
	return a.path.DerivationPath()
}

func (a *Address) IsExternal() bool {
	return a.path.DerivationPath().Charge() == mhda.ChargeExternal
}

func (a *Address) AddressIndex() mhda.AddressIndex {
	return a.path.DerivationPath().AddressIndex()
}
func (a *Address) IsHardenedAddress() bool {
	return a.path.DerivationPath().IsHardenedAddress()
}

func (a *Address) Network() mhda.CoinType {
	return a.path.DerivationPath().Coin()
}

func (a *Address) Account() mhda.AccountIndex {
	return a.path.DerivationPath().Account()
}

func (a *Address) IsW3() bool {
	return a.flags&flagW3Enabled != 0
}

func (a *Address) SetW3() {
	a.flags = a.flags | flagW3Enabled
}

func (a *Address) UnsetW3() {
	a.flags = a.flags &^ flagW3Enabled
}
