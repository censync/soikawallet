module github.com/censync/soikawallet

go 1.18

require (
	github.com/btcsuite/btcd v0.23.4
	github.com/btcsuite/btcd/btcutil v1.1.3
	github.com/censync/go-airgap v0.0.0-20230822093833-f9f616954fbc
	github.com/censync/go-i18n v1.1.0
	github.com/censync/go-mhda v0.0.0-20230907105920-0f9d9a7b9913
	github.com/censync/go-zbar v0.0.0-20230729001432-1a756b569f38
	github.com/ethereum/go-ethereum v1.13.0
	github.com/gdamore/tcell/v2 v2.6.0
	github.com/gorilla/websocket v1.5.0
	github.com/rivo/tview v0.0.0-20230208211350-7dfff1ce7854
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/stretchr/testify v1.8.1
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/crypto v0.12.0
)

require (
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/bits-and-blooms/bitset v1.5.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.2 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.10.0 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set/v2 v2.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/ethereum/c-kzg-4844 v0.3.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/exp v0.0.0-20230810033253-352e893a4cad // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/term v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/tools v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)

// replace github.com/censync/go-airgap => ../go-airgap
// replace github.com/censync/go-mhda => ../go-mhda
replace github.com/rivo/tview => ../tview
