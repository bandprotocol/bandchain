module Proof = {
  type t = {
    jsonProof: Js.Json.t,
    evmProofBytes: JsBuffer.t,
  };

  let decodeProof = json =>
    JsonUtils.Decode.{
      jsonProof: json |> at(["result", "jsonProof"], json => json),
      evmProofBytes: json |> at(["result", "evmProofBytes"], string) |> JsBuffer.fromHex,
    };
};

let get = (~requestId: int, ~pollInterval=?, ()) => {
  let json = Axios.use({j|d3n/proof/$requestId|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Proof.decodeProof);
};
