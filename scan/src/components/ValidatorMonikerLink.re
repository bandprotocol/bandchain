module Styles = {
  open Css;

  let container = w =>
    style([display(`flex), cursor(`pointer), width(w), alignItems(`center)]);
};

[@react.component]
let make =
    (
      ~validatorAddress: Address.t,
      ~moniker: string,
      ~identity=?,
      ~weight=Text.Regular,
      ~size=Text.Md,
      ~underline=false,
      ~width=`auto,
    ) => {
  <Link
    className={Styles.container(width)}
    route={Route.ValidatorIndexPage(validatorAddress, ProposedBlocks)}>
    {switch (identity) {
     | Some(identity') =>
       <> <Avatar moniker identity=identity' /> <HSpacing size=Spacing.sm /> </>
     | None => React.null
     }}
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
  </Link>;
};
