#!/bin/sh

# Gold price endpoint
url="https://www.freeforexapi.com/api/live?pairs=USDXAU"

rate=$(curl -s -X GET $url -H "accept: application/json" | jq '.rates' | jq '.USDXAU' | jq '.rate')
echo $(bc <<< "1.0/$rate")

