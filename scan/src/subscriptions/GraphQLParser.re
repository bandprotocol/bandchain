let int64 = json => json |> Js.Json.decodeNumber |> Belt.Option.getExn |> int_of_float;
let bytea = json =>
  json
  |> Js.Json.decodeString
  |> Belt.Option.getExn
  |> Js.String.substr(~from=2)
  |> JsBuffer.fromHex;
