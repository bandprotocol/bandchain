type t =
  | Mnemonic(Mnemonic.t)
  | Ledger(Ledger.t);

let createFromMnemonic = mnemonic => {
  // Just use arbitrary rpcUrl, chainID beacuase they didn't use in BandWeb3
  let bandChain = BandWeb3.network("rpcUrl", "chainID");
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
