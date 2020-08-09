# ICON Bridge

## Setup

```
virtualenv -p python3 .
source bin/activate
pip install tbears
```

### Run test

```
tbears test bridge
```

### Encoding (OBI)

- `validators_bytes`:
  An Array of validators with voting power which is a parameter of function `on_install` and function `update_validator_powers`

  #### Struct

  ```
      PyObi(
          """
          [
              {
                  pubkey:bytes,
                  power:u64
              }
          ]
          """
      )
  ```

  #### Example

  ```
      input: [
            {
                pubkey: "a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5",
                power: 100
            },
            {
                pubkey: "724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2",
                power: 100
            },
            {
                pubkey: "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479",
                power: 100
            },
            {
                pubkey: "d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f2c67aaeee2bc8d00ff253a326eb04de8ce88eea1224c2bb8ca69296a1c753f53",
                power: 100
            }
      ]
      output: 0000000400000040a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5000000000000006400000040724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2000000000000006400000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479000000000000006400000040d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f2c67aaeee2bc8d00ff253a326eb04de8ce88eea1224c2bb8ca69296a1c753f530000000000000064
  ```

- `multi_store_bytes`: The concatenation of stores hash in Bandchain which is one of the parameters of a function `relay_oracle_state`

  #### Struct

  ```
    PyObi("""bytes""")
  ```

  ### Example

  ```
  //  acc_to_gov_stores_merkle_hash = 10d1bb5a23a0ed3668df0b2a4257b46225db7391c56e940ad95536d0e5a75d79
  //  main_and_mint_stores_merkle_hash = f4be0eeff9a2e52e11da7591c029cee08ecfe59435dbd731679905a94c2bb53b
  //  oracle_iavl_state_hash = 8fdacf87d1d46d85ee2585ca5a304a2a499080b8c7460f5efac7d053911fbee6
  //  params_stores_merkle_hash = b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9
  //  slashing_to_upgrade_stores_merkle_hash = df898ce21b5981ca9994759bb2b77204496e9d062a080b52c2b38b8d8df0f552

  input: 10d1bb5a23a0ed3668df0b2a4257b46225db7391c56e940ad95536d0e5a75d79f4be0eeff9a2e52e11da7591c029cee08ecfe59435dbd731679905a94c2bb53b8fdacf87d1d46d85ee2585ca5a304a2a499080b8c7460f5efac7d053911fbee6b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9df898ce21b5981ca9994759bb2b77204496e9d062a080b52c2b38b8d8df0f552
  output:  0000014010d1bb5a23a0ed3668df0b2a4257b46225db7391c56e940ad95536d0e5a75d79f4be0eeff9a2e52e11da7591c029cee08ecfe59435dbd731679905a94c2bb53b8fdacf87d1d46d85ee2585ca5a304a2a499080b8c7460f5efac7d053911fbee6b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9df898ce21b5981ca9994759bb2b77204496e9d062a080b52c2b38b8d8df0f552
  ```
