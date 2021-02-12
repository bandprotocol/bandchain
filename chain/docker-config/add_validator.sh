echo "usage ketchup faculty bench jewel rocket latin absurd decide field party reunion cook entry scout scene miss box memory museum decorate guide few verify" \
    | bandcli keys add supplier --recover --keyring-backend test

echo "y" | bandcli tx staking create-validator \
  --amount 100000000loki \
  --commission-max-change-rate 0.010000000000000000 \
  --commission-max-rate 0.200000000000000000 \
  --commission-rate 0.100000000000000000 \
  --chain-id odin \
  --from supplier \
  --moniker oracle-validator \
  --pubkey odinvalconspub1addwnpepqge86lvslkpfk0rlz0ah9tat0vntx8yele36hhfpflehfehydlutkvdvhfm \
  --min-self-delegation 1 \
  --keyring-backend test