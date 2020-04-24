type t =
  | Mnemonic(Mnemonic.t)
  | Ledger(Ledger.t);

let createFromMnemonic = mnemonic => {
  Mnemonic(Mnemonic.create(mnemonic));
};

let createFromLedger = () => {
  let%Promise ledger = Ledger.create();
  Promise.ret(Ledger(ledger));
};

let getAddressAndPubKey =
  fun
  | Mnemonic(x) => x |> Mnemonic.getAddressAndPubKey |> Promise.ret
  | Ledger(x) => x |> Ledger.getAddressAndPubKey;

// TODO: (string) => JsBuffer.t
// let sign = msg =>
//   fun
//   | Mnemonic(x) => JsBuffer.from([||]) |> Promise.ret
//   | Ledger => msg |> JsBuffer.fromHex |> Promise.ret;
