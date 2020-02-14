#!/bin/bash
source ~/.profile
cd owasm/res

for f in *; do
    if [ -d "$f" ]; then
        wasm-pack build $f
        cp ./$f/target/wasm32-unknown-unknown/release/*.wasm ./
    fi
done

cd ../../
go test ./...
