type t;

[@bs.val] external from: array(int) => t = "Buffer.from";

[@bs.val] external _from: (string, string) => t = "Buffer.from";
let fromHex = hexstr => hexstr->HexUtils.normalizeHexString->_from("hex");

let fromBase64 = hexstr => _from(hexstr, "base64");

[@bs.send] external _toString: (t, string) => string = "toString";
[@bs.send] external toString: t => string = "toString";

let toHex = (~with0x=false, buf) => (with0x ? "0x" : "") ++ buf->_toString("hex");

let toBase64 = buf => buf->_toString("base64");

[@bs.val] external toArray: t => array(int) = "Array.from";

let base64ToHex = base64str => base64str->fromBase64->toHex;
let hexToBase64 = hexstr => hexstr->fromHex->toBase64;
let arrayToHex = arr => arr->from->toHex;
let hexToArray = hexstr => hexstr->fromHex->toArray;
let arrayToBase64 = arr => arr->from->toBase64;
let base64ToArray = base64str => base64str->fromBase64->toArray;
