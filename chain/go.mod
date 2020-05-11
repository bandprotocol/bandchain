module github.com/bandprotocol/bandchain/chain

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200423152229-f1fdde5d1b18
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/ethereum/go-ethereum v1.9.13
	github.com/gin-gonic/gin v1.6.3
	github.com/gogo/protobuf v1.3.1
	github.com/gorilla/mux v1.7.4
	github.com/jinzhu/gorm v1.9.12
	github.com/kyokomi/emoji v2.2.2+incompatible
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/mattn/go-shellwords v1.0.10
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/perlin-network/life v0.0.0-20191203030451-05c0e0f7eaea
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.13.3
	github.com/tendermint/tendermint v0.33.4
	github.com/tendermint/tm-db v0.5.1
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.1
