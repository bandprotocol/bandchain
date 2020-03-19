type alignment =
  | Start
  | Center
  | End;

module Styles = {
  open Css;

  let col = style([margin4(~top=`zero, ~right=Spacing.xs, ~left=Spacing.xs, ~bottom=`zero)]);
  let colSize = sz => style([flex(`num(sz))]);
  let alignSelf =
    Belt.Option.mapWithDefault(
      _,
      style([]),
      fun
      | Start => style([alignSelf(`flexStart)])
      | Center => style([alignSelf(`center)])
      | End => style([alignSelf(`flexEnd)]),
    );

  let justifyContent =
    Belt.Option.mapWithDefault(
      _,
      style([]),
      fun
      | Start => style([justifyContent(`flexStart)])
      | Center => style([justifyContent(`center)])
      | End => style([justifyContent(`flexEnd)]),
    );
  let alignItems =
    Belt.Option.mapWithDefault(
      _,
      style([]),
      fun
      | Start => style([alignItems(`flexStart)])
      | Center => style([alignItems(`center)])
      | End => style([alignItems(`flexEnd)]),
    );
};

[@react.component]
let make = (~size=?, ~alignSelf=?, ~alignItems=?, ~justifyContent=?, ~children) => {
  <div
    className={Css.merge([
      Styles.col,
      size->Belt.Option.mapWithDefault("", Styles.colSize),
      Styles.alignSelf(alignSelf),
      Styles.justifyContent(justifyContent),
      Styles.alignItems(alignItems),
    ])}>
    children
  </div>;
};
