type alignSelf =
  | FlexStart
  | Center
  | FlexEnd;

module Styles = {
  open Css;

  let col = style([margin4(~top=`px(0), ~right=Spacing.xs, ~left=Spacing.xs, ~bottom=`px(0))]);
  let colSize = sz => style([flex(`num(sz))]);
  let alignSelf =
    Belt.Option.mapWithDefault(
      _,
      style([]),
      fun
      | FlexStart => style([alignSelf(`flexStart)])
      | Center => style([alignSelf(`center)])
      | FlexEnd => style([alignSelf(`flexEnd)])
    );
};

[@react.component]
let make = (~size=?, ~alignSelf=?, ~children) => {
  <div
    className={Css.merge([
      Styles.col,
      size->Belt.Option.mapWithDefault("", Styles.colSize),
      Styles.alignSelf(alignSelf),
    ])}>
    children
  </div>;
};
