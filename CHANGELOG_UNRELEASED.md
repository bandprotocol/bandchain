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

- (feat) [\#2829](https://github.com/bandprotocol/bandchain/pull/2829) Added `price_symbols` endpoint
- (bugs) [\#2828](https://github.com/bandprotocol/bandchain/pull/2828) Add missing height on `multi_request_search` and `request_prices` endpoints
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

- (impv) [\#2798](https://github.com/bandprotocol/bandchain/pull/2798) Added the warning msg to max button on send/delegate msg
- (impv) [\#2766](https://github.com/bandprotocol/bandchain/pull/2766) Revamp undelegate, redelegate, withdraw reward and vote modal
- (chore) [\#2767](https://github.com/bandprotocol/bandchain/pull/2767) Revamp send token modal
- (chore) [\#2749](https://github.com/bandprotocol/bandchain/pull/2749) Revamp delegate modal
- (impv) [\#2727](https://github.com/bandprotocol/bandchain/pull/2727) Revamp connect modal
- (impv) [\#2784](https://github.com/bandprotocol/bandchain/pull/2784) Revamp every state of confirmation popup
- (bugs) [\#2783](https://github.com/bandprotocol/bandchain/pull/2783) Validate address prefix
- (bugs) [\#2758](https://github.com/bandprotocol/bandchain/pull/2758) Fixed the wrapped button issue on safari
- (impv) [\#2754](https://github.com/bandprotocol/bandchain/pull/2754) Polish UI, remove unused
- (impv) [\#2747](https://github.com/bandprotocol/bandchain/pull/2747) Patch for fast sync and new GuanYu DB
- (impv) [\#2740](https://github.com/bandprotocol/bandchain/pull/2740) Add UI test for all pages
- (impv) [\#2731](https://github.com/bandprotocol/bandchain/pull/2731) Add no transaction placeholder and fix mobile margin
- (impv) [\#2728](https://github.com/bandprotocol/bandchain/pull/2728) Add delegation, redelegation, undelegation and withdraw reward UI test
- (impv) [\#2725](https://github.com/bandprotocol/bandchain/pull/2725) Show validator name if voter is validator
- (impv) [\#2719](https://github.com/bandprotocol/bandchain/pull/2719) Implement query Band supply
- (bugs) [\#2716](https://github.com/bandprotocol/bandchain/pull/2716) Fixed no decimal balance bug on account index
  page
- (impv) [\#2714](https://github.com/bandprotocol/bandchain/pull/2714) Handle NotFound case for Datasources and Oracle-scripts page
- (chore) [\#2713](https://github.com/bandprotocol/bandchain/pull/2713) Remove redundant emptyContainer styles
- (chore) [\#2695](https://github.com/bandprotocol/bandchain/pull/2695) Remove unused images and replace with icon.
- (impv) [\#2693](https://github.com/bandprotocol/bandchain/pull/2693) Created the new button component and patched to all buttons
- (impv) [\#2689](https://github.com/bandprotocol/bandchain/pull/2689) Styled the top part of account index page
- (bugs) [\#2687](https://github.com/bandprotocol/bandchain/pull/2687) Fixed NaN number (Urgent)
- (impv) [\#2684](https://github.com/bandprotocol/bandchain/pull/2684) Adjust validator voted from 250 to 100
- (impv) [\#2676](https://github.com/bandprotocol/bandchain/pull/2676) Setup cypress to scan, added new sendToken testcase
- (bugs) [\#2673](https://github.com/bandprotocol/bandchain/pull/2673) Fixed sorting on moniker with emoji
- (impv) [\#2672](https://github.com/bandprotocol/bandchain/pull/2672) Hid the proposal desc from tx index page table
- (impv) [\#2671](https://github.com/bandprotocol/bandchain/pull/2671) Adjusted tooltip width and added webapi
- (impv) [\#2670](https://github.com/bandprotocol/bandchain/pull/2670) Format the count number
- (feat) [\#2669](https://github.com/bandprotocol/bandchain/pull/2669) Add active validator's rank on Validator Index Page
- (bugs) [\#2666](https://github.com/bandprotocol/bandchain/pull/2666) Handle rate limit msg when decoding
- (impv) [\#2663](https://github.com/bandprotocol/bandchain/pull/2663) Fixed and updated `ChainIDBadge`
- (impv) [\#2644](https://github.com/bandprotocol/bandchain/pull/2644) Added `netlify.toml` configuration
- (feat) [\#2594](https://github.com/bandprotocol/bandchain/pull/2594) Added meta og tag to scan

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
- (impv) [\#2802](https://github.com/bandprotocol/bandchain/pull/2802) pyband: Transaction class, add test cases
- (feat) [\#2799](https://github.com/bandprotocol/bandchain/pull/2799) pyband: Implemented Message class
- (impv) [\#2789](https://github.com/bandprotocol/bandchain/pull/2789) pyband: get chain id from specific rest endpoint.
- (impv) [\#2739](https://github.com/bandprotocol/bandchain/pull/2739) pyband: fix client raise error when get fail
- (impv) [\#2652](https://github.com/bandprotocol/bandchain/pull/2652) pyband: use string instead of class annotation for Python3.6
- (bugs) [\#2651](https://github.com/bandprotocol/bandchain/pull/2651) pyband: fix bug get latest block

### MISC
