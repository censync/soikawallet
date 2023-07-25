package dto

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
	DerivationPaths []string
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
	DerivationPath string
}

type UnsetAddressW3DTO struct {
	DerivationPath string
}

type GetAccountsByNetworkDTO struct {
	NetworkType uint32
}

type GetAddressesByAccountDTO struct {
	NetworkType  uint32
	AccountIndex uint32
}

type GetAddressTokensByPathDTO struct {
	DerivationPath string
}

type GetAddressTokensBalanceByPathDTO struct {
	DerivationPath string
}

type GetGasCalculatorConfigDTO struct {
	DerivationPath string
	To             string
	Value          float64
	Standard       uint8
	Contract       string
}

type GetTokenAllowanceDTO struct {
	DerivationPath string
	To             string
	Value          float64
	Standard       uint8
	Contract       string
}

type SendTokensDTO struct {
	DerivationPath string
	To             string
	Value          float64
	GasTipCap      uint64
	GasFeeCap      uint64
	Standard       uint8
	Contract       string
}

type GetTxReceiptDTO struct {
	NetworkType uint32
	NodeIndex   uint32
	Hash        string
}

type GetRPCListByIndexDTO struct {
	Index       uint32
	NetworkType uint32
}

type GetRPCListByNetworkDTO struct {
	NetworkType uint32
}

type SetRPCLinkedAccountDTO struct {
	NetworkType  uint32
	AccountIndex uint32
	NodeIndex    uint32
}

type RemoveRPCLinkedAccountDTO struct {
	NetworkType  uint32
	AccountIndex uint32
}

type GetRPCLinkedAccountCountDTO struct {
	NetworkType uint32
	NodeIndex   uint32
}

type AddRPCDTO struct {
	NetworkType uint32
	Title       string
	Endpoint    string
}

type RemoveRPCDTO struct {
	NetworkType uint32
	Index       uint32
}

type GetRPCInfoDTO struct {
	NetworkType uint32
	NodeIndex   uint32
}

type FlushKeysDTO struct {
	Force bool
}

type GetFiatCurrencyDTO struct {
	NetworkType uint32
}
