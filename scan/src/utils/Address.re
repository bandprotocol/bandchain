type t =
  | Address(string); // string is hex (without 0x)

exception WrongPrefixAddress(string);

let validatePrefix = bech32Address => {
  let prefix = bech32Address->Bech32.prefixGet;
  prefix == "band" || prefix == "bandvaloper" || prefix == "bandvalconspub";
};

let fromBech32 = bech32str => {
  let decoded = bech32str->Bech32.decode;
  validatePrefix(decoded)
    ? Address(decoded->Bech32.wordsGet->Bech32.fromWords->JsBuffer.arrayToHex)
    : raise(WrongPrefixAddress("Address is not correct prefix."));
};

let fromBech32Opt = bech32str => {
  let decodedOpt = bech32str |> Bech32.decodeOpt;
  let%Opt decoded = decodedOpt;

  validatePrefix(decoded)
    ? Some(Address(decoded->Bech32.wordsGet->Bech32.fromWords->JsBuffer.arrayToHex)) : None;
};

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
