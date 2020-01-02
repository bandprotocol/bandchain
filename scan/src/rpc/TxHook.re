module Msg = {
  module OracleReport = {
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
    | OracleReport(OracleReport.t);

  let decode = json =>
    JsonUtils.Decode.(
      switch (json |> field("type", string)) {
      | "zoracle/Report" => OracleReport(json |> field("value", OracleReport.decode))
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
