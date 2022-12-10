package dict

import (
	"github.com/censync/go-i18n"
)

var english = i18n.Dictionary{
	"ui.button": {
		"ok":                "Ok",
		"quit":              "Quit",
		"wallet_create":     "Create wallet",
		"wallet_restore":    "Restore wallet",
		"qr_show_hq":        "Show HQ",
		"finish":            "Finish",
		"next":              "Next",
		"back":              "Back",
		"generate_mnemonic": "Generate mnemonic",
		"copy_to_clipboard": "Copy",
		"create":            "Create",
		"add":               "Add",
		"remove":            "Remove",
	},
	"ui.label": {
		"splash_option_airgap":    "Create wallet with Soika AirGap application",
		"splash_option_no_airgap": "Create wallet without safe AirGap key storage and offline transactions signing",
		"info":                    "The most secure open source cryptocurrency non-custodial wallet with AirGap support",
	},
}
