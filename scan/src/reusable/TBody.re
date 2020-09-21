module Styles = {
  open Css;

  let container = (pv, ph) =>
    style([
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
      display(`flex),
      alignItems(`center),
      padding2(~v=pv, ~h=ph),
      overflow(`hidden),
    ]);
  let containerBase = (pv, ph) =>
    style([
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
      padding2(~v=pv, ~h=ph),
      overflow(`hidden),
    ]);
  let minHeight = height => style([minHeight(`px(height))]);
};

[@react.component]
let make = (~minHeight=45, ~children, ~paddingV=`px(10), ~paddingH=`zero) => {
  <div
    className={Css.merge([Styles.container(paddingV, paddingH), Styles.minHeight(minHeight)])}>
    children
  </div>;
};
module Grid = {
  [@react.component]
  let make = (~children, ~paddingV=`px(15), ~paddingH=`zero) => {
    <div className={Css.merge([Styles.containerBase(paddingV, paddingH)])}> children </div>;
  };
};
