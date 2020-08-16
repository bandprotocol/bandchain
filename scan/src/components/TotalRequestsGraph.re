module Styles = {
  open Css;

  let card =
    style([
      backgroundColor(Colors.white),
      height(`percent(100.)),
      borderRadius(`px(4)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      Media.mobile([height(`px(300))]),
    ]);
};

[@react.component]
let make = () => {
  <div className=Styles.card> <Text value="Total Requests" /> </div>;
};
