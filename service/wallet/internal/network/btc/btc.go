package btc

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/censync/soikawallet/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// base58.Encode()

type BTC struct {
	*types.BaseNetwork
}

func NewBTC(baseNetwork *types.BaseNetwork) *BTC {
	return &BTC{BaseNetwork: baseNetwork}
}

func (b *BTC) Address(pub *ecdsa.PublicKey) string {
	serializedAddr := crypto.FromECDSAPub(pub)
	addr, err := btcutil.NewAddressPubKey(serializedAddr, &chaincfg.MainNetParams)
	if err != nil {
		return `undefined`
	}
	return addr.AddressPubKeyHash().String()
}
