module Styles = {
  open Css;

  let button =
    style([
      display(`flex),
      height(`px(30)),
      justifyContent(`center),
      alignItems(`center),
      border(`px(1), `solid, Colors.blueGray3),
      borderRadius(`px(4)),
      cursor(`pointer),
      alignSelf(`center),
      outline(`zero, `none, Colors.white),
      opacity(1.),
      transition(~duration=400, "all"),
      disabled([cursor(`default), opacity(0.5)]),
    ]);
};

[@react.component]
let make = (~onClick, ~disabled) => {
  <button className=Styles.button onClick disabled> <Text value="MAX" /> </button>;
};
