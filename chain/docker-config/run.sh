#!/bin/bash

cp ./genesis.json ~/.bandd/config/genesis.json
bandd init $1 --chain-id bandchain

pwd

cp ./$1/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp ./$1/node_key.json ~/.bandd/config/node_key.json

bandd start --rpc.laddr tcp://0.0.0.0:26657 \
    --p2p.persistent_peers 11392b605378063b1c505c0ab123f04bd710d7d7@172.18.0.11:26656,0851086afcd835d5a6fb0ffbf96fcdf74fec742e@172.18.0.12:26656,63808bd64f2ec19acb2a494c8ce8467c595f6fba@172.18.0.14:26656,7b58b086dd915a79836eb8bfa956aeb9488d13b0@172.18.0.13:26656
