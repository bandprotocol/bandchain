type pos_t =
  | Title
  | Subtitle
  | Text;

let prefixFontSize =
  fun
  | Title => Text.Xxl
  | Subtitle => Text.Lg
  | Text => Text.Md;

let addressFontSize =
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

  let container = style([display(`flex), cursor(`pointer)]);

  let pointerEvents =
    fun
    | Title => style([pointerEvents(`none)])
    | Subtitle
    | Text => style([pointerEvents(`auto)]);
};

[@react.component]
let make = (~address, ~position=Text) => {
  let noPrefixAddress = address |> Address.toBech32 |> Js.String.sliceToEnd(~from=4);

  <div
    className={Css.merge([Styles.container, Styles.pointerEvents(position)])}
    onClick={_ => Route.redirect(Route.AccountIndexPage(address, Route.AccountTransactions))}>
    <Text
      value="band"
      size={position |> prefixFontSize}
      weight=Text.Semibold
      code=true
      spacing={position |> letterSpacing}
    />
    <Text
      value=noPrefixAddress
      size={position |> addressFontSize}
      weight=Text.Regular
      spacing={position |> letterSpacing}
      code=true
      nowrap=true
      block=true
      ellipsis=true
    />
  </div>;
};
