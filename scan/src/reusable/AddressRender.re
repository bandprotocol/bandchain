type is_pointer =
  | Pointer
  | NonPointer;

type position =
  | Title
  | Subtitle
  | Text(is_pointer)
  | Nav;

module Styles = {
  open Css;

  let container = style([display(`flex), cursor(`pointer)]);

  let pointerEvents =
    fun
    | Title
    | Subtitle
    | Nav
    | Text(Pointer) => style([pointerEvents(`auto)])
    | Text(NonPointer) => style([pointerEvents(`none)]);

  let prefix = style([fontWeight(`num(600))]);

  let font =
    fun
    | Title => style([fontSize(`px(18)), lineHeight(`px(24))])
    | Subtitle => style([fontSize(`px(14)), lineHeight(`px(28)), letterSpacing(`em(0.02))])
    | Text(_) => style([fontSize(`px(12)), lineHeight(`px(16))])
    | Nav => style([fontSize(`px(10)), lineHeight(`px(14))]);

  let base =
    style([overflow(`hidden), textOverflow(`ellipsis), whiteSpace(`nowrap), display(`block)]);

  let copy = style([width(`px(15)), marginLeft(`px(10)), cursor(`pointer)]);
};

[@react.component]
let make = (~address, ~position=Text(Pointer), ~validator=false, ~copy=false) => {
  let prefix = validator ? "bandvaloper" : "band";

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
    <span className={Css.merge([Styles.font(position), Styles.base, Text.Styles.code])}>
      <span className=Styles.prefix> {prefix |> React.string} </span>
      {noPrefixAddress |> React.string}
    </span>
    {copy
       ? <img
           src=Images.copy
           className=Styles.copy
           onClick={_ => {
             validator
               ? Copy.copy(address |> Address.toOperatorBech32)
               : Copy.copy(address |> Address.toBech32)
           }}
         />
       : React.null}
  </div>;
};
