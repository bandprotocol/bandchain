type t =
  | Hash(string); // string is hex (without 0x)

let fromHex = hexstr =>
  if (hexstr->String.sub(0, 2) == "0x") {
    Hash(hexstr->String.lowercase_ascii->String.sub(2, hexstr->String.length - 2));
  } else {
    Hash(hexstr->String.lowercase_ascii);
  };

let fromBase64 = base64str => base64str->JsBuffer.base64ToHex->fromHex;

let toBase64 =
  fun
  | Hash(hash) => hash->JsBuffer.hexToBase64;

let toHex =
  fun
  | Hash(hash) => hash;
