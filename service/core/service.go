// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

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
	Version() string
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
