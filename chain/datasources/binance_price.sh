#!/bin/sh

symbol=$1

# Cryptocurrency price endpoint
url="https://api.binance.com/api/v1/depth?symbol=${symbol}USDT&limit=5"

# Performs data fetching and stores result to variable
result=$(curl -s -X GET $url -H "accept: application/json")

bid=$(jq -r ".bids | .[0] | .[0]" <<< "$result")
ask=$(jq -r ".asks | .[0] | .[0]" <<< "$result")

if [[ "$bid" == "null" || "$ask" == "null" ]];
then
    exit 1
else
    bc -l <<< "($bid+$ask)/2.0" | awk '{printf "%.8f\n", $0}'
fi
