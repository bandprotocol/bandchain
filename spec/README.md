<div align="center">
  <h1>📫 BandChain - Decentralized Oracle Network</h1>
  <b>WORKING DRAFT</b>
</div>

[![](https://img.shields.io/badge/chat-on%20Telegram%20💬-blue.svg)](https://t.me/bandprotocol)
[![](https://img.shields.io/badge/chat-on%20Discord%20🤖-violet.svg)](https://discord.gg/es9CK4)

This live document outlines implementation plans and research directions for [Band Protocol](https://bandprotocol.com)'s decentralized data delivery network. The implementations of the core protocol and its supporting tools are available for public inspection on [Github](https://github.com/bandprotocol/bandchain).

# Introduction

The majority of existing smart contract platforms, while support trustless execution of arbitrary programs, lack access to real-world data. This renders smart contracts not as useful. BandChain connects public blockchains with off-chain information, with the following design goals.

1. **Speed and Scalability:** The system must be able to serve a large quantity of data requests to multiple public blockchains with minimal latency and high throughput. The expected response time must be in the order of seconds.

2. **Cross-Chain Compatibility:** The system must be blockchain-agnostic and able to serve data to most available public blockchains. Verification of data authenticity on the target blockchains must be efficient and trustless.

3. **Data Flexibility:** The system must be generic and able to support different ways to fetch and aggregate data, including both permissionless, publicly available data and information guarded by centralized parties.

BandChain archives the aforementioned goals with a blockchain specifically built for off-chain data curation. The blockchain supports generic data requests and on-chain aggregations with WebAssembly-powered oracle scripts. Oracle results on BandChain blockchain can be sent across to other blockchains via the [Inter-Blockchain Communication protocol (IBC)](https://cosmos.network/ibc/) or customized one-way bridges with minimal latency.

# Table of Contents

The specification is split into multiple sub-documents. See the following for the complete list.

1. [Terminology](./1_terminology.md) - Explanation of terms used in BandChain specification.
2. [System Overview](./2_system_overview.md) - Overview of the network architecture.
3. [Token Economics](./3_token_economics.md) - BAND token and its economic inside of BandChain.
4. [Blockchain Parameters](./4_blockchain_parameters) - Global parameters and the proposed initial values.
5. [Protocol Messages](./5_protocol_messages.md) - List of BandChain-specific messages.
6. [IBC Packets](./6_ibc_packets.md) - Complete list of packets recognized by Bandchain.
7. [Lite-Client Protocol](./7_lite_client_protocol) - TODO
8. [Oracle WebAssembly (Owasm)](./8_oracle_webassembly.md) - TODO
9. [Input/Output Encoding](./9_input_output_encoding.md) - TODO
