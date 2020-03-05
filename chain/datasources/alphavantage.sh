#!/bin/sh

symbol=$1
api_key=$2

url="https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=$symbol&apikey=$api_key"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".[\"Global Quote\"][\"05. price\"]"
