type t = {
  address: Address.t,
  pubKey: PubKey.t,
  wallet: Wallet.t,
};

type a =
  | Connect(Wallet.t, Address.t, PubKey.t)
  | Disconnect
  | SendRequest(ID.OracleScript.t, JsBuffer.t, Js.Promise.t(TxCreator.response_t) => unit);

let reducer = state =>
  fun
  | Connect(wallet, address, pubKey) => Some({wallet, pubKey, address})
  | Disconnect => None
  | SendRequest(oracleScriptID, calldata, callback) =>
    switch (state) {
    | Some({address, wallet, pubKey}) =>
      callback(
        {
          let%Promise rawTx =
            TxCreator.createRawTx(
              address,
              [|Request(oracleScriptID, calldata, "4", "4", address, "")|],
            );
          let%Promise signature = Wallet.sign(TxCreator.sortAndStringify(rawTx), wallet);
          let signedTx =
            TxCreator.createSignedTx(
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
    {
      "value": (state->Belt.Option.map(({address}) => address), dispatch),
      "children": children,
    },
  );
};
