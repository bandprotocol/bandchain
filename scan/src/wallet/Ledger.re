type t = {
  transport: LedgerJS.transport_t,
  app: LedgerJS.t,
  path: array(int),
  prefix: string,
};

type ledger_app_t =
  | Cosmos
  | BandChain;

let getPath = (ledgerApp, accountIndex) => {
  switch (ledgerApp) {
  | Cosmos => [|44, 118, 0, 0, accountIndex|]
  | BandChain => [|44, 494, 0, 0, accountIndex|]
  };
};

let getAppName =
  fun
  | Cosmos => "Cosmos"
  | BandChain => "BandChain";

// TODO: hard-coded minimum version
let getRequiredVersion =
  fun
  | Cosmos => "1.5.0"
  | BandChain => "2.12.0";

let create = (ledgerApp, accountIndex) => {
  // TODO: handle interaction timeout later
  let timeout = 10000;
  let path = getPath(ledgerApp, accountIndex);
  let prefix = "band";

  let%Promise transport =
    Os.isWindow()
      ? LedgerJS.createTransportWebHID(timeout) : LedgerJS.createTransportWebUSB(timeout);

  let app = LedgerJS.createApp(transport);
  let%Promise pubKeyInfo = LedgerJS.publicKey(app, path);
  let%Promise appInfo = LedgerJS.appInfo(app);
  let%Promise version = LedgerJS.getVersion(app);

  let LedgerJS.{major, minor, patch, test_mode, device_locked} = version;
  let userVersion = {j|$major.$minor.$patch|j};
  let requiredAppName = getAppName(ledgerApp);
  let requiredVersion = getRequiredVersion(ledgerApp);

  // 36864(0x9000) will return if there is no error.
  // TODO: improve handle error
  // Validatate step
  // 1. Check return code of pubKeyInfo
  // 2. If pass, then check app version
  // 3. If pass, then check test_mode
  if (pubKeyInfo.return_code != 36864) {
    if (appInfo.appName != requiredAppName) {
      let appName = appInfo.appName;
      Js.Console.log({j|App name is not $requiredAppName. (Current is $appName)|j});
      Js.Promise.reject(Not_found);
    } else if (device_locked) {
      Js.Console.log3("Device is locked", pubKeyInfo, version);
      Js.Promise.reject(Not_found);
    } else {
      Js.Console.log(pubKeyInfo.error_message);
      Js.Promise.reject(Not_found);
    };
  } else if (!Semver.gte(userVersion, requiredVersion)) {
    Js.Console.log({j|Cosmos app version must >= $requiredVersion (Current is $userVersion)|j});
    Js.Promise.reject(Not_found);
  } else if (test_mode) {
    Js.Console.log3("test mode is not supported", pubKeyInfo, version);
    Js.Promise.reject(Not_found);
  } else {
    Promise.ret({transport, app, path, prefix});
  };
};

let getAddressAndPubKey = x => {
  let prefix = "band";
  let responsePromise = LedgerJS.getAddressAndPubKey(x.app, x.path, prefix);
  let%Promise response = responsePromise;

  if (response.return_code != 36864) {
    Js.Console.log(response.error_message);
    Js.Promise.reject(Not_found);
  } else {
    Promise.ret((
      response.bech32_address |> Address.fromBech32,
      response.compressed_pk |> JsBuffer.from |> JsBuffer.toHex |> PubKey.fromHex,
    ));
  };
};

let sign = (x, message) => {
  let responsePromise = LedgerJS.sign(x.app, x.path, message);
  let%Promise response = responsePromise;
  response.signature |> Secp256k1.signatureImport |> JsBuffer.from |> Promise.ret;
};
