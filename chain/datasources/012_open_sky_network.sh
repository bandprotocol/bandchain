#!/bin/bash

flight_op=$1
airport=$2
icao24=$3
begin=$4
end=$5

url="https://opensky-network.org/api/flights/$flight_op?airport=$airport&begin=$begin&end=$end"

# Performs data fetching and parses the result
if curl -s -X GET $url -H "accept: application/json" | jq -er "." | grep -q "$icao24"; then
  echo "true";
else
  echo "false";
fi
