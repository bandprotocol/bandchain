### Band Oracle 2

## Prepare environment

1. Install PostgresSQL `brew install postgresql`
2. Install Golang
3. Install Rust
4. run `cd owasm/chaintests/bitcoin_block_count/`
5. run `wasm-pack build .`
6. `make install` in chain directory
7. Open 3 tabs on cmd

## How to install and run Band Oracle 2

1. Open first cmd tab for running the BandChain
2. Open second cmd tab for running the Band Oracle 2
3. Open third cmd tab for running the BandChian CLI

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

### How to run Band Oracle 2

1. Go to chain directory
2. run `chmod +x scripts/start_bandoracled2.sh` to change the access permission of start_bandoracled2.script
3. run `./scripts/start_bandoracled2.sh validator [number of reporter]` to start Band Oracle D2

### Try to request data BandChain

After we have `BandChain` and `Band Oracle 2` running, now we can request data on BandChain.
Example of requesting data on BandChain

```
bandcli tx oracle request 1 -c 030000004254436400000000000000 1 1  --from requester --chain-id bandchain --gas 3000000 --keyring-backend test  --from requester
```
