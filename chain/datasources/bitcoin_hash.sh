#!/bin/sh

height=$1

url="http://api.blockcypher.com/v1/btc/main/blocks/$height?txstart=1&limit=1"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".[\"hash\"]"
