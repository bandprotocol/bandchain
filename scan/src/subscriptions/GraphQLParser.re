let int64 = json => json |> Js.Json.decodeNumber |> Belt.Option.getExn |> int_of_float;
let buffer = json =>
  json
  |> Js.Json.decodeString
  |> Belt.Option.getExn
  |> Js.String.substr(~from=2)
  |> JsBuffer.fromHex;

let time = json => json |> int64 |> MomentRe.momentWithUnix;

let hash = json =>
  json |> Js.Json.decodeString |> Belt.Option.getExn |> Js.String.substr(~from=2) |> Hash.fromHex;
