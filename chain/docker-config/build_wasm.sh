#!/bin/bash

echo $1
cd $1
pwd
wasm-pack build
