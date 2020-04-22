type t;

type wrapped_msg_t;

// type signed_msg_t;

// Start of TODO: abstract out later (hard-coded for MsgRequest for now)
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

type tx_t = {
  fee: StdMsgRequest.fee_t,
  memo: string,
  msg: array(StdMsgRequest.msg_t),
  signatures: array(signature_t),
};

type signed_msg_t = {
  mode: string,
  tx: tx_t,
};
// End of TODO

type account_result_t = {
  accountNumber: int,
  sequence: int,
};

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

[@bs.send] external newStdMsgSend: (t, StdMsgSend.t) => wrapped_msg_t = "newStdMsg";

[@bs.send] external newStdMsgRequest: (t, StdMsgRequest.t) => wrapped_msg_t = "newStdMsg";

[@bs.send] external sign: (t, wrapped_msg_t, JsBuffer.t, string) => signed_msg_t = "sign";

[@bs.send] external _broadcast: (t, signed_msg_t) => Js.Promise.t(Js.Json.t) = "broadcast";

let getAccounts = (instance, address) => {
  let%Promise rawResult = instance->_getAccounts(address);

  Promise.ret(
    JsonUtils.Decode.{
      accountNumber: rawResult |> at(["result", "value", "account_number"], int),
      sequence: rawResult |> at(["result", "value", "sequence"], int),
    },
  );
};

let addPublicKeyToSignedMsg: signed_msg_t => unit = [%bs.raw
  {|
function(signedMsg) {
  for (const sig of signedMsg.tx.signatures) {
    // sha256("tendermint/PubKeySecp256k1") = f8ccea**eb5ae987**ea423e6cc0e94297a53bd6862df3b3a02a6f6fc89250308760
    // We use eb5ae987 AND 21 (= 21 base 16 = 33 = pubkey size)
    sig.public_key = Buffer.from(
      'eb5ae98721' + Buffer.from(sig.pub_key.value, 'base64').toString('hex'), 'hex'
    ).toString('base64')
  }
}
  |}
];

let broadcast = (instance, signedMsg) => {
  addPublicKeyToSignedMsg(signedMsg);
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

// TODO: abstract out later
let createSignedMsgRequest = (msgRequest: StdMsgRequest.t, signature, pubKey, mode) => {
  {
    mode,
    tx: {
      fee: msgRequest.fee,
      memo: msgRequest.memo,
      msg: msgRequest.msgs,
      signatures: [|
        {
          pub_key: {
            type_: "tendermint/PubKeySecp256k1",
            value: pubKey |> JsBuffer.fromHex |> JsBuffer.toBase64,
          },
          public_key: "eb5ae98721" ++ pubKey |> JsBuffer.fromHex |> JsBuffer.toBase64,
          signature,
        },
      |],
    },
  };
};

// //EXAMPLE
// let cosmos = network("https://d3n.bandprotocol.com/rest", "bandchain");

// cosmos->setPath("m/44'/494'/0'/0/0");
// cosmos->setBech32MainPrefix("band");

// let test = _ => Js.Console.log("TEST");

// let address =
//   cosmos->getAddress(
//     "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic",
//   );

// let ecPairPriv =
//   cosmos->getECPairPriv(
//     "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic",
//   );

// {
//   let%Promise data = cosmos->getAccounts(address);
//   Js.Console.log2("yo22222", data);

//   //msgSend
//   let msgSend = StdMsgSend.create(
//     ~fromAddress="band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
//     ~toAddress="band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
//     ~sendAmount=100,
//     ~feeAmount=5000,
//     ~accountNumber=data.accountNumber,
//     ~sequence=data.sequence,
//   );
//   //msgRequest
//   let msgRequest =
//      ~oracleScriptID=1,
//      ~calldata="RVRI" |> JsBuffer.fromBase64,
//      ~requestedValidatorCount=4,
//      ~sufficientValidatorCount=4,
//      ~sender"band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
//      ~feeAmount=1000000,
//      ~gas=300000,
//      ~accountNumber=data.accountNumber,
//      ~sequence=data.sequence,
//     );

//   let wrappedMsg = cosmos->newStdMsgRequest(msgRequest);
//   let signedMsg = cosmos->sign(wrappedMsg, ecPairPriv);
//   let%Promise res = cosmos->broadcast(signedMsg);
//   Js.Console.log(res);

//   Promise.ret(1);
// };
