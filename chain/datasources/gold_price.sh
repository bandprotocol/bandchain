#!/bin/sh

# Gold price endpoint
url="https://www.freeforexapi.com/api/live?pairs=USDXAU"

code=$(curl -s -X GET $url -H "accept: application/json" | jq '.code')
if [[ $code -eq 200 ]];
then
    rate=$(curl -s -X GET $url -H "accept: application/json" | jq '.rates' | jq '.USDXAU' | jq '.rate')
    echo $(bc <<< "1.0/$rate")
else
    exit 1
fi
