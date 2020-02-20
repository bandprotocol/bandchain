#!/bin/sh

symbol=$1

# Cryptocurrency price endpoint: https://min-api.cryptocompare.com/documentation
url="https://min-api.cryptocompare.com/data/price?fsym=$symbol&tsyms=USD"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".USD"
