# ⛓ Decentralized Data Delivery Network (D3N)

<div style="text-align:center"><img style="width:400px; margin: 20px 0" src="https://i.imgur.com/Saw9Dtx.png" /></div>


[![Join the chat at https://t.me/bandprotocol](https://img.shields.io/badge/chat-on%20Telegram%20💬-blue.svg)](https://t.me/bandprotocol)

This live document outlines implementation plans and research directions for [Band Protocol](https://bandprotocol.com)'s decentralized data delivery network (D3N). The implementation is available for public inspection [here](https://github.com/bandprotocol/d3n).

## Introduction to D3N

If you have heard of Band Protocol, you probably know that we are tackling the blockchain oracle problem. The majority of existing smart contract platforms, while support trustless executation of arbitrary programs, lack access to real-world data. This renders smart contracts not as useful. D3N connects public blockchains with off-chain information, with the following design goals.

1. **Speed and Scalability:** The system must be able to serve a large quantity of data requests to multiple public blockchains with minimal latency.
2. **Cross-Chain Compatibility:** The system must be blockchain-agnostic platform and can can serve data to all available public blockchains.
3. **Data Flexibility:** The system must be generic and able to support different ways to fetch and aggregate data.

## How Do Data Requests in D3N Work?

D3N is a high-performance public blockchain that enables participants to request for arbitrary off-chain data and computation. Built on top [Tendermint](https://tendermint.com/) and [Cosmos SDK](http://cosmos.network/), the blockchain utilizes [BFT](https://en.wikipedia.org/wiki/Byzantine_fault) consensus algorithm to reach immediate finality upon getting confirmations from more than 2/3 of block validators. In D3N, block validators are chosen based on the number of stakes they have. In addition to signing blocks, block validators are also responsible for reponding to data requests submitted to the system.

![](https://i.imgur.com/EIIl1MY.jpg)

D3N's main functionality is that it allows arbitrary data request on the blockchain. The lifecycle of the process is described below. Note that other blockchain aspects (such as asset transfer, stake, slash) of D3N are similar to other Cosmos-like blockchains. 

1. A participant submits a data *request* transaction containing an oWASM script that specifies the data he or she wants to obtain.
2. Once the transaction is confirmed (expected to be within one second), block validators inspect the script and performs off-chain computations to retrieve data required to compute the final result.
3. Block validators submit *report* transactions to the blockchain.
4. After passing the timing threshold to report data, the blockchain computes the final result of the data request using data points reported by block validators.
5. The final result becomes available on the blockchain's state tree. The data can be sent to other blockchains through [IBC](https://cosmos.network/docs/spec/ibc/) or blockchain bridges.

## Data Request Specification

In D3N, a data request is written as a [Turing-complete](https://en.wikipedia.org/wiki/Turing_completeness) program that can be complied to a specifialized variant of [WebAssembly](https://webassembly.org/), called oWASM. oWASM bytecodes are similar to that of standard WASM bytescodes, with the following added constraints.

1. The program must export *three* toplevel functions: `__alocate`, `__prepare`, and `__execute`. The specification of each function is described below.
2. Arrays are passed around as an `i32` integer representing the memory location of the data. Data has its first 4 bytes encoding the length of the array in [little-endian](https://en.wikipedia.org/wiki/Endianness) format followed by the actual data.
3. All features that may lead to non-determinism are disabled.

### Allocation Function

`__allocate` function takes one `i32` argument, representing the size of data, in bytes, that the caller would like to allocate. The function allocates a contiguous memory space of the requested size, and returns the location of the first byte as an `i32`. 

### Preparation Function

Upon confirming the data request transaction, D3N runs `__prepare` function of the request script. The function takes no arguments and returns a `u8` array that JSON-encodes data required to perform execution.

### Execution Function

`__execute` is called once the reporting period ends. The function takes two arguments: an array of `u8` arrays, each representing the JSON-encoded data output that each block validator reports, and an `i32` integer representing the total number of validators eligible to report data. The function must return a `u8` array that is the final result of the data request.

<!-- ## User Types in D3N

### 1) Data Requesters

Data Requesters are the primary end-users of the D3N ecosystem. They are the dApp users who need access to oracle data. 

### 2) Data Providers

Data Providers are the parties who get data from off-chain and report it on to D3N. They are elected by xxxx

### 3) Block Validators

Top N entities by the number of votes are eligible for producing and voting for blocks in the D3N system. Block validators are rewarded by together a valid block which is consent and signed by at least ⅔ of all the validators.

In addition to signing for blocks, block validators are responsible for 
 -->
 
## BAND Native Token

D3N utilizes is native token BAND to incenvize the block validators to produce new blocks and responds to data requests. Band tokens can be held by accounts on D3N and can be transferred to a different account with `Transfer` transaction.

### Inflationary Schedule

As a sovereign chain, D3N uses inflrationary model incentivize participants to stake to secure the network. Anyone can apply to become block validators if given sufficient votes from the community in terms of stakes. Stakers that support active validators enjoy the inflated supply of Band tokens. Inflation rate ranges from 7% to 20% with the staking target of 70% of total BAND supply.

### Transaction Fees

To prevent denial-of-service attacks, sending a transaction to D3N requires the transaction signer to burn a certain amount of Band tokens. (TODO: Specify transaction fee)

| <!--         | Transaction Type     | Gas Usage |
| ------------ | -------------------- |
| Transfer     | 100                  |
| Store Code   | 1000 + 10 per byte   |
| Data Request | 500 + execution cost |
| Data Report  | 0                    | -->       |


## Inter-Chain Architecture (Bridge)

In essense, D3N is an independent blockchain network built specifically for data governance and delivery. Data requests are made on D3N and then transported to another blockchain to use.

This design requires that all the blockchains that request data from D3N must implement native mechanism to validate data to ensure that it's valid and come from D3N network.

There are two methods to transport data to the other blockchain
1. **Validator reported data:** data requester commission D3N's block validator to report data to the targeted blockchain via auction. Only the blockchins supported by at least one validator can receive data via this method.
2. **Requester reported data:** data requester submits the data to targeted blockchain itself.


## Protocol Transactions

This section explains the set of transactions available to be broadcasted on D3N network. Note that the majority of the messages are from Cosmos standard modules, with added transactions to support D3N's oracle features.

### MsgAllocateCode 

| Attribute | Type             | Detail                              |
| --------- | ---------------- | ----------------------------------- |
| Code      | `[]byte`         | Compiled oWASM script               |
| Proposer  | `sdk.AccAddress` | Proposer address                    |
| Deposit   | `sdk.Coins`      | Initial BAND deposit for code space |

### MsgEndorseCode

`MsgEndorseCode`

| Attribute | Type             | Detail |
| --------- | ---------------- | ------ |
| CodeHash  | `[]byte`         |        |
| Endorser  | `sdk.AccAddress` |        |
| Amount    | `sdk.Coins`      |        |

### MsgReleaseCode

`MsgWithdrawCode`

### MsgRequest

### MsgReport

| Name      | Type             | Detail                         |
| --------- | ---------------- | ------------------------------ |
| RequestID | `uint64`         | Request ID                     |
| Data      | `[]byte`         | Raw bytes represent query data |
| Validator | `sdk.ValAddress` | Provider address               |

### MsgSend
### MsgMultiSend

### MsgCreateValidator
### MsgEditValidator
### MsgDelegate
### MsgBeginRedelegate
### MsgUndelegate
### MsgSetWithdrawAddress
### MsgWithdrawDelegatorReward
### MsgWithdrawValidatorCommission
### MsgSubmitEvidence
### MsgUnjail
### MsgSubmitProposal
### MsgDeposit
### MsgVote

<!-- ## Data Sources

### Web Oracle Request Module

### TCD Module

## Data Consumption

To receive and validate the correctness of data from D3N, the blockchain must xxx

## Fee Model

The fee model is currently

## Future Research



- **Economic of Data Providers**: TODO
- **TCD and Dataset Tokens**: TODO -->