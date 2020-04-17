type is_pointer =
  | Pointer
  | NonPointer;

type position =
  | Title
  | Subtitle
  | Text(is_pointer)
  | Nav;

let prefixFontSize =
  fun
  | Title => Text.Xxl
  | Subtitle => Text.Lg
  | Text(_) => Text.Md
  | Nav => Text.Sm;

let addressFontSize =
  fun
  | Title => Text.Xxl
  | Subtitle => Text.Lg
  | Text(_) => Text.Md
  | Nav => Text.Sm;

let lineHeight =
  fun
  | Title => Text.Px(24)
  | Subtitle => Text.Px(18)
  | Text(_) => Text.Px(16)
  | Nav => Text.Px(14);

let letterSpacing =
  fun
  | Title
  | Text(_)
  | Nav => Text.Unset
  | Subtitle => Text.Em(0.02);

module Styles = {
  open Css;

  let container = style([display(`flex), cursor(`pointer)]);

  let pointerEvents =
    fun
    | Title
    | Subtitle
    | Text(Pointer)
    | Nav => style([pointerEvents(`auto)])
    | Text(NonPointer) => style([pointerEvents(`none)]);
};

[@react.component]
let make = (~address, ~position=Text(Pointer), ~validator=false) => {
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
      height={position |> lineHeight}
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
      height={position |> lineHeight}
    />
  </div>;
};
