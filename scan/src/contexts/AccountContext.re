type t = {
  address: Address.t,
  privKey: JsBuffer.t,
};

type a =
  | Connect(string)
  | Disconnect
  | SendRequest(ID.OracleScript.t, JsBuffer.t, Js.Promise.t(BandWeb3.response_t) => unit);

let bandchain = BandWeb3.network(Env.rpc, "bandchain");
bandchain->BandWeb3.setPath("m/44'/494'/0'/0/0");
bandchain->BandWeb3.setBech32MainPrefix("band");

let reducer = state =>
  fun
  | Connect(mnemonic) => {
      let newAddress = bandchain |> BandWeb3.getAddress(_, mnemonic) |> Address.fromBech32;
      let newPrivKey = bandchain |> BandWeb3.getECPairPriv(_, mnemonic);
      Some({address: newAddress, privKey: newPrivKey});
    }
  | Disconnect => None
  | SendRequest(oracleScriptID, calldata, callback) =>
    switch (state) {
    | Some({address, privKey}) =>
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
              ~expiration=20,
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
      );
      state;
    | None =>
      callback(Promise.ret(BandWeb3.Unknown));
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
