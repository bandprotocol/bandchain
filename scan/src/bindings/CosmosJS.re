type t;

type tx_response_t = {
  txHash: Hash.t,
  rawLog: string,
  success: bool,
};

type response_t =
  | Tx(tx_response_t)
  | Unknown;

[@bs.module "@cosmostation/cosmosjs"] external network: (string, string) => t = "network";

[@bs.send] external setPath: (t, string) => unit = "setPath";

[@bs.send] external setBech32MainPrefix: (t, string) => unit = "setBech32MainPrefix";

[@bs.send] external getAddress: (t, string) => string = "getAddress";

[@bs.send] external getECPairPriv: (t, string) => JsBuffer.t = "getECPairPriv";

[@bs.send] external _getAccounts: (t, string) => Js.Promise.t(Js.Json.t) = "getAccounts";

[@bs.send] external _broadcast: (t, TxCreator.t) => Js.Promise.t(Js.Json.t) = "broadcast";

let broadcast = (instance, signedMsg) => {
  let%Promise rawResponse = instance->_broadcast(signedMsg);

  Promise.ret(
    Tx(
      JsonUtils.Decode.{
        txHash: rawResponse |> at(["txhash"], string) |> Hash.fromHex,
        rawLog: rawResponse |> at(["raw_log"], string),
        success: rawResponse |> optional(field("logs", _ => ())) |> Belt_Option.isSome,
      },
    ),
  );
};
