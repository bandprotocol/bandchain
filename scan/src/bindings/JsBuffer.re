type t;

[@bs.val] external from: array(int) => t = "Buffer.from";

[@bs.val] external _from: (string, string) => t = "Buffer.from";
let fromHex = hexstr =>
  if (hexstr->String.sub(0, 2) == "0x") {
    hexstr->String.lowercase_ascii->String.sub(2, hexstr->String.length - 2)->_from("hex");
  } else {
    _from(hexstr->String.lowercase_ascii, "hex");
  };

let fromBase64 = hexstr => _from(hexstr, "base64");

[@bs.send] external _toString: (t, string) => string = "toString";
let toHex = buf => buf->_toString("hex");
let toBase64 = buf => buf->_toString("base64");

[@bs.val] external toArray: t => array(int) = "Array.from";

let base64ToHex = base64str => base64str->fromBase64->toHex;
let hexToBase64 = hexstr => hexstr->fromHex->toBase64;
let arrayToHex = arr => arr->from->toHex;
let hexToArray = hexstr => hexstr->fromHex->toArray;
let arrayToBase64 = arr => arr->from->toBase64;
let base64ToArray = base64str => base64str->fromBase64->toArray;
