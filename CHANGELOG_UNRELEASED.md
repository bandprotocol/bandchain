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

- (impv) [\#2607](https://github.com/bandprotocol/bandchain/pull/2607) Query proof with specific height parameter

### Yoda

- (impv) [\#2631](https://github.com/bandprotocol/bandchain/pull/2631) Add request timeout for each data source execution

### Emitter & Flusher

- (impv) [\#2553](https://github.com/bandprotocol/bandchain/pull/2553) fast-sync: emit `unbonding`, `delegation` and `redelegation` from start state
- (impv) [\#2558](https://github.com/bandprotocol/bandchain/pull/2558) fast-sync: emit all oracle module from start state
- (bugs) [\#2601](https://github.com/bandprotocol/bandchain/pull/2601) Downgrade Kafka go
- (bugs) [\#2600](https://github.com/bandprotocol/bandchain/pull/2600) Fix bug handle new transaction

### Scan

- (impv) [\#2621](https://github.com/bandprotocol/bandchain/pull/2621) Wired up data for proposal vote
- (bugs) [\#2606](https://github.com/bandprotocol/bandchain/pull/2606) Fix overflow text on request index page
- (impv) [\#2604](https://github.com/bandprotocol/bandchain/pull/2604) Moved the proposal route to wenchang route
- (feat) [\#2603](https://github.com/bandprotocol/bandchain/pull/2603) Added `guanyu-poa` on chain id
- (impv) [\#2599](https://github.com/bandprotocol/bandchain/pull/2599) Polish style on revamp GuanYu part 2
- (feat) [\#2598](https://github.com/bandprotocol/bandchain/pull/2598) Implemented VoteSub, Vote breakdown table and wire up
- (feat) [\#2597](https://github.com/bandprotocol/bandchain/pull/2597) Implemented VoteMsg and modal action
- (impv) [\#2596](https://github.com/bandprotocol/bandchain/pull/2596) Fix total deposit, deposit amount type
- (bugs) [\#2595](https://github.com/bandprotocol/bandchain/pull/2595) Handle `client_id` optional case
- (feat) [\#2594](https://github.com/bandprotocol/bandchain/pull/2594) Added meta og tag to scan
- (impv) [\#2593](https://github.com/bandprotocol/bandchain/pull/2593) Polish UI
- (feat) [\#2592](https://github.com/bandprotocol/bandchain/pull/2592) Added new voting overview and results box to proposal index page
- (feat) [\#2591](https://github.com/bandprotocol/bandchain/pull/2591) Implemented the deposit overview, depositors table
- (feat) [\#2578](https://github.com/bandprotocol/bandchain/pull/2578) Implemented the top part for proposal index page, markdown component
- (bugs) [\#2576](https://github.com/bandprotocol/bandchain/pull/2576) Fix overflow value on History Bonded Token, implement HistoryOracleParser and unit test.
- (impv) [\#2573](https://github.com/bandprotocol/bandchain/pull/2573) Added tooltip text to each place which is lorem
- (feat) [\#2570](https://github.com/bandprotocol/bandchain/pull/2570) Created ProposalSub and ProposalHomepage, and also implemented the route for both home and index page
- (bugs) [\#2568](https://github.com/bandprotocol/bandchain/pull/2568) Updated reporters subscription
- (impv) [\#2566](https://github.com/bandprotocol/bandchain/pull/2566) Implemented TxIndexpage with new theme
- (impv) [\#2563](https://github.com/bandprotocol/bandchain/pull/2563) Implemented the revamp block index page
- (impv) [\#2561](https://github.com/bandprotocol/bandchain/pull/2561) Implemented the new layout for Validator Homepage
- (impv) [\#2560](https://github.com/bandprotocol/bandchain/pull/2560) Refactor, improve Block Home Page for revamp GuanYu version
- (impv) [\#2557](https://github.com/bandprotocol/bandchain/pull/2557) Fix and clean up copy button
- (impv) [\#2552](https://github.com/bandprotocol/bandchain/pull/2552) Wire up related data source, handle nullable timestamp
- (impv) [\#2550](https://github.com/bandprotocol/bandchain/pull/2550) Polish request index page
- (bugs) [\#2527](https://github.com/bandprotocol/bandchain/pull/2527) Fix search bugs on ds & os home page
- (impv) [\#2516](https://github.com/bandprotocol/bandchain/pull/2516) Fix loading width and bg color
- (impv) [\#2499](https://github.com/bandprotocol/bandchain/pull/2499) Implemented Total request chart
- (impv) [\#2455](https://github.com/bandprotocol/bandchain/pull/2455) Implement full copy non evm proof
- (impv) [\#2363](https://github.com/bandprotocol/bandchain/pull/2363) Added loading state on AccountIndexPage, fixed some layout to make everything consistance.

### Bridges

- (feat) [\#2629](https://github.com/bandprotocol/bandchain/pull/2629) Add contracts and interfaces for Band Standard dataset. Refactored `bridges/evm/` directory to me interface and library files into their own subfolders
- (feat) [\#2620](https://github.com/bandprotocol/bandchain/pull/2620) Bridge & CacheBridge with relay multiple
- (feat) [\#2548](https://github.com/bandprotocol/bandchain/pull/2548) Added `BridgeProxy` contract
- (feat) [\#2385](https://github.com/bandprotocol/bandchain/pull/2385) Add icon bridge

### Runtime

- (bugs) [\#2590](https://github.com/bandprotocol/bandchain/pull/2590) Fix verify executable size failed

### Owasm

### Oracle Binary Encoding (OBI)

- (feat) [\#2546](https://github.com/bandprotocol/bandchain/pull/2546) Implement encode / decode for fix length array

### Helpers

- (impv) [\#2577](https://github.com/bandprotocol/bandchain/pull/2577) Wallet module completion

### MISC
