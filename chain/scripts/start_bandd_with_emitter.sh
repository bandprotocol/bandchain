cp ./docker-config/single-validator/priv_validator_key.json ~/.bandd/config/priv_validator_key.json
cp ./docker-config/single-validator/node_key.json ~/.bandd/config/node_key.json

# Delete old kafka topic
kafka-topics --zookeeper localhost:2181 --delete --topic test

# Create kafka topic
kafka-topics --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic test

# Init table from flusher
source ../flusher/venv/bin/activate
python ../flusher/main.py init bandchain test --db localhost:5432/my_db


# start bandchain
bandd start --with-emitter test \
  --rpc.laddr tcp://0.0.0.0:26657 --pruning=nothing
