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
      code: JsBuffer.t,
      owner: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        code: json |> field("code", string) |> JsBuffer.fromBase64,
        owner: json |> field("owner", string) |> Address.fromBech32,
      };
  };

  module Request = {
    type t = {
      codeHash: Hash.t,
      params: JsBuffer.t,
      reportPeriod: int,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        codeHash: json |> field("codeHash", string) |> Hash.fromBase64,
        params: json |> field("params", string) |> JsBuffer.fromBase64,
        reportPeriod: json |> field("reportPeriod", intstr),
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module Report = {
    type t = {
      requestId: int,
      data: JsBuffer.t,
      validator: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        requestId: json |> field("requestID", intstr),
        data: json |> field("data", string) |> JsBuffer.fromBase64,
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

  let decodeTx = json =>
    JsonUtils.Decode.{
      sender:
        json
        |> at(["tx", "value", "signatures"], list(Signature.decode))
        |> Belt_List.getExn(_, 0)
        |> ((firstSignature: Signature.t) => firstSignature.pubKey)
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

let atHash = txHash => {
  let txHashHex = txHash->Hash.toHex;
  let json = Axios.use({j|txs/$txHashHex|j}, ());
  json |> Belt.Option.map(_, Tx.decodeTx);
};

let atHeight = (height, ~page=1, ~limit=25, ~pollInterval=?, ()) => {
  let json = Axios.use({j|txs?tx.height=$height&page=$page&limit=$limit|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Tx.decodeTxs);
};

let latest = (~page=1, ~limit=10, ~pollInterval=?, ()) => {
  let json = Axios.use({j|d3n/txs/latest?page=$page&limit=$limit|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Tx.decodeTxs);
};

let withCodehash = (~codeHash, ~page=1, ~limit=10, ~pollInterval=?, ()) => {
  let codeHashHex = codeHash->Hash.toHex;
  let json =
    Axios.use(
      {j|txs?request.codehash=$codeHashHex&page=$page&limit=$limit|j},
      ~pollInterval?,
      (),
    );
  json |> Belt.Option.map(_, Tx.decodeTxs);
};
