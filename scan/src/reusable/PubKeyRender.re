type pos_t =
  | Title
  | Subtitle
  | Text;

let prefixFontSize =
  fun
  | Title => Text.Xxl
  | Subtitle => Text.Lg
  | Text => Text.Md;

let pubKeyFontSize =
  fun
  | Title => Text.Xxl
  | Subtitle => Text.Lg
  | Text => Text.Md;

let lineHeight =
  fun
  | Title => Text.Px(23)
  | Subtitle => Text.Px(18)
  | Text => Text.Px(16);

module Styles = {
  open Css;

  let container = display_ => style([display(display_), wordBreak(`breakAll)]);
};

[@react.component]
let make = (~pubKey, ~position=Text, ~alignLeft=false, ~display=`flex) => {
  let noPrefixAddress = pubKey |> PubKey.toBech32 |> Js.String.sliceToEnd(~from=14);

  <div className={Styles.container(display)}>
    <Text
      value={"bandvalconspub" ++ noPrefixAddress}
      size={position |> pubKeyFontSize}
      code=true
      align=?{alignLeft ? None : Some(Text.Right)}
    />
  </div>;
};
