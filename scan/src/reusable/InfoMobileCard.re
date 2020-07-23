type t =
  | Address(Address.t, int)
  | Height(ID.Block.t)
  | Count(int)
  | Float(float)
  | Timestamp(MomentRe.Moment.t)
  | TxHash(Hash.t, int)
  | Validator(Address.t, string, string)
  | Messages(Hash.t, list(TxSub.Msg.t), bool, string)
  | Loading(int);

module Styles = {
  open Css;
  let vFlex = style([display(`flex), alignItems(`center)]);
  let addressContainer = w => {
    style([width(`px(w))]);
  };
};

[@react.component]
let make = (~info) => {
  switch (info) {
  | Address(address, width) =>
    <div className={Styles.addressContainer(width)}>
      <AddressRender address position=AddressRender.Text clickable=true />
    </div>
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
  | Float(value) =>
    <Text value={value |> Format.fPretty} size=Text.Md spacing={Text.Em(0.02)} code=true />
  | Timestamp(time) => <Timestamp time size=Text.Md weight=Text.Regular code=true />
  | Validator(address, moniker, identity) =>
    <ValidatorMonikerLink
      validatorAddress=address
      moniker
      size=Text.Md
      identity
      width={`percent(100.)}
    />
  | TxHash(txHash, width) => <TxLink txHash width size=Text.Lg />
  | Messages(txHash, messages, success, errMsg) =>
    <TxMessages txHash messages success errMsg width=360 />
  | Loading(width) => <LoadingCensorBar width height=21 />
  };
};
