package dto

import mhda "github.com/censync/go-mhda"

type InitWalletDTO struct {
	Mnemonic          string
	Passphrase        string
	SkipMnemonicCheck bool
}

type GenerateMnemonicDTO struct {
	BitSize  int
	Language string
}

type AddAddressesDTO struct {
	MhdaPaths []string
}

type AddLabelDTO struct {
	LabelType uint8
	Title     string
}

type RemoveLabelDTO struct {
	LabelType uint8
	Index     uint32
}

type SetLabelLinkDTO struct {
	LabelType uint8
	Index     uint32
	Path      string
}

type RemoveLabelLinkDTO struct {
	LabelType uint8
	Path      string
}

type RemoveLabelLinkedAccountDTO struct {
	Index uint32
}

type SetAddressW3DTO struct {
	MhdaPath string
}

type UnsetAddressW3DTO struct {
	ChainKey mhda.ChainKey
}

type GetAccountsByNetworkDTO struct {
	ChainKey mhda.ChainKey
}

type GetAddressesByAccountDTO struct {
	ChainKey     mhda.ChainKey
	AccountIndex uint32
}

type GetAddressTokensByPathDTO struct {
	MhdaPath string
}

type GetAddressTokensBalanceByPathDTO struct {
	DerivationPath string
}

type GetGasCalculatorConfigDTO struct {
	Operation string
	MhdaPath  string
	To        string
	Value     float64
	Standard  uint8
	Contract  string
}

type GetTokenAllowanceDTO struct {
	MhdaPath string
	To       string
	Value    float64
	Standard uint8
	Contract string
}

type SendTokensDTO struct {
	MhdaPath  string
	To        string
	Value     float64
	Gas       uint64
	GasTipCap uint64
	GasFeeCap uint64
	Standard  uint8
	Contract  string
}

type GetTxReceiptDTO struct {
	ChainKey  mhda.ChainKey
	NodeIndex uint32
	Hash      string
}

type GetRPCListByIndexDTO struct {
	ChainKey mhda.ChainKey
	Index    uint32
}

type GetRPCListByNetworkDTO struct {
	ChainKey mhda.ChainKey
}

type SetRPCLinkedAccountDTO struct {
	ChainKey     mhda.ChainKey
	AccountIndex uint32
	NodeIndex    uint32
}

type RemoveRPCLinkedAccountDTO struct {
	ChainKey     mhda.ChainKey
	AccountIndex uint32
}

type GetRPCLinkedAccountCountDTO struct {
	ChainKey  mhda.ChainKey
	NodeIndex uint32
}

type AddRPCDTO struct {
	ChainKey mhda.ChainKey
	Title    string
	Endpoint string
}

type RemoveRPCDTO struct {
	ChainKey mhda.ChainKey
	Index    uint32
}

type GetRPCInfoDTO struct {
	ChainKey  mhda.ChainKey
	NodeIndex uint32
}

type FlushKeysDTO struct {
	Force bool
}

type GetFiatCurrencyDTO struct {
	ChainKey mhda.ChainKey
}
