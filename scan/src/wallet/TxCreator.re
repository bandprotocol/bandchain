type msg_send_t = {
  to_address: string,
  from_address: string,
  amount: array(Coin.t),
};

type msg_request_t = {
  oracleScriptID: string,
  calldata: string,
  requestedValidatorCount: string,
  sufficientValidatorCount: string,
  sender: string,
  clientID: string,
};

type amount_t = {
  amount: string,
  denom: string,
};

type fee_t = {
  amount: array(amount_t),
  gas: string,
};

type msg_input_t =
  | Send(Address.t, Address.t, Coin.t)
  | Request(ID.OracleScript.t, JsBuffer.t, string, string, Address.t, string);

type msg_payload_t = {
  [@bs.as "type"]
  type_: string,
  value: Js.Json.t,
};

type account_result_t = {
  accountNumber: int,
  sequence: int,
};

type pub_key_t = {
  [@bs.as "type"]
  type_: string,
  value: string,
};

type signature_t = {
  pub_key: pub_key_t,
  public_key: string,
  signature: string,
};

type raw_tx_t = {
  msgs: array(msg_payload_t),
  chain_id: string,
  fee: fee_t,
  memo: string,
  account_number: string,
  sequence: string,
};

type signed_tx_t = {
  fee: fee_t,
  memo: string,
  msg: array(msg_payload_t),
  signatures: array(signature_t),
};

type t = {
  mode: string,
  tx: signed_tx_t,
};

type tx_response_t = {
  txHash: Hash.t,
  rawLog: string,
  success: bool,
};

type response_t =
  | Tx(tx_response_t)
  | Unknown;

let getAccountInfo = address => {
  let url = Env.rpc ++ "/auth/accounts/" ++ (address |> Address.toBech32);
  let%Promise info = Axios.get(url);
  let data = info##data;
  Promise.ret(
    JsonUtils.Decode.{
      accountNumber: data |> at(["result", "value", "account_number"], int),
      sequence: data |> at(["result", "value", "sequence"], int),
    },
  );
};

let sortAndStringify: raw_tx_t => string = [%bs.raw
  {|
  function sortAndStringify(obj) {
    function sortObject(obj) {
      if (obj === null) return null;
      if (typeof obj !== "object") return obj;
      if (Array.isArray(obj)) return obj.map(sortObject);
      const sortedKeys = Object.keys(obj).sort();
      const result = {};
      sortedKeys.forEach(key => {
        result[key] = sortObject(obj[key])
      });
      return result;
    }

    return JSON.stringify(sortObject(obj));
  }
|}
];

let createMsg = (msg: msg_input_t): msg_payload_t => {
  let msgType =
    switch (msg) {
    | Send(_) => "cosmos-sdk/MsgSend"
    | Request(_) => "oracle/Request"
    };

  let msgValue =
    switch (msg) {
    | Send(fromAddress, toAddress, coins) =>
      Js.Json.stringifyAny({
        to_address: toAddress |> Address.toBech32,
        from_address: fromAddress |> Address.toBech32,
        amount: [|coins|],
      })
      |> Belt_Option.getExn
      |> Js.Json.parseExn
    | Request(
        ID.OracleScript.ID(oracleScriptID),
        calldata,
        requestedValidatorCount,
        sufficientValidatorCount,
        sender,
        clientID,
      ) =>
      Js.Json.stringifyAny({
        oracleScriptID: oracleScriptID |> string_of_int,
        calldata: calldata |> JsBuffer.toBase64,
        requestedValidatorCount,
        sufficientValidatorCount,
        sender: sender |> Address.toBech32,
        clientID,
      })
      |> Belt_Option.getExn
      |> Js.Json.parseExn
    };
  {type_: msgType, value: msgValue};
};

// TODO: Reme hardcoded values
let createRawTx = (address, msgs) => {
  let%Promise accountInfo = getAccountInfo(address);
  Promise.ret({
    msgs: msgs->Belt_Array.map(createMsg),
    chain_id: "bandchain",
    fee: {
      amount: [|{amount: "1000000", denom: "uband"}|],
      gas: "3000000",
    },
    memo: "",
    account_number: accountInfo.accountNumber |> string_of_int,
    sequence: accountInfo.sequence |> string_of_int,
  });
};

let createSignedTx = (~signature, ~pubKey, ~tx: raw_tx_t, ~mode, ()) => {
  let oldPubKey = {type_: "tendermint/PubKeySecp256k1", value: pubKey |> PubKey.toBase64};
  let newPubKey = "eb5ae98721" ++ (pubKey |> PubKey.toHex) |> JsBuffer.hexToBase64;
  let signedTx = {
    fee: tx.fee,
    memo: tx.memo,
    msg: tx.msgs,
    signatures: [|{pub_key: oldPubKey, public_key: newPubKey, signature}|],
  };
  {mode, tx: signedTx};
};

let broadcast = signedTx => {
  /* TODO: FIX THIS MESS */
  let convert: t => Js.t('a) = [%bs.raw {|
function(data) {return {...data};}
  |}];

  let%Promise rawResponse = Axios.postData(Env.rpc ++ "/txs", convert(signedTx));
  let response = rawResponse##data;
  Promise.ret(
    Tx(
      JsonUtils.Decode.{
        txHash: response |> at(["txhash"], string) |> Hash.fromHex,
        rawLog: response |> at(["raw_log"], string),
        success: response |> optional(field("logs", _ => ())) |> Belt_Option.isSome,
      },
    ),
  );
};
