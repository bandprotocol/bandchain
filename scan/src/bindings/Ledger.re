type t;

type transport_t;

type response_t = {
  bech32_address: string,
  return_code: string,
  error_message: string,
  compressed_pk: array(int),
};

type sign_response_t = {
  return_code: string,
  error_message: string,
  signature: array(int),
};

type addr_pubkey_t = {
  address: Address.t,
  pubKey: string,
};

[@bs.module "@ledgerhq/hw-transport-webusb"] [@bs.scope "default"] [@bs.val]
external createTransportU2F: unit => Js.Promise.t(transport_t) = "create";

[@bs.module "ledger-cosmos-js"] [@bs.new] external createApp: transport_t => t = "default";
[@bs.send]
external _getAddressAndPubKey: (t, array(int), string) => Js.Promise.t(response_t) =
  "getAddressAndPubKey";
[@bs.send] external _sign: (t, array(int), string) => Js.Promise.t(sign_response_t) = "sign";

[@bs.module "secp256k1"] external signatureImport: array(int) => array(int) = "signatureImport";

let getApp = () => {
  let%Promise transport = createTransportU2F();
  Promise.ret(createApp(transport));
};

let getAddressAndPubKey = () => {
  let%Promise app = getApp();

  //  TODO: remove hard-coded path later
  let path = [|44, 118, 0, 0, 0|];
  let prefix = "band";
  let responsePromise = _getAddressAndPubKey(app, path, prefix);
  let%Promise response = responsePromise;

  Promise.ret({
    address: response.bech32_address |> Address.fromBech32,
    pubKey: response.compressed_pk |> JsBuffer.from |> JsBuffer.toHex,
  });
};

let sign = message => {
  let%Promise app = getApp();

  //  TODO: remove hard-coded path later
  let path = [|44, 118, 0, 0, 0|];
  let responsePromise = _sign(app, path, message);
  let%Promise response = responsePromise;
  response.signature |> signatureImport |> JsBuffer.from |> Promise.ret;
};
