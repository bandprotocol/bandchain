module Styles = {
  open Css;

  let button =
    style([
      backgroundColor(Colors.gray4),
      padding2(~h=`px(8), ~v=`px(4)),
      display(`flex),
      width(`px(110)),
      borderRadius(`px(6)),
      cursor(`pointer),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(20, 32, 184, 0.2))),
    ]);

  let withHeight = style([maxHeight(`px(12))]);
};

[@react.component]
let make = (~link, ~description) => {
  <a href=link target="_blank" rel="noopener">
    <div className=Styles.button>
      <img src=Images.externalLink className=Styles.withHeight />
      <HSpacing size=Spacing.sm />
      <Text value=description size=Text.Sm block=true color=Colors.gray7 nowrap=true />
    </div>
  </a>;
};
