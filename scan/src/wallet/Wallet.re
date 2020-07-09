type t =
  | Mnemonic(Mnemonic.t)
  | Ledger(Ledger.t);

let createFromMnemonic = mnemonic => {
  Mnemonic(Mnemonic.create(mnemonic));
};

let createFromLedger = (ledgerApp, accountIndex) => {
  let%Promise ledger = Ledger.create(ledgerApp, accountIndex);
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

let disconnect =
  fun
  | Mnemonic(_) => ()
  | Ledger({transport}) => transport |> LedgerJS.close;
