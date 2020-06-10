module Styles = {
  open Css;
  let withWidth = w => style([display(`flex), maxWidth(`px(w)), cursor(`pointer)]);
};

[@react.component]
let make = (~txHash: Hash.t, ~width: int, ~size=Text.Md, ~weight=Text.Medium) => {
  <Link className={Styles.withWidth(width)} route={Route.TxIndexPage(txHash)}>
    <Text
      block=true
      code=true
      spacing={Text.Em(0.02)}
      value={txHash |> Hash.toHex(~upper=true)}
      weight
      ellipsis=true
      size
      color=Colors.gray7
    />
  </Link>;
};
