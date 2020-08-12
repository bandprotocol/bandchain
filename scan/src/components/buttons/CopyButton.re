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

  // Code

  let buttonCode = w =>
    style([
      backgroundColor(Colors.white),
      border(`px(1), `solid, Colors.bandBlue),
      borderRadius(`px(4)),
      cursor(`pointer),
      padding2(~v=`px(5), ~h=`px(10)),
      width(`px(w)),
    ]);
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

module Code = {
  [@react.component]
  let make = (~data, ~title, ~width=105) => {
    <a
      className={Css.merge([
        Styles.buttonCode(width),
        CssHelper.flexBox(~align=`center, ~justify=`center, ()),
      ])}
      onClick={_ => {Copy.copy(data |> JsBuffer.toUTF8)}}>
      <img src=Images.copy className=Styles.withHeight />
      <HSpacing size=Spacing.sm />
      <Text value=title size=Text.Md block=true color=Colors.bandBlue nowrap=true />
    </a>;
  };
};
