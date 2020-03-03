<div align="center">
  <img width="300" src="assets/d3n_banner.png" />
</div>

# Decentralized Data Delivery Network (D3N)

This repository is a [monorepo] containing the reference implementation of D3N and its various supporting tools. See below for the breakdown and explanation of each module. README for each of the modules.

## Table of Contents

| Module               | Description                                               |
| -------------------- | --------------------------------------------------------- |
| [`chain`](chain)     | 🔗 D3N blockchain reference implementation                |
| [`bridges`](bridges) | 📡 Lite client bridges on other smart contract platforms  |
| [`owasm`](owasm)     | 🔮 WebAssembly library for writing oracle scripts         |
| [`scan`](scan)       | 🔍 Web interface to explore D3N network                   |
| [`spec`](spec)       | 📖 D3N research and specification knowledge base          |
| [`studio`](studio)   | 🎬 In-browser IDE for testing and deploying owasm scripts |

## Running with Docker

There are 2 ways to run bandchian

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

## License & Contributing

All modules are licensed under the terms of the Apache 2.0 License unless otherwise specified in the LICENSE file at module's root.

We highly encourage participation from the community to help with D3N development. If you are interested in developing with D3N or have suggestion for protocol improvements, please open an issue, submit a pull request, or [drop as a line].

[monorepo]: https://en.wikipedia.org/wiki/Monorepo
[drop as a line]: mailto:connect@bandprotocol.com
