module Styles = {
  open Css;

  let errorContainer =
    style([
      padding(`px(10)),
      color(Colors.red5),
      backgroundColor(Colors.red1),
      border(`px(1), `solid, Colors.red5),
      borderRadius(`px(4)),
      marginBottom(`px(24)),
      selector("> i", [marginRight(`px(8))]),
    ]);
};

type log_t = {message: string};

type err_t = {log: option(string)};

let decodeLog = json => JsonUtils.Decode.{message: json |> field("message", string)};

let decode = json => JsonUtils.Decode.{log: json |> optional(field("log", string))};

let parseErr = msg => {
  exception WrongNetwork(string);
  switch (Env.network) {
  | "GUANYU" => msg
  | "WENCHANG"
  | "GUANYU38" =>
    let err =
      {
        let%Opt json = msg |> Json.parse;
        let%Opt x = json |> Js.Json.decodeArray;
        let%Opt y = x->Belt.Array.get(0);
        let%Opt logStr = (y |> decode).log;
        let%Opt logJson = logStr |> Json.parse;
        let log = logJson |> decodeLog;
        Opt.ret(log.message);
      }
      |> Belt.Option.getWithDefault(_, msg);
    "Error: " ++ err;
  | _ => raise(WrongNetwork("Incorrect or unspecified NETWORK environment variable"))
  };
};

module Full = {
  [@react.component]
  let make = (~msg) => {
    <div className={Css.merge([Styles.errorContainer, CssHelper.flexBox(~wrap=`nowrap, ())])}>
      <Icon name="fal fa-exclamation-circle" size=14 color=Colors.red5 />
      <Text value={msg |> parseErr} size=Text.Lg spacing={Text.Em(0.02)} breakAll=true />
    </div>;
  };
};

module Mini = {
  [@react.component]
  let make = (~msg) => {
    <Text value={msg |> parseErr} code=true size=Text.Sm breakAll=true />;
  };
};
