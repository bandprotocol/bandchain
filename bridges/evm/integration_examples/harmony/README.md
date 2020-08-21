# Harmony Example Deploy

## Project Structure

- [**../../contracts/**](../../contracts): Contains the contracts associated with this project
- **.env**: Contains the environment variables used in the project
- **truffle-config.js**: Responsible for passing the environment variables from `.env` into our project
- **index.js**: file for easily sending a request to the `relayAndSafe` function in our `ReceiverMock` contract

## Configuration

### RPC Enpoint

We'll be using Harmony's testnet endpoint throughout this example. The URL for the endpoint is

```text
https://api.s0.b.hmny.io
```

The full list of endpoints can be found [here](https://docs.harmony.one/developers/harmony-networks/explorer-and-rpc-endpoints).

### Environment Variables

For a deploy to any of Harmony's networks, we need three pieces of information:

- `<NETOWRK>_PRIVATE_KEY`: The private key of the account we're deploying from
- `<NETWORK>_MNEMONIC`: The deploying account's mnenonic
- `<NETWORK>_URL`: The URL of the network we want to deploy to

Where `<NETWORK>` is one either `LOCAL`, `TESTNET`, or `MAINNET`. In our case, these pieces of information is stored in the `.env` file.

## Build

Before we can deploy, we first need to install the necessary dependencies and compile our contract.

```bash
yarn install
yarn truffle compile
```

## Deployment

#### Compile and Deploy Your Contract

Now that our contracts have been compiled, we proceed to deploying it to Harmony's testnet.

```bash
yarn truffle deploy --network=testnet --reset
```

### Check Contract Deploy

To check that your contract has been successfully deployed, run the following command:

```bash
yarn truffle networks
```

If the contract deployment was successful, you should something similar to:

```bash
Network: local (id: 2)
  No contracts deployed.

Network: mainnet0 (id: 1)
  No contracts deployed.

# Successful contract deploy will show some info here
Network: testnet (id: 2)
  BridgeWithCache: 0x176fb19ddAf98f639C55A2008623B581821E7B42
  Migrations: 0x1290AbE9bd7123c56a1244e5ABb689a8Fbef7DD9
  ReceiverMock: 0x2F23d4C4Ae0C4B13808EEf8CD8FF99eDf59c2343
```

### Test ReceiverMock contract

To sending relayAndSafe transaction, run the following command:

```bash
node index.js
```

This will send a request calling ReceiverMock.sol's `relayAndSafe` function. Please see the [`index.js`](/harmony/index.js) for details on the request process, payload, etc.

### Harmony JSON-RPC API

To retrieve information from the Harmony blockchain, we will be using Harmony's JSON-RPC API ([doc](https://docs.harmony.one/home/developers/api)). The API specification is quite similar [Ethereum's](https://github.com/ethereum/wiki/wiki/JSON-RPC), for those who are familiar with it.

#### Transaction Verification

First let's make a request to check whether the previous transaction was successful. To do this we'll use the [`hmy_gettransactionReceipt`](https://docs.harmony.one/home/developers/api/methods/transaction-related-methods/hmy_gettransactionreceipt) method. When using curl, this call has the following structure:

```bash
curl --location --request POST 'https://api.s0.b.hmny.io' \
--header 'Content-Type: application/json' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "jsonrpc":"2.0",
    "method":"hmy_getTransactionReceipt",
    "params":[
    "0xdb9e715485432bc84a1d0f8bc4ea001b4b5c4cc4659ab4bace4abe1d59d93d14"],
    "id":1
}'
```

where the hex value inside `params` is the transaction hash of the transaction we want to look up.

The response return will look something like:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "blockHash": "0xc1c3572dcfedb6207c8396e0bbb9fe8254b56def988e5bad8b59751395b647e8",
    "blockNumber": "0xff494",
    "contractAddress": null,
    "cumulativeGasUsed": "0x6b44d",
    "from": "one18t4yj4fuutj83uwqckkvxp9gfa0568uc48ggj7",
    "gasUsed": "0x6b44d",
    "logs": [],
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "shardID": 0,
    "status": "0x1",
    "to": "one1enza0gdgwxguz9kxn3fyq5rzdqn8kc76zdfwu7",
    "transactionHash": "0xdb9e715485432bc84a1d0f8bc4ea001b4b5c4cc4659ab4bace4abe1d59d93d14",
    "transactionIndex": "0x0"
  }
}
```

If the transaction is successful, `status` will have a value of `0x1` as seen here.

#### Retrieving A Contract's Variable Value

If instead we want to query the value of a variable in a specific contract, we would instead use the [`hmy_call`](https://docs.harmony.one/home/developers/api/methods/contract-related-methods/hmy_call) specification. In this case we will check if the `latestReq` variable in `relayAndSafe` has been set as expected after being called earlier.

```bash
curl --location --request POST 'https://api.s0.b.hmny.io' \
--header 'Content-Type: application/json' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "hmy_call",
    "params": [
        {
            "to": "0xccc5D7A1A87191C116c69c5240506268267B63dA",
            "data": "0x8a0d3c31"
        },
        "latest"
    ],
    "id": 1
}'
```

The data we need to pass in are:

- `to`: the address of the contract we want to query (as shown with `truffle networks`)
- `data`: The first 4 bytes of the hex returned when running variable name append with "()" (e.g. "`latestReq()`") through a Keccak-256 hash function ([online tool](https://emn178.github.io/online-tools/keccak_256.html))

If the variable is successfully set, the response body should have the following structure

```json
{
   "jsonrpc":"2.0",
   "id":1,
   "result":"0x00000000000000000000000000000000000000000000000000000000000000a00
   00000000000000000000000000000000000000000000000000000000000000100000000000000
   000000000000000000000000000000000000000000000000e0000000000000000000000000000
   00000000000000000000000000000000000040000000000000000000000000000000000000000
   00000000000000000000000400000000000000000000000000000000000000000000000000000
   0000000000962616e642074657374000000000000000000000000000000000000000000000000
   0000000000000000000000000000000000000000000000000000000000001e303330303030303
   0343235343433363430303030303030303030303030300000"
}
```

where the value of `result` is non-zero.

## Resources

### Code

- Harmony's official GitHub [repo](https://github.com/harmony-one)
- Harmony JavaScript [SDK](https://github.com/harmony-one/sdk)

### Documentation

- Harmony [documentation](https://docs.harmony.one/home/)
- Ethereum's JSON-RPC page on their [wiki](https://github.com/ethereum/wiki/wiki/JSON-RPC)

### Tools

- Online Keccak-256 [hasher](https://emn178.github.io/online-tools/keccak_256.html)
