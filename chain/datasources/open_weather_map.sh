#!/bin/bash

country=$1
main_field=$2
sub_field=$3

url="https://api.openweathermap.org/data/2.5/weather?q=$country&appid=ac7c05361f8f91652eab609377134ab7"

# Performs data fetching and parses the result
curl -s -X GET $url -H "accept: application/json" | jq -er ".[\"$main_field\"][\"$sub_field\"]"
