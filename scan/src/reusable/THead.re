module Styles = {
  open Css;

  let container = height_ =>
    style([
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
      display(`flex),
      alignItems(`center),
      height(`px(height_)),
    ]);
  let containerBase =
    style([
      backgroundColor(Colors.blueGray1),
      padding2(~v=`zero, ~h=`px(24)),
      selector("> div", [height(`px(52))]),
    ]);
};

[@react.component]
let make = (~children, ~height=30) => {
  <div className={Styles.container(height)}> children </div>;
};
module Grid = {
  [@react.component]
  let make = (~children) => {
    <div className=Styles.containerBase> children </div>;
  };
};
