type coin_amount_t = {
  value: list(Coin.t),
  hasDenom: bool,
};

type t =
  | Address(Address.t, int, bool)
  | Height(ID.Block.t)
  | Coin(coin_amount_t)
  | Count(int)
  | Float(float, option(int))
  | Percentage(float, option(int))
  | Timestamp(MomentRe.Moment.t)
  | TxHash(Hash.t, int)
  | Validator(Address.t, string, string)
  | Messages(Hash.t, list(TxSub.Msg.t), bool, string)
  | PubKey(PubKey.t)
  | Badge(TxSub.Msg.badge_theme_t)
  | Loading(int)
  | Text(string)
  | Nothing;

module Styles = {
  open Css;
  let vFlex = style([display(`flex), alignItems(`center)]);
  let addressContainer = w => {
    style([width(`px(w))]);
  };
  let badge = color =>
    style([
      display(`inlineFlex),
      padding2(~v=`px(5), ~h=`px(10)),
      backgroundColor(color),
      borderRadius(`px(15)),
    ]);
};

[@react.component]
let make = (~info) => {
  switch (info) {
  | Address(address, width, isValidator) =>
    <div className={Styles.addressContainer(width)}>
      <AddressRender address position=AddressRender.Text clickable=true validator=isValidator />
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
  | Float(value, digits) =>
    <Text
      value={value |> Format.fPretty(~digits?)}
      size=Text.Md
      spacing={Text.Em(0.02)}
      code=true
    />
  | Percentage(value, digits) =>
    <Text
      value={value |> Format.fPercent(~digits?)}
      size=Text.Md
      spacing={Text.Em(0.02)}
      code=true
    />
  | Coin({value, hasDenom}) =>
    <AmountRender coins=value pos={hasDenom ? AmountRender.TxIndex : Fee} />
  | Text(text) =>
    <Text
      value=text
      size=Text.Lg
      weight=Text.Semibold
      code=true
      spacing={Text.Em(0.02)}
      nowrap=true
      ellipsis=true
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
  | PubKey(publicKey) => <PubKeyRender pubKey=publicKey />
  | TxHash(txHash, width) => <TxLink txHash width size=Text.Lg />
  | Messages(txHash, messages, success, errMsg) =>
    <TxMessages txHash messages success errMsg width=360 />
  | Badge({text, textColor, bgColor}) =>
    <div className={Styles.badge(bgColor)}>
      <Text value=text size=Text.Xs spacing={Text.Em(0.07)} color=textColor />
    </div>
  | Loading(width) => <LoadingCensorBar width height=21 />
  | Nothing => React.null
  };
};
