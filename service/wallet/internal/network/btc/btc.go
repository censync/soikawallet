package btc

import (
	"crypto/ecdsa"
	"errors"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/censync/soikawallet/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// base58.Encode()

type Bitcoin struct {
	*types.BaseNetwork
	//clients map[uint32]*rpcclient.Client
}

func NewBTC(baseNetwork *types.BaseNetwork) *Bitcoin {
	return &Bitcoin{BaseNetwork: baseNetwork /*, clients: map[uint32]*rpcclient.Client{}*/}
}

// https://en.bitcoin.it/wiki/Original_Bitcoin_client/API_calls_list
/*
func (b *Bitcoin) getClient(nodeId uint32) (*rpcclient.Client, error) {
	var err error
	if b.clients[nodeId] != nil {
		return b.clients[nodeId], nil
	} else {
		//b.client[nodeId], err = ethclient.Dial(e.DefaultRPC().Endpoint())
		connCfg := &rpcclient.ConnConfig{
			Host:         b.DefaultRPC().Endpoint(),
			HTTPPostMode: true,
			DisableTLS:   false,
		}
		b.clients[nodeId], err = rpcclient.New(connCfg, nil)
		return b.clients[nodeId], err
	}
	return nil, errors.New("not implemented yet")
}
*/

func (b *Bitcoin) Address(pub *ecdsa.PublicKey) string {
	serializedAddr := crypto.CompressPubkey(pub)
	addr, err := btcutil.NewAddressPubKey(serializedAddr, &chaincfg.MainNetParams)
	if err != nil {
		return `undefined`
	}
	return addr.AddressPubKeyHash().String()
}

func (b *Bitcoin) GetBalance(ctx *types.RPCContext) (float64, error) {
	/*client, err := b.getClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}
	balance, err := client.GetBalance(ctx.CurrentAccount())
	if err != nil {
		return 0, err
	}
	return balance.ToBTC(), nil*/
	return 0, errors.New("not implemented yet")
}

func (b *Bitcoin) GetTokenBalance(ctx *types.RPCContext, contract string, decimals int) (*big.Float, error) {
	return nil, nil
}

func (b *Bitcoin) GetGasConfig(ctx *types.RPCContext, args ...interface{}) (map[string]uint64, error) {
	return nil, nil
}
