DIR=`dirname "$0"`

mkdir -p pkg/owasm/res

# release mode
cd ../owasm/chaintests

for f in *; do
    if [ -d "$f" ]; then
        RUSTFLAGS='-C link-arg=-s' cargo build --target wasm32-unknown-unknown --release --package $f
        cp ../target/wasm32-unknown-unknown/release/$f.wasm ../../chain/pkg/owasm/res
    fi
done
cd ../../chain

# debug mode
cd ../owasm/chaindebugtests
for f in *; do
    if [ -d "$f" ]; then
        RUSTFLAGS='-C link-arg=-s' cargo build --target wasm32-unknown-unknown --package $f
        cp ../target/wasm32-unknown-unknown/debug/$f.wasm ../../chain/pkg/owasm/res
    fi
done

cd ../../chain
go test ./...

rm -rf ~/.band*
dropdb my_db
createdb my_db

# initial new node
bandd init validator --chain-id bandchain --oracle band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs
echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | bandcli keys add validator --recover --keyring-backend test  --coin-type 494
echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | bandcli keys add requester --recover --keyring-backend test  --coin-type 494

# add accounts to genesis
bandd add-genesis-account validator 10000000000000uband --keyring-backend test
bandd add-genesis-account requester 10000000000000uband --keyring-backend test

# genesis configurations
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

# register initial validators
bandd gentx \
    --amount 100000000uband \
    --node-id 11392b605378063b1c505c0ab123f04bd710d7d7 \
    --pubkey bandvalconspub1addwnpepq06h7wvh5n5pmrejr6t3pyn7ytpwd5c0kmv0wjdfujs847em8dusjl96sxg \
    --name validator \
    --keyring-backend test

# collect genesis transactions
bandd collect-gentxs
cp ./docker-config/single-validator/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp ./docker-config/single-validator/node_key.json ~/.bandd/config/node_key.json

# start bandchain
bandd start --with-db "postgres: port=5432 user=postgres dbname=my_db sslmode=disable" --rpc.laddr tcp://0.0.0.0:26657   --pruning=nothing
