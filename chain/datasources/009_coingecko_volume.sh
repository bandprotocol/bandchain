#!/bin/sh

symbol=$1

# Cryptocurrency volume endpoint: https://www.coingecko.com/api/documentations/v3
url="https://api.coingecko.com/api/v3/coins/$symbol/market_chart?vs_currency=usd&days=1"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".total_volumes | .[-1] | .[1]"
