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

  let logo = style([width(`px(12))]);

  let clickable = style([cursor(`pointer)]);
};

// TODO: we'll clean up this later.
[@react.component]
let make = (~data, ~title, ~width=105) => {
  <div
    className={Styles.button(width)}
    onClick={_ => {Copy.copy(data |> JsBuffer.toHex(~with0x=false))}}>
    <img src=Images.copy className=Styles.logo />
    <HSpacing size=Spacing.sm />
    <Text value=title size=Text.Sm block=true color=Colors.bandBlue nowrap=true />
  </div>;
};

module Modern = {
  [@react.component]
  let make = (~data, ~title, ~width=105, ~py=5, ~px=10, ~pySm=py, ~pxSm=px) => {
    let (copied, setCopy) = React.useState(_ => false);
    <div
      className={Css.merge([
        CssHelper.btn(~variant=Outline, ~px, ~py, ~pxSm, ~pySm, ()),
        CssHelper.flexBox(~align=`center, ~justify=`center, ()),
        Styles.clickable,
      ])}
      onClick={_ => {
        Copy.copy(data);
        setCopy(_ => true);
        let _ = Js.Global.setTimeout(() => setCopy(_ => false), 700);
        ();
      }}>
      {copied
         ? <img src=Images.tickIcon className=Styles.logo />
         : <img src=Images.copy className=Styles.logo />}
      <HSpacing size=Spacing.sm />
      <Text value=title size=Text.Md block=true color=Colors.bandBlue nowrap=true />
    </div>;
  };
};
