Step to generate your gentx Message

0. Make sure you can run `bandd` and `bandcli`

1. Generate genesis file + node_key + validator private key

```
bandd init <moniker> --chain-id bandchain
```

After run this command, It will generate 2 file that you need to keep in `$HOME/.bandd/config`

- node_key.json
- priv_validator_key.json

2. Generate operator address using in our chain

```
bandcli keys add <key_name>
```

Fill your password
Please remember your mnemonic for this key

3. Add your validator address to genesis state for gentx

```
bandd add-genesis-account $(bandcli keys show <key_name> -a) 1uband
```

4. Generate your create validator msg

```
bandd gentx --amount 1uband --name <key_name>
```

Your gentx result will store in `$HOME/.bandd/config/gentx` and upload a file in this directory to us.
Your name file will be `gentx-........json`
