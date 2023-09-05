package wallet

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	airgap "github.com/censync/go-airgap"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/wallet/internal/network"
	"github.com/censync/soikawallet/service/wallet/meta"
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/currencies"
	"github.com/censync/soikawallet/util/seed"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

const (
	hardenedKeyStart = uint32(0x80000000) // 2^31
	fiatTitle        = "USD"
	fiatSymbol       = "$"
)

type Wallet struct {
	// instanceId compressed public key for root key, used for identify device instance
	instanceId []byte
	rootKey    *hdkeychain.ExtendedKey
	// addresses key: mhda_nss of address => val: address
	meta           *meta.Meta
	currenciesFiat *currencies.FiatCurrencies
}

func (s *Wallet) getNetworkProvider(ctx *types.RPCContext) (types.NetworkAdapter, error) {
	return network.WithContext(ctx)
}

func (s *Wallet) getRPCProvider(chainKey mhda.ChainKey) types.RPCAdapter {
	return network.Get(chainKey)
}

// Init initializes static instance of wallet with mnemonic and optional passphrase.
// If result is successful, will be returned base58 encoded compressed root public key.
func (s *Wallet) Init(dto *dto.InitWalletDTO) (string, error) {
	var err error
	dto.Mnemonic = strings.TrimSpace(dto.Mnemonic)
	dto.Passphrase = strings.TrimSpace(dto.Passphrase)

	// Check for singleton
	if s.instanceId != nil {
		return "", errors.New("wallet already initialized")
	}

	// SkipMnemonicCheck flag used only for testing vectors
	if !dto.SkipMnemonicCheck {
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

	*s = Wallet{
		instanceId:     masterPubKey.SerializeCompressed(),
		rootKey:        masterKey,
		meta:           meta.InitMeta(),
		currenciesFiat: currencies.NewFiatCurrencies(fiatTitle, fiatSymbol),
	}
	return s.getInstanceId(), nil
}

func (s *Wallet) getInstanceId() string {
	return base58.Encode(s.instanceId)
}

func (s *Wallet) GetTxReceipt(dto *dto.GetTxReceiptDTO) (map[string]interface{}, error) {
	ctx := types.NewRPCContext(mhda.ChainKey(dto.ChainKey), dto.NodeIndex)
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}
	return provider.TxGetReceipt(ctx, dto.Hash)
}

func (s *Wallet) GetAccountsByNetwork(dto *dto.GetAccountsByNetworkDTO) []*resp.AccountResponse {
	accountsIndex := map[mhda.AccountIndex]bool{}

	for _, addr := range s.meta.Addresses() {
		if addr.MHDA().Chain().Key() == dto.ChainKey {
			accountsIndex[addr.MHDA().DerivationPath().Account()] = true
		}
	}

	accounts := make([]*resp.AccountResponse, 0)

	for accountIndex := range accountsIndex {
		// Deprecated
		// 	accountPath, err := types.CreateAccountPath(types.CoinType(dto.ChainKey), accountIndex)
		// 	if err != nil {
		// 		continue
		//	}

		// Fix account path
		accounts = append(accounts, &resp.AccountResponse{
			// Path:        accountPath.String(),
			ChainKey: dto.ChainKey,
			Account:  accountIndex,
			// Label:       s.meta.GetAccountLabel(accountPath.String()),
		})
	}

	return accounts
}

func (s *Wallet) FlushKeys(dto *dto.FlushKeysDTO) {
	s.rootKey = nil
	/*for key := range s.addresses {
		if dto.Force || !s.addresses[key].isKeyDelivered {
			s.addresses[key].key.Free()
			s.addresses[key].key = nil
		}
	}*/
}

func (s *Wallet) ExportMeta() (*resp.AirGapMessage, error) {
	data, err := s.meta.MarshalJSON()
	if err != nil {
		return nil, err
	}
	airgapMsg := airgap.NewAirGap(airgap.VersionDefault, s.instanceId).
		CreateMessage().
		AddOperation(types.OpMetaWallet, data)
	chunks, err := airgapMsg.MarshalB64Chunks()
	if err != nil {
		return nil, err
	}
	return &resp.AirGapMessage{
		Chunks: chunks,
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
