let withCommas = value =>
  value
  |> Js.String.split(".")
  |> Js.Array.mapi((part, idx) =>
       if (idx == 0) {
         part
         |> Js.String.split("")
         |> Js.Array.reverseInPlace
         |> Js.Array.reducei(
              (acc, ch, idx) => idx mod 3 == 0 && idx != 0 ? ch ++ "," ++ acc : ch ++ acc,
              "",
            );
       } else {
         "." ++ part;
       }
     )
  |> Js.Array.reduce((a, b) => a ++ b, "");

let fPretty = (~digits=?, value) => {
  switch (digits) {
  | Some(digits') => withCommas(value->Js.Float.toFixedWithPrecision(~digits=digits'))
  | None =>
    withCommas(
      if (value >= 1000000.) {
        value->Js.Float.toFixedWithPrecision(~digits=0);
      } else if (value > 100.) {
        value->Js.Float.toFixedWithPrecision(~digits=2);
      } else if (value > 1.) {
        value->Js.Float.toFixedWithPrecision(~digits=4);
      } else {
        value->Js.Float.toFixedWithPrecision(~digits=6);
      },
    )
  };
};

let fCurrency = value =>
  if (value >= 1e9) {
    "$" ++ (value /. 1e9 |> fPretty(~digits=2)) ++ "B";
  } else if (value >= 1e6) {
    "$" ++ (value /. 1e6 |> fPretty(~digits=2)) ++ "M";
  } else if (value >= 1e3) {
    "$" ++ (value /. 1e3 |> fPretty(~digits=2)) ++ "K";
  } else {
    "$" ++ (value |> fPretty(~digits=2));
  };

let fPercentChange = value =>
  (value > 0. ? "+" : "") ++ value->Js.Float.toFixedWithPrecision(~digits=2) ++ "%";

let fPercent = (~digits=?, value) => {
  (
    switch (digits) {
    | Some(digits') => withCommas(value->Js.Float.toFixedWithPrecision(~digits=digits'))
    | None =>
      withCommas(
        if (value > 1.) {
          value->Js.Float.toFixedWithPrecision(~digits=2);
        } else {
          value->Js.Float.toFixedWithPrecision(~digits=6);
        },
      )
    }
  )
  ++ "%";
};

let iPretty = value => withCommas(value->string_of_int);
