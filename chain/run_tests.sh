#!/bin/bash

mkdir -p pkg/owasm/res

# release mode
cd ../owasm/chaintests

for f in *; do
    if [ -d "$f" ]; then
        RUSTFLAGS='-C link-arg=-s' cargo build --target wasm32-unknown-unknown --release --package $f
        cp ../target/wasm32-unknown-unknown/release/$f.wasm ../../chain/pkg/owasm/res
    fi
done
cd ../../chain

# debug mode
cd ../owasm/chaindebugtests
for f in *; do
    if [ -d "$f" ]; then
        RUSTFLAGS='-C link-arg=-s' cargo build --target wasm32-unknown-unknown --package $f
        cp ../target/wasm32-unknown-unknown/debug/$f.wasm ../../chain/pkg/owasm/res
    fi
done

cd ../../chain
go test ./...
