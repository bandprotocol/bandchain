#!/bin/bash

bandd init validator --chain-id=bandchain

cp /zoracle/docker-config/single-validator/priv_validator_state.json ~/.bandd/data/priv_validator_state.json
cp /zoracle/docker-config/single-validator/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp /zoracle/docker-config/single-validator/node_key.json ~/.bandd/config/node_key.json
cp /zoracle/docker-config/single-validator/genesis.json ~/.bandd/config/genesis.json

# add cors in config.toml
cd ~/.bandd/config/
sed 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' config.toml > config_tmp.toml
mv config_tmp.toml config.toml
cd /zoracle/

bandd start --rpc.laddr tcp://0.0.0.0:26657
