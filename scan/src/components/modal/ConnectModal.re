module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      justifyContent(`center),
      width(`percent(100.)),
      height(`percent(100.)),
    ]);
};

[@react.component]
let make = (~value) => {
  <div className=Styles.container> <Text value /> </div>;
};
