package btc

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/censync/soikawallet/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// base58.Encode()

type Bitcoin struct {
	*types.BaseNetwork
}

func NewBTC(baseNetwork *types.BaseNetwork) *Bitcoin {
	return &Bitcoin{BaseNetwork: baseNetwork}
}

func (b *Bitcoin) Address(pub *ecdsa.PublicKey) string {
	serializedAddr := crypto.CompressPubkey(pub)
	addr, err := btcutil.NewAddressPubKey(serializedAddr, &chaincfg.MainNetParams)
	if err != nil {
		return `undefined`
	}
	return addr.AddressPubKeyHash().String()
}

func (b *Bitcoin) GetBalance(ctx *types.RPCContext) (float64, error) {
	return 0, nil
}

func (b *Bitcoin) GetTokenBalance(ctx *types.RPCContext, contract string, decimals int) (*big.Float, error) {
	return nil, nil
}

func (b *Bitcoin) GetGasConfig(ctx *types.RPCContext, args ...interface{}) (map[string]uint64, error) {
	return nil, nil
}
