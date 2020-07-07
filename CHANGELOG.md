<!--
(feat): New feature
(impv): Improvement / Enhancement
(docs): Documentation
(bugs): Bug fixes
(chore): Chore/cleanup work
-->

# Changelog

## [v1.0.4](https://github.com/bandprotocol/bandchain/releases/tag/v1.0.4)

- Bump Cosmos-SDK to version 0.38.5
- Introduce active oracle validators concept.
- Kick start `Emitter` extended app that emit some messages to Kafka server.
- `Flusher` service that consume messages from Kafka server to populate database.

### Chain (Consensus)

- (impv) [\#2148](https://github.com/bandprotocol/bandchain/pull/2148) Add more unit tests and fix nil retdata and nil calldata.
- (impv) [\#2143](https://github.com/bandprotocol/bandchain/pull/2143) Separate resolve functions into success/failure/expired cases.
- (patch) [\#2135](https://github.com/bandprotocol/bandchain/pull/2135) Bump Cosmos-SDK to version 0.38.5.
- (impv) [\#2137](https://github.com/bandprotocol/bandchain/pull/2137) Update SetResult (Resolve) keeper signature and remove extra events.
- (impv) [\#2128](https://github.com/bandprotocol/bandchain/pull/2128) Add expired request check when report and change type of `requestTime`.
- (feat) [\#2131](https://github.com/bandprotocol/bandchain/pull/2131) Add validator status querier.
- (bug) [\#2125](https://github.com/bandprotocol/bandchain/pull/2125) Fix request with duplicate external id and empty raw request bug.
- (impv) [\#2093](https://github.com/bandprotocol/bandchain/pull/2093) Add missing pieces on app.go + some refactor and comments.
- (feat) [\#2111](https://github.com/bandprotocol/bandchain/pull/2111), [\#2117](https://github.com/bandprotocol/bandchain/pull/2117) Introduce the notion of active validators who are performing oracle tasks.
- (chore) [\#2084](https://github.com/bandprotocol/bandchain/pull/2084) Rename ValidatorReportInto to ReportInfo.
- (impv) [\#2074](https://github.com/bandprotocol/bandchain/pull/2074) Use rolling block hash as seed for validator sampling.
- (impv) [\#2080](https://github.com/bandprotocol/bandchain/pull/2080), [\#2107](https://github.com/bandprotocol/bandchain/pull/2107) Set count state on genesis.go and remove default values from getters.

### Chain (Not consensus)

- (impv) [\#2161](https://github.com/bandprotocol/bandchain/pull/2161) Update default genesis parameters for oracle and distr.
- (impv) [\#2151](https://github.com/bandprotocol/bandchain/pull/2151) Fix `exec_env` code more consistent.
- (impv/chore) [\#2147](https://github.com/bandprotocol/bandchain/pull/2147) Add abci begin block test for reward allocation.
- (impv) [\#2146](https://github.com/bandprotocol/bandchain/pull/2146) Various minor improvements and unit tests.
- (chore) [\#2139](https://github.com/bandprotocol/bandchain/pull/2139) Remove db in favor of emitter.
- (chore) [\#2129](https://github.com/bandprotocol/bandchain/pull/2129) Add unit tests for validator status keeper.
- (chore) [\#2130](https://github.com/bandprotocol/bandchain/pull/2130) Add ABCI begin block rolling seed test.
- (chore) [\#2126](https://github.com/bandprotocol/bandchain/pull/2126) More test cleanups in request and result keepers.
- (impv/chore) [\#2124](https://github.com/bandprotocol/bandchain/pull/2124) Add genesis ds and os and move same test logic to testapp.
- (impv) [\#2113](https://github.com/bandprotocol/bandchain/pull/2113) Add tests in types and deactivate event for activate flow.
- (feat) [\#2114](https://github.com/bandprotocol/bandchain/pull/2114) Add more unit test coverage and enhance code comments in pkg.
- (bug) [\#2110](https://github.com/bandprotocol/bandchain/pull/2074) Set `bandoracled` max capacity of event subscription channel.
- (impv) [\#2104](https://github.com/bandprotocol/bandchain/pull/2104) Update default gas/size consensus params and clean up cmd code.
- (chore) [\#2082](https://github.com/bandprotocol/bandchain/pull/2082) Reorder, reword, and remove unused error codes.
- (chore) [\#2060](https://github.com/bandprotocol/bandchain/pull/2060) Remove unused /bandchain/file endpoints and custom swagger from bandcli REST.

### Emitter & Flusher

- (impv) [\#2144](https://github.com/bandprotocol/bandchain/pull/2144) Add validator votes table.
- (impv) [\#2138](https://github.com/bandprotocol/bandchain/pull/2138) Docker deployment and handle jail event and unjail msg.
- (chore) [\#2108](https://github.com/bandprotocol/bandchain/pull/2108) Add script to run bandchain with emitter and flusher locally.
- (impv) [\#2132](https://github.com/bandprotocol/bandchain/pull/2132) Implement emitter handler for bank messages.
- (impv) [\#2136](https://github.com/bandprotocol/bandchain/pull/2136) Add reward field in validators table.
- (impv) [\#2121](https://github.com/bandprotocol/bandchain/pull/2121) Add handle edit validator msg for emitter
- (impv) [\#2118](https://github.com/bandprotocol/bandchain/pull/2118) Implement emitter handle MsgDelegate, MsgUndelegate and MsgBeginRedelegate.
- (impv) [\#2109](https://github.com/bandprotocol/bandchain/pull/2109) Add Validator table and handle create validator message from genesis file and tx for emitter and flusher.
- (impv) [\#2072](https://github.com/bandprotocol/bandchain/pull/2072) Handle resolve request for emitter/flusher.
- (impv) [\#2115](https://github.com/bandprotocol/bandchain/pull/2115) Implement emitter handle MsgEditDataSource/OracleScript.
- (impv) [\#2106](https://github.com/bandprotocol/bandchain/pull/2106) Implement emitter handle MsgCreateDataSource/OracleScript.
- (feat) [\#2022](https://github.com/bandprotocol/bandchain/pull/2022) Initial implementation of BandChain emitter/flusher.

### Scan

- (impv) [\#2005](https://github.com/bandprotocol/bandchain/pull/2005) Show max commission rate and max commission change on validator index page
- (feat) [\#2001](https://github.com/bandprotocol/bandchain/pull/2001) Add `Countup.js` to animate balance and reward
- (chore) [\#1987](https://github.com/bandprotocol/bandchain/pull/1987) Remove `request_tab_t` in `Route.re`
- (impv) [\#1971](https://github.com/bandprotocol/bandchain/pull/1971) Add loading state on Tx Index Page
- (impv) [\#1947](https://github.com/bandprotocol/bandchain/pull/1947) Fix fade out on modal

### Bridges

- (feat) [\#2055](https://github.com/bandprotocol/bandchain/pull/2055) Implement BridgeWithCache to keep the latest response for any unique request packet

### Owasm

- (feat) [\#2150](https://github.com/bandprotocol/bandchain/pull/2150) Increase wasm gas limit and consume gas when reading / writing.

### Oracle Binary Encoding (OBI)

- (impv) [#1947](https://github.com/bandprotocol/bandchain/pull/2065) Remove obi.js build process

### MISC

- (chore) [\#2105](https://github.com/bandprotocol/bandchain/pull/2105) Add pull request template to describe PR.
- (chore) [\#2068](https://github.com/bandprotocol/bandchain/pull/2068) Remove `band-consumer` from repository.

## [v1.0.3-alpha](https://github.com/bandprotocol/bandchain/releases/tag/v1.0.3-alpha)

### Chain

- (impv) [\#2057](https://github.com/bandprotocol/bandchain/pull/2057) Change wasm execute gas to 100000.
- (chore) [\#2053](https://github.com/bandprotocol/bandchain/pull/2053) Remove data sources and oracle scripts from repo.
- (chore) [\#2056](https://github.com/bandprotocol/bandchain/pull/2056) Remove IBCInfo from current v0.38 release.
- (impv) [\#2052](https://github.com/bandprotocol/bandchain/pull/2052) Improve proof endpoint and patch evm bridge contract
- (impv) [\#1934](https://github.com/bandprotocol/bandchain/pull/1934) Wireup go-owasm with blockchain properly.
- (impv) [\#2050](https://github.com/bandprotocol/bandchain/pull/2050) Return proper HTTP status codes on REST endpoints.
- (bugs) [\#2042](https://github.com/bandprotocol/bandchain/pull/2042) Add request and resolve time on failed and expired requests
- (bug) [\#2047](https://github.com/bandprotocol/bandchain/pull/2047) Fix request search ordering in 0.38.
- (bug) [\#2046](https://github.com/bandprotocol/bandchain/pull/2043) Use dash for bandcli report-info.
- (impv) [\#2043](https://github.com/bandprotocol/bandchain/pull/2043) Add full raw requests information in request struct.
- (chore) [\#2040](https://github.com/bandprotocol/bandchain/pull/2040) + [\#2044](https://github.com/bandprotocol/bandchain/pull/2044) Set HomeFlag to /tmp for SimApp.
- (patch) [\#2037](https://github.com/bandprotocol/bandchain/pull/2037) Patch Multistore proof to new structure tree of 0.38.
- (impv) [\#2021](https://github.com/bandprotocol/bandchain/pull/2021) Update chain test when execute bad wasm and get result from go-owasm runtime.
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
- (bugs) [\#1996](https://github.com/bandprotocol/bandchain/pull/1996) Fix id bug on `RequestSub.re`
- (impv) [\#1986](https://github.com/bandprotocol/bandchain/pull/1986) Add autofocus input on submittx modal
- (feat) [\#1985](https://github.com/bandprotocol/bandchain/pull/1985) Add redelegate button and submit transaction modal
- (chore) [\#1982](https://github.com/bandprotocol/bandchain/pull/1982) Update network on `ChainIDBadge.re`
- (impv) [\#1970](https://github.com/bandprotocol/bandchain/pull/1970) Update oracle script bridge code generator to use Obi standard
- (impv) [\#1958](https://github.com/bandprotocol/bandchain/pull/1958) Add human-readable error when broadcast tx, fix withdraw reward msg on guanyu
- (impv) [\#1940](https://github.com/bandprotocol/bandchain/pull/1940) Trim Address input, remind user when sending token to themself
- (impv) [\#1939](https://github.com/bandprotocol/bandchain/pull/1939) Use Format.re on user balance
- (feat) [\#1938](https://github.com/bandprotocol/bandchain/pull/1938) Add validator's image from identity
- (impv) [\#1928](https://github.com/bandprotocol/bandchain/pull/1928) Add chainID for guanyu-devnet
- (impv) [\#1925](https://github.com/bandprotocol/bandchain/pull/1925) Removed input field for unused oracle script input
- (impv) [\#1906](https://github.com/bandprotocol/bandchain/pull/1906/files) Add searchbar on validator home page

### Bridges

- (bug) [\#2049](https://github.com/bandprotocol/bandchain/pull/2049) Fix signed data prefix not consistency bug.
- (patch) [\#2041](https://github.com/bandprotocol/bandchain/pull/2041) Patch Multistore to new structure tree of 0.38.

### Owasm

- (impv) [\#2004](https://github.com/bandprotocol/bandchain/pull/2004) Implement recursive and memory test in go-owasm .
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

- (feat) [\#2058](https://github.com/bandprotocol/bandchain/pull/2058) Implement faucet service for testnet.

## [v1.0.2-alpha](https://github.com/bandprotocol/bandchain/releases/tag/v1.0.2-alpha)

### Chain

- (impv) [\#1753](https://github.com/bandprotocol/bandchain/pull/1753) Make sure bandoracled cache data source content
- (impv) [\#1920](https://github.com/bandprotocol/bandchain/pull/1920) Cleaned up default data sources and oracle scripts info
- (feat) [\#1917](https://github.com/bandprotocol/bandchain/pull/1917) Add memo to tx table in cacher.
- (impv/chore) [\#1893](https://github.com/bandprotocol/bandchain/pull/1893) Cleanup and Add genesis command to add data sources / oracle scripts.
- (feat) [\#1910](https://github.com/bandprotocol/bandchain/pull/1910) Implement request-search REST and CLI endpoint.
- (impv) [\#1903](https://github.com/bandprotocol/bandchain/pull/1903) Store result in store instead of result hash.
- (feat) [\#1905](https://github.com/bandprotocol/bandchain/pull/1905) Add request + response packet to request rest endpoint.
- (impv) [\#1901](https://github.com/bandprotocol/bandchain/pull/1901) Added `moniker` field to `delegations_view` table.
- (feat) [\#1879](https://github.com/bandprotocol/bandchain/pull/1879) Keep `Request` and `Report` around. We avoid premature optimization at the moment.
- (feat) [\#1873](https://github.com/bandprotocol/bandchain/pull/1873) Add `in-before-resolve` field in `Report` data structure.
- (impv) [\#1880](https://github.com/bandprotocol/bandchain/pull/1880) Oracled handle failed data source execution.
- (feat) [\#1875](https://github.com/bandprotocol/bandchain/pull/1875) Add CLI and REST query interface for request.
- (chore) [\#1869](https://github.com/bandprotocol/bandchain/pull/1869) Update new schema and source code url for all oracle scripts.
- (chore) [\#1864](https://github.com/bandprotocol/bandchain/pull/1864) Remove unused query types.
- (impv) [\#1792](https://github.com/bandprotocol/bandchain/pull/1792) Request data message handler test.

### Scan

- (imprv) [\#1924](https://github.com/bandprotocol/bandchain/pull/1924) Shorten bonded token amount on validator home page
- (bug) [\#1918](https://github.com/bandprotocol/bandchain/pull/1918) Fix UI overlap, word, chainIDBadge
- (feat) [\#1913](https://github.com/bandprotocol/bandchain/pull/1913) Add uptime bar on Validator Home Page
- (imprv) [\#1911](https://github.com/bandprotocol/bandchain/pull/1911) Use validator moniker instead of validator address on Account page.
- (imprv) [\#1904](https://github.com/bandprotocol/bandchain/pull/1904) Added helper tooltips.
- (imprv) [\#1900](https://github.com/bandprotocol/bandchain/pull/1900) Shorten marketcap amount on landing page
- (feat) [\#1888](https://github.com/bandprotocol/bandchain/pull/1888) Added OBI bindings and patched scan to use OBI standard
- (bug) [\#1861](https://github.com/bandprotocol/bandchain/pull/1861) Fix name, endpoint of guanyu-devnet on chain-id selection
- (feat) [\#1856](https://github.com/bandprotocol/bandchain/pull/1856) Add sorting on Validator Home Page.

### Owasm

- (feat) [\#2023](https://github.com/bandprotocol/bandchain/pull/2023) Remove wabt dependencies from Rust code.
- (feat) [\#1919](https://github.com/bandprotocol/bandchain/pull/1919) Add Makefile entry to allow compiling dylib and so for go-owasm via docker.
- (feat) [\#1907](https://github.com/bandprotocol/bandchain/pull/1907) Implement OBISchema derive for generate schema of Input and Output struct
- (feat) [\#1858](https://github.com/bandprotocol/bandchain/pull/1858) Add go-owasm to BandChain monorepo.

### MISC

- (chore) [\#1876](https://github.com/bandprotocol/bandchain/pull/1876)Update docker script fix build failed when using go-owasm local package
