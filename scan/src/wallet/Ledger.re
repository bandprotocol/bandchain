type t = {
  app: LedgerJS.t,
  path: array(int),
  prefix: string,
};

let create = () => {
  // TODO: handle interaction timeout later
  let timeout = 10000;
  let path = [|44, 118, 0, 0, 0|];
  let prefix = "band";
  let%Promise transport = LedgerJS.createTransportWebUSB(timeout);

  let app = LedgerJS.createApp(transport);
  let%Promise pubKeyInfo = LedgerJS.publicKey(app, path);
  let%Promise appInfo = LedgerJS.appInfo(app);
  // TODO: check version
  // let%Promise version = LedgerJS.getVersion(app);

  // 36864(0x9000) will return if there is no error.
  // TODO: improve handle error
  if (pubKeyInfo.return_code != 36864) {
    if (pubKeyInfo.return_code == 28160) {
      Js.Console.log2("pubKeyInfo", pubKeyInfo);
      Js.Promise.reject(Not_found);
    } else if (appInfo.appName != "Cosmos") {
      Js.Console.log2("appInfo", appInfo);
      Js.Promise.reject(Not_found);
    } else {
      Js.Console.log(pubKeyInfo.error_message);
      Js.Promise.reject(Not_found);
    };
  } else {
    Promise.ret({app, path, prefix});
  };
};

let getApp = () => {
  // TODO: handle interaction timeout later
  let timeout = 10000;
  let%Promise transport = LedgerJS.createTransportWebUSB(timeout);
  Promise.ret(LedgerJS.createApp(transport));
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
