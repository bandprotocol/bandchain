<--
(feat): New feature
(impv): Improvement / Enhancement
(docs): Documentation
(bugs): Bug fixes
(chore): Chore/cleanup work
-->

# Changelog

## [Unreleased]

### Chain

- (impv) [\#1753](https://github.com/bandprotocol/bandchain/pull/1753) Make sure bandoracled cache data source content
- (impv) [\#1920](https://github.com/bandprotocol/bandchain/pull/1920) Cleaned up default data sources and oracle scripts info
- (feat) [\#1917](https://github.com/bandprotocol/bandchain/pull/1917) Add memo to tx table in cacher.
- (impv/chore) [\#1893](https://github.com/bandprotocol/bandchain/pull/1893) Cleanup and Add genesis command to add data sources / oracle scripts.
- (feat) [\#1910](https://github.com/bandprotocol/bandchain/pull/1910) Implement request-search REST and CLI endpoint.
- (impv) [\#1903](https://github.com/bandprotocol/bandchain/pull/1903) Store result in store instead of result hash.
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
- (imprv) [\#1904](https://github.com/bandprotocol/bandchain/pull/1904) Added helper tooltips.
- (imprv) [\#1900](https://github.com/bandprotocol/bandchain/pull/1900) Shorten marketcap amount on landing page
- (feat) [\#1888](https://github.com/bandprotocol/bandchain/pull/1888) Added OBI bindings and patched scan to use OBI standard
- (bug) [\#1861](https://github.com/bandprotocol/bandchain/pull/1861) Fix name, endpoint of guanyu-devnet on chain-id selection
- (feat) [\#1856](https://github.com/bandprotocol/bandchain/pull/1856) Add sorting on Validator Home Page.

### Bridges

### Owasm

- (feat) [\#1919](https://github.com/bandprotocol/bandchain/pull/1919) Add Makefile entry to allow compiling dylib and so for go-owasm via docker.
- (feat) [\#1907](https://github.com/bandprotocol/bandchain/pull/1907) Implement OBISchema derive for generate schema of Input and Output struct
- (feat) [\#1858](https://github.com/bandprotocol/bandchain/pull/1858) Add go-owasm to BandChain monorepo.

### Oracle Binary Encoding (OBI)

### MISC

- (chore) [\#1876](https://github.com/bandprotocol/bandchain/pull/1876)Update docker script fix build failed when using go-owasm local package
