#!/bin/bash

DIR=`dirname "$0"`

# remove old genesis
rm -rf ~/.band*

# initial new node
bandd init validator --chain-id bandchain

# make uband staking denom
cat ~/.bandd/config/genesis.json \
    | python3 -c 'import json; import sys; genesis = json.loads(sys.stdin.read()); genesis["app_state"]["staking"]["params"]["bond_denom"] = "uband"; genesis["consensus_params"]["validator"]["pub_key_types"] = ["secp256k1"]; genesis["app_state"]["crisis"]["constant_fee"]["denom"] = "uband"; genesis["app_state"]["gov"]["deposit_params"]["min_deposit"][0]["denom"] = "uband"; print(json.dumps(genesis))' \
    > ~/.bandd/config/genesis.json.temp

mv ~/.bandd/config/genesis.json.temp ~/.bandd/config/genesis.json

# create acccounts
expect $DIR/../add-account.exp \
    validator \
    12345678 \
    "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \

# add accounts to genesis
bandd add-genesis-account $(bandcli keys show validator -a) 10000000000000uband

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
    --name validator \
    --ip 172.18.0.11

# collect genesis transactions
bandd collect-gentxs

# copy genesis to the proper location!
cp ~/.bandd/config/genesis.json $DIR/genesis.json
