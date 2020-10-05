module Styles = {
  open Css;
  let height = he => style([height(he)]);
  let display = dp => style([display(dp)]);
  let justifyContent = jc => style([justifyContent(jc)]);
  let alignItems = ai => style([alignItems(ai)]);
  let flexDirection = fd => style([flexDirection(fd)]);
  let backgroundColor = bc => style([backgroundColor(bc)]);
};

[@react.component]
let make =
    (
      ~height=`px(300),
      ~display=`flex,
      ~justifyContent=`center,
      ~alignItems=`center,
      ~flexDirection=`column,
      ~backgroundColor=Colors.white,
      ~children,
    ) => {
  <div
    className={Css.merge([
      Styles.height(height),
      Styles.display(display),
      Styles.justifyContent(justifyContent),
      Styles.alignItems(alignItems),
      Styles.flexDirection(flexDirection),
      Styles.backgroundColor(backgroundColor),
    ])}>
    children
  </div>;
};
