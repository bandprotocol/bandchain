type t;

type wrapped_msg_t;

type signed_msg_t;

type account_result_t = {
  accountNumber: string,
  sequence: string,
};

type response_t = {txHash: Hash.t};

[@bs.module "@cosmostation/cosmosjs"] external network: (string, string) => t = "network";

[@bs.send] external setPath: (t, string) => unit = "setPath";

[@bs.send] external setBech32MainPrefix: (t, string) => unit = "setBech32MainPrefix";

[@bs.send] external getAddress: (t, string) => string = "getAddress";

[@bs.send] external getECPairPriv: (t, string) => JsBuffer.t = "getECPairPriv";

[@bs.send] external _getAccounts: (t, string) => Js.Promise.t(Js.Json.t) = "getAccounts";

[@bs.send] external newStdMsgSend: (t, StdMsgSend.t) => wrapped_msg_t = "newStdMsg";

[@bs.send] external newStdMsgRequest: (t, StdMsgRequest.t) => wrapped_msg_t = "newStdMsg";

[@bs.send] external sign: (t, wrapped_msg_t, JsBuffer.t) => signed_msg_t = "sign";

[@bs.send] external _broadcast: (t, signed_msg_t) => Js.Promise.t(Js.Json.t) = "broadcast";

let getAccounts = (instance, address) => {
  let%Promise rawResult = instance->_getAccounts(address);

  Promise.ret(
    JsonUtils.Decode.{
      accountNumber: rawResult |> at(["result", "value", "account_number"], string),
      sequence: rawResult |> at(["result", "value", "sequence"], string),
    },
  );
};

let broadcast = (instance, signedMsg) => {
  let%Promise rawResponse = instance->_broadcast(signedMsg);

  Promise.ret(
    JsonUtils.Decode.{txHash: rawResponse |> at(["txhash"], string) |> Hash.fromHex},
  );
};

// let cosmos = network("http://localhost:8010", "bandchain");

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

//   // let sendMsg = {
//   //   StdMsgSend.msgs: [|
//   //     {
//   //       type_: "cosmos-sdk/MsgSend",
//   //       value: {
//   //         amount: [|{amount: "777", denom: "uband"}|],
//   //         from_address: address,
//   //         to_address: "band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun",
//   //       },
//   //     },
//   //   |],
//   //   chain_id: "bandchain",
//   //   fee: {
//   //     amount: [|{amount: "5000", denom: "uband"}|],
//   //     gas: "200000",
//   //   },
//   //   memo: "",
//   //   account_number: data.accountNumber,
//   //   sequence: data.sequence,
//   // };
//   let msgRequest = {
//     StdMsgRequest.msgs: [|
//       {
//         type_: "zoracle/Request",
//         value: {
//           oracleScriptID: "1",
//           calldata: "RVRI",
//           requestedValidatorCount: "1",
//           sufficientValidatorCount: "1",
//           expiration: "20",
//           prepareGas: "20000",
//           executeGas: "150000",
//           sender: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs",
//         },
//       },
//     |],
//     chain_id: "bandchain",
//     fee: {
//       amount: [|{amount: "5000", denom: "uband"}|],
//       gas: "3000000",
//     },
//     memo: "",
//     account_number: data.accountNumber,
//     sequence: data.sequence,
//   };

//   let wrappedMsg = cosmos->newStdMsgRequest(msgRequest);
//   let signedMsg = cosmos->sign(wrappedMsg, ecPairPriv);
//   let%Promise res = cosmos->broadcast(signedMsg);
//   Js.Console.log(res);

//   Promise.ret(1);
// };
