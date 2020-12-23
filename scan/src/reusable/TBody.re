module Styles = {
  open Css;

  let containerBase = (pv, ph) =>
    style([
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, `num(0.05))),
      ),
      backgroundColor(white),
      marginBottom(`px(1)),
      padding2(~v=pv, ~h=ph),
      overflow(`hidden),
    ]);
  let minHeight = height => style([minHeight(`px(height))]);
};

[@react.component]
let make = (~children, ~paddingV=`px(15), ~paddingH=`zero) => {
  <div className={Css.merge([Styles.containerBase(paddingV, paddingH)])}> children </div>;
};
