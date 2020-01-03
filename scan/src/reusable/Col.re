module Styles = {
  open Css;

  let col = style([margin4(~top=`px(0), ~right=Spacing.xs, ~left=Spacing.xs, ~bottom=`px(0))]);
  let colSize = sz => style([flex(`num(sz))]);
};

[@react.component]
let make = (~size=?, ~children) => {
  <div className={Css.merge([Styles.col, size->Belt.Option.mapWithDefault("", Styles.colSize)])}>
    children
  </div>;
};
