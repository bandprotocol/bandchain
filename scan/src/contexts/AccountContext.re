type t = {
  address: Address.t,
  privKey: JsBuffer.t,
};

type a =
  | Connect(string)
  | Disconnect;

let bandchain = BandWeb3.network("http://localhost:8010", "bandchain");
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

  React.createElement(
    React.Context.provider(context),
    {
      "value": (state->Belt.Option.map(({address}) => address), dispatch),
      "children": children,
    },
  );
};
