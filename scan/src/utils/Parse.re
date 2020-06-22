let getBandAmount = (maxValue, amount) => {
  let maxValueBigInt = maxValue |> Big.fromFloat;
  switch (amount |> Big.fromString) {
  | amount =>
    let uband = amount |> Big.times(Config.ubandMult);
    let roundedUband = uband->Big.round(~dp=0, ());
    Big.eq(uband, roundedUband)
      ? {
        switch (Big.lte(uband, maxValueBigInt), Big.gt(uband, Big.fromInt(0))) {
        | (true, true) => Result.Ok(uband->Big.toFixed(~dp=0, ())->Int64.of_string)
        | (false, _) => Err("Insufficient Amount")
        | (_, false) => Err("Amount must be more than 0")
        };
      }
      : Err("Maximum precision is 6");
  | exception _ => Err("Invalid value")
  };
};

let address = addr => {
  switch (Address.fromBech32Opt(addr->String.trim)) {
  | Some(address) => Result.Ok(address)
  | None => Err("Invalid address")
  };
};
