type t =
  | BlockHeight(ID.Block.t)
  | BlockTXNCount(int)
  | Timestamp(MomentRe.Moment.t)
  | Validator(Address.t, string, string)
  | Loading(int);

module Styles = {
  open Css;

  let vFlex = style([display(`flex), alignItems(`center)]);
  let addressContainer = maxwidth_ => style([alignItems(`center), maxWidth(`px(maxwidth_))]);
  let datasourcesContainer = style([display(`flex), alignItems(`center), flexWrap(`wrap)]);
  let oracleScriptContainer = style([display(`flex), width(`px(240))]);
  let validatorsContainer = style([display(`flex), flexDirection(`column), flexWrap(`wrap)]);
  let marginRightOnly = size => style([marginRight(`px(size))]);
};

[@react.component]
let make = (~info) => {
  switch (info) {
  | BlockHeight(height) =>
    <div className=Styles.vFlex> <TypeID.Block id=height position=TypeID.Subtitle /> </div>
  | BlockTXNCount(value) =>
    <Text
      value={value |> Format.iPretty}
      size=Text.Md
      weight=Text.Semibold
      spacing={Text.Em(0.02)}
      code=true
    />
  | Timestamp(time) =>
    <Timestamp time=time size=Text.Md weight=Text.Regular code=true />
  | Validator(address, moniker, identity) =>
    <ValidatorMonikerLink
      validatorAddress=address
      moniker
      size=Text.Md
      identity
      width={`percent(100.)}
    />
  | Loading(width) => <LoadingCensorBar width height=22 />
  };
};
