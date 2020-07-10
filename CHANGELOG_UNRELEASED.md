<!--
(feat): New feature
(impv): Improvement / Enhancement
(docs): Documentation
(bugs): Bug fixes
(chore): Chore/cleanup work
-->

## [Unreleased]

### Chain (Consensus)

### Chain (Non-consensus)

### Yoda

- (bugs) [\#2193](https://github.com/bandprotocol/bandchain/pull/2193) Fix bug start_docker.sh and start_yoda.sh.
- (feat) [\#2190](https://github.com/bandprotocol/bandchain/pull/2190) Add yoda command to bandchain.

### Emitter & Flusher

- (impv) [\#2196](https://github.com/bandprotocol/bandchain/pull/2196) Fix create view table command
- (impv) [\#2186](https://github.com/bandprotocol/bandchain/pull/2186) Add temporary view tables.
- (bugs) [\#2192](https://github.com/bandprotocol/bandchain/pull/2192) Add `commission_amount` field to Withdraw Commission Reward.
- (bugs) [\#2179](https://github.com/bandprotocol/bandchain/pull/2179) Add parse bytes to convert nil slices to empty slices.
- (impv) [\#2181](https://github.com/bandprotocol/bandchain/pull/2181) Change all foreign key that refers from tx_hash to id.
- (bugs) [\#2191](https://github.com/bandprotocol/bandchain/pull/2191) Add `reward_amount` field to Withdraw Reward Msg.
- (impv) [\#2177](https://github.com/bandprotocol/bandchain/pull/2177) Add field validator id for validators table.
- (impv) [\#2178](https://github.com/bandprotocol/bandchain/pull/2178) Add field id in accounts table.
- (bugs) [\#2180](https://github.com/bandprotocol/bandchain/pull/2180) Fix bug data source and oracle script id from genesis.
- (bugs) [\#2182](https://github.com/bandprotocol/bandchain/pull/2182) Flusher: change external id to BigInteger.
- (bugs) [\#2183](https://github.com/bandprotocol/bandchain/pull/2183) Add data_sources_id FK in raw_requests table.
- (impv) [\#2184](https://github.com/bandprotocol/bandchain/pull/2184) Add `validator_moniker` on `add_reporter` and `remove_reporter` msg.
- (impv) [\#2142](https://github.com/bandprotocol/bandchain/pull/2142) Add account transactions table.
- (impv) [\#2170](https://github.com/bandprotocol/bandchain/pull/2170) Add validator status field in table.
- (impv) [\#2169](https://github.com/bandprotocol/bandchain/pull/2169) Add unbonding and redelegation table.
- (impv) [\#2160](https://github.com/bandprotocol/bandchain/pull/2160) Add missing foreign key in report table.

### Scan

- (bug) [\#2119](https://github.com/bandprotocol/bandchain/pull/2119) Fix matching order on search.
- (bug) [\#2061](https://github.com/bandprotocol/bandchain/pull/2061) Fix wrong delegation count on account index

### Bridges

- (impv) [\#67](https://github.com/bandprotocol/bandchain/pull/2175) Patched bridge contracts to use Solidity version 0.6.11

### Owasm

### Oracle Binary Encoding (OBI)

### Helpers

### MISC

- (impv) [\#2195](https://github.com/bandprotocol/bandchain/pull/2195) Remove lib and update executable as base64 encoded.

