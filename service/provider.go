package service

import (
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/wallet"
	"github.com/censync/soikawallet/types"
)

var walletInstance WalletAdapter = &wallet.Wallet{}

func API() WalletAdapter {
	return walletInstance
}

type WalletAdapter interface {
	Init(dto *dto.InitWalletDTO) error
	GenerateMnemonic(dto *dto.GenerateMnemonicDTO) (string, error)

	SendTokens(dto *dto.SendTokensDTO) (txId string, err error)
	GetTxReceipt(dto *dto.GetTxReceiptDTO) (map[string]interface{}, error)
	// GetAllAccounts() []types.AccountIndex
	GetAccountsByCoin(dto *dto.GetAccountsByCoinDTO) []types.AccountIndex
	GetInstanceId() string

	GetAccountLabels() map[uint32]string
	GetAddressLabels() map[uint32]string
	AddLabel(dto *dto.AddLabelDTO) (uint32, error)
	RemoveLabel(dto *dto.RemoveLabelDTO) error

	AllRPC(dto *dto.GetRPCListByCoinDTO) map[uint32]*types.RPC
	AddRPC(dto *dto.AddRPCDTO) (uint32, error)
	RemoveRPC(dto *dto.RemoveRPCDTO) error

	AddAddresses(dto *dto.AddAddressesDTO) ([]*responses.AddressResponse, error)
	GetAddressesByAccount(dto *dto.GetAddressesByAccountDTO) []*responses.AddressResponse
	// GetAllAddresses() []*types.AddressResponse
	GetAddressTokensByPath(dto *dto.GetAddressTokensByPathDTO) (tokens map[string]float64, err error)

	// Nodes
	AccountLinkRPCSet(dto *dto.SetRPCLinkedAccountDTO) error
	RemoveAccountLinkRPC(dto *dto.RemoveRPCLinkedAccountDTO) error
	GetRPCLinkedAccountCount(dto *dto.GetRPCLinkedAccountCountDTO) int
	GetRPCInfo(dto *dto.GetRPCInfoDTO) (map[string]interface{}, error)
}
