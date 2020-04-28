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

let sign = msg =>
  fun
  | Mnemonic(x) => x |> Mnemonic.sign(_, msg) |> Promise.ret
  | Ledger(x) => x |> Ledger.sign(_, msg);
