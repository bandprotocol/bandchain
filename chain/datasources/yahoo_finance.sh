#!/bin/sh

symbol=$1

url="https://finance.yahoo.com/quote/$symbol"

# Performs data fetching
raw=$(curl -s -X GET $url)
remove_newline_raw=$(cat <<< $raw | awk '{print}' ORS='')
result=$(echo "$remove_newline_raw" | awk -F'root.App.main = ' '{print $2}' | awk -F'; }\\\(this\\\)' '{print $1}')
echo "$result" | jq -er '.["context"]["dispatcher"]["stores"]["QuoteSummaryStore"]["price"]["regularMarketPrice"]["raw"]'
