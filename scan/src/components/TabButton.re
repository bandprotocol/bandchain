module Styles = {
  open Css;

  let buttonContainer = active =>
    style([
      backgroundColor(active ? Colors.brightPurple : Colors.white),
      border(`px(1), `solid, active ? Colors.brightPurple : Colors.lightGray),
      borderRadius(`px(6)),
      height(`px(35)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      padding2(~v=Spacing.xs, ~h=Spacing.sm),
    ]);
  // let buttonInner = style([position(`absolute)]);
};

[@react.component]
let make = (~active, ~text) => {
  <div className={Styles.buttonContainer(active)}>
    <Text
      value=text
      weight=Text.Semibold
      size=Text.Md
      color={active ? Colors.white : Colors.darkGrayText}
    />
  </div>;
};
