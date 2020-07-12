module Decode = {
  include Json.Decode;
  let intstr = string |> map(int_of_string);
  let uamount = string |> map(float_of_string);
  let moment = string |> map(MomentRe.moment);
  let floatstr = string |> map(float_of_string);
  let intWithDefault = (key, json) =>
    json |> optional(key(int)) |> Belt.Option.getWithDefault(_, 0);
  let bufferWithDefault = (key, json) =>
    json |> optional(key(string)) |> Belt.Option.getWithDefault(_, "") |> JsBuffer.fromBase64;
  let strWithDefault = (key, json) =>
    json |> optional(key(string)) |> Belt.Option.getWithDefault(_, "");
};
