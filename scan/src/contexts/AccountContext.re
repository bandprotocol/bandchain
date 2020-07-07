type t = {
  address: Address.t,
  pubKey: PubKey.t,
  wallet: Wallet.t,
  chainID: string,
};

type a =
  | Connect(Wallet.t, Address.t, PubKey.t, string)
  | Disconnect
  | SendRequest(ID.OracleScript.t, JsBuffer.t, Js.Promise.t(TxCreator.response_t) => unit);

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
  | SendRequest(oracleScriptID, calldata, callback) =>
    switch (state) {
    | Some({address, wallet, pubKey, chainID}) =>
      callback(
        {
          let%Promise rawTx =
            TxCreator.createRawTx(
              ~address,
              // Client id can't be an empty string (""), so we need to add "from_scan"
              // TODO: Make this more intuitive
              ~msgs=[|Request(oracleScriptID, calldata, "4", "4", address, "from_scan")|],
              ~chainID,
              ~gas="700000",
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
