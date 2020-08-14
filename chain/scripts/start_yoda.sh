#!/bin/bash

rm -rf ~/.yoda

# config chain id
yoda config chain-id bandchain

# add validator to yoda config
yoda config validator $(bandcli keys show $1 -a --bech val --keyring-backend test)

# setup execution endpoint
yoda config executor "rest:https://3hdt5gnbr6.execute-api.ap-southeast-1.amazonaws.com/live/py-execution?timeout=10s"

# setup broadcast-timeout to yoda config
yoda config broadcast-timeout "30s"

# setup rpc-poll-interval to yoda config
yoda config rpc-poll-interval "1s"

# setup max-try to yoda config
yoda config max-try 5

echo "y" | bandcli tx oracle activate --from $1 --keyring-backend test

# wait for activation transaction success
sleep 2

for i in $(eval echo {1..$2})
do
  # add reporter key
  yoda keys add reporter$i
done

# send band tokens to reporters
echo "y" | bandcli tx multi-send 1000000uband $(yoda keys list -a) --from $1 --keyring-backend test

# wait for sending band tokens transaction success
sleep 2

# add reporter to bandchain
echo "y" | bandcli tx oracle add-reporters $(yoda keys list -a) --from $1 --keyring-backend test

# wait for addding reporter transaction success
sleep 2

# run yoda
yoda run
