package dto

type InitWalletDTO struct {
	Mnemonic        string
	Passphrase      string
	SkipPrefixCheck bool
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

type GetAccountsByCoinDTO struct {
	CoinType uint32
}

type GetAddressesByAccountDTO struct {
	CoinType     uint32
	AccountIndex uint32
}

type GetAddressTokensByPathDTO struct {
	DerivationPath string
}

type GetAddressTokensBalanceByPathDTO struct {
	DerivationPath string
}

type SendTokensDTO struct {
	DerivationPath string
	To             string
	Count          string
}

type GetTxReceiptDTO struct {
	DerivationPath string
	Hash           string
}

type GetRPCListByIndexDTO struct {
	Index    uint32
	CoinType uint32
}

type GetRPCListByCoinDTO struct {
	CoinType uint32
}

type SetRPCLinkedAccountDTO struct {
	CoinType     uint32
	AccountIndex uint32
	NodeIndex    uint32
}

type RemoveRPCLinkedAccountDTO struct {
	CoinType     uint32
	AccountIndex uint32
}

type GetRPCLinkedAccountCountDTO struct {
	CoinType  uint32
	NodeIndex uint32
}

type AddRPCDTO struct {
	CoinType uint32
	Title    string
	Endpoint string
}

type RemoveRPCDTO struct {
	CoinType uint32
	Index    uint32
}

type GetRPCInfoDTO struct {
	CoinType  uint32
	NodeIndex uint32
}

type FlushKeysDTO struct {
	Force bool
}
