module Styles = {
  open Css;

  let button = w =>
    style([
      backgroundColor(Colors.blue1),
      padding2(~h=`px(8), ~v=`px(4)),
      display(`flex),
      width(`px(w)),
      borderRadius(`px(6)),
      cursor(`pointer),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(20, 32, 184, 0.2))),
    ]);

  let withHeight = style([maxHeight(`px(12))]);
};

[@react.component]
let make = (~data, ~title, ~width=105) => {
  <div
    className={Styles.button(width)}
    onClick={_ => {Copy.copy(data |> JsBuffer.toHex(~with0x=false))}}>
    <img src=Images.copy className=Styles.withHeight />
    <HSpacing size=Spacing.sm />
    <Text value=title size=Text.Sm block=true color=Colors.bandBlue nowrap=true />
  </div>;
};
