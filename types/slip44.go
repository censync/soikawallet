package types

import "sort"

const (
	Ethereum = CoinType(60)
	Tron     = CoinType(195)
	Polygon  = CoinType(966)
)

var (
	registeredCoins = map[string]CoinType{
		`Ethereum`: Ethereum,
		`Tron`:     Tron,
		`Polygon`:  Polygon,
	}

	registeredCoinsIndexes = map[CoinType]string{}
	registeredCoinsTypes   []CoinType
	registeredCoinsNames   []string
)

func init() {
	for name, coinType := range registeredCoins {
		registeredCoinsIndexes[coinType] = name
		registeredCoinsTypes = append(registeredCoinsTypes, coinType)
		registeredCoinsNames = append(registeredCoinsNames, name)
	}
	sort.Strings(registeredCoinsNames)
}

func GetCoins() map[string]CoinType {
	return registeredCoins
}

func GetCoinNames() []string {
	return registeredCoinsNames
}

func GetCoinTypes() []CoinType {
	return registeredCoinsTypes
}

func GetCoinByName(str string) CoinType {
	if coinType, ok := registeredCoins[str]; ok {
		return coinType
	} else {
		return 0
	}
}

func GetCoinNameByIndex(coinType CoinType) string {
	return registeredCoinsIndexes[coinType]
}

func IsCoinExists(val CoinType) bool {
	for i := range registeredCoinsTypes {
		if registeredCoinsTypes[i] == val {
			return true
		}
	}
	return false
}
