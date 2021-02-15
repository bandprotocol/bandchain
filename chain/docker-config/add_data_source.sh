echo "y" | bandcli tx oracle create-data-source \
  --name "mock data source" \
  --description "mock data source with 'Hello, World!'" \
  --script ./mock.py \
  --owner odin1nnfeguq30x6nwxjhaypxymx3nulyspsuja4a2x \
  --from supplier \
  --chain-id odin \
  --gas auto \
  --keyring-backend test

echo "y" | bandcli tx oracle create-oracle-script \
  --name "mock oracle script" \
  --description "mock oracle script with 'Hello, World!' querying 1st DS" \
  --script ./mock.wasm \
  --owner odin1nnfeguq30x6nwxjhaypxymx3nulyspsuja4a2x \
  --from supplier \
  --chain-id odin \
  --gas auto \
  --keyring-backend test