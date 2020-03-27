// This function receive any string and fliter the latest part
// that is hex and then make it lowercase.
// In other words, just strip prefix out and then lowercase.
// Please see HexUtils_test.re for examples.
let normalizeHexString = hexstr =>
  hexstr
  ->Js.Re.exec_("[0-9a-fA-F]+$"->Js.Re.fromString, _)
  ->Belt_Option.mapWithDefault([||], result =>
      result->Js.Re.captures->Belt_Array.keepMap(Js.toOption)
    )
  ->Belt_Array.get(0)
  ->Belt_Option.getWithDefault(_, "")
  ->String.lowercase_ascii;
