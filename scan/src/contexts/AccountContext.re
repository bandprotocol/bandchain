type t = {
  address: Address.t,
  pubKey: PubKey.t,
  wallet: Wallet.t,
};

type a =
  | Connect(Wallet.t, Address.t, PubKey.t)
  | Disconnect
  | SendRequest(ID.OracleScript.t, JsBuffer.t, Js.Promise.t(BandWeb3.response_t) => unit)
  | SendRequestWithLedger;

let bandchain = BandWeb3.network(Env.rpc, "bandchain");
bandchain->BandWeb3.setPath("m/44'/494'/0'/0/0");
bandchain->BandWeb3.setBech32MainPrefix("band");

let reducer = state =>
  fun
  | Connect(wallet, address, pubKey) => Some({wallet, pubKey, address})
  | Disconnect => None
  | SendRequest(oracleScriptID, calldata, callback) =>
    switch (state) {
    | Some({address, wallet}) =>
      switch (wallet) {
      | Mnemonic({privKey}) =>
        callback(
          {
            let%Promise {accountNumber, sequence} =
              bandchain->BandWeb3.getAccounts(address |> Address.toBech32);
            let msgRequest =
              StdMsgRequest.create(
                oracleScriptID,
                ~calldata,
                ~requestedValidatorCount=4,
                ~sufficientValidatorCount=4,
                ~sender=address,
                ~feeAmount=1000000,
                ~gas=3000000,
                ~accountNumber=accountNumber |> string_of_int,
                ~sequence=sequence |> string_of_int,
              );
            let wrappedMsg = bandchain->BandWeb3.newStdMsgRequest(msgRequest);
            let signedMsg = bandchain->BandWeb3.sign(wrappedMsg, privKey, "block");
            let%Promise res = bandchain->BandWeb3.broadcast(signedMsg);

            Promise.ret(res);
          },
        )
      | Ledger(_) => Js.Console.log("Send request via ledger on WIP")
      };

      state;
    | None =>
      callback(Promise.ret(BandWeb3.Unknown));
      state;
    }
  | SendRequestWithLedger =>
    switch (state) {
    | Some(_) =>
      // let _ = {
      //   // TODO: 1. save address to state
      //   //       2. handle error when ledger doesn't connect.
      //   let%Promise {address, pubKey} = LedgerJS.getAddressAndPubKey();
      //   let%Promise {accountNumber, sequence} =
      //     bandchain->BandWeb3.getAccounts(address |> Address.toBech32);
      //   let msgRequest =
      //     StdMsgRequest.create(
      //       ID.OracleScript.ID(2),
      //       ~calldata=JsBuffer.fromBase64("AwAAAEJUQ2QAAAAAAAAA"),
      //       ~requestedValidatorCount=4,
      //       ~sufficientValidatorCount=4,
      //       ~sender=address,
      //       ~feeAmount=1000000,
      //       ~gas=3000000,
      //       ~accountNumber=accountNumber |> string_of_int,
      //       ~sequence=sequence |> string_of_int,
      //     );
      //   let stringifiedMsg = msgRequest |> StdMsgRequest.sortAndStringify;
      //   let%Promise signature = LedgerJS.sign(stringifiedMsg);

      //   let signBase64 = signature |> JsBuffer.toBase64;

      //   let signedMsg = BandWeb3.createSignedMsgRequest(msgRequest, signBase64, pubKey, "block");
      //   Js.Console.log2("signedMsg", signedMsg);
      //   let%Promise res = bandchain->BandWeb3.broadcast(signedMsg);
      //   Js.Console.log(res);

      //   Promise.ret();
      // };
      state
    | None => state
    };

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let (state, dispatch) = React.useReducer(reducer, None);

  React.createElement(
    React.Context.provider(context),
    {
      "value": (state->Belt.Option.map(({address}) => address), dispatch),
      "children": children,
    },
  );
};
