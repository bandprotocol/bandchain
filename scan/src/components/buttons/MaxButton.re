module Styles = {
  open Css;

  let button =
    style([
      display(`flex),
      height(`px(37)),
      justifyContent(`center),
      background(Colors.white),
      alignItems(`center),
      border(`px(1), `solid, Colors.bandBlue),
      borderRadius(`px(4)),
      cursor(`pointer),
      alignSelf(`center),
      outline(`zero, `none, Colors.white),
      opacity(1.),
      transition(~duration=400, "all"),
      disabled([cursor(`default), opacity(0.5)]),
      padding(`px(10)),
    ]);
};

[@react.component]
let make = (~onClick, ~disabled) => {
  <button className=Styles.button onClick disabled>
    <Text value="Max" color=Colors.bandBlue weight=Text.Medium />
  </button>;
};
