package wallet

import (
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/wallet/internal/airgap"
	"github.com/censync/soikawallet/service/wallet/internal/network"
	"github.com/censync/soikawallet/service/wallet/meta"
	"github.com/censync/soikawallet/types"
	"github.com/google/uuid"
	"strings"
)

const (
	hardenedKeyStart = uint32(0x80000000) // 2^31
)

type Wallet struct {
	instanceId uuid.UUID
	bip44Key   *hdkeychain.ExtendedKey
	addresses  map[string]*address
	meta       *meta.Meta
}

func (s *Wallet) getNetworkProvider(ctx *types.RPCContext) (types.NetworkAdapter, error) {
	return network.WithContext(ctx)
}

func (s *Wallet) getRPCProvider(coinType types.CoinType) types.RPCAdapter {
	return network.Get(coinType)
}

func (s *Wallet) SendTokens(dto *dto.SendTokensDTO) (txId string, err error) {
	addressPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return ``, err
	}

	addr, err := s.address(addressPath)

	if err != nil {
		return ``, err
	}

	if addr.key == nil {
		return ``, nil
	}

	ctx := types.NewRPCContext(addr.CoinType(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return "", err
	}

	return provider.TxSendBase(ctx, ``, addr.key.Get())
}

func (s *Wallet) GetTxReceipt(dto *dto.GetTxReceiptDTO) (map[string]interface{}, error) {
	dto.DerivationPath = strings.TrimSpace(dto.DerivationPath)
	dto.Hash = strings.TrimSpace(dto.Hash)

	addressPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return nil, err
	}

	addr, err := s.address(addressPath)

	if err != nil {
		return nil, err
	}
	ctx := types.NewRPCContext(addr.CoinType(), addr.nodeIndex)
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}
	return provider.TxGetReceipt(ctx, dto.Hash)
}

/*
func (s *Wallet) GetAllCoins() []types.CoinType {
	var coins []types.CoinType
	for coin := range s.addresses {
		coins = append(coins, coin)
	}
	return coins
}

func (s *Wallet) GetAllAccounts() []types.AccountIndex {
	var accounts []types.AccountIndex
	for coin := range s.coins {
		for account := range s.coins[coin] {
			accounts = append(accounts, account)
		}
	}
	return accounts
}*/

func (s *Wallet) GetAccountsByCoin(dto *dto.GetAccountsByCoinDTO) []types.AccountIndex {
	accountsIndex := map[types.AccountIndex]bool{}

	for _, addr := range s.addresses {
		if addr.Path().Coin() == types.CoinType(dto.CoinType) {
			accountsIndex[addr.Path().Account()] = true
		}
	}

	accounts := make([]types.AccountIndex, 0)

	for accountIndex := range accountsIndex {
		accounts = append(accounts, accountIndex)
	}

	return accounts
}

func (s *Wallet) GetInstanceId() string {
	return s.instanceId.String()
}

func (s *Wallet) FlushKeys(dto *dto.FlushKeysDTO) {
	s.bip44Key = nil
	for key := range s.addresses {
		if dto.Force || !s.addresses[key].staticKey {
			s.addresses[key].key.Free()
			s.addresses[key].key = nil
		}
	}
}

func (s *Wallet) ExportMeta() (*resp.AirGapMessageResponse, error) {
	data, err := s.MarshalJSON()
	if err != nil {
		return nil, err
	}
	chunks, err := airgap.NewChunks(data, 192)
	if err != nil {
		return nil, err
	}
	return &resp.AirGapMessageResponse{
		Chunks: chunks.ChunksBase64(),
	}, nil
}

func (s *Wallet) MarshalJSON() ([]byte, error) {
	var strPaths []string
	// TODO: Add internal
	addresses := s.GetAllAddresses()
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
