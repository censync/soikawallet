package types

import "sort"

const (
	Bitcoin  = NetworkType(1)
	Ethereum = NetworkType(60)
	Tron     = NetworkType(195)
	Polygon  = NetworkType(966)
	BSC      = NetworkType(9006)
)

var (
	registeredNetworks = map[string]NetworkType{
		`Bitcoin`:  Bitcoin,
		`Ethereum`: Ethereum,
		`Tron`:     Tron,
		`Polygon`:  Polygon,
		`BSC`:      BSC,
	}

	registeredNetworksIndexes = map[NetworkType]string{}
	registeredNetworksTypes   []NetworkType
	registeredNetworksNames   []string
)

func init() {
	for name, networkType := range registeredNetworks {
		registeredNetworksIndexes[networkType] = name
		registeredNetworksTypes = append(registeredNetworksTypes, networkType)
		registeredNetworksNames = append(registeredNetworksNames, name)
	}
	sort.Strings(registeredNetworksNames)
}

func GetNetworksNames() []string {
	return registeredNetworksNames
}

func GetNetworks() []NetworkType {
	return registeredNetworksTypes
}

func GetNetworkByName(str string) NetworkType {
	if networkType, ok := registeredNetworks[str]; ok {
		return networkType
	} else {
		return 0
	}
}

func GetNetworkNameByIndex(networkType NetworkType) string {
	return registeredNetworksIndexes[networkType]
}

func IsNetworkExists(val NetworkType) bool {
	for i := range registeredNetworksTypes {
		if registeredNetworksTypes[i] == val {
			return true
		}
	}
	return false
}
