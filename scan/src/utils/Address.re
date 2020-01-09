type t =
  | Address(string); // string is hex (without 0x)

let fromBech32 = bech32str => {
  Address(bech32str->Bech32.decode->Bech32.wordsGet->Bech32.fromWords->JsBuffer.arrayToHex);
};

let fromHex = hexstr =>
  if (hexstr->String.sub(0, 2) == "0x") {
    Address(hexstr->String.lowercase_ascii->String.sub(2, 40));
  } else {
    Address(hexstr->String.lowercase_ascii);
  };

let toHex =
  fun
  | Address(hexstr) => hexstr;

let bech32ToHex = bech32str => bech32str->fromBech32->toHex;

let toOperatorBech32 =
  fun
  | Address(hexstr) =>
    hexstr |> JsBuffer.hexToArray |> Bech32.toWords |> Bech32.encode("bandvaloper");

let toBech32 =
  fun
  | Address(hexstr) => hexstr |> JsBuffer.hexToArray |> Bech32.toWords |> Bech32.encode("band");

let hexToOperatorBech32 = hexstr => hexstr->fromHex->toOperatorBech32;
let hexToBech32 = hexstr => hexstr->fromHex->toBech32;
