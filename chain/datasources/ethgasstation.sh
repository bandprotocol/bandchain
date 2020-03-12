#!/bin/sh

type=$1

url="https://ethgasstation.info/json/ethgasAPI.json"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".[\"$type\"]"
