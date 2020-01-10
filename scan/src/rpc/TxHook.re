module Coin = {
  type t = {
    denom: string,
    amount: float,
  };

  let decodeCoin = json =>
    JsonUtils.Decode.{
      denom: json |> field("denom", string),
      amount: json |> field("amount", uamount),
    };
};

module Msg = {
  module Send = {
    type t = {
      fromAddress: Address.t,
      toAddress: Address.t,
      amount: list(Coin.t),
    };

    let decode = json =>
      JsonUtils.Decode.{
        fromAddress: json |> field("from_address", string) |> Address.fromBech32,
        toAddress: json |> field("to_address", string) |> Address.fromBech32,
        amount: json |> field("amount", list(Coin.decodeCoin)),
      };
  };

  module Store = {
    type t = {
      code: string,
      owner: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        code: json |> field("code", string),
        owner: json |> field("owner", string) |> Address.fromBech32,
      };
  };

  module Request = {
    type t = {
      codeHash: Hash.t,
      params: string,
      reportPeriod: int,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        codeHash: json |> field("codeHash", string) |> Hash.fromBase64,
        params: json |> field("params", string),
        reportPeriod: json |> field("reportPeriod", intstr),
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module Report = {
    type t = {
      requestId: int,
      data: string,
      validator: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        requestId: json |> field("requestID", intstr),
        data: json |> field("data", string),
        validator: json |> field("validator", string) |> Address.fromBech32,
      };
  };

  type t =
    | Unknown
    | Send(Send.t)
    | Store(Store.t)
    | Request(Request.t)
    | Report(Report.t);

  let decode = json =>
    JsonUtils.Decode.(
      switch (json |> field("type", string)) {
      | "cosmos-sdk/MsgSend" => Send(json |> field("value", Send.decode))
      | "zoracle/Store" => Store(json |> field("value", Store.decode))
      | "zoracle/Request" => Request(json |> field("value", Request.decode))
      | "zoracle/Report" => Report(json |> field("value", Report.decode))
      | _ => Unknown
      }
    );
};

module Signature = {
  type t = {
    pubKey: PubKey.t,
    pubKeyType: string,
    signature: JsBuffer.t,
  };

  let decode = json =>
    JsonUtils.Decode.{
      pubKey: json |> at(["pub_key", "value"], string) |> PubKey.fromBase64,
      pubKeyType: json |> at(["pub_key", "type"], string),
      signature: json |> field("signature", string) |> JsBuffer.fromBase64,
    };
};

module Tx = {
  type t = {
    sender: Address.t,
    blockHeight: int,
    hash: Hash.t,
    timestamp: MomentRe.Moment.t,
    gasWanted: int,
    gasUsed: int,
    messages: list(Msg.t),
  };

  let getFirstSignerAddress = (sigsList: list(Signature.t)) => {
    let sigsArr = sigsList |> Belt_List.toArray;
    sigsArr[0].pubKey;
  };

  let decodeTx = json =>
    JsonUtils.Decode.{
      sender:
        json
        |> at(["tx", "value", "signatures"], list(Signature.decode))
        |> getFirstSignerAddress
        |> PubKey.toAddress,
      blockHeight: json |> field("height", intstr),
      hash: json |> field("txhash", string) |> Hash.fromHex,
      timestamp: json |> field("timestamp", moment),
      gasWanted: json |> field("gas_wanted", intstr),
      gasUsed: json |> field("gas_used", intstr),
      messages: json |> at(["tx", "value", "msg"], list(Msg.decode)),
    };

  let decodeTxs = json => JsonUtils.Decode.(json |> field("txs", list(decodeTx)));
};

let at_hash = tx_hash => {
  let tx_hash_hex_str = tx_hash->Hash.toHex;
  let json = Axios.use({j|txs/$tx_hash_hex_str|j}, ());
  json |> Belt.Option.map(_, Tx.decodeTx);
};

let at_height = (height, ~page=1, ~limit=25, ~pollInterval=?, ()) => {
  let json = Axios.use({j|txs?tx.height=$height&page=$page&limit=$limit|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Tx.decodeTxs);
};

let latest = (~page=1, ~limit=10, ~pollInterval=?, ()) => {
  let json = Axios.use({j|d3n/txs/latest?page=$page&limit=$limit|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Tx.decodeTxs);
};

let with_code_hash = (~code_hash, ~page=1, ~limit=10, ~pollInterval=?, ()) => {
  let code_hash_hex_str = code_hash->Hash.toHex;
  let json =
    Axios.use(
      {j|txs?request.codehash=$code_hash_hex_str&page=$page&limit=$limit|j},
      ~pollInterval?,
      (),
    );
  json |> Belt.Option.map(_, Tx.decodeTxs);
};
