#!/bin/bash

DIR=`dirname "$0"`

# remove old genesis
rm -rf ~/.band*

# initial new node
bandd init validator --chain-id bandchain --zoracle band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs

# create acccounts
expect $DIR/../add-account.exp \
    validator \
    12345678 \
    "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \

expect $DIR/../add-account.exp \
    requester \
    12345678 \
    "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \

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
    --ip 172.18.0.15 \
    --keyring-backend test

# collect genesis transactions
bandd collect-gentxs

# copy genesis to the proper location!
cp ~/.bandd/config/genesis.json $DIR/genesis.json
