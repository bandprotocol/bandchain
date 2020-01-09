type t =
  | PubKey(string); // string is hex (without 0x)

let fromHex = hexstr =>
  if (hexstr->String.sub(0, 2) == "0x") {
    PubKey(hexstr->String.lowercase_ascii->String.sub(2, hexstr->String.length - 2));
  } else {
    PubKey(hexstr->String.lowercase_ascii);
  };

let fromBech32 = bech32str =>
  PubKey(bech32str->Bech32.decode->Bech32.wordsGet->Bech32.fromWords->JsBuffer.arrayToHex);

let fromBase64 = base64str => base64str->JsBuffer.base64ToHex->fromHex;

let toAddress =
  fun
  | PubKey(hexstr) =>
    hexstr
    ->JsBuffer.hexToArray
    ->Belt_Array.sliceToEnd(-33)
    ->Sha256.digest
    ->RIPEMD160.hexDigest
    ->Address.fromHex;

let toHex =
  fun
  | PubKey(hexstr) => hexstr;

let toPubKeyHexOnly =
  fun
  | PubKey(hexstr) =>
    hexstr->JsBuffer.hexToArray->Belt_Array.sliceToEnd(-33)->JsBuffer.arrayToHex;

let toBech32 =
  fun
  | PubKey(hexstr) =>
    hexstr |> JsBuffer.hexToArray |> Bech32.toWords |> Bech32.encode("bandvalconspub");
