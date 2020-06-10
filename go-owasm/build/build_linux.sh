#!/bin/bash

cargo build --release
cp target/release/deps/libgo_owasm.so api
