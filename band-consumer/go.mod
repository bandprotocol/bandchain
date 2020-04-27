module github.com/bandprotocol/band-consumer

go 1.13

require (
	github.com/bandprotocol/bandchain/chain v0.0.0-20200413032603-6ae7f6e32df7
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200406170659-df5badaf4c2b
	github.com/gorilla/mux v1.7.4
	github.com/otiai10/copy v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.3
	github.com/tendermint/tm-db v0.5.1
)

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
