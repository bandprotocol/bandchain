#!/bin/bash

bandcli config chain-id odin

echo "usage ketchup faculty bench jewel rocket latin absurd decide field party reunion cook entry scout scene miss box memory museum decorate guide few verify" \
    | bandcli keys add $1 --recover --keyring-backend test


# Create system validator
echo "y" | bandcli tx staking create-validator \
  --amount 100000000loki \
  --commission-max-change-rate 0.010000000000000000 \
  --commission-max-rate 0.200000000000000000 \
  --commission-rate 0.100000000000000000 \
  --from $1 \
  --moniker oracle-validator \
  --pubkey odinvalconspub1addwnpepqge86lvslkpfk0rlz0ah9tat0vntx8yele36hhfpflehfehydlutkvdvhfm \
  --min-self-delegation 1 \
  --broadcast-mode block \
  --keyring-backend test \
  --node $2

# Create data source and oracle script
echo "y" | bandcli tx oracle create-data-source \
  --name "mock data source" \
  --description "mock data source with 'Hello, World!'" \
  --script /home/stepan/Projects/band/GeoDB-Limited/odin-developer-edition/data-source-scripts/geo-data-v3.py \
  --owner odin1nnfeguq30x6nwxjhaypxymx3nulyspsuja4a2x \
  --from $1 \
  --gas auto \
  --broadcast-mode block \
  --keyring-backend test \
  --node $2

echo "y" | bandcli tx oracle create-oracle-script \
  --name "mock oracle script" \
  --description "mock oracle script with 'Hello, World!' going on 1 DS" \
  --script /home/stepan/Projects/band/GeoDB-Limited/odin-developer-edition/oracle-scripts/geo_data_source_v3.wasm \
  --owner odin1nnfeguq30x6nwxjhaypxymx3nulyspsuja4a2x \
  --from $1 \
  --gas auto \
  --broadcast-mode block \
  --keyring-backend test \
  --node $2

./chain/scripts/start_yoda_odin.sh $1 $2