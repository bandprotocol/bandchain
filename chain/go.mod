module github.com/bandprotocol/d3n/chain

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.37.4
	github.com/ethereum/go-ethereum v1.9.7
	github.com/gin-gonic/gin v1.5.0
	github.com/gorilla/mux v1.7.0
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/iavl v0.12.4
	github.com/tendermint/tendermint v0.32.8
	github.com/tendermint/tm-db v0.2.0
	github.com/wasmerio/go-ext-wasm v0.0.0-20191113152408-05371a4e2fe5
	golang.org/x/crypto v0.0.0-20191122220453-ac88ee75c92c // indirect
	golang.org/x/net v0.0.0-20191124235446-72fef5d5e266 // indirect
)

replace github.com/tendermint/tendermint => github.com/bandprotocol/tendermint v0.32.9-0.20200116035555-cfd3a61367ea
