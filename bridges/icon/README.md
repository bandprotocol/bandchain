# ICON Bridge

## Setup

```
virtualenv -p python3 .
source bin/activate
pip install tbears
```

## Run test

```
tbears test bridge
```

## Deployment

#### Bridge

```
tbears deploy bridge -k mumustorekey.key -c bridge/tbears_cli_config_testnet.json
```

#### ReceiverMock

```
tbears deploy receiver_mock -k mumustorekey.key -c receiver_mock/deploy_config_testnet.json
```

#### CacheConsumerMock

```
tbears deploy cache_consumer_mock -k mumustorekey.key -c cache_consumer_mock/deploy_config_testnet.json
```

## Interaction

#### Bridge

- relay (only relay the newer oracle state)
  ```
  tbears sendtx ./bridge/example_send_relay.json  -k mumustorekey.key -u https://bicon.net.solidwallet.io/api/v3 -n 3
  ```
- relay_and_verify
  ```
  tbears sendtx ./bridge/example_send_relay_and_verify.json  -k mumustorekey.key -u https://bicon.net.solidwallet.io/api/v3 -n 3
  ```

#### ReceiverMock

- relay_and_safe
  ```
  tbears sendtx ./receiver_mock/example_send_relay_and_safe.json -k mumustorekey.key -u https://bicon.net.solidwallet.io/api/v3 -n 3
  ```

#### CacheConsumerMock

- consume_cache

```
tbears sendtx ./cache_consumer_mock/example_send_consume.json  -k mumustorekey.key -u https://bicon.net.solidwallet.io/api/v3 -n 3
```

## Example deployed contracts on the testnet

#### Bridge contract

