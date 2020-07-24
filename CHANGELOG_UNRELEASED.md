<!--
(feat): New feature
(impv): Improvement / Enhancement
(docs): Documentation
(bugs): Bug fixes
(chore): Chore/cleanup work
-->

## [Unreleased]

### Chain (Consensus)

- (chain) [\#2333](https://github.com/bandprotocol/bandchain/pull/2333) Upgrade to Cosmos-SDK version 0.39.1.

### Chain (Non-consensus)

- (impv) [\#2332](https://github.com/bandprotocol/bandchain/pull/2232) Emit gas used as an attribute event during prepare and execute.
- (bugs) [\#2297](https://github.com/bandprotocol/bandchain/pull/2297) Update json key name of request and response packets.

### Yoda

- (impv) [\#2307](https://github.com/bandprotocol/bandchain/pull/2307) Add Yoda configurable timeout duration.

### Emitter & Flusher

- (impv) [\#2319](https://github.com/bandprotocol/bandchain/pull/2319) Add index on blocks table and swap order of primary key of validator_votes table.
- (impv) [\#2302](https://github.com/bandprotocol/bandchain/pull/2302) Add offset check before sync flusher.
- (bugs) [\#2298](https://github.com/bandprotocol/bandchain/pull/2298) Fix bug `accumulated_commission` in `emitSetValidator`.
- (bugs) [\#2295](https://github.com/bandprotocol/bandchain/pull/2295) Truncate `accumulated_commission` precision.

### Scan

- (impv) [\#2334](https://github.com/bandprotocol/bandchain/pull/2334) Implemented the sorting function on validator homepage's mobile layout.
- (impv) [\#2330](https://github.com/bandprotocol/bandchain/pull/2330) Fixed share_percentage decoder in DelegationSub
- (impv) [\#2317](https://github.com/bandprotocol/bandchain/pull/2317) Implemented account Index Page (Mobile)
- (impv) [\#2315](https://github.com/bandprotocol/bandchain/pull/2315) Improved how to pass account type on the AddressRender component
- (impv) [\#2312](https://github.com/bandprotocol/bandchain/pull/2312) Implemented the BlockIndexPage layout for mobile view
- (impv) [\#2316](https://github.com/bandprotocol/bandchain/pull/2316) Implemented the ValidatorHomePage layout for mobile view
- (impv) [\#2310](https://github.com/bandprotocol/bandchain/pull/2310) Implemented the TxIndexpage layout for mobile view
- (impv) [\#2305](https://github.com/bandprotocol/bandchain/pull/2305) Implement the TxHomepage layout for mobile view and adjusted the pagination on mobile view.
- (impv) [\#2313](https://github.com/bandprotocol/bandchain/pull/2313) Added commision amount on Account Index Page
- (feat) [\#2294](https://github.com/bandprotocol/bandchain/pull/2294) Implemented top part of `ValidatorIndexPage` for mobile
- (impv) [\#2299](https://github.com/bandprotocol/bandchain/pull/2299) Update the latest transactions table for mobile version.
- (bugs) [\#2290](https://github.com/bandprotocol/bandchain/pull/2290) Fix average block time calculation on `ValidatorHomePage` when using new cacher
- (feat) [\#2296](https://github.com/bandprotocol/bandchain/pull/2296) Implemented delegators and proposed blocks table in `ValidatorIndexPage` mobile version

### Bridges

### Runtime

### Owasm

### Oracle Binary Encoding (OBI)

### Helpers

- (feat) [\#2301](https://github.com/bandprotocol/bandchain/pull/2301) Add `pyband` initial implementation.

### MISC

- (chore) [\#2279](https://github.com/bandprotocol/bandchain/pull/2279) Update `yoda` README.
