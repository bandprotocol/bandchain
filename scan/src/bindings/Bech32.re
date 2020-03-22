[@bs.deriving abstract]
type decoded_t = {
  prefix: string,
  words: array(int),
};

let safeDecode: string => (decoded_t, bool) = [%bs.raw
  {|
function(str) {
  const bech32 = require("bech32");
  try {
    const result = bech32.decode(str);
    return [result,true];
  } catch(_) {
    return [{prefix:"",words:[]}, false];
  }
}
  |}
];

let decodeOpt = str =>
  switch (str->safeDecode) {
  | (result, true) => Some(result)
  | _ => None
  };

[@bs.module "bech32"] [@bs.val] external fromWords: array(int) => array(int) = "fromWords";
[@bs.module "bech32"] [@bs.val] external toWords: array(int) => array(int) = "toWords";

[@bs.module "bech32"] [@bs.val] external decode: string => decoded_t = "decode";
[@bs.module "bech32"] [@bs.val] external encode: (string, array(int)) => string = "encode";
