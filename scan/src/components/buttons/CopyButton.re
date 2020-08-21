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

  let buttonCode = (~w, ~py, ~px, ()) =>
    style([
      backgroundColor(Colors.white),
      border(`px(1), `solid, Colors.bandBlue),
      borderRadius(`px(4)),
      cursor(`pointer),
      padding2(~v=`px(py), ~h=`px(px)),
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

module Modern = {
  [@react.component]
  let make = (~data, ~title, ~width=105, ~py=5, ~px=10) => {
    let (copied, setCopy) = React.useState(_ => false);
    <a
      className={Css.merge([
        Styles.buttonCode(~w=width, ~px, ~py, ()),
        CssHelper.flexBox(~align=`center, ~justify=`center, ()),
      ])}
      onClick={_ => {
        Copy.copy(data);
        setCopy(_ => true);
        let _ = Js.Global.setTimeout(() => setCopy(_ => false), 700);
        ();
      }}>
      {copied
         ? <img src=Images.tickIcon className=Styles.withHeight />
         : <img src=Images.copy className=Styles.withHeight />}
      <HSpacing size=Spacing.sm />
      <Text value=title size=Text.Md block=true color=Colors.bandBlue nowrap=true />
    </a>;
  };
};
