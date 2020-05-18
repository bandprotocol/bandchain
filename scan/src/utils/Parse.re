let parseBandAmount = amount => {
  let%Opt amountFloat = float_of_string_opt(amount);

  amountFloat *. 1e6 == Js.Math.floor_float(amountFloat *. 1e6)
    ? Some(amountFloat *. 1e6) : None;
};
