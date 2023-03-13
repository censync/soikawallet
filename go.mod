module github.com/censync/soikawallet

go 1.18

require (
	github.com/btcsuite/btcd v0.23.4
	github.com/btcsuite/btcd/btcutil v1.1.3
	github.com/censync/go-i18n v1.1.0
	github.com/ecies/go/v2 v2.0.4
	github.com/ethereum/go-ethereum v1.10.26
	github.com/gdamore/tcell/v2 v2.6.0
	github.com/google/uuid v1.3.0
	github.com/rivo/tview v0.0.0-20230208211350-7dfff1ce7854
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/stretchr/testify v1.8.0
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/crypto v0.6.0
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rjeczalik/notify v0.9.3 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/rivo/tview => ../tview
