#!/bin/sh

symbol=$1

# Cryptocurrency price endpoint: https://www.coingecko.com/api/documentations/v3
url="https://api.coingecko.com/api/v3/simple/price?ids=$symbol&vs_currencies=usd"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".[\"$symbol\"].usd"
