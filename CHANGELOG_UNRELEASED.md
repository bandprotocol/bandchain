<!--
(feat): New feature
(impv): Improvement / Enhancement
(docs): Documentation
(bugs): Bug fixes
(chore): Chore/cleanup work
-->

## [Unreleased]

### Chain (Consensus)

- (bugs) [\#2251](https://github.com/bandprotocol/bandchain/pull/2251)  go-owasm: Allow OEI to read nil external data

### Chain (Non-consensus)

- (impv) [\#2232](https://github.com/bandprotocol/bandchain/pull/2218) CLI/REST for query active oracle validators.

### Yoda

- (impv) [\#2247](https://github.com/bandprotocol/bandchain/pull/2247) Use max data size configurations from oracle module
- (impv) [\#2249](https://github.com/bandprotocol/bandchain/pull/2249) docker: Run test program during initialization.
- (feat) [\#2218](https://github.com/bandprotocol/bandchain/pull/2218) Implement MultiExec to combine multiple executors.

### Emitter & Flusher

- (feat) [\#2248](https://github.com/bandprotocol/bandchain/pull/2248) Add identity column on delegation view table.
- (bugs) [\#2273](https://github.com/bandprotocol/bandchain/pull/2273) Fix bug update delegators table after withdraw reward.
- (bugs) [\#2255](https://github.com/bandprotocol/bandchain/pull/2255) Fix bug `reward_amount` and `commission_amount` in extra field.
- (bugs) [\#2252](https://github.com/bandprotocol/bandchain/pull/2252) `handle_set_validator` get wrong validator id.
- (impv) [\#2250](https://github.com/bandprotocol/bandchain/pull/2250) Add account id in `validators` table.
- (feat) [\#2246](https://github.com/bandprotocol/bandchain/pull/2246) Implement handle all resolve proposal status.
- (feat) [\#2242](https://github.com/bandprotocol/bandchain/pull/2242) Implement handle MsgVote for emitter and flusher.
- (feat) [\#2241](https://github.com/bandprotocol/bandchain/pull/2241) Implement handle MsgDeposit for emitter and flusher.
- (impv) [\#2240](https://github.com/bandprotocol/bandchain/pull/2240) Emit Kafka msg for staking genesis state.
- (feat) [\#2238](https://github.com/bandprotocol/bandchain/pull/2238) Implement handle MsgSubmitProposal for emitter and flusher.

### Scan

- (feat) [\#2276](https://github.com/bandprotocol/bandchain/pull/2276) Implemented mobile version of homepage top part
- (impv) [\#2245](https://github.com/bandprotocol/bandchain/pull/2245) Added more features and patch for guanyu testnet
- (impv) [\#2237](https://github.com/bandprotocol/bandchain/pull/2237/files) Add validator's oracle status.
- (bugs) [\#2236](https://github.com/bandprotocol/bandchain/pull/2236) Fixed uptime query on ValidatorIndexPage
- (impv) [\#2235](https://github.com/bandprotocol/bandchain/pull/2235) added `unbonding` and `redelegate` tabs to account index page.
- (impv) [\#2234](https://github.com/bandprotocol/bandchain/pull/2234) Added support for `received` transaction on scan
- (impv) [\#2228](https://github.com/bandprotocol/bandchain/pull/2228) Add `expired` request status to Scan
- (impv) [\#2203](https://github.com/bandprotocol/bandchain/pull/2203/files) Patch request, report subscription for new cacher.
- (impv) [\#2199](https://github.com/bandprotocol/bandchain/pull/2199) Remove proposed blocks count & adjust ui
- (impv) [\#2176](https://github.com/bandprotocol/bandchain/pull/2176/files) Fixed delegations/unbonding subs, avg blocktime and validator uptime to work with new cacher.

### Bridges

### Runtime

- (impv) [\#2233](https://github.com/bandprotocol/bandchain/pull/2233) Remove max timeout variable.
- (impv) [\#2230](https://github.com/bandprotocol/bandchain/pull/2230) Update lambda function follow Remote data source executor.

### Owasm

- (impv) [\#2231](https://github.com/bandprotocol/bandchain/pull/2231) Maintain gas used and gas limit in VMConfig.

### Oracle Binary Encoding (OBI)

### Helpers

### MISC
