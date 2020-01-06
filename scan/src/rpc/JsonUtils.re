module Decode = {
  include Json.Decode;
  let intstr = string |> map(int_of_string);
  let uamount = string |> map(float_of_string) |> map(v => v /. 1e6);
  let moment = string |> map(MomentRe.moment);
};
