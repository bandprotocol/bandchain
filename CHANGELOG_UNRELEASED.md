<!--
(feat): New feature
(impv): Improvement / Enhancement
(docs): Documentation
(bugs): Bug fixes
(chore): Chore/cleanup work
-->

## [Unreleased]

### Chain

- (chore) [\#2040](https://github.com/bandprotocol/bandchain/pull/2040) Set HomeFlag to /tmp for SimApp.
- (patch) [\#2037](https://github.com/bandprotocol/bandchain/pull/2037) Patch Multistore proof to new structure tree of 0.38.
- (impv) [\#1892](https://github.com/bandprotocol/bandchain/pull/1892) Rewrite execution environment for wasmer
- (patch) [\#1999](https://github.com/bandprotocol/bandchain/pull/1999) Patch to Cosmos-SDK 0.38.4
- (impv) [\#1981](https://github.com/bandprotocol/bandchain/pull/1981) Remove gas price refund for report transactions.
- (feat) [\#1872](https://github.com/bandprotocol/bandchain/pull/1872) Add Owasm executor on Google Cloud function
- (feat) [\#1951](https://github.com/bandprotocol/bandchain/pull/1951) Add REST endpoint to get genesis file.
- (feat) [\#1908](https://github.com/bandprotocol/bandchain/pull/1908) Support google cloud function REST executor
- (feat) [\#1929](https://github.com/bandprotocol/bandchain/pull/1929) Add report info of validator on CLI and REST endpoint.
- (chore) [\#1933](https://github.com/bandprotocol/bandchain/pull/1933) Improve code quality and conciseness in x/oracle/keeper.
- (bugs) [\#1935](https://github.com/bandprotocol/bandchain/pull/1935) Fix wrong order on call search latest request.
- (feat) [\#1927](https://github.com/bandprotocol/bandchain/pull/1927) Add reporter list of validator on CLI and REST endpoint.
- (impv) [\#1891](https://github.com/bandprotocol/bandchain/pull/1891) Add end-to-end request flow test.

### Scan

- (impv) [\#2039](https://github.com/bandprotocol/bandchain/pull/2039) Forward max_span_size from module to go-owasm
- (impv) [\#2008](https://github.com/bandprotocol/bandchain/pull/2008) Add target validator's address on redelegate
- (feat) [\#2001](https://github.com/bandprotocol/bandchain/pull/2001) Add `Countup.js` to animate balance and reward
- (bugs) [\#1996](https://github.com/bandprotocol/bandchain/pull/1996) Fix id bug on `RequestSub.re`
- (chore) [\#1987](https://github.com/bandprotocol/bandchain/pull/1987) Remove `request_tab_t` in `Route.re`
- (impv) [\#1986](https://github.com/bandprotocol/bandchain/pull/1986) Add autofocus input on submittx modal
- (feat) [\#1985](https://github.com/bandprotocol/bandchain/pull/1985) Add redelegate button and submit transaction modal
- (chore) [\#1982](https://github.com/bandprotocol/bandchain/pull/1982) Update network on `ChainIDBadge.re`
- (impv) [\#1970](https://github.com/bandprotocol/bandchain/pull/1970) Update oracle script bridge code generator to use Obi standard
- (impv) [\#1958](https://github.com/bandprotocol/bandchain/pull/1958) Add human-readable error when broadcast tx, fix withdraw reward msg on guanyu
- (impv) [\#1947](https://github.com/bandprotocol/bandchain/pull/1947) Fix fade out on modal
- (impv) [\#1940](https://github.com/bandprotocol/bandchain/pull/1940) Trim Address input, remind user when sending token to themself
- (impv) [\#1939](https://github.com/bandprotocol/bandchain/pull/1939) Use Format.re on user balance
- (feat) [\#1938](https://github.com/bandprotocol/bandchain/pull/1938) Add validator's image from identity
- (impv) [\#1928](https://github.com/bandprotocol/bandchain/pull/1928) Add chainID for guanyu-devnet
- (impv) [\#1925](https://github.com/bandprotocol/bandchain/pull/1925) Removed input field for unused oracle script input
- (impv) [\#1906](https://github.com/bandprotocol/bandchain/pull/1906/files) Add searchbar on validator home page

### Bridges

### Owasm

- (impv) [\#2026](https://github.com/bandprotocol/bandchain/pull/2026) Remove get_calldata_size and get_external_data_size from OEI.
- (chore) [\#2025](https://github.com/bandprotocol/bandchain/pull/2025) Cleanup code documentation and error messages
- (feat) [\#1998](https://github.com/bandprotocol/bandchain/pull/1998) Implement safe guard to check import and export from wasm
- (feat) [\#1994](https://github.com/bandprotocol/bandchain/pull/1994) Implement inject stack height when compile wasm code.
- (feat) [\#1988](https://github.com/bandprotocol/bandchain/pull/1988) Implement safe guard to check memory limit from wasm
- (feat) [\#1937](https://github.com/bandprotocol/bandchain/pull/1937) Allow gas configuration from golang world to wasmer.
- (impv) [\#1936](https://github.com/bandprotocol/bandchain/pull/1936) Return error if writing beyond the span capacity.
- (impv) [#\1941](https://github.com/bandprotocol/bandchain/pull/1941) Add error code standard for wasm compilation.
- (impv) [#\1941](https://github.com/bandprotocol/bandchain/pull/1941) Fix how to build share object in Linux.
- (feat) [\#1922](https://github.com/bandprotocol/bandchain/pull/1922) Add wat to wasm function.

### Oracle Binary Encoding (OBI)

- (impv) [\#2027](https://github.com/bandprotocol/bandchain/pull/2027) Make Go-OBI natively support encoding or decoding multiple values.

### Helpers

- (feat) [\#1963](https://github.com/bandprotocol/bandchain/pull/1963) Add Bandchain.js

### MISC
