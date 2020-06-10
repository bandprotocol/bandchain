#!/bin/sh

symbol=$1

# Cryptocurrency volume endpoint: https://min-api.cryptocompare.com/documentation
url="https://min-api.cryptocompare.com/data/symbol/histoday?fsym=$symbol&tsym=USD&limit=1"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".Data | .[0] | .total_volume_total"
