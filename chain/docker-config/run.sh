#!/bin/bash

bandd init $1 --chain-id=bandchain

cp /zoracle/docker-config/$1/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp /zoracle/docker-config/$1/priv_validator_state.json ~/.bandd/data/priv_validator_state.json
cp /zoracle/docker-config/$1/node_key.json ~/.bandd/config/node_key.json
cp /zoracle/docker-config/genesis.json ~/.bandd/config/genesis.json

# add cors in config.toml
cd ~/.bandd/config/
sed 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' config.toml > config_tmp.toml
mv config_tmp.toml config.toml
cd /zoracle/

if [ "$1" == "query-node" ];then
    mkdir ~/.banddb
    bandd start --rpc.laddr tcp://0.0.0.0:26657 --with-db "postgres: host=172.18.0.88 port=5432 user=postgres dbname=postgres password=postgrespassword sslmode=disable" \
    --p2p.persistent_peers cf65f9b5f290e3eef62dbf721397bc2e3fd47ecd@wenchang-testnet2-alice.node.bandchain.org:26656,c5d042cca13c34ee5318a4e5fc49dd7950f0a1da@wenchang-testnet2-bob.node.bandchain.org:26656 --uptime-look-back 20
else
    bandd start --rpc.laddr tcp://0.0.0.0:26657 \
    --p2p.persistent_peers 11392b605378063b1c505c0ab123f04bd710d7d7@172.18.0.11:26656,0851086afcd835d5a6fb0ffbf96fcdf74fec742e@172.18.0.12:26656,63808bd64f2ec19acb2a494c8ce8467c595f6fba@172.18.0.14:26656,7b58b086dd915a79836eb8bfa956aeb9488d13b0@172.18.0.13:26656
fi

