build:
	@go build -ldflags "-s -w" -o ./soikawallet ./cmd/tui/main.go

build-testnet:
	@go build -tags testnet -o ./soikawallet ./cmd/tui/main.go

gen-chainlink:
	@wget https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol
	solc --overwrite --abi AggregatorV3Interface.sol -o .
	abigen --abi=AggregatorV3Interface.abi --pkg=chainlink --out=./service/wallet/internal/oracle/chainlink/aggregator_v3_interface.go
clean:
	@rm -f AggregatorV3Interface*
