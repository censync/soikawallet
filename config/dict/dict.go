package dict

import (
	"github.com/censync/go-i18n"
)

const (
	defaultLanguage = `en`
)

var dictCollection = i18n.DictionaryCollection{
	"en": english,
}

func GetTr(language string) *i18n.Translator {
	err := i18n.Init(defaultLanguage, dictCollection)
	if err != nil {
		panic("loading dictionaries error")
	}
	return i18n.New(language)
}
