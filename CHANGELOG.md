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

- (impv) [\#1901](https://github.com/bandprotocol/bandchain/pull/1901) Added `moniker` field to `delegations_view` table.
- (feat) [\#1879](https://github.com/bandprotocol/bandchain/pull/1879) Keep `Request` and `Report` around. We avoid premature optimization at the moment.
- (feat) [\#1873](https://github.com/bandprotocol/bandchain/pull/1873) Add `in-before-resolve` field in `Report` data structure.
- (impv) [\#1880](https://github.com/bandprotocol/bandchain/pull/1880) Oracled handle failed data source execution.
- (feat) [\#1875](https://github.com/bandprotocol/bandchain/pull/1875) Add CLI and REST query interface for request.
- (chore) [\#1869](https://github.com/bandprotocol/bandchain/pull/1869) Update new schema and source code url for all oracle scripts.
- (chore) [\#1864](https://github.com/bandprotocol/bandchain/pull/1864) Remove unused query types.
- (impv) [\#1792](https://github.com/bandprotocol/bandchain/pull/1792) Request data message handler test.

### Scan

- (imprv) [\#1911](https://github.com/bandprotocol/bandchain/pull/1911) Use validator moniker instead of validator address on Account page.
- (imprv) [\#1900](https://github.com/bandprotocol/bandchain/pull/1900) Shorten marketcap amount on landing page
- (bug) [\#1861](https://github.com/bandprotocol/bandchain/pull/1861) Fix name, endpoint of guanyu-devnet on chain-id selection

### Bridges

### Owasm

- (feat) [\#1858](https://github.com/bandprotocol/bandchain/pull/1858) Add go-owasm to BandChain monorepo.

### Oracle Binary Encoding (OBI)

### MISC

- (chore) [\#1876](https://github.com/bandprotocol/bandchain/pull/1876)Update docker script fix build failed when using go-owasm local package
