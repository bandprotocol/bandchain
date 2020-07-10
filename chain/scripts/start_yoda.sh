#!/bin/bash

rm -rf ~/.yoda

# config chain id
yoda config chain-id bandchain

# add validator to yoda config
yoda config validator $(bandcli keys show $1 -a --bech val --keyring-backend test)

echo "y" | bandcli tx oracle activate --from $1 --keyring-backend test

# wait for activation transaction success
sleep 2

for i in $(eval echo {1..$2})
do
  # add reporter key
  yoda keys add reporter$i

  # send band tokens to reporter
  echo "y" | bandcli tx send $1 $(yoda keys show reporter$i) 1000000uband --keyring-backend test

  # wait for sending band tokens transaction success
  sleep 2

  # add reporter to bandchain
  echo "y" | bandcli tx oracle add-reporter $(yoda keys show reporter$i) --from $1 --keyring-backend test

  # wait for addding reporter transaction success
  sleep 2
done

# run yoda
yoda run
