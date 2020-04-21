# Scan

## Installation

```
yarn install
yarn global add local-cors-proxy
```

## Download GraphQL Schema

```
npx get-graphql-schema  https://d3n.bandprotocol.com/v1/graphql -j > graphql_schema.json
```

## Running App Development

In 2 separate tabs:

```sh
# First tab
yarn bsb -make-world -w -ws _ # ReasonML compiler

# Second tab
RPC_URL=https://d3n.bandprotocol.com/rest GRAPHQL_URL=wss://d3n.bandprotocol.com/v1/graphql LAMBDA_URL=<LAMBDA_URL> FAUCET_URL=https://d3n.bandprotocol.com/faucet/request yarn parcel index.html --no-cache # Serve to localhost:1234
```

## Build production

```sh
RPC_URL=https://d3n.bandprotocol.com/rest GRAPHQL_URL=wss://d3n.bandprotocol.com/v1/graphql LAMBDA_URL=<LAMBDA_URL> FAUCET_URL=https://d3n.bandprotocol.com/faucet/request yarn build
```
