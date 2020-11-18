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

- (feat) [\#2886](https://github.com/bandprotocol/bandchain/pull/2886) Added `pending_requests` querier
- (feat) [\#2870](https://github.com/bandprotocol/bandchain/pull/2870) fast-sync: emit all gov module from start state
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

- (impv) [\#2896](https://github.com/bandprotocol/bandchain/pull/2896) Improved yoda to handle pending requests
- (impv) [\#2857](https://github.com/bandprotocol/bandchain/pull/2857) Improved yoda gas estimation function

### Emitter & Flusher

- (bugs) [\#2641](https://github.com/bandprotocol/bandchain/pull/2641) Fix bug flusher when update validator and remove reporter
- (impv) [\#2572](https://github.com/bandprotocol/bandchain/pull/2572) cdb: Implemented view table for track vote statistic

### Scan

- (bugs) [\#2879](https://github.com/bandprotocol/bandchain/pull/2879) Fixed the address width on small desktop screen
- (bugs) [\#2868](https://github.com/bandprotocol/bandchain/pull/2868) Fixed click outside icon bug
- (impv) [\#2851](https://github.com/bandprotocol/bandchain/pull/2851) Add ESC key event listener
- (impv) [\#2850](https://github.com/bandprotocol/bandchain/pull/2850) Disable withdraw rewards if reward is zero
- (impv) [\#2858](https://github.com/bandprotocol/bandchain/pull/2858) Added click outside function to user account panel
- (impv) [\#2798](https://github.com/bandprotocol/bandchain/pull/2798) Added the warning msg to max button on send/delegate msg
- (impv) [\#2784](https://github.com/bandprotocol/bandchain/pull/2784) Revamp every state of confirmation popup
- (bugs) [\#2783](https://github.com/bandprotocol/bandchain/pull/2783) Validate address prefix
- (chore) [\#2767](https://github.com/bandprotocol/bandchain/pull/2767) Revamp send token modal
- (impv) [\#2766](https://github.com/bandprotocol/bandchain/pull/2766) Revamp undelegate, redelegate, withdraw reward and vote modal
- (impv) [\#2763](https://github.com/bandprotocol/bandchain/pull/2763) Implemented the data msg cat on tx index page
- (bugs) [\#2758](https://github.com/bandprotocol/bandchain/pull/2758) Fixed the wrapped button issue on safari
- (chore) [\#2755](https://github.com/bandprotocol/bandchain/pull/2755) Implemented Oracle tx msg restructure on txtable
- (impv) [\#2754](https://github.com/bandprotocol/bandchain/pull/2754) Polish UI, remove unused
- (chore) [\#2749](https://github.com/bandprotocol/bandchain/pull/2749) Revamp delegate modal
- (impv) [\#2748](https://github.com/bandprotocol/bandchain/pull/2748) Implemented the token msg category on txindex page (both success and fail msgs)
- (impv) [\#2747](https://github.com/bandprotocol/bandchain/pull/2747) Patch for fast sync and new GuanYu DB
- (chore) [\#2742](https://github.com/bandprotocol/bandchain/pull/2742) Validator and proposal message restructure
- (impv) [\#2740](https://github.com/bandprotocol/bandchain/pull/2740) Add UI test for all pages
- (impv) [\#2731](https://github.com/bandprotocol/bandchain/pull/2731) Add no transaction placeholder and fix mobile margin
- (impv) [\#2728](https://github.com/bandprotocol/bandchain/pull/2728) Add delegation, redelegation, undelegation and withdraw reward UI test
- (impv) [\#2727](https://github.com/bandprotocol/bandchain/pull/2727) Revamp connect modal
- (impv) [\#2725](https://github.com/bandprotocol/bandchain/pull/2725) Show validator name if voter is validator
- (impv) [\#2723](https://github.com/bandprotocol/bandchain/pull/2723) Re-factor to Msg structure
- (impv) [\#2719](https://github.com/bandprotocol/bandchain/pull/2719) Implement query Band supply
- (bugs) [\#2716](https://github.com/bandprotocol/bandchain/pull/2716) Fixed no decimal balance bug on account index
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

- (feat) [\#2529](https://github.com/bandprotocol/bandchain/pull/2529) Add harmony to bridge evm
- (docs) [\#2691](https://github.com/bandprotocol/bandchain/pull/2691) Add simple price db example for doc
- (feat) [\#2632](https://github.com/bandprotocol/bandchain/pull/2632) Add aggregator contract for ICON bridge
- (impv) [\#2626](https://github.com/bandprotocol/bandchain/pull/2626) Icon bridge fixed from auditing process

### Runtime

- (impv) [\#2494](https://github.com/bandprotocol/bandchain/pull/2494) Patch google cloud function

### Owasm

### Oracle Binary Encoding (OBI)

### Helpers

- (feat) [\#2889](https://github.com/bandprotocol/bandchain/pull/2889) Bandchain.js: Started to create the client module with mock request test
- (feat) [\#2872](https://github.com/bandprotocol/bandchain/pull/2872) Bandchain.js: Added new Address class and added more fn on PublicKey class
- (feat) [\#2865](https://github.com/bandprotocol/bandchain/pull/2865) bandchain.js: Add Private Key and verify on Public Key on Wallet
- (impv) [\#2863](https://github.com/bandprotocol/bandchain/pull/2863) bandchain.js: Add Github action
- (impv) [\#2862](https://github.com/bandprotocol/bandchain/pull/2862) bandchain.js: Implement `getSignData`, `getTxData` and all `with_*` methods on Transaction module
- (feat) [\#2855](https://github.com/bandprotocol/bandchain/pull/2855) pyband: Implemented with_auto fn on tx module
- (impv) [\#2835](https://github.com/bandprotocol/bandchain/pull/2835) pyband: Add msg delegate
- (impv) [\#2830](https://github.com/bandprotocol/bandchain/pull/2830) pyband: Add msg send
- (impv) [\#2838](https://github.com/bandprotocol/bandchain/pull/2838) pyband: refactor get_latest_block on client module
- (impv) [\#2826](https://github.com/bandprotocol/bandchain/pull/2826) pyband: Add Pyband test on Github Action
- (impv) [\#2807](https://github.com/bandprotocol/bandchain/pull/2807) pyband: Refactor and wrote tests for client and send tx
- (impv) [\#2803](https://github.com/bandprotocol/bandchain/pull/2803) pyband: Fix typing and add PrivateKey.from_hex on Wallet
- (impv) [\#2802](https://github.com/bandprotocol/bandchain/pull/2802) pyband: Transaction class, add test cases
- (feat) [\#2799](https://github.com/bandprotocol/bandchain/pull/2799) pyband: Implemented Message class
- (impv) [\#2789](https://github.com/bandprotocol/bandchain/pull/2789) pyband: get chain id from specific rest endpoint.
- (impv) [\#2739](https://github.com/bandprotocol/bandchain/pull/2739) pyband: fix client raise error when get fail
- (impv) [\#2652](https://github.com/bandprotocol/bandchain/pull/2652) pyband: use string instead of class annotation for Python3.6
- (bugs) [\#2651](https://github.com/bandprotocol/bandchain/pull/2651) pyband: fix bug get latest block

### MISC
