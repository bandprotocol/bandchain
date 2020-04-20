#!/bin/sh

length=$1

url="https://qrng.anu.edu.au/API/jsonI.php?length=$length&type=uint8"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".[\"data\"]"
