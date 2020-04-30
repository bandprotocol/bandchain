[@bs.module "secp256k1"]
external publicKeyCreate: (JsBuffer.t, bool) => array(int) = "publicKeyCreate";

[@bs.module "secp256k1"] external signatureImport: array(int) => array(int) = "signatureImport";

type sign_info_t = {signature: array(int)};

[@bs.module "secp256k1"] external ecdsaSign: (JsBuffer.t, JsBuffer.t) => sign_info_t = "ecdsaSign";
