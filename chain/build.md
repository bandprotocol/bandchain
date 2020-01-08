# Building and running the application

## Building the `band` application

If you want to build the `band` application in this repo to see the functionalities, **Go 1.13.0+** is required .

Add some parameters to environment is necessary if you have never used the `go mod` before.

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
echo "export GO111MODULE=on" >> ~/.bash_profile
source ~/.bash_profile
```

Now, you can install and run the application.

```
# Clone the source of the tutorial repository
git clone git@github.com:bandprotocol/d3n.git
cd d3n/chain/
```

```bash
# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
bandd help
bandcli help
```

## Running the live network and using the commands

To initialize configuration and a `genesis.json` file for your application and an account for the transactions, start by running:

> _*NOTE*_: In the below commands addresses are pulled using terminal utilities. You can also just input the raw strings saved from creating keys, shown below. The commands require [`jq`](https://stedolan.github.io/jq/download/) to be installed on your machine.

> _*NOTE*_: If you have run the tutorial before, you can start from scratch with a `bandd unsafe-reset-all` or by deleting both of the home folders `rm -rf ~/.req*`

```bash
# Initialize configuration files and genesis file
  # moniker is the name of your node
bandd init <moniker> --chain-id bandchain


# Copy the `Address` output here and save it for later use
bandcli keys add owner

# Copy the `Address` output here and save it for later use
bandcli keys add validator1

# Copy the `Address` output here and save it for later use
bandcli keys add validator2

# Add both accounts, with band to the genesis file
bandd add-genesis-account $(bandcli keys show owner -a) 100000000uband
bandd add-genesis-account $(bandcli keys show validator1 -a) 100000000uband
bandd add-genesis-account $(bandcli keys show validator2 -a) 100000000uband

# Configure your CLI to eliminate need for chain-id flag
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

bandd gentx --name owner <or your key_name>
```

After you have generated a genesis transcation, you will have to input the gentx into the genesis file, so that your band chain is aware of the validators. To do so, run:

`bandd collect-gentxs`

and to make sure your genesis file is correct, run:

`bandd validate-genesis`

You can now start `bandd` by calling `bandd start`. You will see logs begin streaming that represent blocks being produced, this will take a couple of seconds.

You have run your first node successfully.

```bash
# First check the accounts to ensure they have funds
bandcli query account $(bandcli keys show owner -a)
bandcli query account $(bandcli keys show validator1 -a)
```

## How to run it

```bash
 # send request
 bandcli tx zoracle request 30 $(xxd -p -c100000000000 ./wasm/res/test_u64.wasm) --from owner --gas 10000000

 # get request by id
 bandcli query zoracle request <reqID>

 # get pending request
 bandcli query zoracle pending_request

 # send report
 bandcli tx zoracle report <reqID> <data> --from <validator>

 bandcli tx zoracle report <reqID> 02000000000000001b000000000000007b22626974636f696e223a7b22757364223a373436392e34397d7d0f000000000000007b22555344223a373531302e32317d --from owner
```
