type t = {
  address: Address.t,
  pubKey: PubKey.t,
  wallet: Wallet.t,
  chainID: string,
};

type send_request_t = {
  oracleScriptID: ID.OracleScript.t,
  calldata: JsBuffer.t,
  callback: Js.Promise.t(TxCreator.response_t) => unit,
  askCount: string,
  minCount: string,
  clientID: string,
  feeLimit: string,
  prepareGas: string,
  executeGas: string,
};

type a =
  | Connect(Wallet.t, Address.t, PubKey.t, string)
  | Disconnect
  | SendRequest(send_request_t);

let reducer = state =>
  fun
  | Connect(wallet, address, pubKey, chainID) => Some({wallet, pubKey, address, chainID})
  | Disconnect => {
      switch (state) {
      | Some({wallet}) => wallet |> Wallet.disconnect
      | None => ()
      };
      None;
    }
  | SendRequest({
      oracleScriptID,
      calldata,
      callback,
      askCount,
      minCount,
      clientID,
      feeLimit,
      prepareGas,
      executeGas,
    }) =>
    switch (state) {
    | Some({address, wallet, pubKey, chainID}) =>
      callback(
        {
          let%Promise rawTx =
            TxCreator.createRawTx(
              ~address,
              ~msgs=[|
                Request(
                  oracleScriptID,
                  calldata,
                  askCount,
                  minCount,
                  address,
                  clientID,
                  {amount: feeLimit, denom: "uband"},
                  prepareGas,
                  executeGas,
                ),
              |],
              ~chainID,
              ~gas="2000000",
              ~feeAmount="0",
              ~memo="send via scan",
              (),
            );
          let%Promise signature = Wallet.sign(TxCreator.sortAndStringify(rawTx), wallet);
          let signedTx =
            TxCreator.createSignedTx(
              ~network=Env.network,
              ~signature=signature |> JsBuffer.toBase64,
              ~pubKey,
              ~tx=rawTx,
              ~mode="block",
              (),
            );
          TxCreator.broadcast(signedTx);
        },
      );

      state;
    | None =>
      callback(Promise.ret(TxCreator.Unknown));
      state;
    };

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let (state, dispatch) = React.useReducer(reducer, None);

  React.createElement(
    React.Context.provider(context),
    {"value": (state, dispatch), "children": children},
  );
};
