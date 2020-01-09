[@bs.module "js-sha256"] [@bs.scope "sha256"] [@bs.val]
external _digest: JsBuffer.t => array(int) = "array";
let digest = arr => arr->JsBuffer.from->_digest;
let hexDigest = arr => arr->digest->JsBuffer.arrayToHex;
