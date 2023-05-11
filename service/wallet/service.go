package wallet

import (
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
)

var walletInstance WalletAdapter = &Wallet{}

func API() WalletAdapter {
	return walletInstance
}

type WalletAdapter interface {
	Init(dto *dto.InitWalletDTO) (string, error)
	GenerateMnemonic(dto *dto.GenerateMnemonicDTO) (string, error)

	SendTokens(dto *dto.SendTokensDTO) (txId string, err error)
	GetTxReceipt(dto *dto.GetTxReceiptDTO) (map[string]interface{}, error)
	// GetAllAccounts() []types.AccountIndex
	GetAccountsByCoin(dto *dto.GetAccountsByCoinDTO) []types.AccountIndex

	GetAccountLabels() map[uint32]string
	GetAddressLabels() map[uint32]string
	AddLabel(dto *dto.AddLabelDTO) (uint32, error)
	RemoveLabel(dto *dto.RemoveLabelDTO) error

	AllRPC(dto *dto.GetRPCListByCoinDTO) map[uint32]*types.RPC
	AddRPC(dto *dto.AddRPCDTO) error
	RemoveRPC(dto *dto.RemoveRPCDTO) error

	AddAddresses(dto *dto.AddAddressesDTO) ([]*resp.AddressResponse, error)
	GetAddressesByAccount(dto *dto.GetAddressesByAccountDTO) []*resp.AddressResponse
	// GetAllAddresses() []*types.AddressResponse
	GetTokensBalancesByPath(dto *dto.GetAddressTokensByPathDTO) (map[string]float64, error)

	// Tokens
	UpsertToken(dto *dto.AddTokenDTO) error
	GetBaseCurrency(dto *dto.GetTokensByNetworkDTO) (*resp.BaseCurrency, error)
	GetTokensByPath(dto *dto.GetAddressTokensByPathDTO) (*resp.AddressTokensListResponse, error)
	GetAllTokensByNetwork(dto *dto.GetTokensByNetworkDTO) (*resp.AddressTokensListResponse, error)
	GetToken(dto *dto.GetTokenDTO) (*resp.TokenConfig, error)

	// nodes
	AccountLinkRPCSet(dto *dto.SetRPCLinkedAccountDTO) error
	RemoveAccountLinkRPC(dto *dto.RemoveRPCLinkedAccountDTO) error
	GetRPCLinkedAccountCount(dto *dto.GetRPCLinkedAccountCountDTO) int
	GetRPCInfo(dto *dto.GetRPCInfoDTO) (map[string]interface{}, error)

	// AirGap
	ExportMeta() (*resp.AirGapMessageResponse, error)
	ExportMetaDebug() ([]byte, error)
}
