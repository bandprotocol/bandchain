#!/bin/sh

tar cvf src.tar Cargo.toml src/*

cat src.tar | base64 -w0 > src.base64

source=`cat src.base64`

curl -XPOST "$HOST/cargo" --data-binary "{ \"opts\": {}, \"tar\": \"$source\" }"

rm src.tar
rm src.base64
