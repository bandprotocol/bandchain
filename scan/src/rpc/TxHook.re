module Msg = {
  module OracleReport = {
    type t = {
      request_id: int,
      data: string,
      validator: string,
    };

    let decode = json =>
      JsonUtils.Decode.{
        request_id: json |> field("requestID", intstr),
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
    block_height: int,
    hash: string,
    timestamp: MomentRe.Moment.t,
    gas_wanted: int,
    gas_used: int,
    messages: list(Msg.t),
  };

  let decode_tx = json =>
    JsonUtils.Decode.{
      block_height: json |> field("height", intstr),
      hash: json |> field("txhash", string),
      timestamp: json |> field("timestamp", moment),
      gas_wanted: json |> field("gas_wanted", intstr),
      gas_used: json |> field("gas_used", intstr),
      messages: json |> at(["tx", "value", "msg"], list(Msg.decode)),
    };

  let decode_txs = json => JsonUtils.Decode.(json |> field("txs", list(decode_tx)));
};

let at_hash = tx_hash => {
  let json = Axios.use({j|txs/$tx_hash|j}, ());
  json |> Belt.Option.map(_, Tx.decode_tx);
};

let at_height = (height, ~page=1, ~limit=25, ~pollInterval=?, ()) => {
  let json = Axios.use({j|txs?tx.height=$height&page=$page&limit=$limit|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Tx.decode_txs);
};
