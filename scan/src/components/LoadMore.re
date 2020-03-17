module Styles = {
  open Css;

  let loadMore =
    style([
      width(`percent(100.)),
      height(`px(28)),
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      cursor(`pointer),
    ]);
};

[@react.component]
let make = (~onClick=_ => ()) => {
  <div className=Styles.loadMore onClick>
    <Text value="LOAD MORE" color=Colors.gray7 weight=Text.Bold size=Text.Sm />
  </div>;
};
