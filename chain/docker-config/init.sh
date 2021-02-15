#!/bin/bash

bandd init $1 --chain-id odin --timeout-commit 5000

cp ./$1/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp ./$1/node_key.json ~/.bandd/config/node_key.json
cp ./genesis.json ~/.bandd/config/genesis.json