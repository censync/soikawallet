package wallet

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/wallet/internal/airgap"
	"github.com/censync/soikawallet/service/wallet/internal/network"
	"github.com/censync/soikawallet/service/wallet/meta"
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/util/seed"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

const (
	hardenedKeyStart = uint32(0x80000000) // 2^31
)

type Wallet struct {
	instanceId []byte
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

func (s *Wallet) Init(dto *dto.InitWalletDTO) (string, error) {
	var err error
	dto.Mnemonic = strings.TrimSpace(dto.Mnemonic)
	dto.Passphrase = strings.TrimSpace(dto.Passphrase)

	if !dto.SkipPrefixCheck {
		err = seed.Check(dto.Mnemonic)
	}

	if err != nil {
		return "", err
	}

	rootSeed := pbkdf2.Key([]byte(dto.Mnemonic), []byte("mnemonic"+dto.Passphrase), 2048, 64, sha512.New)

	masterKey, err := generateKeyFromSeed(&rootSeed)

	if err != nil {
		return "", errors.New("cannot initialize master key")
	}

	masterPubKey, err := masterKey.ECPubKey()

	if err != nil {
		return "", errors.New("cannot initialize master pub key")
	}

	bip44Key, err := masterKey.Derive(hardenedKeyStart + 44)
	if err != nil {
		return "", errors.New("cannot initialize BIP-44 key")
	}
	masterPubKey.SerializeCompressed()
	*s = Wallet{
		instanceId: masterPubKey.SerializeCompressed(), //(masterPubKey.SerializeCompressed()),
		bip44Key:   bip44Key,
		addresses:  map[string]*address{},
		meta:       meta.InitMeta(),
	}
	return s.getInstanceId(), nil
}

func (s *Wallet) getInstanceId() string {
	return base58.Encode(s.instanceId)
}

func (s *Wallet) SendTokens(dto *dto.SendTokensDTO) (txId string, err error) {
	dto.To = strings.TrimSpace(dto.To)
	addressPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return ``, err
	}

	addr, err := s.address(addressPath)

	if err != nil {
		return ``, err
	}

	if addr.key == nil {
		return ``, errors.New("empty key for sign, use airgap option")
	}

	if len(dto.To) < 4 {
		return ``, errors.New("incorrect recipient address")
	}

	ctx := types.NewRPCContext(addr.CoinType(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return "", err
	}

	return provider.TxSendBase(ctx, dto.To, dto.Value, addr.key.Get())
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

func (s *Wallet) GetAccountsByCoin(dto *dto.GetAccountsByCoinDTO) []*resp.AccountResponse {
	accountsIndex := map[types.AccountIndex]bool{}

	for _, addr := range s.addresses {
		if addr.Path().Coin() == types.CoinType(dto.CoinType) {
			accountsIndex[addr.Path().Account()] = true
		}
	}

	accounts := make([]*resp.AccountResponse, 0)

	for accountIndex := range accountsIndex {
		accountPath, err := types.CreateAccountPath(types.CoinType(dto.CoinType), accountIndex)
		if err != nil {
			continue
		}
		accounts = append(accounts, &resp.AccountResponse{
			Path:     accountPath.String(),
			CoinType: accountPath.Coin(),
			Account:  accountPath.Account(),
			Label:    s.meta.GetAccountLabel(accountPath.String()),
		})
	}

	return accounts
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
	data, err := s.meta.MarshalJSON()
	if err != nil {
		return nil, err
	}
	a := airgap.Restore(s.instanceId).
		CreateMessage().
		AddOperation(types.OpMetaWallet, data)
	chunks, err := airgap.NewChunks(a.Bytes(), 192)
	if err != nil {
		return nil, err
	}
	return &resp.AirGapMessageResponse{
		Chunks: chunks.ChunksBase64(),
	}, nil
}

func (s *Wallet) ExportMetaDebug() ([]byte, error) {
	data, err := s.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return data, nil
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