[`cx947d7eb9766c4a59a127967574add8bf22304c8e`](https://bicon.tracker.solidwallet.io/contract/cx947d7eb9766c4a59a127967574add8bf22304c8e#readcontract)

#### ReceiverMock contract

[`cxc4b59b0def8bbc2306c77d0e9d11075407a9f70e`](https://bicon.tracker.solidwallet.io/contract/cxc4b59b0def8bbc2306c77d0e9d11075407a9f70e#readcontract)

#### CacheConsumerMock contract

[`cx71e33c94c5827af2039578fb7ad30044326e4348`](https://bicon.tracker.solidwallet.io/contract/cxc4b59b0def8bbc2306c77d0e9d11075407a9f70e#readcontract)

## Encoding (OBI)

- `validators_bytes`:
  An Array of validators with voting power which is a parameter of function `on_install` and function `update_validator_powers`

  #### Encoding/Decoding

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
      encoded_input: 0000000400000040a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5000000000000006400000040724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2000000000000006400000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479000000000000006400000040d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f2c67aaeee2bc8d00ff253a326eb04de8ce88eea1224c2bb8ca69296a1c753f530000000000000064
  ```

- `multi_store_bytes`: The concatenation of stores hash in Bandchain which is one of the parameters of a function `relay_oracle_state`

  #### Encoding/Decoding

  ```
    bytes
  ```

  ### Example

  ```
  //  acc_to_gov_stores_merkle_hash = 10d1bb5a23a0ed3668df0b2a4257b46225db7391c56e940ad95536d0e5a75d79
  //  main_and_mint_stores_merkle_hash = f4be0eeff9a2e52e11da7591c029cee08ecfe59435dbd731679905a94c2bb53b
  //  oracle_iavl_state_hash = 8fdacf87d1d46d85ee2585ca5a304a2a499080b8c7460f5efac7d053911fbee6
  //  params_stores_merkle_hash = b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9
  //  slashing_to_upgrade_stores_merkle_hash = df898ce21b5981ca9994759bb2b77204496e9d062a080b52c2b38b8d8df0f552

  input: 10d1bb5a23a0ed3668df0b2a4257b46225db7391c56e940ad95536d0e5a75d79f4be0eeff9a2e52e11da7591c029cee08ecfe59435dbd731679905a94c2bb53b8fdacf87d1d46d85ee2585ca5a304a2a499080b8c7460f5efac7d053911fbee6b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9df898ce21b5981ca9994759bb2b77204496e9d062a080b52c2b38b8d8df0f552
  ```

- `block_merkle_part_bytes`: The concatenation of extra merkle parts to compute block hash which is one of the parameters of a function `relay_oracle_state`

  #### Encoding/Decoding

  ```
    bytes
  ```

  ### Example

  ```
    // versionAndChainIdHash = 32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e
    // timeHash = 8f7beb334b2ff0c5c176144edbbfe9426cdb41fea466d43d0b89e2515f9e5a51
    // lastBlockIDAndOther = 70e961a346803fa2198490f4ded2683e31384b2a37725e719ae7ca164eb22ceb
    // nextValidatorHashAndConsensusHash = 004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d9
    // lastResultsHash = 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d
    // evidenceAndProposerHash = 0efe3e12f46363c7779140d4ce659925db52f19053e114d7cc4efd666b37f79f

  input: 32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e8f7beb334b2ff0c5c176144edbbfe9426cdb41fea466d43d0b89e2515f9e5a5170e961a346803fa2198490f4ded2683e31384b2a37725e719ae7ca164eb22ceb004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d0efe3e12f46363c7779140d4ce659925db52f19053e114d7cc4efd666b37f79f
  ```

- `signatures_bytes`: The OBI encoded of an array of signatures signed on the block which is one of the parameters of a function `relay_oracle_state`

  #### Encoding/Decoding

  ```
    PyObi(
        """
        [
            {
                r: bytes,
                s: bytes,
                v: u8,
                signed_data_prefix: bytes,
                signed_data_suffix: bytes
            }
        ]
        """
    )
  ```

  ### Example

  ```
  input: [
            {
                r: "628716ac49023de84adddddcbef8007c2e41e5b58306ce87a0afad5447bc6210"
                s: "2f520db2bff3003d5612e03b7aaa99472164c73922a977af95e1ffc2a67c53b4"
                v: 27
                signed_data_prefix: "6e080211840500000000000022480a20"
                signed_data_suffix: "12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510cbe1cfbe02320962616e64636861696e"
            },
            {
                r: "ff2ba7e2bd2175827997c706451b5da768b6873d7ba4129fc6ee54e62ba9c593"
                s: "3c7f7e5b08d1733d430658431545c9a2f57e4641b3b4cd52e567f27be9485e60"
                v: 28
                signed_data_prefix: "6e080211840500000000000022480a20"
                signed_data_suffix: "12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a3f0e5be02320962616e64636861696e"
            },
            {
                r: "5a2f66b4d62d905b98277cd2807a324f0651340e80ae0249e500beb5ddcdce11"
                s: "3c1ed3d960b19e0ca7d321215874c6e91407ae3d2748f2e3b617fad833c30b6d"
                v: 27
                signed_data_prefix: "6e080211840500000000000022480a20"
                signed_data_suffix: "12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a394a6bf02320962616e64636861696e"
            },
        ]
  encoded_input:  0000000300000020628716ac49023de84adddddcbef8007c2e41e5b58306ce87a0afad5447bc6210000000202f520db2bff3003d5612e03b7aaa99472164c73922a977af95e1ffc2a67c53b41b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510cbe1cfbe02320962616e64636861696e00000020ff2ba7e2bd2175827997c706451b5da768b6873d7ba4129fc6ee54e62ba9c593000000203c7f7e5b08d1733d430658431545c9a2f57e4641b3b4cd52e567f27be9485e601c000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a3f0e5be02320962616e64636861696e000000205a2f66b4d62d905b98277cd2807a324f0651340e80ae0249e500beb5ddcdce11000000203c1ed3d960b19e0ca7d321215874c6e91407ae3d2748f2e3b617fad833c30b6d1b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a394a6bf02320962616e64636861696e
  ```

- `encode_packet`: The OBI encode of request packet and response packet which is one of the parameters of a function `verify_oracle_data`

  #### Encoding/Decoding

  ```
    PyObi(
        """
        {
            req: {
                client_id: string,
                oracle_script_id: u64,
                calldata: bytes,
                ask_count: u64,
                min_count: u64
            },
            res: {
                client_id: string,
                request_id: u64,
                ans_count: u64,
                request_time: u64,
                resolve_time: u64,
                resolve_status: u32,
                result: bytes
            }
        }
        """
    )
  ```

  ### Example

  ```
  input: {
        req: {
            client_id: "from_scan",
            oracle_script_id: 13,
            params: "0000000342544300000003555344000000046d65616e0000000000000064",
            ans_count: 4,
            min_count: 4,
        },
        res: {
            client_id: "from_scan",
            request_id: 9,
            ans_count: 4,
            request_time: 1593001602,
            resolve_time: 1593001606,
            resolve_status: 2,
            result: 0000000000000064,
        }
  }
  encoded_input:  0000000966726f6d5f7363616e000000000000000d0000001e0000000342544300000003555344000000046d65616e0000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000090000000000000004000000005ef34682000000005ef3468600000002000000080000000000000064
  ```

- `iavl_merkle_paths`: The merkle proof that shows how the data leave is part of the oracle iAVL which is one of the parameters of a function `verify_oracle_data`

  #### Encoding/Decoding

  ```
    PyObi(
        """
        [
            {
                is_data_on_right: bool,
                subtree_height: u8,
                subtree_size: u64,
                subtree_version: u64,
                sibling_hash: bytes
            }
        ]
        """
    )
  ```

  ### Example

  ```
  input: [
    {
        is_data_on_right: true,
        subtree_height: 1,
        subtree_size: 2,
        subtree_version: 444,
        sibling_hash: "1d3927941aec0602e1fcab94b4d3e57da16a2455165cb097adfdfdf34a3e1688",
    },
    {
        is_data_on_right: true,
        subtree_height: 1,
        subtree_size: 2,
        subtree_version: 444,
        sibling_hash: "1d3927941aec0602e1fcab94b4d3e57da16a2455165cb097adfdfdf34a3e1688",
    }
  ]

  encoded_input: 000000020101000000000000000200000000000001bc000000201d3927941aec0602e1fcab94b4d3e57da16a2455165cb097adfdfdf34a3e16880102000000000000000400000000000001bc0000002017480c9a0faca2c2ecef5fdc7627b11f898f5652f3022e8d0b2da40356d6d8e4
  ```

- `proof`: The encoded data for oracle state relay and data verification which is the parameter of function `relay_and_verify` and function `relay`

  #### Encoding/Decoding

  ```
    PyObi(
        """
        {
            block_height: u64,
            multi_store: bytes,
            merkle_parts: bytes,
            signatures: bytes,
            encoded_packet: bytes,
            version: u64,
            iavl_merkle_paths: bytes
        }
        """
    )
  ```

  ### Example

  ```
  input: {
        block_height: 1412,
        multi_store: "811032f41a5918cfc860ebcea1b7678564c622b4e3f04e28c9d8b195f04661e11e06ea167a89de60abf52c60350793c3823a0a4b9681f59daeb8f20942b8bcf3ff2b888760ffa9efa09e35763d8a5ba3ec794d00bcd1f9a901fa007f1f16e2ebb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9c6507de937080730702a2c39ad8aa4666fd24e12a9496e4546a0bd4031e6d4ee",
        block_merkle_parts: "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e2a448c96300f1cb90567d815c6b2f6ec0bba29615fd88deab57299c76f0b30fa3ebf2264c2c941df347eb97f55b4adcd4784f4f0a8e66d9f2e3fda413115e085004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d0efe3e12f46363c7779140d4ce659925db52f19053e114d7cc4efd666b37f79f",
        signatures: "0000000300000020628716ac49023de84adddddcbef8007c2e41e5b58306ce87a0afad5447bc6210000000202f520db2bff3003d5612e03b7aaa99472164c73922a977af95e1ffc2a67c53b41b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510cbe1cfbe02320962616e64636861696e00000020ff2ba7e2bd2175827997c706451b5da768b6873d7ba4129fc6ee54e62ba9c593000000203c7f7e5b08d1733d430658431545c9a2f57e4641b3b4cd52e567f27be9485e601c000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a3f0e5be02320962616e64636861696e000000205a2f66b4d62d905b98277cd2807a324f0651340e80ae0249e500beb5ddcdce11000000203c1ed3d960b19e0ca7d321215874c6e91407ae3d2748f2e3b617fad833c30b6d1b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a394a6bf02320962616e64636861696e"
        encode_packet: "0000000966726f6d5f7363616e000000000000000d0000001e0000000342544300000003555344000000046d65616e0000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000090000000000000004000000005ef34682000000005ef3468600000002000000080000000000000064"
        version: 1409,
        iavl_paths: "000000020101000000000000000200000000000001bc000000201d3927941aec0602e1fcab94b4d3e57da16a2455165cb097adfdfdf34a3e16880102000000000000000400000000000001bc0000002017480c9a0faca2c2ecef5fdc7627b11f898f5652f3022e8d0b2da40356d6d8e4",
    }

  encoded_input: 0000000000000584000000a0811032f41a5918cfc860ebcea1b7678564c622b4e3f04e28c9d8b195f04661e11e06ea167a89de60abf52c60350793c3823a0a4b9681f59daeb8f20942b8bcf3ff2b888760ffa9efa09e35763d8a5ba3ec794d00bcd1f9a901fa007f1f16e2ebb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9c6507de937080730702a2c39ad8aa4666fd24e12a9496e4546a0bd4031e6d4ee000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e2a448c96300f1cb90567d815c6b2f6ec0bba29615fd88deab57299c76f0b30fa3ebf2264c2c941df347eb97f55b4adcd4784f4f0a8e66d9f2e3fda413115e085004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d0efe3e12f46363c7779140d4ce659925db52f19053e114d7cc4efd666b37f79f000001e40000000300000020628716ac49023de84adddddcbef8007c2e41e5b58306ce87a0afad5447bc6210000000202f520db2bff3003d5612e03b7aaa99472164c73922a977af95e1ffc2a67c53b41b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510cbe1cfbe02320962616e64636861696e00000020ff2ba7e2bd2175827997c706451b5da768b6873d7ba4129fc6ee54e62ba9c593000000203c7f7e5b08d1733d430658431545c9a2f57e4641b3b4cd52e567f27be9485e601c000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a3f0e5be02320962616e64636861696e000000205a2f66b4d62d905b98277cd2807a324f0651340e80ae0249e500beb5ddcdce11000000203c1ed3d960b19e0ca7d321215874c6e91407ae3d2748f2e3b617fad833c30b6d1b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a394a6bf02320962616e64636861696e000000840000000966726f6d5f7363616e000000000000000d0000001e0000000342544300000003555344000000046d65616e0000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000090000000000000004000000005ef34682000000005ef346860000000200000008000000000000006400000000000005810000014800000006010100000000000000020000000000000581000000208d09a67a6aebb1498f3d104ae86c23d0c8e2e86dcb0a96f1088aca50b6d5e4030102000000000000000300000000000005810000002021c6cca538be5e9673479ba5694886335cdbbd0fcccd8d62376728c1133f00270103000000000000000500000000000005810000002057c9a0c6d2089b3911dc3bfd173ccd88922ed48030686ba3b3f8e4325b5f3cac0104000000000000000c00000000000005810000002014dc21a48e352ec5d08b52473da16cc6ca515e9e94124ea68b84a011b5da400c010500000000000000170000000000000581000000200fcacf91a581cf180802122abb9d02a67590fe38c7aea5baaa53416e91bf30e201070000000000000043000000000000058300000020c0062b4ff613ecbcf16e37a4423961afebf0302a208702d67e583622ad4a3dda
  ```
