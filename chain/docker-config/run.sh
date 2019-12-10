#!/bin/bash

bandd init $1 --chain-id=bandchain

cp /zoracle/docker-config/$1/config.toml ~/.bandd/config/config.toml
cp /zoracle/docker-config/$1/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp /zoracle/docker-config/$1/node_key.json ~/.bandd/config/node_key.json
cp /zoracle/docker-config/genesis.json ~/.bandd/config/genesis.json

bandd start --rpc.laddr tcp://0.0.0.0:26657 & go run cmd/provider/main.go
