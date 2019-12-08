#!bin/bash

# remove old genesis
rm -rf ~/.band*

# initial new node
bandd init node-validator --chain-id bandchain

# create acccount
bandcli keys add owner
bandcli keys add user
# bandcli keys add bob --no-backup

bandd add-genesis-account $(bandcli keys show owner -a) 99000000000000stake
bandd add-genesis-account $(bandcli keys show user -a) 1000000000000stake

bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

bandd gentx --name owner

bandd collect-gentxs
bandd validate-genesis

bandd start
