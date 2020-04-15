type t;

type transport_t;

type response_t = {
  bech32_address: string,
  return_code: string,
  error_message: string,
  major: string,
  minor: string,
  patch: string,
  test_mode: string,
  device_locked: string,
};

type sign_response_t = {
  return_code: string,
  error_message: string,
  signature: JsBuffer.t,
};

[@bs.module "@ledgerhq/hw-transport-webusb"] [@bs.scope "default"] [@bs.val]
external createTransportU2F: unit => Js.Promise.t(transport_t) = "create";

[@bs.module "ledger-cosmos-js"] [@bs.new] external createApp: transport_t => t = "default";
[@bs.send]
external getAddressAndPubKey: (t, array(int), string) => Js.Promise.t(response_t) =
  "getAddressAndPubKey";
[@bs.send] external sign: (t, array(int), string) => Js.Promise.t(sign_response_t) = "sign";

let getApp = () => {
  let%Promise transport = createTransportU2F();
  Js.Promise.resolve(createApp(transport));
};

let getAddress = () => {
  let%Promise app = getApp();

  // TODO: remove hard-coded later
  let path = [|44, 118, 4, 0, 2|];
  let prefix = "cosmos";
  let responsePromise = getAddressAndPubKey(app, path, prefix);
  let%Promise response = responsePromise;
  Js.Promise.resolve(response.bech32_address);
};

let signExampleTx = () => {
  let%Promise app = getApp();

  //  TODO: remove hard-coded path later
  let path = [|44, 118, 0, 0, 0|];
  let message = {f|{"account_number":"6571","chain_id":"cosmoshub-2","fee":{"amount":[{"amount":"5000","denom":"uatom"}],"gas":"200000"},"memo":"Delegated with Ledger from union.market","msgs":[{"type":"cosmos-sdk/MsgDelegate","value":{"amount":{"amount":"1000000","denom":"uatom"},"delegator_address":"cosmos102hty0jv2s29lyc4u0tv97z9v298e24t3vwtpl","validator_address":"cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7"}}],"sequence":"0"}|f};
  let responsePromise = sign(app, path, message);
  let%Promise response = responsePromise;
  Js.Promise.resolve(response);
};
