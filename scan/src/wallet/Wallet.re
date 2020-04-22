type t =
  | Mnemonic(Mnemonic.t)
  | Ledger(Ledger.t);

let createFromMnemonic = mnemonic => {
  // Just use arbitrary rpcUrl, chainID beacuase they didn't use in CosmosJS
  let bandChain = CosmosJS.network("rpcUrl", "chainID");
  bandChain->CosmosJS.setPath("m/44'/494'/0'/0/0");
  bandChain->CosmosJS.setBech32MainPrefix("band");
  let privKey = bandChain |> CosmosJS.getECPairPriv(_, mnemonic);
  Mnemonic({bandChain, mnemonic, privKey});
};

let createFromLedger = () => {
  // TODO: handle interaction timeout later
  let timeout = 10000;
  let path = [|44, 118, 0, 0, 0|];
  let prefix = "band";
  let%Promise transport = LedgerJS.createTransportWebUSB(timeout);

  let app = LedgerJS.createApp(transport);
  let%Promise pubKeyInfo = LedgerJS.publicKey(app, path);
  let%Promise appInfo = LedgerJS.appInfo(app);
  let%Promise version = LedgerJS.getVersion(app);

  // TODO: check version
  Js.Console.log(version);
  if (pubKeyInfo.return_code == 28160) {
    Js.Console.log2("pubKeyInfo", pubKeyInfo);
    Js.Promise.reject(Not_found);
  } else if (appInfo.appName != "Cosmos") {
    Js.Console.log2("appInfo", appInfo);
    Js.Promise.reject(Not_found);
  } else {
    Promise.ret(Ledger({app, path, prefix}));
  };
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
