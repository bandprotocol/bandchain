#!/bin/bash

DIR=`dirname "$0"`

# remove old genesis
rm -rf ~/.band*

# initial new node
bandd init node-validator --chain-id bandchain

# make uband staking denom
cat ~/.bandd/config/genesis.json \
    | python3 -c 'import json; import sys; genesis = json.loads(sys.stdin.read()); genesis["app_state"]["staking"]["params"]["bond_denom"] = "uband"; genesis["consensus_params"]["validator"]["pub_key_types"] = ["secp256k1"]; genesis["app_state"]["crisis"]["constant_fee"]["denom"] = "uband"; genesis["app_state"]["gov"]["deposit_params"]["min_deposit"][0]["denom"] = "uband"; genesis["app_state"]["mint"]["params"]["mint_denom"] = "uband"; print(json.dumps(genesis))' \
    > ~/.bandd/config/genesis.json.temp

mv ~/.bandd/config/genesis.json.temp ~/.bandd/config/genesis.json

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
bandd add-genesis-account $(bandcli keys show validator1 -a) 10000000000000uband
bandd add-genesis-account $(bandcli keys show validator2 -a) 10000000000000uband
bandd add-genesis-account $(bandcli keys show validator3 -a) 10000000000000uband
bandd add-genesis-account $(bandcli keys show validator4 -a) 10000000000000uband

bandd add-genesis-account $(bandcli keys show requester -a) 10000000000000uband

# genesis configurations
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

# register initial validators
echo "12345678" | bandd gentx \
    --amount 100000000uband \
    --node-id 11392b605378063b1c505c0ab123f04bd710d7d7 \
    --pubkey bandvalconspub1addwnpepq06h7wvh5n5pmrejr6t3pyn7ytpwd5c0kmv0wjdfujs847em8dusjl96sxg \
    --name validator1 \
    --ip 172.18.0.11

echo "12345678" | bandd gentx \
    --amount 100000000uband \
    --node-id 0851086afcd835d5a6fb0ffbf96fcdf74fec742e \
    --pubkey bandvalconspub1addwnpepqfey4c5ul6m5juz36z0dlk8gyg6jcnyrvxm4werkgkmcerx8fn5g2gj9q6w \
    --name validator2 \
    --ip 172.18.0.12

echo "12345678" | bandd gentx \
    --amount 100000000uband \
    --node-id 7b58b086dd915a79836eb8bfa956aeb9488d13b0 \
    --pubkey bandvalconspub1addwnpepqwj5l74gfj8j77v8st0gh932s3uyu2yys7n50qf6pptjgwnqu2arxkkn82m \
    --name validator3 \
    --ip 172.18.0.13

echo "12345678" | bandd gentx \
    --amount 100000000uband \
    --node-id 63808bd64f2ec19acb2a494c8ce8467c595f6fba \
    --pubkey bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj \
    --name validator4 \
    --ip 172.18.0.14

# collect genesis transactions
bandd collect-gentxs

# change monikers
cat ~/.bandd/config/genesis.json \
    | python3 -c 'import json; import sys; genesis = json.loads(sys.stdin.read()); genesis["app_state"]["genutil"]["gentxs"][0]["value"]["msg"][0]["value"]["description"]["moniker"] = "node-validator-1"; genesis["app_state"]["genutil"]["gentxs"][1]["value"]["msg"][0]["value"]["description"]["moniker"] = "node-validator-2"; genesis["app_state"]["genutil"]["gentxs"][2]["value"]["msg"][0]["value"]["description"]["moniker"] = "node-validator-3"; genesis["app_state"]["genutil"]["gentxs"][3]["value"]["msg"][0]["value"]["description"]["moniker"] = "node-validator-4"; print(json.dumps(genesis))' \
    > ~/.bandd/config/genesis.json.temp

mv ~/.bandd/config/genesis.json.temp ~/.bandd/config/genesis.json

# copy genesis to the proper location!
cp ~/.bandd/config/genesis.json $DIR/genesis.json
