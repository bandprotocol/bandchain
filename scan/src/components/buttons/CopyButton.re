module Styles = {
  open Css;

  let logo = style([width(`px(12))]);

  let clickable = style([cursor(`pointer)]);
};

[@react.component]
let make = (~data, ~title, ~width=105, ~py=5, ~px=10, ~pySm=py, ~pxSm=px) => {
  let (copied, setCopy) = React.useState(_ => false);
  <Button
    variant=Button.Outline
    onClick={_ => {
      Copy.copy(data);
      setCopy(_ => true);
      let _ = Js.Global.setTimeout(() => setCopy(_ => false), 700);
      ();
    }}
    py
    px
    pySm
    pxSm>
    <div className={CssHelper.flexBox(~align=`center, ~justify=`center, ())}>
      {copied
         ? <img src=Images.tickIcon className=Styles.logo />
         : <img src=Images.copy className=Styles.logo />}
      <HSpacing size=Spacing.sm />
      <Text value=title size=Text.Md block=true color=Colors.bandBlue nowrap=true />
    </div>
  </Button>;
};
