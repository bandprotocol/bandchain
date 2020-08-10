module Styles = {
  open Css;
  let searchContainer =
    style([
      display(`flex),
      alignItems(`center),
      position(`relative),
      before([
        backgroundImage(`url(Images.searchGray)),
        contentRule(`text("")),
        width(`px(15)),
        height(`px(15)),
        backgroundRepeat(`noRepeat),
        display(`block),
        backgroundPositions([`center, `center]),
        position(`absolute),
        top(`percent(5.)),
      ]),
    ]);
  let searchBar =
    style([
      backgroundColor(Colors.transparent),
      borderRadius(`zero),
      border(`zero, `none, Colors.white),
      borderBottom(`px(1), `solid, Colors.gray8),
      placeholder([color(Colors.blueGray3)]),
      paddingLeft(`px(20)),
      paddingBottom(`px(10)),
      focus([outlineStyle(`none)]),
      width(`percent(100.)),
      maxWidth(`px(300)),
    ]);
};

[@react.component]
let make = (~placeholder, ~onChange, ~debounce=500) => {
  let (changeValue, setChangeValue) = React.useState(_ => "");

  React.useEffect1(
    () => {
      let timeoutId = Js.Global.setTimeout(() => onChange(_ => changeValue), debounce);
      Some(() => Js.Global.clearTimeout(timeoutId));
    },
    [|changeValue|],
  );

  <div className=Styles.searchContainer>
    <input
      type_="text"
      className=Styles.searchBar
      placeholder
      onChange={event => {
        let newVal = ReactEvent.Form.target(event)##value |> String.lowercase_ascii;
        setChangeValue(_ => newVal);
      }}
    />
  </div>;
};
