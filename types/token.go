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
	activesTokenStandards = map[CoinType][]TokenStandard{
		Ethereum: {TokenERC20, TokenERC721, TokenERC1155},
		Tron:     {TokenTRC20, TokenTRC10},
		Polygon:  {TokenERC20, TokenERC721, TokenERC1155},
	}
	registeredTokenStandardNames = map[CoinType][]string{}
	registeredTokenIndexes       = map[string]TokenStandard{}
)

type TokenStandard uint8

func init() {
	for tokenStandard, tokenStandardName := range registeredTokenStandards {
		registeredTokenIndexes[tokenStandardName] = tokenStandard
	}

	for coinType, activeTokenStandards := range activesTokenStandards {
		names := make([]string, 0)
		for _, tokenStandard := range activeTokenStandards {
			names = append(names, registeredTokenStandards[tokenStandard])
		}
		//sort.Strings(names)
		registeredTokenStandardNames[coinType] = names
	}
}

func GetTokenStandardNames(coinType CoinType) []string {
	return registeredTokenStandardNames[coinType]
}

func GetTokenStandards(coinType CoinType) []TokenStandard {
	return activesTokenStandards[coinType]
}

func GetTokenStandByName(str string) TokenStandard {
	if tokenStandard, ok := registeredTokenIndexes[str]; ok {
		return tokenStandard
	} else {
		return 0
	}
}
