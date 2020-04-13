#!/bin/sh

# Cryptocurrency price endpoint
url="https://api.binance.com/api/v1/depth?symbol=ATOMUSDT&limit=5"

# Performs data fetching and stores result to variable
result=$(curl -s -X GET $url -H "accept: application/json")

bid=$(echo "$result" | jq -r ".bids | .[0] | .[0]")
ask=$(echo "$result" | jq -r ".asks | .[0] | .[0]")

if [[ "$bid" == "null" || "$ask" == "null" ]];
then
    exit 1
else
    echo "($bid+$ask)/2.0" | bc -l | awk '{printf "%.8f\n", $0}'
fi
