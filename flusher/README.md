
## Prepare environment

1. Install Java `brew cask install java`
2. Install Kafka `brew install kafka`
3. Start Zookeeper services `brew services start zookeeper`
4. Start Kafka services `brew services start kafka`
5. run `python3 -m venv venv && source venv/bin/activate`
6. run `pip install -r requirements.txt`
7. If you have openssl problem while install dependencies run `brew install openssl && export LIBRARY_PATH=$LIBRARY_PATH:/usr/local/opt/openssl/lib/`
8. `make install` in chain directory
9. Open 2 tabs on cmd

### How to run BandChain with emitter on development mode

1. Go to chain directory
2. run `chmod +x scripts/generate_genesis.sh` to change the access permission of start_bandd_with_emitter.script
3. run `./scripts/generate_genesis.sh` to start generate genesis
4. run `chmod +x scripts/start_bandd_with_emitter.sh` to change the access permission of start_bandd_with_emitter.script
5. run `./scripts/start_bandd_with_emitter.sh` to start BandChain


### How to run flusher

1. run `source venv/bin/activate`
2. run `python main.py sync --db localhost:5432/my_db`

### Log consumer

kafka-console-consumer --bootstrap-server localhost:9092 --topic <topic_name> --from-beginning

### Troubleshooting

If you experience problems while trying to start Kafka, you can try to examine Kafka and ZooKeeper logs using the following commands:

1. run `tail -f /usr/local/var/log/zookeeper/zookeeper.log`
2. run `tail -f /usr/local/var/log/kafka/kafka.log`
