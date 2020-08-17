### Yoda

## Prepare environment

1. Install PostgresSQL `brew install postgresql`
2. Install Golang
3. Install Rust
4. Install Docker
5. run `cd owasm/chaintests/bitcoin_block_count/`
6. run `wasm-pack build .`
7. `make install` in chain directory
8. Open 3 tabs on cmd
9. run `docker pull bandprotocol/runtime:1.0.2`

## How to install and run Yoda

1. Open first cmd tab for running the BandChain
2. Open second cmd tab for running the Yoda
3. Open third cmd tab for running the Bandchain CLI

### How to run BandChain on development mode

1. Go to chain directory
2. Setup your PostgresSQL user, port and database name on `start_bandd.sh`
3. run `chmod +x scripts/start_bandd.sh` to change the access permission of start_bandd.script
4. run `./scripts/start_bandd.sh` to start BandChain
5. If fail, try owasm pack build then run script again.

```
cd ../owasm/chaintests/bitcoin_block_count/
wasm-pack build .
cd ../../../chain
```

### How to run Yoda

1. Go to chain directory
2. run `chmod +x scripts/start_yoda.sh` to change the access permission of start_yoda.script
3. run `./scripts/start_yoda.sh validator [number of reporter]` to start Yoda

### Try to request data BandChain

After we have `BandChain` and `Yoda` running, now we can request data on BandChain.
Example of requesting data on BandChain

```
bandcli tx oracle request 1 -c 0000000342544300000000000003e8 1 1  --from requester --chain-id bandchain --gas 3000000 --keyring-backend test  --from requester
```
