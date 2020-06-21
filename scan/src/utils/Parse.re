let getBandAmount = (maxValue, amount) => {
  switch (float_of_string_opt(amount)) {
  | Some(amountFloat) =>
    let uband = amountFloat *. 1e6;
    uband == Js.Math.floor_float(uband)
      ? {
        switch (uband <= maxValue, uband > 0.) {
        | (true, true) => Result.Ok(uband)
        | (false, _) => Err("Insufficient Amount")
        | (_, false) => Err("Amount must be more than 0")
        };
      }
      : Err("Maximum precision is 4");
  | None => Err("Invalid value")
  };
};

let address = addr => {
  switch (Address.fromBech32Opt(addr->String.trim)) {
  | Some(address) => Result.Ok(address)
  | None => Err("Invalid address")
  };
};
