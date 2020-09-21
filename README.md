# BandChain - Decentralized Data Delivery Network

This repository is a [monorepo] containing the reference implementation of BandChain and its various supporting tools. See below for the breakdown and explanation of each module. README for each of the modules.

## Table of Contents

| Module                 | Description                                               |
| ---------------------- | --------------------------------------------------------- |
| [`chain`](chain)       | 🔗 BandChain blockchain reference implementation          |
| [`bridges`](bridges)   | 📡 Lite client bridges on other smart contract platforms  |
| [`lambda`](lambda)     | 👷‍♂️ AWS Lambda package for running data source executables |
| [`go-owasm`](go-owasm) | 🐀 Go library for executing oracle scripts with Wasmer    |
| [`helpers`](helpers)   | 🔪 Client-side utility libraries                          |
| [`obi`](obi)           | 📦 Oracle binary encoding implementations                 |
| [`owasm`](owasm)       | 🔮 WebAssembly library for writing oracle scripts         |
| [`scan`](scan)         | 🔍 Web interface to explore D3N network                   |

## Running with Docker

There are 2 ways to run bandchain

#### Run on 4 validators

```
./chain/docker-config/generate-genesis.sh && docker-compose up multi-validator
```

#### Run 1 validator (for development)

```
./chain/docker-config/single-validator/generate-genesis.sh && docker-compose up single-validator
```

#### (Optional) Run Owasm Studio

```
docker-compose up <single or multi> owasm-studio
```

#### Tear down

```
docker-compose down -v
```

## Running a Validator Node

[📚 Guide to Becoming a Validator](https://medium.com/bandprotocol/bandchain-wenchang-testnet-2-how-to-join-as-a-validator-76bc4180ddd7)

## License & Contributing

All modules are licensed under the terms of the Apache 2.0 License unless otherwise specified in the LICENSE file at module's root.

We highly encourage participation from the community to help with D3N development. If you are interested in developing with D3N or have suggestion for protocol improvements, please open an issue, submit a pull request, or [drop as a line].

[monorepo]: https://en.wikipedia.org/wiki/Monorepo
[drop as a line]: mailto:connect@bandprotocol.com
