// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

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
		"accept":            "Accept",
		"reject":            "Reject",
	},
	"ui.tab": {
		"wizard": "Wizard",
		"bulk":   "Bulk",
	},
	"ui.label": {
		"entropy":                 "Entropy",
		"language":                "Language",
		"mnemonic":                "Mnemonic",
		"passphrase":              "Password",
		"addresses":               "Addresses",
		"options":                 "Options",
		"choose_chain":            "Select network",
		"derivation_type":         "Derivation type",
		"derivation_path":         "Derivation path",
		"use_hardened":            "Use hardened",
		"row_add":                 "Add row",
		"row_remove":              "Remove row",
		"account":                 "Account",
		"charge":                  "Charge",
		"index":                   "Index",
		"cannot_create_addr":      "Cannot create address: {err}",
		"splash_option_airgap":    "Create wallet with Soika AirGap application",
		"splash_option_no_airgap": "Create wallet without safe AirGap key storage and offline transactions signing",
		"info":                    "The most secure open source cryptocurrency non-custodial wallet with AirGap support",
		"terms_of_use":            "Choose life\nChoose non-custodial wallet\nChoose network\nChoose private node\nChoose token\nChoose operation\nChoose recipient\nChoose safety",
		"privacy_policy":          "We don't collect any user data",
		"terms_of_use_testnet":    "This is UNSAFE BUILD FOR DEVELOPERS ONLY\nDO NOT USE IT FOR OPERATIONS WITH MAINNET",
		"privacy_policy_testnet":  "This is TESTING VERSION FOR DEVELOPERS\nThis application version may collect debug information",
	},
}
