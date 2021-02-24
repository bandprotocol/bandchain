#!/bin/bash

rm -rf ~/.yoda

# config chain id
yoda config chain-id odin

# add validator to yoda config
yoda config validator $(bandcli keys show $1 -a --bech val --keyring-backend test)

# setup execution endpoint
yoda config executor "rest:https://iv3lgtv11a.execute-api.ap-southeast-1.amazonaws.com/live/master?timeout=10s"

# setup broadcast-timeout to yoda config
yoda config broadcast-timeout "5m"

# setup rpc-poll-interval to yoda config
yoda config rpc-poll-interval "1s"

# setup max-try to yoda config
yoda config max-try 5

yoda config node $2

echo "y" | bandcli tx oracle activate --from $1 --keyring-backend test --broadcast-mode block --node $2

# wait for activation transaction success
sleep 10

yoda keys add reporter

# send band tokens to reporters
echo "y" | bandcli tx multi-send 1000000loki $(yoda keys list -a) --from $1 --keyring-backend test --broadcast-mode block --node $2

# wait for sending band tokens transaction success
sleep 10

# add reporter to bandchain
echo "y" | bandcli tx oracle add-reporters $(yoda keys list -a) --from $1 --keyring-backend test --broadcast-mode block --node $2

# wait for addding reporter transaction success
sleep 10

# run yoda
yoda run
