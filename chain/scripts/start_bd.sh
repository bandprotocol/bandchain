cp ./docker-config/single-validator/priv_validator_key.json ~/.bd/config/priv_validator_key.json
cp ./docker-config/single-validator/node_key.json ~/.bd/config/node_key.json

# start bandchain
bd start --rpc.laddr tcp://0.0.0.0:1237
