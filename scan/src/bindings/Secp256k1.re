[@bs.module "secp256k1"]
external publicKeyCreate: (JsBuffer.t, bool) => JsBuffer.t = "publicKeyCreate";

[@bs.module "secp256k1"] external signatureImport: array(int) => array(int) = "signatureImport";
