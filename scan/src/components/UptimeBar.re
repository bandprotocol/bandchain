module Styles = {
  open Css;

  let bar =
    style([
      display(`flex),
      width(`percent(100.)),
      height(`px(15)),
      border(`px(1), `solid, Colors.blueGray4),
      borderRadius(`px(6)),
      padding(`px(1)),
      justifyContent(`flexStart),
    ]);

  let innerBar = (percent, color) =>
    style([backgroundColor(color), width(`percent(percent)), borderRadius(`px(6))]);
};

[@react.component]
let make = (~percent) => {
  let color =
    if (percent == 100.) {
      Colors.purple8;
    } else if (percent < 100. && percent >= 79.) {
      Colors.purple4;
    } else {
      Colors.purple2;
    };

  <div className=Styles.bar> <div className={Styles.innerBar(percent, color)} /> </div>;
};
