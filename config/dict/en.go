package dict

import (
	"github.com/censync/go-i18n"
)

var english = i18n.Dictionary{
	"tui.button": {
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
		"accept":            "Accept",
		"reject":            "Reject",
	},
	"tui.label": {
		"splash_option_airgap":    "Create wallet with Soika AirGap application",
		"splash_option_no_airgap": "Create wallet without safe AirGap key storage and offline transactions signing",
		"info":                    "The most secure open source cryptocurrency non-custodial wallet with AirGap support",
		"terms_of_use":            "Choose life\nChoose non-custodial wallet\nChoose network\nChoose private node\nChoose token\nChoose operation\nChoose recipient\nChoose safety",
		"privacy_policy":          "We don't collect any user data",
		"terms_of_use_testnet":    "This is UNSAFE BUILD FOR DEVELOPERS ONLY\nDO NOT USE IT FOR OPERATIONS WITH MAINNET",
		"privacy_policy_testnet":  "This is TESTING VERSION FOR DEVELOPERS\nThis application version may collect debug information",
	},
}
