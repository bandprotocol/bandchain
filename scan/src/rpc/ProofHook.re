module Proof = {
  type t = {
    proof: Js.Json.t,
    evmProofBytes: JsBuffer.t,
  };

  let decodeProof = json =>
    JsonUtils.Decode.{
      proof: json |> at(["result", "proof"], json => json),
      evmProofBytes: json |> at(["result", "evmProofBytes"], string) |> JsBuffer.fromHex,
    };
};

let get = (requestId: ID.Request.t) => {
  let (json, reload) = AxiosHooks.useWithReload({j|oracle/proof/$requestId|j});
  (json |> Belt.Option.map(_, Proof.decodeProof), reload);
};
