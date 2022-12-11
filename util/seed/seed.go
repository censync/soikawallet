package seed

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

const (
	Entropy128 = 128
	Entropy160 = 160
	Entropy192 = 192
	Entropy224 = 224
	Entropy256 = 256
)

var mnemonicCount = map[int]int{
	12: Entropy128,
	15: Entropy160,
	18: Entropy192,
	21: Entropy224,
	24: Entropy256,
}

func EntropyList() []string {
	var result []string
	for _, entropy := range mnemonicCount {
		result = append(result, strconv.Itoa(entropy))
	}
	sort.Strings(result)
	return result
}

// 128 - 12, 160 - 15, 192 - 18, 224 - 21, 256 - 24
func Check(mnemonic string) error {
	entropy := entropyByMnemonic(mnemonic)

	if entropy == 0 {
		return errors.New("undefined entropy")
	}

	if !checkDuplicates(mnemonic) {
		return errors.New("mnemonic has duplicates prefix")
	}
	return nil
}

func entropyByMnemonic(mnemonic string) int {
	wordsCount := len(strings.Fields(mnemonic))
	return mnemonicCount[wordsCount]
}

func checkDuplicates(str string) bool {
	var index = map[string]bool{}
	arr := strings.Fields(str)
	for idx := range arr {
		prefix := substr(arr[idx], 4)
		if _, ok := index[prefix]; !ok {
			index[prefix] = true
		} else {
			return false
		}
	}
	return true
}

func substr(input string, length int) string {
	asRunes := []rune(input)

	if 0 >= len(asRunes) {
		return ""
	}

	if length > len(asRunes) {
		length = len(asRunes) - 0
	}

	return string(asRunes[0 : 0+length])
}
