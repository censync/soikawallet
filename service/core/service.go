package core

import (
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
)

var walletInstance CoreAdapter = &Wallet{}

func API() CoreAdapter {
	return walletInstance
}

type Core interface {
	Init(dto *dto.InitWalletDTO) (string, error)
	GenerateMnemonic(dto *dto.GenerateMnemonicDTO) (string, error)
	GetTokensBalancesByAddress(dto *dto.GetAddressTokensByPathDTO) (map[string]float64, error)

	GetAccountsByNetwork(dto *dto.GetAccountsByNetworkDTO) []*resp.AccountResponse

	// Chain operations
	GetAllEvmW3Chains() []*resp.ChainInfo
}

type Address interface {
	SetAddressW3(dto *dto.SetAddressW3DTO) error
	UnsetAddressW3(dto *dto.SetAddressW3DTO) error
	AddAddresses(dto *dto.AddAddressesDTO) ([]*resp.AddressResponse, error)
	GetAddressesByAccount(dto *dto.GetAddressesByAccountDTO) []*resp.AddressResponse
	GetTokensByPath(dto *dto.GetAddressTokensByPathDTO) (*resp.AddressTokensListResponse, error)
}

type AirGap interface {
	ProcessAirGapMessage(dto *dto.AirGapMessageDTO) (string, error)
	ExportMeta() (*resp.AirGapMessage, error)
	ExportMetaDebug() ([]byte, error)
}

type Operation interface {
	GetAllowance(dto *dto.GetTokenAllowanceDTO) (uint64, error)
	ApproveTokens(dto *dto.SendTokensDTO) (string, error)
	SendTokens(dto *dto.SendTokensDTO) (string, error)
	SendTokensPrepare(dto *dto.SendTokensDTO) (*resp.AirGapMessage, error)
	GetTxReceipt(dto *dto.GetTxReceiptDTO) (map[string]interface{}, error)
}

type Currencies interface {
	UpdateFiatCurrencies() map[string]float64
	GetFiatCurrency(dto *dto.GetFiatCurrencyDTO) (float64, string, string)
	GetGasCalculatorConfig(dto *dto.GetGasCalculatorConfigDTO) (*resp.CalculatorConfig, error)
}

type Meta interface {
	Label
	RPC
}

type Label interface {
	AddLabel(dto *dto.AddLabelDTO) (uint32, error)
	RemoveLabel(dto *dto.RemoveLabelDTO) error
	SetLabelLink(dto *dto.SetLabelLinkDTO) error
	RemoveLabelLink(dto *dto.RemoveLabelLinkDTO) error
	GetAccountLabels() map[uint32]string
	GetAddressLabels() map[uint32]string
}

type RPC interface {
	AllRPC(dto *dto.GetRPCListByNetworkDTO) map[uint32]*types.RPC
	AddRPC(dto *dto.AddRPCDTO) error
	RemoveRPC(dto *dto.RemoveRPCDTO) error

	AccountLinkRPCSet(dto *dto.SetRPCLinkedAccountDTO) error
	RemoveAccountLinkRPC(dto *dto.RemoveRPCLinkedAccountDTO) error
	GetRPCLinkedAccountCount(dto *dto.GetRPCLinkedAccountCountDTO) int
	GetRPCInfo(dto *dto.GetRPCInfoDTO) (map[string]interface{}, error)
	ExecuteRPC(dto *dto.ExecuteRPCRequestDTO) ([]byte, error)
}

type Token interface {
	UpsertToken(dto *dto.AddTokenDTO) error
	GetBaseCurrency(dto *dto.GetTokensByNetworkDTO) (*resp.BaseCurrency, error)
	GetAllTokensByNetwork(dto *dto.GetTokensByNetworkDTO) (*resp.AddressTokensListResponse, error)
	GetToken(dto *dto.GetTokenDTO) (*resp.TokenConfig, error)
}

type CoreAdapter interface {
	Core
	Address
	AirGap
	Operation
	Currencies
	Token
	Meta
}
