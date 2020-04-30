type t =
  | Address(string); // string is hex (without 0x)

let fromBech32 = bech32str => {
  Address(bech32str->Bech32.decode->Bech32.wordsGet->Bech32.fromWords->JsBuffer.arrayToHex);
};

let fromBech32Opt = bech32str =>
  bech32str
  |> Bech32.decodeOpt
  |> Belt.Option.map(_, x => Address(x->Bech32.wordsGet->Bech32.fromWords->JsBuffer.arrayToHex));

let fromHex = hexstr => Address(hexstr->HexUtils.normalizeHexString);

let toHex = (~with0x=false, ~upper=false) =>
  fun
  | Address(hexstr) => {
      let lowercase = (with0x ? "0x" : "") ++ hexstr;
      upper ? lowercase->String.uppercase_ascii : lowercase;
    };

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

let isEqual = (Address(hexstr1), Address(hexst2)) => hexstr1 == hexst2;
