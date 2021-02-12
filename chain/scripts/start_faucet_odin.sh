rm -rf ~/.faucet
faucet config chain-id odin
faucet config port 5005

faucet keys add worker

echo "y" | bandcli tx send supplier $(faucet keys show worker) 1000000000000loki --keyring-backend test

faucet run