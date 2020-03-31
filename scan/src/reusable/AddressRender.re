type is_pointer =
  | IsPointer
  | NonPointer;

type pos_t =
  | Title
  | Subtitle
  | Text(is_pointer);

let prefixFontSize =
  fun
  | Title => Text.Xxl
  | Subtitle => Text.Lg
  | Text(_) => Text.Md;

let addressFontSize =
  fun
  | Title => Text.Xxl
  | Subtitle => Text.Lg
  | Text(_) => Text.Md;

let lineHeight =
  fun
  | Title => Text.Px(23)
  | Subtitle => Text.Px(18)
  | Text(_) => Text.Px(16);

let letterSpacing =
  fun
  | Title
  | Text(_) => Text.Unset
  | Subtitle => Text.Em(0.02);

module Styles = {
  open Css;

  let container = style([display(`flex), cursor(`pointer)]);

  let pointerEvents =
    fun
    | Title
    | Subtitle
    | Text(IsPointer) => style([pointerEvents(`auto)])
    | Text(NonPointer) => style([pointerEvents(`none)]);
};

[@react.component]
let make = (~address, ~position=Text(IsPointer), ~validator=false) => {
  let noPrefixAddress =
    validator
      ? address |> Address.toOperatorBech32 |> Js.String.sliceToEnd(~from=11)
      : address |> Address.toBech32 |> Js.String.sliceToEnd(~from=4);

  <div
    className={Css.merge([
      Styles.container,
      Styles.pointerEvents(validator ? Text(NonPointer) : position),
    ])}
    onClick={_ => Route.redirect(Route.AccountIndexPage(address, Route.AccountTransactions))}>
    <Text
      value={validator ? "bandvaloper" : "band"}
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
