type t =
  | Height(ID.Block.t)
  | Count(int)
  | Timestamp(MomentRe.Moment.t)
  | Validator(Address.t, string, string)
  | Loading(int);

module Styles = {
  open Css;
  let vFlex = style([display(`flex), alignItems(`center)]);
};

[@react.component]
let make = (~info) => {
  switch (info) {
  | Height(height) =>
    <div className=Styles.vFlex> <TypeID.Block id=height position=TypeID.Subtitle /> </div>
  | Count(value) =>
    <Text
      value={value |> Format.iPretty}
      size=Text.Md
      weight=Text.Semibold
      spacing={Text.Em(0.02)}
      code=true
    />
  | Timestamp(time) => <Timestamp time size=Text.Md weight=Text.Regular code=true />
  | Validator(address, moniker, identity) =>
    <ValidatorMonikerLink
      validatorAddress=address
      moniker
      size=Text.Md
      identity
      width={`percent(100.)}
    />
  | Loading(width) => <LoadingCensorBar width height=21 />
  };
};