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

let fPretty = value =>
  withCommas(
    if (value > 1000000.) {
      value->Js.Float.toFixedWithPrecision(~digits=0);
    } else if (value > 100.) {
      value->Js.Float.toFixedWithPrecision(~digits=2);
    } else if (value > 1.) {
      value->Js.Float.toFixedWithPrecision(~digits=4);
    } else {
      value->Js.Float.toFixedWithPrecision(~digits=6);
    },
  );

let fPercent = value =>
  (value > 0. ? "+" : "") ++ value->Js.Float.toFixedWithPrecision(~digits=2) ++ "%";

let iPretty = value => withCommas(value->string_of_int);
