let int64 = json => json |> Js.Json.decodeNumber |> Belt.Option.getExn |> int_of_float;
let string = json => json |> Js.Json.decodeString |> Belt.Option.getExn;
let stringExn = (stringOpt: option(string)) => stringOpt |> Belt_Option.getExn;
let buffer = str => str |> Js.String.substr(~from=2) |> JsBuffer.fromHex;

let timeS = json => {
  json
  |> Js.Json.decodeNumber
  |> Belt.Option.getExn
  |> int_of_float
  |> MomentRe.momentWithUnix
  |> MomentRe.Moment.defaultUtc;
};

let timeMS = json => {
  json
  |> Js.Json.decodeNumber
  |> Belt.Option.getExn
  |> MomentRe.momentWithTimestampMS
  |> MomentRe.Moment.defaultUtc;
};

let timestamp = json =>
  json
  |> Js.Json.decodeString
  |> Belt.Option.getExn
  |> MomentRe.momentWithFormat(_, "YYYY-MM-DDTHH:mm:ss.SSSSSS")
  |> MomentRe.Moment.defaultUtc;

let timestampWithDefault = jsonOpt =>
  jsonOpt
  |> Belt_Option.flatMap(_, x => Some(timestamp(x)))
  |> Belt.Option.getWithDefault(_, MomentRe.momentNow());

let optionBuffer = Belt_Option.map(_, buffer);

let optionTimeMS = Belt_Option.map(_, timeMS);

let optionTimeS = Belt_Option.map(_, timeS);

let optionTimeSExn = timeSOpt => timeSOpt |> Belt_Option.getExn |> timeS;

let bool = json => json |> Js.Json.decodeBoolean |> Belt.Option.getExn;

let hash = str => str |> Js.String.substr(~from=2) |> Hash.fromHex;

let coinRegEx = "([0-9]+)([a-z][a-z0-9/]{2,31})" |> Js.Re.fromString;

let intToCoin = int_ => int_ |> float_of_int |> Coin.newUBANDFromAmount;

let coin = json => {
  json |> Js.Json.decodeNumber |> Belt_Option.getExn |> Coin.newUBANDFromAmount;
};

let coinExn = jsonOpt => {
  jsonOpt
  |> Belt_Option.flatMap(_, Js.Json.decodeNumber)
  |> Belt.Option.getExn
  |> Coin.newUBANDFromAmount;
};
let coinWithDefault = jsonOpt => {
  jsonOpt
  |> Belt_Option.flatMap(_, Js.Json.decodeNumber)
  |> Belt.Option.getWithDefault(_, 0.0)
  |> Coin.newUBANDFromAmount;
};

let coins = str =>
  str
  |> Js.String.split(",")
  |> Belt_List.fromArray
  |> Belt_List.keepMap(_, coin =>
       if (coin == "") {
         None;
       } else {
         let result = coin |> Js.Re.exec_(coinRegEx) |> Belt_Option.getExn |> Js.Re.captures;
         Some({
           Coin.denom: result[2] |> Js.Nullable.toOption |> Belt_Option.getExn,
           amount: result[1] |> Js.Nullable.toOption |> Belt_Option.getExn |> float_of_string,
         });
       }
     );

let addressExn = jsonOpt => jsonOpt |> Belt_Option.getExn |> Address.fromBech32;

let numberWithDefault = jsonOpt =>
  jsonOpt |> Belt_Option.flatMap(_, Js.Json.decodeNumber) |> Belt.Option.getWithDefault(_, 0.0);

let floatWithDefault = jsonOpt =>
  jsonOpt |> Belt_Option.flatMap(_, Js.Json.decodeNumber) |> Belt.Option.getWithDefault(_, 0.);

let floatExn = json => {
  json |> Js.Json.decodeNumber |> Belt.Option.getExn;
};
