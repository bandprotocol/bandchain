type t = {
  address: Address.t,
  privKey: JsBuffer.t,
};

type a =
  | Connect(string)
  | Disconnect;

let lcdURL = "http://localhost:8010";
let chainID = "bandchain";
let path = "m/44'/494'/0'/0/0";
let bech32Prefix = "band";

let reducer = _ =>
  fun
  | Connect(mnemonic) => {
      let bandchain = BandWeb3.network(lcdURL, chainID);
      bandchain->BandWeb3.setPath(path);
      bandchain->BandWeb3.setBech32MainPrefix(bech32Prefix);

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
