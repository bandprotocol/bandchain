# OWASM Studio Server

A simple server for compiling OWASM code from OWASM Studio

## Running service locally

Before running it first time, you will need to setup the dependencies:

1. rustup, see https://www.rustup.rs/
2. rustc nightly channel: `rustup toolchain install nightly`
3. and wasm target: `rustup target add wasm32-unknown-unknown`
4. wasm-gc: `cargo install wasm-gc`
5. wasm-bindgen: `cargo install wasm-bindgen-cli`
6. rustfmt nightly channel: `rustup component add rustfmt-preview --toolchain nightly`

```sh
yarn
AWS_ACCESS_KEY=<key> AWS_SECRET_ACCESS_KEY=<secret> yarn start
```

By default it will run on "0.0.0.0:8082" address. Use `PORT` environment variable to change it.
