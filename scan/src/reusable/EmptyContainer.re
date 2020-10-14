module Styles = {
  open Css;
  let emptyContainer =
    style([
      justifyContent(`center),
      alignItems(`center),
      flexDirection(`column),
      width(`percent(100.)),
    ]);

  let height = he => style([height(he)]);
  let display = dp => style([display(dp ? `flex : `none)]);
  let backgroundColor = bc => style([backgroundColor(bc)]);
};

[@react.component]
let make = (~height=`px(300), ~display=true, ~backgroundColor=Colors.white, ~children) => {
  <div
    className={Css.merge([
      Styles.emptyContainer,
      Styles.height(height),
      Styles.display(display),
      Styles.backgroundColor(backgroundColor),
    ])}>
    children
  </div>;
};
