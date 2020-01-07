module Styles = {
  open Css;

  let vFlex =
    style([
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      justifyContent(`center),
      height(`px(600)),
    ]);
};

[@react.component]
let make = () => {
  <div className=Styles.vFlex>
    <Text value="Block Home Page" size=Text.Xxl weight=Text.Bold nowrap=true />
  </div>;
};
