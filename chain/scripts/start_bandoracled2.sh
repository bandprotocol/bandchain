#!/bin/bash

rm -rf ~/.oracled

# config chain id
bandoracled2 config chain-id bandchain

# add validator to bandoracled2 config
bandoracled2 config validator bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec

for i in $(eval echo {1..$1})
do
  bandoracled2 keys add reporter$i
  # send band tokens to reporters
  echo "y" | bandcli tx send validator $(bandoracled2 keys show reporter$i) 1000000uband --keyring-backend test

  # wait for sending transaction success
  sleep 2

  # add reporter to bandchain
  echo "y" | bandcli tx oracle add-reporter $(bandoracled2 keys show reporter$i) --from validator --keyring-backend test

  # wait for addding reporter transaction success
  sleep 2
done

# run bandoracled2
bandoracled2 run