let getBandAmount = amount => {
  let%Opt amountFloat = float_of_string_opt(amount);

  // TODO: return with error if it cannot parse
  let uband = amountFloat *. 1e6;
  uband == Js.Math.floor_float(uband) ? Some(uband) : None;
};
