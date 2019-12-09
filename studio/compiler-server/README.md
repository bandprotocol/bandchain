# Rustc service for WebAssembly Studio

This is heroku rustc compiler microservice for WebAssembly Studio. It is also a slug builder.

## WASM-RUN Warning
`app/wasm-run` is built for MacOS only. To use in linux, please contact swit@bandprotocol.com.

## Running service locally

The app/server can be run from the "app" folder:

```sh
cd app
nodemon .
```
  
By default it will run on "0.0.0.0:8082" address. Use `PORT` environment variable to change it.

Before running it first time, you will need to setup the dependencies:

1. rustup, see https://www.rustup.rs/
2. rustc nightly channel: `rustup toolchain install nightly`
3. and wasm target: `rustup target add wasm32-unknown-unknown --toolchain nightly`
4. wasm-gc: `cargo install wasm-gc`
5. wasm-bindgen: `cargo install wasm-bindgen-cli`
6. rustfmt nightly channel: `rustup component add rustfmt-preview --toolchain nightly`

See also [Rust for the Web](https://www.hellorust.com/setup/wasm-target/) for details.

## Change configuration in WebAssembly Studio

The "rustc" and  "cargo" endpoints addresses can be located at the https://github.com/wasdk/WebAssemblyStudio/blob/master/config.json and be changed locally.

## Building and installing the slug

Run `APP=<your-heroku-app> make publish` to fully build and setup slug at your app.
