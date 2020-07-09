type t;

type transport_t;

type addr_pukey_response_t = {
  bech32_address: string,
  return_code: int,
  error_message: string,
  compressed_pk: array(int),
};

type pubkey_response_t = {
  return_code: int,
  error_message: string,
  compressed_pk: PubKey.t,
};

type version_t = {
  return_code: int,
  error_message: string,
  test_mode: bool,
  major: int,
  minor: int,
  patch: int,
  device_locked: bool,
};

type app_info_t = {
  return_code: int,
  error_message: string,
  appName: string,
  appVersion: string,
};

type sign_response_t = {
  return_code: string,
  error_message: string,
  signature: array(int),
};

[@bs.module "@ledgerhq/hw-transport-webhid"] [@bs.scope "default"] [@bs.val]
external createTransportWebHID: int => Js.Promise.t(transport_t) = "create";

[@bs.module "@ledgerhq/hw-transport-webusb"] [@bs.scope "default"] [@bs.val]
external createTransportWebUSB: int => Js.Promise.t(transport_t) = "create";

[@bs.module "ledger-cosmos-js"] [@bs.new] external createApp: transport_t => t = "default";
[@bs.send]
external getAddressAndPubKey: (t, array(int), string) => Js.Promise.t(addr_pukey_response_t) =
  "getAddressAndPubKey";
[@bs.send] external publicKey: (t, array(int)) => Js.Promise.t(pubkey_response_t) = "publicKey";
[@bs.send] external sign: (t, array(int), string) => Js.Promise.t(sign_response_t) = "sign";
[@bs.send] external getVersion: t => Js.Promise.t(version_t) = "getVersion";
[@bs.send] external appInfo: t => Js.Promise.t(app_info_t) = "appInfo";
// TODO: It should return promise
[@bs.send] external close: transport_t => unit = "close";
