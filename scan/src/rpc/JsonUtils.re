module Decode = {
  include Json.Decode;

  let intstr = string |> map(int_of_string);
  let moment = string |> map(MomentRe.moment);
};
