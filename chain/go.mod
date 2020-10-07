module github.com/bandprotocol/bandchain/chain

go 1.13

require (
	github.com/bandprotocol/bandchain/go-owasm v0.0.0-00010101000000-000000000000
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/ethereum/go-ethereum v1.9.19
	github.com/gin-gonic/gin v1.6.3
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gorilla/mux v1.7.4
	github.com/hashicorp/golang-lru v0.5.4
	github.com/kyokomi/emoji v2.2.4+incompatible
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/oasisprotocol/oasis-core/go v0.0.0-20200730171716-3be2b460b3ac
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/rakyll/statik v0.1.7 // indirect
	github.com/segmentio/kafka-go v0.4.5
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.14.0
	github.com/tendermint/tendermint v0.33.8
	github.com/tendermint/tm-db v0.5.1
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.1

replace github.com/bandprotocol/bandchain/go-owasm => ../go-owasm
