module Styles = {
  open Css;

  let avatar = width_ => style([width(`px(width_)), borderRadius(`percent(50.))]);
};

let decodeThem = json =>
  json |> JsonUtils.Decode.at(["pictures", "primary", "url"], JsonUtils.Decode.string);

let decode = json =>
  json
  |> JsonUtils.Decode.field("them", JsonUtils.Decode.array(decodeThem))
  |> Belt.Array.get(_, 0);

module Placeholder = {
  [@react.component]
  let make = (~moniker, ~width) =>
    <img
      src={j|https://ui-avatars.com/api/?rounded=true&size=128&name=$moniker&color=9714B8&background=F3CEFD|j}
      className={Styles.avatar(width)}
    />;
};

module Keybase = {
  [@react.component]
  let make = (~identity, ~moniker, ~width) =>
    {
      let resOpt =
        AxiosHooks.use(
          {j|https://keybase.io/_/api/1.0/user/lookup.json?key_suffix=$identity&fields=pictures|j},
        );
      let%Opt res = resOpt;

      switch (res |> decode) {
      | Some(url) => Some(<img src=url className={Styles.avatar(width)} />)
      | None =>
        // Log for debug
        Js.Console.log3("none", identity, res);
        Some(<Placeholder moniker width />);
      | exception err =>
        // Log for debug
        Js.Console.log3(identity, res, err);
        Some(<Placeholder moniker width />);
      };
    }
    |> Belt.Option.getWithDefault(_, <LoadingCensorBar width height={width - 4} radius=100 />);
};

[@react.component]
let make = (~moniker, ~identity, ~width=25) =>
  React.useMemo1(
    () => identity != "" ? <Keybase identity moniker width /> : <Placeholder moniker width />,
    [|identity|],
  );
