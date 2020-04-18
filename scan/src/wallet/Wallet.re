type t =
  | Mnemonic(Mnemonic.t)
  | Ledger(Ledger.t);

let createFromMnemonic = (rpcUrl, mnemonic) => {
  let bandChain = BandWeb3.network(rpcUrl, "bandchain");
  bandChain->BandWeb3.setPath("m/44'/494'/0'/0/0");
  bandChain->BandWeb3.setBech32MainPrefix("band");
  Mnemonic({bandChain, mnemonic});
};

let createFromLedger = path => {
  Promise.ret(Ledger(path));
};

let getAddress =
  fun
  | Mnemonic(x) => x |> Mnemonic.getAddress |> Promise.ret
  | Ledger(x) => x |> Ledger.getAddress;

// let sign = msg =>
//   fun
//   | Mnemonic(x) => JsBuffer.from([||]) |> Promise.ret
//   | Ledger => msg |> JsBuffer.fromHex |> Promise.ret;
