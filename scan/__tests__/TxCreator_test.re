open Jest;

open TxCreator;

open PubKey;

open Expect;

let pubKey =
  PubKey("eb5ae9872103a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33");
let signature = "CHdU7pVFBLl4GqWvMNlyOh5fdoOagkf3MSf5UfY7DzAEzVX2YOUZpbEKuBDDvEGDTc3u0Pl42zE04GLpSfQzOw";

describe("expect TxCreator to give the correct message", () => {
  test("should be able to create correct message for MsgSend", () => {
    expect({
      mode: "block",
      tx: {
        msg: [|
          {
            type_: "cosmos-sdk/MsgSend",
            value:
              Js.Json.stringifyAny({
                amount: [|{amount: "100.", denom: "uband"}|],
                from_address: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
                to_address: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
              })
              |> Belt_Option.getExn
              |> Js.Json.parseExn,
          },
        |],
        fee: {
          amount: [|{amount: "10000", denom: "uband"}|],
          gas: "1000000",
        },
        memo: "",
        signatures: [|
          {
            pub_key:
              Js.Json.object_(
                Js.Dict.fromList([
                  ("type", Js.Json.string("tendermint/PubKeySecp256k1")),
                  ("value", Js.Json.string(pubKey |> PubKey.toBase64)),
                ]),
              ),
            public_key: "eb5ae98721" ++ (pubKey |> PubKey.toHex) |> JsBuffer.hexToBase64,
            signature,
          },
        |],
      },
    })
    |> toEqual(
         createSignedTx(
           ~network="WENCHANG",
           ~signature,
           ~pubKey,
           ~mode="block",
           ~tx={
             fee: {
               amount: [|{amount: "10000", denom: "uband"}|],
               gas: "1000000",
             },
             memo: "",
             chain_id: "bandchain",
             account_number: "2",
             sequence: "2",
             msgs: [|
               {
                 type_: "cosmos-sdk/MsgSend",
                 value:
                   Js.Json.stringifyAny({
                     amount: [|{amount: "100.", denom: "uband"}|],
                     from_address: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
                     to_address: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
                   })
                   |> Belt_Option.getExn
                   |> Js.Json.parseExn,
               },
             |],
           },
           (),
         ),
       )
  });
  test("should be able to create correct message for MsgRequest", () => {
    expect({
      mode: "block",
      tx: {
        msg: [|
          {
            type_: "oracle/Request",
            value:
              Js.Json.stringifyAny({
                oracle_script_id: "1",
                calldata: "RVRI",
                ask_count: "4",
                min_count: "4",
                sender: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
                client_id: "",
                fee_limit: [|{amount: "100.", denom: "uband"}|],
                prepare_gas: "30000",
                execute_gas: "50000",
              })
              |> Belt_Option.getExn
              |> Js.Json.parseExn,
          },
        |],
        fee: {
          amount: [|{amount: "10000", denom: "uband"}|],
          gas: "1000000",
        },
        memo: "",
        signatures: [|
          {
            pub_key:
              Js.Json.string("eb5ae98721" ++ (pubKey |> PubKey.toHex) |> JsBuffer.hexToBase64),
            public_key: "eb5ae98721" ++ (pubKey |> PubKey.toHex) |> JsBuffer.hexToBase64,
            signature,
          },
        |],
      },
    })
    |> toEqual(
         createSignedTx(
           ~network="GUANYU",
           ~signature,
           ~pubKey,
           ~mode="block",
           ~tx={
             fee: {
               amount: [|{amount: "10000", denom: "uband"}|],
               gas: "1000000",
             },
             memo: "",
             chain_id: "bandchain",
             account_number: "2",
             sequence: "2",
             msgs: [|
               {
                 type_: "oracle/Request",
                 value:
                   Js.Json.stringifyAny({
                     oracle_script_id: "1",
                     calldata: "RVRI",
                     ask_count: "4",
                     min_count: "4",
                     sender: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
                     client_id: "",
                     fee_limit: [|{amount: "100.", denom: "uband"}|],
                     prepare_gas: "30000",
                     execute_gas: "50000",
                   })
                   |> Belt_Option.getExn
                   |> Js.Json.parseExn,
               },
             |],
           },
           (),
         ),
       )
  });
});
