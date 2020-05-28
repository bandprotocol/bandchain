type position =
  | Title
  | Subtitle
  | Text
  | Nav;

module Styles = {
  open Css;

  let container = style([display(`flex), cursor(`pointer)]);

  let clickable = active =>
    active ? style([pointerEvents(`auto)]) : style([pointerEvents(`none)]);

  let prefix = style([fontWeight(`num(600))]);

  let font =
    fun
    | Title => style([fontSize(`px(18)), lineHeight(`px(24))])
    | Subtitle => style([fontSize(`px(14)), lineHeight(`px(28)), letterSpacing(`em(0.02))])
    | Text => style([fontSize(`px(12)), lineHeight(`px(16))])
    | Nav => style([fontSize(`px(10)), lineHeight(`px(14))]);

  let base =
    style([overflow(`hidden), textOverflow(`ellipsis), whiteSpace(`nowrap), display(`block)]);

  let copy = style([width(`px(15)), marginLeft(`px(10)), cursor(`pointer)]);
};

[@react.component]
let make = (~address, ~position=Text, ~validator=false, ~copy=false, ~clickable=true) => {
  let prefix = validator ? "bandvaloper" : "band";

  let noPrefixAddress =
    validator
      ? address |> Address.toOperatorBech32 |> Js.String.sliceToEnd(~from=11)
      : address |> Address.toBech32 |> Js.String.sliceToEnd(~from=4);

  <>
    <div
      className={Css.merge([Styles.container, Styles.clickable(clickable)])}
      onClick={_ =>
        Route.redirect(
          validator
            ? Route.ValidatorIndexPage(address, Route.ProposedBlocks)
            : Route.AccountIndexPage(address, Route.AccountTransactions),
        )
      }>
      <span className={Css.merge([Styles.font(position), Styles.base, Text.Styles.code])}>
        <span className=Styles.prefix> {prefix |> React.string} </span>
        {noPrefixAddress |> React.string}
      </span>
    </div>
    {copy
       ? <>
           {switch (position) {
            | Title => <HSpacing size=Spacing.md />
            | _ => <HSpacing size=Spacing.sm />
            }}
           <CopyRender
             width={
               switch (position) {
               | Title => 15
               | _ => 12
               }
             }
             message={
               validator ? address |> Address.toOperatorBech32 : address |> Address.toBech32
             }
           />
         </>
       : React.null}
  </>;
};
