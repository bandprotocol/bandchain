module Styles = {
  open Css;

  let buttonContainer = active =>
    style([
      backgroundColor(active ? Colors.purple3 : Colors.white),
      border(`px(1), `solid, active ? Colors.purple3 : Colors.gray4),
      borderRadius(`px(6)),
      height(`px(35)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      padding2(~v=Spacing.xs, ~h=Spacing.md),
      cursor(`pointer),
    ]);
};

[@react.component]
let make = (~active, ~text, ~route) => {
  <div className={Styles.buttonContainer(active)} onClick={_ => route |> Route.redirect}>
    <Text
      value=text
      weight=Text.Semibold
      size=Text.Md
      color={active ? Colors.white : Colors.gray6}
    />
  </div>;
};
