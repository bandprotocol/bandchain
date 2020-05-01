[@bs.deriving abstract]
type decoded_t = {
  prefix: string,
  words: array(int),
};

[@bs.module "bech32"] [@bs.val] external fromWords: array(int) => array(int) = "fromWords";
[@bs.module "bech32"] [@bs.val] external toWords: array(int) => array(int) = "toWords";

[@bs.module "bech32"] [@bs.val] external decode: string => decoded_t = "decode";
[@bs.module "bech32"] [@bs.val] external encode: (string, array(int)) => string = "encode";

let decodeOpt = str =>
  switch (str->decode) {
  | result => Some(result)
  | exception _ => None
  };
