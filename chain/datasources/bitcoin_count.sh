#!/bin/sh

url="https://blockchain.info/q/getblockcount"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er "."
