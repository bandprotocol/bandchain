type t = {
  app: LedgerJS.t,
  path: array(int),
  prefix: string,
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

// TODO:
// let sign = message => {
//   let%Promise app = getApp();

//   //  TODO: remove hard-coded path later
//   let path = [|44, 118, 0, 0, 0|];
//   let responsePromise = LedgerJS.sign(app, path, message);
//   let%Promise response = responsePromise;
//   response.signature |> LedgerJS.signatureImport |> JsBuffer.from |> Promise.ret;
// };
