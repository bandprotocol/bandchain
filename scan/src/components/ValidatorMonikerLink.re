module Styles = {
  open Css;

  let container = style([display(`flex), cursor(`pointer)]);
};

[@react.component]
let make = (~validator: ValidatorHook.Validator.t) => {
  <div
    className={Css.merge([Styles.container])}
    onClick={_ =>
      Route.redirect(Route.ValidatorIndexPage(validator.operatorAddress, ProposedBlocks))
    }>
    <Text
      value={validator.moniker}
      color=Colors.gray7
      code=true
      weight=Text.Regular
      spacing={Text.Em(0.02)}
      block=true
      size=Text.Md
      nowrap=true
      ellipsis=true
    />
  </div>;
};
