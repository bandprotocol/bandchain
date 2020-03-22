#!/bin/sh

RAW_BASE=$1
RAW_QUOTE=$2



URL="https://api.coingecko.com/api/v3/simple/price?ids=$BASE&vs_currencies=$QUOTE"
KEY=".$BASE.usd"

# Performs data fetching and parses the result
curl -s -X GET $URL -H "accept: application/json" | jq -r ".[\"$BASE\"].usd"
