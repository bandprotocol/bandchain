type t = {
  address: Address.t,
  privKey: JsBuffer.t,
};

type context_value_t = {
  address: option(Address.t),
  sendRequest:
    React.callback(
      BandScan.ID.OracleScript.t,
      BandScan.JsBuffer.t => Js.Promise.t(BandScan.BandWeb3.response_t),
    ),
};

type a =
  | Connect(string)
  | Disconnect;

let bandchain = BandWeb3.network("https://d3n.bandprotocol.com/rest", "bandchain");
bandchain->BandWeb3.setPath("m/44'/494'/0'/0/0");
bandchain->BandWeb3.setBech32MainPrefix("band");

let reducer = _ =>
  fun
  | Connect(mnemonic) => {
      let newAddress = bandchain |> BandWeb3.getAddress(_, mnemonic) |> Address.fromBech32;
      let newPrivKey = bandchain |> BandWeb3.getECPairPriv(_, mnemonic);
      Some({address: newAddress, privKey: newPrivKey});
    }
  | Disconnect => None;

let context = React.createContext(ContextHelper.default);

[@react.component]
let make = (~children) => {
  let (state, dispatch) = React.useReducer(reducer, None);

  let sendRequest =
    React.useCallback1(
      (oracleScriptID, calldata) =>
        switch (state) {
        | Some({address, privKey}) =>
          let%Promise data = bandchain->BandWeb3.getAccounts(address |> Address.toBech32);
          let msgRequest =
            StdMsgRequest.create(
              oracleScriptID,
              ~calldata,
              ~requestedValidatorCount=4,
              ~sufficientValidatorCount=3,
              ~expiration=20,
              ~prepareGas=20000,
              ~executeGas=150000,
              ~sender=address,
              ~feeAmount=1000000,
              ~gas=3000000,
              ~accountNumber=data.accountNumber,
              ~sequence=data.sequence,
            );

          let wrappedMsg = bandchain->BandWeb3.newStdMsgRequest(msgRequest);
          let signedMsg = bandchain->BandWeb3.sign(wrappedMsg, privKey, "block");
          let%Promise res = bandchain->BandWeb3.broadcast(signedMsg);

          Promise.ret(res);
        | None => Promise.ret(BandWeb3.Unknown)
        },
      [|state->Belt_Option.mapWithDefault("", ({privKey}) => privKey |> JsBuffer.toHex)|],
    );

  let contextValue: context_value_t = {
    address: state->Belt.Option.map(({address}) => address),
    sendRequest,
  };

  React.createElement(
    React.Context.provider(context),
    {"value": (contextValue, dispatch), "children": children},
  );
};
