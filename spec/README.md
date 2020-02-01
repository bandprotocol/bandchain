<div align="center">
  <h1>📫 Decentralized Data Delivery Network</h1>
  <img width="300" src="/assets/d3n_banner.png" />
  <br />
  <b>WORKING DRAFT</b>
</div>

[![](https://img.shields.io/badge/chat-on%20Telegram%20💬-blue.svg)](https://t.me/bandprotocol)

This live document outlines implementation plans and research directions for [Band Protocol](https://bandprotocol.com)'s decentralized data delivery network (D3N). The implementations of the core protocol and its supporting tools are available for public inspection on [Github](https://github.com/bandprotocol/d3n).

# Introduction
The majority of existing smart contract platforms, while support trustless executation of arbitrary programs, lack access to real-world data. This renders smart contracts not as useful. D3N connects public blockchains with off-chain information, with the following design goals.

1. **Speed and Scalability:** The system must be able to serve a large quantity of data requests to multiple public blockchains with minimal latency.
2. **Cross-Chain Compatibility:** The system must be blockchain-agnostic and able to serve data to most available public blockchains.
3. **Data Flexibility:** The system must be generic and able to support different ways to fetch and aggregate data, including both permissionless publicly-available information and information guarded by centralized parties.


D3N archives the aforementioned goals with a blockchain specifically built for offchain data curation. The blockchain supports generic data requests and on-chain resolves with WebAssembly-powered oracle scripts. Data on D3N blockchain can be sent across via one-way bridges to other blockchains with minimal latency.

# D3N System Overview
D3N is a high-performance public blockchain that allows anyone to request for arbitrary off-chain data and computation. Built on top [Tendermint](https://tendermint.com/) and [Cosmos SDK](http://cosmos.network/), the blockchain utilizes [BFT](https://en.wikipedia.org/wiki/Byzantine_fault) consensus algorithm to reach immediate finality upon getting confirmations from a sufficient number of block validators. In addition to signing blocks, block validators are responsible for reponding to data requests submitted to the system.

## Terminology
Before diving further into the blockchain protocol and how it works, we take the opportunity here to explain various terms that we will later use to explain the protocol.

### Data Sources
[sec:term-data-sources]: #Data-Sources

A data source is the most fundamental unit of the oracle system. It describes the procedure to retrieve a raw data point from a primary source and the fee associated with one data query. In D3N, a data source can be registered into the system by anyone. Once registered, a data source can be *owned* or *unowned*. An owned data source can be changed and upgraded by its owner who can also claims the collected fees, while an unowned data source is immutable and cannot be changed.

Note that even though an unowned data source cannot be changed, it can still be controlled by centralized parties if its procedure depends on centralized sources. Below is some examples of data source procedures written in [Bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell)). 

**Example:** A script to retrieve a cryptocurrency price from [CoinGecko](https://www.coingecko.com/) with one argument: the currency identifier. The script assumes that [cURL](https://en.wikipedia.org/wiki/CURL) and [jq](https://github.com/stedolan/jq) are available on the host and the host is connected to the Internet.

```bash
#!/bin/sh

# Cryptocurrency price endpoint: https://www.coingecko.com/api/documentations/v3
URL="https://api.coingecko.com/api/v3/simple/price?ids=$1&vs_currencies=usd"
KEY=".$1.usd"

# Performs data fetching and parses the result
curl -s -X GET $URL -H "accept: application/json" | jq -r ".[\"$1\"].usd"
```

**Example:** A script to resolve a given hostname into IP addresses. The script assumes that [getent](https://en.wikipedia.org/wiki/Getent) and [awk](https://en.wikipedia.org/wiki/AWK) are available on the host and the host is connected to the DNS network.

```bash
#!/bin/sh

getent hosts $1 | awk '{ print $1 }'
```

### Oracle Scripts
[sec:term-oracle-scripts]: #Oracle-Scripts

An oracle script is a executable program that encodes (1) the set of raw data requests to the data sources it needs and (2) the way to aggregate raw data reports into the final result. An oracle script may depend on an input data, making it reusable without needing to rewrite the whole program. It may also depend on other oracle scripts, making oracle scripts composable just like [smart contracts](https://en.wikipedia.org/wiki/Smart_contract).

**Example**: A [psudocode](https://en.wikipedia.org/wiki/Pseudocode) showing an example oracle script to fetch cryptocurrency prices from multiple sources: [CoinGecko](https://coinmarketcap.com/), [CryptoCompare](https://www.cryptocompare.com/), and [CoinMarketCap](https://coinmarketcap.com/). The code assumes that data sources for the three sources are available and reports the average among all sources from all data reporters as the final result.

```python
# 1st Phase. Emits raw data requests that the oracle script needs.
def prepare(symbol):
    request(get_px_from_coin_gecko, symbol)
    request(get_px_from_crypto_compare, symbol)
    request(get_px_from_coin_market_cap, symbol)

# 2nd Phase. Aggregates raw data reports into the final result.
def aggregate(symbol, number_of_reporters):
    data_report_count = 0
    price_sum = 0.0
    for reporter_index in range(number_of_reporters):
        for data_source in (
            get_px_from_coin_gecko,
            get_px_from_crypto_compare,
            get_px_from_coin_market_cap,
        ):
            price_sum = receive(reporter_index, data_source, symbol)
            data_report_count += 1
    return price_sum / data_report_count
```

### Data Requests
[sec:term-data-request]: #DataRequests

A data request is a transaction from a user to perform a data query based on an oracle script. A data request transaction specifies the oracle script to execute, the parameters to the script, and [other security parameters][sec:msg-request-data].

### Raw Data Requests
[sec:term-raw-data-request]: #RawDataRequests

Raw data requests are requests to primary sources emitted while an oracle script is being executed in the first phase. It essentially consists of a data source's procedure and the associated parameters. Raw data requests are expected to be resolved by D3N block validators in the form of raw data reports. 

### Raw Data Reports
[sec:term-raw-data-report]: #RawDataReport

Raw data reports are the results from resolving raw data requests by D3N block validators. The raw reports are submitted to D3N. Once sufficient number of reports are collected, they will be used in the oracle script's second phase to compute the final result of the data request.

### Oracle Data Proof
[sec:term-oracle-data-proof]: #OracleDataProof

Once the aggregation is complete, the final result of the data request is stored permanently in D3N's global state. As similar to most other blockchains, the whole state of D3N can be represented as a [Merkle root hash](https://en.wikipedia.org/wiki/Merkle_tree). An oracle data proof is a merkle proof that shows the existence of the final result of the data request with other information related to it, including the oracle script hash, the parameters, the time of execution, etc.

## Data Request Lifecycle
What is unique about D3N is that it natively supports external data queries. All other aspects (such as asset transfering, staking, slashing, etc) are similar to other Cosmos-like blockchains. The lifecycle of data requests is described below.

**Before data requests can be made**

1. [Data sources][sec:term-data-sources] related to the request must be published to the system.
2. The [oracle script][sec:term-oracle-scripts] that encodes the data request must also be published.

**For each individual data request**

![](https://i.imgur.com/EIIl1MY.jpg)

1. A participant sends [data request][sec:term-data-request] to the network by broadcasting [MsgRequestData][sec:msg-request-data]. The message includes the oracle script it wants to invoke and other additional parameters.
2. Once the transaction is confirmed, the oracle script's preparation function will be run in a decentralized manner. The function will emit the set of [raw data requests][sec:term-raw-data-request] necessary to continue the oracle script's execution.
4. D3N block validators inspect the raw data requests and execute the associated data sources' procedures as required by the raw data requests.
5. Block validators sends [raw data reports][sec:term-raw-data-report] to the network by broadcasting [MsgReportData][sec:msg-report-data]. The raw data reports are stored temporarily on the blockchain.
6. Once sufficient block validators (as specified in the data request's security parametes) have reported the raw data points, D3N resumes the execution of the oracle script to aggregate raw data points into the final result, which will be stored on the blockchain permanently.
7. The final result becomes available on the blockchain's state tree. The data can be sent to other blockchains through D3N's [inter-blockchain communication][sec:inter-chain-architecture].

# Token Economics
Below is the proposed token economics of D3N. Note that the blockchain will be implemented in such the way that all the numbers are adjustable through on-chain governance.

## BAND Token
D3N utilizes is native token BAND to incentivize the block validators to produce new blocks and submit responses to data requests. BAND token holders will be able to stake to become a block validators or delegate their staking power to another block validator to earn portion of the collected fees and inflationary rewards. Additionally, BAND token holders can participate in the governance of D3N blockchain.

## Inflationary Rewards
As a sovereign blockchain, D3N uses inflationary model to incentivize participants to stake and contribute to the security of the network. Anyone can apply to become block validators given sufficient votes from the community in terms of bonded tokens. Stakers that support active validators enjoy the inflated supply of Band tokens. Inflation rate ranges from 7% to 20% with the staking target of 70% of the total BAND supply.

## Transaction Fees
BANDs must be paid as the transaction fee for a transaction to be processed in the network. Each block validator can subjectively set the minimum gas fee for processing the transaction, and choose whatever transactions it wants to include in the block it is proposing, as long as the total gas limit is not exceeded.

## Community Development Funds
2% of block rewards are diverted to community fund pool. The funds are intended to promote long-term sustainability of the ecosystem. These funds can also be distributed in accordance with the decisions made by the governance system.

## Token Burn Events
50% of the gas fee collected in terms of BANDs from all of the transactions that are processed on the D3N blockchain are burned permanently. The mechanism forces the supply of BAND to contract as the utilization of the network grows, ensuring that BAND's token value is correlated with network utilization.

## Stablecoins
While BAND tokens are available as fee tokens to pay for data source queries. In reality, data vendors may prefer to accept payment in [stablecoins](https://en.wikipedia.org/wiki/Stablecoin). D3N will support stablecoins natively to faciliate payments in stablecoins. At this stage, there are two main mechanics to support stablecoins in consideration 

- **BAND-backed stablecoins**: Similar to [MakerDAO](https://makerdao.com/) on Ethereum, the system can support minting new stablecoins (i.e. DAI counterparts) using BAND tokens as collateral. The same liquidation mechanism applies here.
- **Bridge from other networks**: Alternatively, D3N can support bridging in stable tokens from other networks, such as DAI, USDC, or Binance USD on [Ethereum](https://ethereum.org/), or Terra on [Terra Network](https://terra.money/).

Both mechanics are viable and the team is actively researching on the long-term economic  and technical consequences of adding them to D3N blockchain.

# Incentivize Non-Public Data Providers
For data that is public available for free, such as Bitcoin transaction formation or weather information provided by a non-profit organization, It is obvious that D3N block validators can access the information without requiring permissions. For information guarded by centralized parties, however, additional layers are required.

## Payment Gateway
For API endpoints that require paying fee to access. D3N supports this by allowing data source owners to set fees that will be paid when an oracle script asks for the data. Once the fee is collected, the blockchain will emit an unforgeable *payment receipt* as a record. Block validators can show the receipt to the API provider to retrieve data. API provider will be able to withdraw the collected fees anytime. Note that the system relies on the API provider acknowledging the payment and serving the data properly. It is putting its reputation at stake and risk losing future revenue if it acts maliciously.

## P2P Authentication Gateway
Being a public and permissionless blockchain, it is naturally tricky for D3N to handle private information. Consider the following situation: a person would like to confirm to a smart contract that she has the credit score beyond 700, according to a credit bureau. There are 

1. How can the smart contract trust that the data actually comes from the intended entity (the credit bureau agency).
2. How can the user guarantee that the private data is not revealed beyond what is necessary.
3. How can the credit bureau verify that the individual gives consent to share the information to the public without revealing the person's identity publicly?

Issue 1) is solved by the fact that D3N relies on multiple parties to query to the same entity. Thus, if the majority of block validators stay honest, data must be valid. Issue 2) is solved by the fact that the data source can attach arbitrary source code to ensure that only necessary data is needed to be commited publicly as raw data reports.

As for issue 3), D3N cannot solve the problem directly. We, however, introduce the concept of peer-to-peer authentication gateway. Each individual is expected to run her own, or give authorization to a host provider, the resolving server that the oracle endpoint can query to verify the request's authencity. In this specific example, once the credit bureau receives the query request, it goes to the user's authentication gateway to confirm the user's consent. The idea is analogously similar to [Solid](https://solid.mit.edu) or [Blockstack](https://blockstack.org/). We are working on the formal specification for this gateway and will update this section once the information becomes available.

# Development Phases
D3N development is publicly available on [Github](https://github.com/bandprotocol/d3n) for public inspection. The technical development is breaking down into different phases. The content below is provided for informational purpose and may change in the future as the software gets developed. Also note that apart from the core protocol, the foundation will also develop related supporting tools, including client side libraries, Owasm IDE, block explorer, etc.

**Phase 0** (Q1 2020)
- Supports token staking.
- Supports on-chain governance parameters.
- Supports asset transfers.

**Phase 1** (tentatively Q2 2020)
- Supports querying public and permissionless data sources.
- Supports gas metering on WebAssembly scripts on-chain.

**Phase 2** (tentatively Q3 2020)
- Supports querying data sources while collecting fees to data source owners.
- Supports function invokes between different oracle scripts.

**Phase 3** (tentatively Q4 2020)
- Supports stablecoins in the ecosystem.
- Supports peer-to-peer authentication gateway.

# Protocol Transactions

## D3N Specific Transactions
Band Protocol is built based on [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) framework, and thus  inherits all transactions available on the Cosmos framework. In addition, the transactions related to performing oracle operations are added. The specification is presented below.

### MsgCreateDataSource
[sec:msg-create-data-source]: #MsgCreateDataSource

Deploys a new data source to D3N network. Once deployed a data source is assigned a unique `int64` identifier which can be used to refer to it forever. Note that if you want data source to stay immutable, simply set `Owner` to be a value you believe that no one on earth knows the matching public key.

**Parameters**
- `Owner: sdk.AccAddress` - The address responsible for maintaining the data source. The owner will be able to withdraw query fees and update.
- `Name: string` - The human-readable string name for this data source. 
- `Fee: sdk.Coins` - The fee that data requester needs to pay per one data query.
- `Executable: []byte` - The content of executable to be run by block validators upon receiving a data request for this data source. The executable can be in any format, as long as it is accepted by the general public.
- `Sender: sdk.AccAddress` - The sender of this transaction. Note that the sender does not need to be the same as the owner.

### MsgEditDataSource
Edits an existing data source given the unique `int64` identifier. The sender must be the owner of the data source for this transaction to succeed.

**Parameters**
- `DataSourceId: int64` - The unique identifier of the data source to edit.
- `Owner: sdk.AccAddress` - The new owner of the data source.
- `Name: sdk.AccAddress` - The new name of the data source.
- `Fee: sdk.Coins` - The new fee of the data source.
- `Executable: []byte` - The new executable content of the data source.
- `Sender: sdk.AccAddress` - The sender of this transaction. Must be the current owner of the data source.

### MsgCreateOracleScript
[sec:msg-create-oracle-script]: #MsgCreateOracleScript

Deploys a new oracle script to the blockchain network. Once deployed, an oracle script is given a unique `int64` identifier and anyone in the network can request for data with it. Similar to data sources, oracle scripts can have an owner who is allowed to patch it.

**Parameters**
- `Owner: sdk.AccAddress` - The address responsible for maintaining the oracle script.
- `Name: string` - The human-readable string name for this oracle script.
- `Code: []byte` - The [Owasm][sec:owasm] compiled binary attached to this oracle script.
- `Sender: sdk.AccAddress` - The sender of this transaction. Note that the sender does not need to be the same as the owner.

### MsgEditOracleScript
Edits an existing oracle script given the unique `int64` identifier. The sender must be the owner of the oracle script for this transaction to succeed.

**Parameters**
- `OracleScriptId: int64` - The unique identifier of the oracle script to edit.
- `Owner: sdk.AccAddress` - The new owner of the oracle script.
- `Name: string` - The new name of the oracle script.
- `Code: []byte` - The new [Owasm][sec:owasm] compiled binary of the oracle script.
- `Sender: sdk.AccAddress` - The sender of this transaction. Must be the current owner of the oracle script.

### MsgRequestData
[sec:msg-request-data]: #MsgRequestData

Requests a new data based on an existing oracle script. A data request will be assigned a unique identifier once the transaction is confirmed. After sufficient block validators report raw data points. The result of data request will be written permanently to the blockchain for further uses.

**Parameters**
- `OracleScriptId: int64` - The unique identifier of the oracle script.
- `Params: []byte` - The data to passed over to the oracle script when invoking its [prepare][sec:owasm-prepare] function.
- `RequestedValidatorCount: int64` - The number of block validators to perform fetching data tasks issued by this oracle script. Setting this higher may make the request safer against attacks or downtimes of validators, but will incur higher gas fees. The value must between 1 and the number of active validators.
- `SufficientValidatorCount: int64` - The number of data reports from block validators needed for this data request to be considered ready for aggregation. The number must be positive and not greater than `RequestValidatorCount` value.
- `Expiration: int64` - The number of blocks that the data request stays active and open for data reports. If the number of validator count is not reached by then. The data request is considered invalid.
- `Sender: sdk.AccAddress` - The sender of this transaction.

### MsgReportData
[sec:msg-report-data]: #MsgReportData

Reports raw data points for the given data request. Each data point corresponds to a data source query issued during the data request script's execution of [`prepare`][sec:owasm-prepare] function.

**Parameters**
- `RequestId: int64` - The unique identifier of the data request.
- `Data: []struct{ externalDataId: int64, data: []byte }` - The array of raw data points. Each element corresponds to a data source query.
- `Sender: sdk.ValAddress` - The sender of this transaction. Must be one of the block validators that are entitied to report data to the data request.

## Gas Fee
Each transaction broadcasted to D3N consumes certain amount of gas, depending on the following three factors.

- **Network bandwidth** - The size of transaction message in bytes. 
- **Instruction count** - The number of WebAssembly operations required to complete the transaction.
- **Memory** - The amount of storage that the transaction uses to store data persistently.

A transaction must specify the maximum gas it intends to spend with the transaction. The transaction is considered failed if it runs out of gas before the transaction is complete. Each block has a maximum gas limit, thus limitting the number of transactions the network can process in a block.

Note that memory is different from the other two in the sense that it continously costs block validators to store the data while the gas is being paid only once. D3N may change gas model by incorporating state rent mechanics to incentivize participants to clean-up unused memory slots in future iterations.

# Inter-Chain Architecture
[sec:inter-chain-architecture]: #Inter-Chain-Architecture

The goal of D3N is to eventually pass over data requested and aggregated on the blockchain to the target location. In this section, we discuss the implementation of the architecture that enables the inter-chain communication.

## Merkle Proofs for Light Client Validation
Data requests and their corresponding final responses are stored persistently in D3N's storage on one [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree) of D3N's multistore app state. The hash of D3N app state on every block is signed by the block validators. Bridges on other blockchains keep track of D3N's validator sets and can verify the validity of app hash. Once the app hash is verified, one can provide proof path to show that a specific piece of information is curated inside of D3N.

The set of validators on D3N blockchain may be continuously changing, and thus must be updated on the target blockchain periodically. Since the signatures of the previous block validators and the merkle root to the state tree that shows the set of current validators are both included in a block header, anyone can relay the information to the target blockchain without relying on a centralized relayer. Due to 21-day unbonding period, the target blockchain must update D3N validator set at least every 21 day.

For the details on how to traverse the iAVL state tree, please see [Merkle proof section on Cosmos whitepaper](https://github.com/cosmos/cosmos/blob/master/WHITEPAPER.md#merkle-tree--proof-specification).

## Reference Implementations
Below is the reference implementation of code to check for (1) sufficient validator signatures when relaying a new block header and (2) correct merkle path when submitting an oracle data point in different platforms. This is an active list and will be updated to reflect the current state.

| Platform        | Implementation                | 
| --------------- | ----------------------------- | 
| [EVM][evm-wiki] | [d3n/bridges/evm][bridge-evm] | 

[evm-wiki]: https://en.wikipedia.org/wiki/Ethereum
[bridge-evm]: https://github.com/bandprotocol/d3n/tree/master/bridges/evm

## Cosmos IBC
Built on top of Cosmos SDK, D3N will naturally support the [Inter-Blockchain Communication (IBC) protocol](https://cosmos.network/ibc/) once it is fully implemented and integrated to the Cosmos SDK codebase. Once integrated, oracle data inside of D3N will be natively accessible to other blockchains that follow the IBC specification without additional overhead. 

# Possible Attack Vectors

## Consensus Level Attacks
As D3N is built on top of Cosmos-SDK and inherits similar security assumptions based on the [Tendermint](https://tendermint.com/) BFT consensus algorithm. In essense, the system stays secure and available as long as not less than 2/3 of the total bonded voting power stays honest. Malicious parties will get slashed or risk losing their token value. For further discussions on different types on attacks and mitigations on the consensus level, see [here](https://github.com/cosmos/cosmos/blob/master/WHITEPAPER.md#appendix).

## Oracle Data Inaccuracy
D3N system on its own does not penalize block validators on reporting "bad" data, since there is no deterministic and fully decentralized way to proof that a validator intentionally tampers with data from the original source.  However, if such attack occurs and data becomes less useful, the value of token is essentially destroyed since there will no longer be any dApps willing to pay for such data. Unbonding period disallows block validators to liquidate their BANDs immediately after the attack. This ensures that block validators suffer the most from the system's collapse. The credible threat of economic loss should be sufficient to discourage community-wide collusion of block validators. Additionally, real-world reputation loss from collusion also serves as an incentive to prevent data providers from acting maliciously. In the future, we also consider imposing token slashing based on on-chain voting to further dis-incentivize dishonest behavior.

In addition, D3N allows data requesters to specify any type of aggregation logic to their oracle scripts. Thus, if they want data to be safe against the attacks from some of the block validators, they may write code that only consider the result usable of it is agreed by most (or even all) validators. Note that the tradeoff is that if some validators go down temporarily, the script might not run to completion despite everyone being honest. This is the flexibility that D3N provides and participants must use their own judgements based on how much they value security vs liveness for their applications.

## State Bloat
D3N does store a huge volume of data, including data sources, oracle scripts, raw data requests/reports, and all query results. As the current gas fee model only charges when the transaction is included in a block without incentivizing users to free-up the space, the blockchain state will grow. In addition, malicious validator may set gas fee to zero and intentionally add more data for everyone to store.

There is a common problem in many blockchain systems, in particular with [Ethreum](https://ethereum.org) and [EOS](https://eos.io). A common research direction is to incorporate state rent model to incentivize users to free-up space or otherwise continuously get charged. The team will be actively researching and monitoring D3N blockchain usage and will upgrade the blockchain should state size become a real issue.

# Oracle WebAssembly (Owasm)
[sec:owasm]: #Oracle-WebAssembly-Owasm

 To support arbritrary oracle requests, D3N allows data requesters to upload [WebAssembly](https://webassembly.org/) binaries with request payloads. Owasm is the standard on top of on top of WebAssembly code for it to be universally understandable as a script for off-chain oracle execution. It specifies the data types to support, the functions the oracle script must expose, and the interface functions that an oracle script can access.

## Data Types
At the protocol level, Owasm only support basic WebAseembly types: `i32`, `i64`, `f32`, and `f64`. In addition, a contiguous span of memory (`bytes`) is represented with an `i32` which is a WebAssembly offset. For the sake of clarity, we will use `i32ptr` for an `i32` that should be intepretted as a memory offset. How to interpret `bytes` into program-specific data structures is left to be defined by users and agreed by other standards.

## Owasm Script Interface (OSI)
Owasm imposes extra. If an oracle script is intended to be a reusable library , it needs not to implement the required functions.

### `prepare() -> ()` 
[sec:owasm-prepare]: #prepare--gt-

**Required** - This function is called once every time a new request is made to the Owasm script via [MsgDataRequest][sec:msg-request-data]. Note that the request parameters of are passed through calldata. The function is expected to call [requestExternalData][sec:owasm-request-external-data] function multiple times depending on the number of data sources it needs.

### `execute(i32, i32) -> (u64ptr)`

**Required** - This function is called after sufficient data points are reported for a data request. Note that the request parameters of are passed through calldata. The function is expected to return the final aggregated data through the call of [finish][sec:owasm-finish] function.

**Parameters**
- `requestedValidatorCount: i32` The number of validators that was specified to 
- `receivedValidatorCount: i32` The number of validators that actually report raw data points prior to this function call.

### `getParametersInfo() -> ()`
**Optional** - Gets the human-readable string representing the specification of this script's parameters. The function invokes [finish][sec:owasm-finish] with the return value.

### `getResultInfo() -> ()`
**Optional** - Gets the human-readable string representing the specification of this script's final result. The function invokes [finish][sec:owasm-finish] with the return value. 

## Oracle Environment Interface (OEI)
The following functions are part of Owasm oracle environment interface, which are accessible to oracle scripts during their executions.

### `finish(i32ptr, i32) -> ()`
[sec:owasm-finish]: #finishi32ptr-i32--gt-

Sets the returning output data for the execution, causing a trap and immediately terminating the execution with success result.

**Parameters**
- `dataOffset: i32ptr` - The memory offset of the returning output data
- `dataLength: i32` - The size of the output data in bytes

### `fail(i32ptr, i32) -> ()`
Terminates the current execution, causing a trap and immediately terminating the execution with failure result.

**Parameters**
- `errMsgOffet: i32ptr` - The memory offset of the returning output data
- `errMsgLength: i32` - The size of the output data in bytes

### `call(i64, i64, i32ptr, i32) -> i32`
[sec:owasm-call]: #calli64-i64-i32ptr-i32--gt-i32

Performs a call into another oracle script. The callee will use a separate memory space. The function will return after the function call is complete, whether with success or failure. The call depth size is limited at 127.

**Parameters**
- `target: i64` - The unique identifier of the oracle script to call into.
- `funcNameOffset: i32ptr` - The starting memory offset of function name.
- `funcNameSize: i32` - The size of function name in bytes.
- `dataOffset: i32ptr` - The starting memory offset of calldata to pass through.
- `dataSize: i32` - The size of calldata in bytes.

**Returns**
- `status: i32` - The status of the call: `0` on success and `1` on failure.

### `callStatic(i64, i64, i32ptr, i32) -> i32`
Performs a call similar to [call][sec:owasm-call] function, but disallows the target script and its children to emit raw data requests.

**Parameters**
- `target: i64` - The unique identifier of the oracle script to call into.
- `funcNameOffset: i32ptr` - The starting memory offset of function name.
- `funcNameSize: i32` - The size of function name in bytes.
- `dataOffset: i32ptr` - The starting memory offset of calldata to pass through.
- `dataSize: i32` - The size of calldata in bytes.

**Returns**
- `status: i32` - The status of the call: `0` on success and `1` on failure.

### `getCallDataSize() -> i32`
Gets the size of calldata set by the current function's caller in bytes. For the outermost function call, the data is always the data request's encoded parameters.

**Returns**
- `dataSize: i32`: - The size of the calldata in bytes.

### `readCallData(i32ptr, i32, i32) -> ()`
Reads the specific part of calldata into the runtime memory. The function will fail if `seekOffset + resultSize > getCallDataSize()`.

**Parameters**
- `resultOffset: i32ptr` - The starting memory offset to copy calldata into.
- `seekOffset: i32` - The starting offset of calldata to start copying.
- `resultSize: i32` - The number of bytes to copy.

### `getReturnDataSize() -> i32`
Gets the size of return data for the last `call` or `callStatic`.

**Returns**
- `dataSize: i32`: - The size of the return data data from the last call in bytes.

### `readReturnData(i32ptr, i32, i32) -> ()`
Reads the specific part of return data into the runtime memory. The function will fail if `seekOffset + resultSize > getReturnDataSize()`. 

**Parameters**
- `resultOffset: i32ptr` - The starting memory offset to copy return data into.
- `seekOffset: i32` - The starting offset of return data to start copying.
- `resultSize: i32` - The number of bytes to copy.

### `requestExternalData(i64, i64, i32ptr, i32)`
[sec:owasm-request-external-data]: #requestExternalDatai64-i64-i32ptr-i32

Issues a raw data request for a specific data source to block validators. The function can only be called during the preparation phase of an oracle script.

**Parameters**
- `dataSourceId: i64` - The unique identifier of the data source to request data.
- `externalDataId: i64` - The unique identifer for this raw data request.
- `dataOffset: i32ptr` - The starting memory offset of the parameter to data source.
- `dataLength: i32` - The size of the parameter in bytes.

### `getExternalDataSize(i64, i32) -> i32`
Gets the size of raw data report for a specific raw data request from a specific block validator.

**Parameters**
- `externalDataId: i64` - The unique identifier of the raw data request.
- `validatorIndex: i32` - The index of block validators. 0 for the first, 1 for the second validator, and so on. Must be less than the number of validators that report data.

### `readExternalData(i64, i32, i32ptr, i32, i32) -> ()`
Reads the specific part of raw data reports into the runtime memory. The function will fail if `seekOffset + resultSize > getExternalDataSize()`. 

**Parameters**
- `externalDataId: i64` - The unique identifier of the raw data request.
- `validatorIndex: i32` - The index of block validators. 0 for the first, 1 for the second validator, and so on. Must be less than the number of validators that report data.
- `resultOffset: i32ptr` - The starting memory offset to copy raw report data into.
- `seekOffset: i32` - The starting offset of raw report data to start copying.
- `resultSize: i32` - The number of bytes to copy.

