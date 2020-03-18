module Styles = {
  open Css;

  let container =
    style([
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
      hover([backgroundColor(Colors.blueGray1)]),
      display(`flex),
      alignItems(`center),
      padding2(~v=`px(10), ~h=`zero),
    ]);

  let minHeight = height => style([minHeight(`px(height))]);
};

[@react.component]
let make = (~minHeight=45, ~children) => {
  <div className={Css.merge([Styles.container, Styles.minHeight(minHeight)])}> children </div>;
};
