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
      fromAddress: string,
      toAddress: string,
      amount: list(Coin.t),
    };

    let decode = json =>
      JsonUtils.Decode.{
        fromAddress: json |> field("from_address", string),
        toAddress: json |> field("to_address", string),
        amount: json |> field("amount", list(Coin.decodeCoin)),
      };
  };

  module Store = {
    type t = {
      code: string,
      owner: string,
    };

    let decode = json =>
      JsonUtils.Decode.{
        code: json |> field("code", string),
        owner: json |> field("owner", string),
      };
  };

  module Request = {
    type t = {
      codeHash: string,
      params: string,
      reportPeriod: int,
      sender: string,
    };

    let decode = json =>
      JsonUtils.Decode.{
        codeHash: json |> field("codeHash", string),
        params: json |> field("params", string),
        reportPeriod: json |> field("reportPeriod", intstr),
        sender: json |> field("sender", string),
      };
  };

  module Report = {
    type t = {
      requestId: int,
      data: string,
      validator: string,
    };

    let decode = json =>
      JsonUtils.Decode.{
        requestId: json |> field("requestID", intstr),
        data: json |> field("data", string),
        validator: json |> field("validator", string),
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

module Tx = {
  type t = {
    blockHeight: int,
    hash: string,
    timestamp: MomentRe.Moment.t,
    gasWanted: int,
    gasUsed: int,
    messages: list(Msg.t),
  };

  let decodeTx = json =>
    JsonUtils.Decode.{
      blockHeight: json |> field("height", intstr),
      hash: json |> field("txhash", string),
      timestamp: json |> field("timestamp", moment),
      gasWanted: json |> field("gas_wanted", intstr),
      gasUsed: json |> field("gas_used", intstr),
      messages: json |> at(["tx", "value", "msg"], list(Msg.decode)),
    };

  let decodeTxs = json => JsonUtils.Decode.(json |> field("txs", list(decodeTx)));
};

let at_hash = tx_hash => {
  let json = Axios.use({j|txs/$tx_hash|j}, ());
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
