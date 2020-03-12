module Styles = {
  open Css;

  let container =
    style([
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
      hover([backgroundColor(Colors.lighterPurple)]),
      display(`flex),
      alignItems(`center),
      padding2(~v=`px(10), ~h=`px(0)),
      minHeight(`px(60)),
    ]);
};

[@react.component]
let make = (~children) => {
  <div className=Styles.container> children </div>;
};
