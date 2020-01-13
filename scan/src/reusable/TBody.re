module Styles = {
  open Css;

  let container = height_ =>
    style([
      height(`px(height_)),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
      cursor(`pointer),
      hover([backgroundColor(Colors.lighterPurple)]),
      display(`flex),
      alignItems(`center),
    ]);
};

[@react.component]
let make = (~children, ~height=60) => {
  <div className={Styles.container(height)}> children </div>;
};
