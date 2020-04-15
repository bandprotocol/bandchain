module Styles = {
  open Css;

  let container = w => style([display(`flex), cursor(`pointer), width(w)]);
};

[@react.component]
let make =
    (
      ~validatorAddress: Address.t,
      ~moniker: string,
      ~weight=Text.Regular,
      ~size=Text.Md,
      ~underline=false,
      ~width=`auto,
    ) => {
  <div
    className={Styles.container(width)}
    onClick={_ => Route.redirect(Route.ValidatorIndexPage(validatorAddress, ProposedBlocks))}>
    <Text
      value=moniker
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
