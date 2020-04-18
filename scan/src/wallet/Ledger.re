type t = {path: array(int)};

let getApp = () => {
  let%Promise transport = LedgerJS.createTransportU2F();
  Promise.ret(LedgerJS.createApp(transport));
};

let getAddress = x => {
  let%Promise app = getApp();

  //  TODO: remove hard-coded path later
  // let path = [|44, 118, 0, 0, 0|];

  let prefix = "band";
  let responsePromise = LedgerJS.getAddressAndPubKey(app, x.path, prefix);
  let%Promise response = responsePromise;

  Promise.ret(response.bech32_address |> Address.fromBech32);
};

let getAddressAndPubKey = x => {
  let%Promise app = getApp();

  //  TODO: remove hard-coded path later
  // let path = [|44, 118, 0, 0, 0|];
  let prefix = "band";
  let responsePromise = LedgerJS.getAddressAndPubKey(app, x.path, prefix);
  let%Promise response = responsePromise;

  Promise.ret(
    LedgerJS.{
      address: response.bech32_address |> Address.fromBech32,
      pubKey: response.compressed_pk |> JsBuffer.from |> JsBuffer.toHex,
    },
  );
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
