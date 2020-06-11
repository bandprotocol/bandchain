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

type status_t =
  | Data(string)
  | NoData
  | Loading;

[@react.component]
let make = (~moniker, ~identity, ~width=25) => {
  let (url, setUrl) = React.useState(_ => Loading);

  React.useEffect1(
    () => {
      if (identity != "") {
        let _ = {
          Axios.get(
            {j|https://keybase.io/_/api/1.0/user/lookup.json?key_suffix=$identity&fields=pictures|j},
          )
          |> Js.Promise.then_(res => {
               switch (res##data |> decode) {
               | Some(url) =>
                 setUrl(_ => Data(url));
                 Js.Promise.resolve();
               | None => Js.Promise.reject(Not_found)
               | exception err => Js.Promise.reject(err)
               }
             })
          |> Js.Promise.catch(err => {
               Js.Console.log(err);
               setUrl(_ => NoData);
               Js.Promise.resolve();
             });
        };
        ();
      } else {
        setUrl(_ => NoData);
      };
      None;
    },
    [|identity|],
  );

  switch (url) {
  | Data(url') => <img src=url' className={Styles.avatar(width)} />
  | NoData =>
    <img
      src={j|https://ui-avatars.com/api/?rounded=true&size=128&name=$moniker&color=fff&background=CA47EB|j}
      className={Styles.avatar(width)}
    />
  | Loading => <LoadingCensorBar width height={width - 4} radius=100 />
  };
};
