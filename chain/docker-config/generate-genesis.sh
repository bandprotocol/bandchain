#!/bin/bash

DIR=`dirname "$0"`

# remove old genesis
rm -rf ~/.band*

# initial new node
bandd init node-validator --chain-id bandchain --zoracle band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs

# create acccounts
expect $DIR/add-account.exp \
    validator1 \
    12345678 \
    "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \

expect $DIR/add-account.exp \
    validator2 \
    12345678 \
    "loyal damage diet label ability huge dad dash mom design method busy notable cash vast nerve congress drip chunk cheese blur stem dawn fatigue" \

expect $DIR/add-account.exp \
    validator3 \
    12345678 \
    "whip desk enemy only canal swear help walnut cannon great arm onion oval doctor twice dish comfort team meat junior blind city mask aware" \

expect $DIR/add-account.exp \
    validator4 \
    12345678 \
    "unfair beyond material banner okay genre camera dumb grit balcony permit room intact code degree execute twin flip half salt script cause demand recipe" \

expect $DIR/add-account.exp \
    requester \
    12345678 \
    "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \

# add accounts to genesis
bandd add-genesis-account validator1 10000000000000uband --keyring-backend test
bandd add-genesis-account validator2 10000000000000uband --keyring-backend test
bandd add-genesis-account validator3 10000000000000uband --keyring-backend test
bandd add-genesis-account validator4 10000000000000uband --keyring-backend test
bandd add-genesis-account requester 10000000000000uband --keyring-backend test

# genesis configurations
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

# create copy of config.toml
cp ~/.bandd/config/config.toml ~/.bandd/config/config.toml.temp

# modify moniker
sed 's/node-validator/node-validator-1/g' ~/.bandd/config/config.toml.temp > ~/.bandd/config/config.toml

# register initial validators
bandd gentx \
    --amount 100000000uband \
    --node-id 11392b605378063b1c505c0ab123f04bd710d7d7 \
    --pubkey bandvalconspub1addwnpepq06h7wvh5n5pmrejr6t3pyn7ytpwd5c0kmv0wjdfujs847em8dusjl96sxg \
    --name validator1 \
    --ip 172.18.0.11 \
    --keyring-backend test

# modify moniker
sed 's/node-validator/node-validator-2/g' ~/.bandd/config/config.toml.temp > ~/.bandd/config/config.toml

bandd gentx \
    --amount 100000000uband \
    --node-id 0851086afcd835d5a6fb0ffbf96fcdf74fec742e \
    --pubkey bandvalconspub1addwnpepqfey4c5ul6m5juz36z0dlk8gyg6jcnyrvxm4werkgkmcerx8fn5g2gj9q6w \
    --name validator2 \
    --ip 172.18.0.12 \
    --keyring-backend test

# modify moniker
sed 's/node-validator/node-validator-3/g' ~/.bandd/config/config.toml.temp > ~/.bandd/config/config.toml

bandd gentx \
    --amount 100000000uband \
    --node-id 7b58b086dd915a79836eb8bfa956aeb9488d13b0 \
    --pubkey bandvalconspub1addwnpepqwj5l74gfj8j77v8st0gh932s3uyu2yys7n50qf6pptjgwnqu2arxkkn82m \
    --name validator3 \
    --ip 172.18.0.13 \
    --keyring-backend test

# modify moniker
sed 's/node-validator/node-validator-4/g' ~/.bandd/config/config.toml.temp > ~/.bandd/config/config.toml

bandd gentx \
    --amount 100000000uband \
    --node-id 63808bd64f2ec19acb2a494c8ce8467c595f6fba \
    --pubkey bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj \
    --name validator4 \
    --ip 172.18.0.14 \
    --keyring-backend test

# remove temp test
rm -rf ~/.bandd/config/config.toml.temp

# collect genesis transactions
bandd collect-gentxs

# copy genesis to the proper location!
cp ~/.bandd/config/genesis.json $DIR/genesis.json
