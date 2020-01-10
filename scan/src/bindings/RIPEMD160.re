type t;

[@bs.module] [@bs.new] external _create: unit => t = "ripemd160";
[@bs.send] external _update: (t, JsBuffer.t) => t = "update";
[@bs.send] external _digest: (t, string) => string = "digest";
let hexDigest = data => _create()->_update(data->JsBuffer.from)->_digest("hex");
let digest = data => data->hexDigest->JsBuffer.hexToArray;
