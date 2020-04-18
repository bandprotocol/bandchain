type t = {
  bandChain: BandWeb3.t,
  mnemonic: string,
};

let getAddress = x => {
  x.bandChain |> BandWeb3.getAddress(_, x.mnemonic) |> Address.fromBech32;
} /* }*/;

// TODO: sign message
// let sign = () => {
//  ""
