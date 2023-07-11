package types

const (
	TokenBase    = TokenStandard(1)
	TokenERC20   = TokenStandard(20)
	TokenERC721  = TokenStandard(21)
	TokenERC777  = TokenStandard(22)
	TokenERC1155 = TokenStandard(23)
	TokenERC4626 = TokenStandard(24)
	TokenTRC10   = TokenStandard(40)
	TokenTRC20   = TokenStandard(41)
	TokenBEP20   = TokenStandard(50)

	ContractBase = `__base`
	ContractZero = `0x0000000000000000000000000000000000000000`
)

var (
	registeredTokenStandards = map[TokenStandard]string{
		TokenBase:    `Base`,
		TokenERC20:   `ERC-20`,
		TokenERC721:  `ERC-771`,
		TokenERC777:  `ERC-777`,
		TokenERC1155: `ERC-1155`,
		TokenERC4626: `ERC-4626`,
		TokenTRC10:   `TRC-10`,
		TokenTRC20:   `TRC-20`,
	}
	activesTokenStandards = map[NetworkType][]TokenStandard{
		Ethereum: {TokenERC20, TokenERC721, TokenERC1155},
		Tron:     {TokenTRC20, TokenTRC10},
		Polygon:  {TokenERC20, TokenERC721, TokenERC1155},
	}
	registeredTokenStandardNames = map[NetworkType][]string{}
	registeredTokenIndexes       = map[string]TokenStandard{}
)

type TokenStandard uint8

func init() {
	for tokenStandard, tokenStandardName := range registeredTokenStandards {
		registeredTokenIndexes[tokenStandardName] = tokenStandard
	}

	for networkType, activeTokenStandards := range activesTokenStandards {
		names := make([]string, 0)
		for _, tokenStandard := range activeTokenStandards {
			names = append(names, registeredTokenStandards[tokenStandard])
		}
		//sort.Strings(names)
		registeredTokenStandardNames[networkType] = names
	}
}

func GetTokenStandardNames(networkType NetworkType) []string {
	return registeredTokenStandardNames[networkType]
}

func GetTokenStandards(networkType NetworkType) []TokenStandard {
	return activesTokenStandards[networkType]
}

func GetTokenStandByName(str string) TokenStandard {
	if tokenStandard, ok := registeredTokenIndexes[str]; ok {
		return tokenStandard
	} else {
		return 0
	}
}
