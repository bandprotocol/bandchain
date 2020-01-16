module Proof = {
  type t = Js.Json.t;

  let decodeProof = JsonUtils.Decode.field("result", json => json);
};

let get = (~requestId: int, ~pollInterval=?, ()) => {
  let json = Axios.use({j|d3n/proof/$requestId|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Proof.decodeProof);
};
