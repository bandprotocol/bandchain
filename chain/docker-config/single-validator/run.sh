#!/bin/bash

bandd init validator --chain-id=bandchain

cp /zoracle/docker-config/single-validator/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp /zoracle/docker-config/single-validator/node_key.json ~/.bandd/config/node_key.json
cp /zoracle/docker-config/single-validator/genesis.json ~/.bandd/config/genesis.json

bandd start --rpc.laddr tcp://0.0.0.0:26657 & go run cmd/provider/main.go
