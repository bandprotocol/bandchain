# OWASM Boilerplate

## Prerequisite

Install

- Rust and Cargo - https://doc.rust-lang.org/cargo/getting-started/installation.html
- Wasm-pack - https://rustwasm.github.io/wasm-pack/installer/

## Test

```
$ cargo test -- --nocapture
```

## Build

```
$ wasm-pack build . --out-name data_request
$ ls pkg/data_request_bg.wasm
```
