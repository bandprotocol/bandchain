module Styles = {
  open Css;

  let col =
    style([
      margin4(
        ~top=`px(0),
        ~right=Spacing.xs,
        ~left=Spacing.xs,
        ~bottom=`px(0),
      ),
    ]);
  let colSize = sz => style([flex(`num(sz))]);
};

[@react.component]
let make = (~size=?, ~nogrow=?, ~children) => {
  <div className={Cn.make([
    Styles.col, 
    size->Cn.mapSome(Styles.colSize), 
  ])}>
    children
  </div>;
};