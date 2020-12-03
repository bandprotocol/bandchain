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

- (impv) [\#2746](https://github.com/bandprotocol/bandchain/pull/2746) Implemented emitter, price, and latest request hooks
- (impv) [\#2789](https://github.com/bandprotocol/bandchain/pull/2789) Added `bandchain/chain_id` endpoint
- (feat) [\#2757](https://github.com/bandprotocol/bandchain/pull/2757) Bring cosmos-hd-path flag
- (bugs) [\#2730](https://github.com/bandprotocol/bandchain/pull/2730) Add Content-Type header on oracle module rest endpoints
- (feat) [\#2718](https://github.com/bandprotocol/bandchain/pull/2718) Added more field in price cache
- (feat) [\#2694](https://github.com/bandprotocol/bandchain/pull/2694) Added pricer to cache latest price
- (feat) [\#2690](https://github.com/bandprotocol/bandchain/pull/2690) Added `multi_request_search`endpoint
- (feat) [\#2653](https://github.com/bandprotocol/bandchain/pull/2653) Added `verify_request` endpoint

### Yoda

### Emitter & Flusher

- (bugs) [\#2641](https://github.com/bandprotocol/bandchain/pull/2641) Fix bug flusher when update validator and remove reporter
- (impv) [\#2572](https://github.com/bandprotocol/bandchain/pull/2572) cdb: Implemented view table for track vote statistic

### Scan

- (bugs) [\#2965](https://github.com/bandprotocol/bandchain/pull/2965) Fixed typo on the transaction modal
- (bugs) [\#2927](https://github.com/bandprotocol/bandchain/pull/2927) Fix default tab on Route's search
- (feat) [\#2854](https://github.com/bandprotocol/bandchain/pull/2854) Implement reinvest

### Bridges

- (docs) [\#2691](https://github.com/bandprotocol/bandchain/pull/2691) Add simple price db example for doc
- (feat) [\#2632](https://github.com/bandprotocol/bandchain/pull/2632) Add aggregator contract for ICON bridge
- (impv) [\#2626](https://github.com/bandprotocol/bandchain/pull/2626) Icon bridge fixed from auditing process

### Runtime

### Owasm

### Oracle Binary Encoding (OBI)

### Helpers

- (impv) [\#2826](https://github.com/bandprotocol/bandchain/pull/2826) pyband: Add Pyband test on Github Action
- (impv) [\#2803](https://github.com/bandprotocol/bandchain/pull/2803) pyband: Fix typing and add PrivateKey.from_hex on Wallet
- (feat) [\#2799](https://github.com/bandprotocol/bandchain/pull/2799) pyband: Implemented Message class
- (impv) [\#2789](https://github.com/bandprotocol/bandchain/pull/2789) pyband: get chain id from specific rest endpoint.
- (impv) [\#2739](https://github.com/bandprotocol/bandchain/pull/2739) pyband: fix client raise error when get fail
- (impv) [\#2652](https://github.com/bandprotocol/bandchain/pull/2652) pyband: use string instead of class annotation for Python3.6
- (bugs) [\#2651](https://github.com/bandprotocol/bandchain/pull/2651) pyband: fix bug get latest block

### MISC
