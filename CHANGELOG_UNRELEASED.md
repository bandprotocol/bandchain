<!--
(feat): New feature
(impv): Improvement / Enhancement
(docs): Documentation
(bugs): Bug fixes
(chore): Chore/cleanup work
-->

## [Unreleased]

### Chain

- (chore) [\#2130](https://github.com/bandprotocol/bandchain/pull/2130) Add ABCI begin block rolling seed test.
- (bug) [\#2075](https://github.com/bandprotocol/bandchain/pull/2075) Add height check when sync on db and fix external id can be zero.
- (bug) [\#2125](https://github.com/bandprotocol/bandchain/pull/2125) Fix request with duplicate external id and empty raw request bug.
- (chore) [\#2126](https://github.com/bandprotocol/bandchain/pull/2126) More test cleanups in request and result keepers.
- (impv) [\#2121](https://github.com/bandprotocol/bandchain/pull/2121) Add handle edit validator msg for emitter
- (impv/chore) [\#2124](https://github.com/bandprotocol/bandchain/pull/2124) Add genesis ds and os and move same test logic to testapp.
- (impv) [\#2113](https://github.com/bandprotocol/bandchain/pull/2113) Add tests in types and deactivate event for activate flow.
- (impv) [\#2109](https://github.com/bandprotocol/bandchain/pull/2109) Add Validator table and handle create validator message from genesis file and tx for emitter and flusher.
- (impv) [\#2093](https://github.com/bandprotocol/bandchain/pull/2093) Add missing pieces on app.go + some refactor and comments.
- (feat) [\#2114](https://github.com/bandprotocol/bandchain/pull/2114) Add more unit test coverage and enhance code comments in pkg.
- (impv) [\#2072](https://github.com/bandprotocol/bandchain/pull/2072) Handle resolve request for emitter/flusher.
- (feat) [\#2111](https://github.com/bandprotocol/bandchain/pull/2111), [\#2117](https://github.com/bandprotocol/bandchain/pull/2117) Introduce the notion of active validators who are performing oracle tasks.
- (bug) [\#2110](https://github.com/bandprotocol/bandchain/pull/2074) Set `bandoracled` max capacity of event subscription channel.
- (impv) [\#2106](https://github.com/bandprotocol/bandchain/pull/2106) Implement emitter handle MsgCreateDataSource/OracleScript.
- (impv) [\#2074](https://github.com/bandprotocol/bandchain/pull/2074) Use rolling block hash as seed for validator sampling.
- (impv) [\#2104](https://github.com/bandprotocol/bandchain/pull/2104) Update default gas/size consensus params and clean up cmd code.
- (chore) [\#2084](https://github.com/bandprotocol/bandchain/pull/2084) Rename ValidatorReportInto to ReportInfo.
- (chore) [\#2082](https://github.com/bandprotocol/bandchain/pull/2082) Reorder, reword, and remove unused error codes.
- (impv) [\#2080](https://github.com/bandprotocol/bandchain/pull/2080), [\#2107](https://github.com/bandprotocol/bandchain/pull/2107) Set count state on genesis.go and remove default values from getters.
- (feat) [\#2022](https://github.com/bandprotocol/bandchain/pull/2022) Initial implementation of BandChain emitter/flusher.
- (chore) [\#2060](https://github.com/bandprotocol/bandchain/pull/2060) Remove unused /bandchain/file endpoints and custom swagger from bandcli REST.

### Scan

- (impv) [\#2005](https://github.com/bandprotocol/bandchain/pull/2005) Show max commission rate and max commission change on validator index page
- (feat) [\#2001](https://github.com/bandprotocol/bandchain/pull/2001) Add `Countup.js` to animate balance and reward
- (chore) [\#1987](https://github.com/bandprotocol/bandchain/pull/1987) Remove `request_tab_t` in `Route.re`
- (impv) [\#1971](https://github.com/bandprotocol/bandchain/pull/1971) Add loading state on Tx Index Page
- (impv) [\#1947](https://github.com/bandprotocol/bandchain/pull/1947) Fix fade out on modal

### Bridges

- (feat) [\#2055](https://github.com/bandprotocol/bandchain/pull/2055) Implement BridgeWithCache to keep the latest response for any unique request packet

### Owasm

### Oracle Binary Encoding (OBI)

- (impv) [#1947](https://github.com/bandprotocol/bandchain/pull/2065) Remove obi.js build process

### Helpers

### MISC

- (chore) [\#2105](https://github.com/bandprotocol/bandchain/pull/2105) Add pull request template to describe PR.
- (chore) [\#2068](https://github.com/bandprotocol/bandchain/pull/2068) Remove `band-consumer` from repository.
