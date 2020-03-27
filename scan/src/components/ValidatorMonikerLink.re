module Styles = {
  open Css;

  let container = style([display(`flex), cursor(`pointer)]);
};

[@react.component]
let make =
    (~validator: ValidatorHook.Validator.t, ~weight=Text.Regular, ~size=Text.Md, ~underline=false) => {
  <div
    className={Css.merge([Styles.container])}
    onClick={_ =>
      Route.redirect(Route.ValidatorIndexPage(validator.operatorAddress, ProposedBlocks))
    }>
    <Text
      value={validator.moniker}
      color=Colors.gray7
      code=true
      weight
      spacing={Text.Em(0.02)}
      block=true
      size
      nowrap=true
      ellipsis=true
      underline
    />
  </div>;
};
