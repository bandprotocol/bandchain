type pos_t =
  | Title
  | Subtitle
  | Text
  | None;

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

let letterSpacing =
  fun
  | Title
  | Text => Text.Unset
  | Subtitle => Text.Em(0.02);

module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      maxWidth(`px(360)),
      justifyContent(`flexEnd),
      wordBreak(`breakAll),
    ]);
};

[@react.component]
let make = (~pubKey, ~position=Text) => {
  let noPrefixAddress =
    pubKey |> PubKey.toBech32 |> Js.String.sliceToEnd(~from=14);

  <div className=Styles.container>
    <Text
      value="bandvalconspub"
      size={position |> prefixFontSize}
      weight=Text.Semibold
      code=true
      spacing={position |> letterSpacing}
      block=true
      nowrap=true
    />
    <Text
      value=noPrefixAddress
      size={position |> pubKeyFontSize}
      weight=Text.Regular
      spacing={position |> letterSpacing}
      code=true
      block=true
      align=Text.Right
      height={Text.Px(16)}
    />
  </div>;
};
