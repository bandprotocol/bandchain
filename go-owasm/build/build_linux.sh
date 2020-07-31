#!/bin/bash

cd go-owasm
cargo build --release
cp target/release/deps/libgo_owasm.so api
