# Scan

## Installation

```
yarn install
yarn global add local-cors-proxy
```

## Download GraphQL Schema

```
npx get-graphql-schema  http://d3n-debug.bandprotocol.com:5433/v1/graphql -j > graphql_schema.json
```

## Running App Development

In 2 separate tabs:

```sh
yarn bsb -make-world -w -ws _ # ReasonML compiler
# Replace https://mock.com/ by the real url and don't forgot the / at the back
RPC_URL=https://mock.com/ GRAPHQL_URL=wss://mock-graphql.com yarn parcel index.html # Serve to localhost:1234
```

## Build production

```sh
# Replace https://mock.com/ by the real url and don't forgot the / at the back
RPC_URL=https://mock.com/ GRAPHQL_URL=wss://mock-graphql.com yarn build
```
