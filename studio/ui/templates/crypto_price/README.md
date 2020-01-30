# OWASM Crypto price boilerplate

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
$ wasm-pack build . --out-name crypto_price
$ ls pkg/crypto_price_bg.wasm
```
