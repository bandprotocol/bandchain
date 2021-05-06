open Css;

type mode_t =
  | Day
  | Dark;

type color_t = Types.Color.t;

type t = {
  text1: color_t,
  text2: color_t,
  mainBg: color_t,
};

let get =
  fun
  | Day => {text1: hex("303030"), text2: hex("7D7D7D"), mainBg: hex("fcfcfc")}
  | Dark => {text1: hex("ffffff"), text2: hex("9A9A9A"), mainBg: hex("000000")};
