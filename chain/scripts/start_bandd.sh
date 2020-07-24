cp ./docker-config/single-validator/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp ./docker-config/single-validator/node_key.json ~/.bandd/config/node_key.json

dropdb my_db
createdb my_db

# start bandchain
bandd start --rpc.laddr tcp://0.0.0.0:26657
