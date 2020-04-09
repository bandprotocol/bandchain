module Styles = {
  open Css;

  let container = (pv, ph) =>
    style([
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
      hover([backgroundColor(Colors.blueGray1)]),
      display(`flex),
      alignItems(`center),
      padding2(~v=pv, ~h=ph),
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
