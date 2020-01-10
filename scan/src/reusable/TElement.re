module Styles = {
  open Css;

  let typeContainer = style([marginRight(`px(20)), maxWidth(`px(210))]);

  let txTypeOval =
    style([
      marginLeft(`px(-2)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(15)),
      padding2(~v=Spacing.xs, ~h=Spacing.sm),
      color(Colors.purple),
      backgroundColor(Colors.purpleLight),
    ]);

  let msgIcon =
    style([
      width(`px(30)),
      height(`px(30)),
      marginTop(`px(5)),
      marginLeft(Spacing.xl),
      marginRight(Spacing.xl),
    ]);

  let hashContainer = style([maxWidth(`px(250))]);
  let feeContainer = style([maxWidth(`px(80))]);
  let timeContainer = style([display(`flex), alignItems(`center), maxWidth(`px(150))]);
  let textContainer = style([display(`flex)]);
};

let renderTxType = txType =>
  <div className=Styles.typeContainer>
    <div className=Styles.txTypeOval> <Text value="DATA REQUEST" size=Text.Xs block=true /> </div>
    <VSpacing size=Spacing.xs />
    <Text value="ETH/USD Price Feed" size=Text.Lg weight=Text.Semibold block=true ellipsis=true />
  </div>;

let renderTxHash = (hash, time) => {
  <div className=Styles.hashContainer>
    <TimeAgos time />
    <VSpacing size={`px(6)} />
    <Text
      block=true
      code=true
      value={hash |> Hash.toHex}
      size=Text.Lg
      weight=Text.Bold
      ellipsis=true
    />
  </div>;
};

let renderHash = hash => {
  <div className=Styles.hashContainer>
    <Text
      block=true
      code=true
      value={hash |> Hash.toHex}
      size=Text.Lg
      weight=Text.Bold
      ellipsis=true
    />
  </div>;
};

let renderAddress = address => {
  <div className=Styles.hashContainer>
    <Text
      block=true
      code=true
      value={address |> Address.toHex}
      size=Text.Lg
      weight=Text.Bold
      ellipsis=true
    />
  </div>;
};

let renderFee = (fee, hasUsd) => {
  <div className=Styles.feeContainer>
    {hasUsd
       ? <>
           <VSpacing size={`px(4)} />
           <Text size=Text.Sm block=true value="$0.002" color=Colors.grayText />
         </>
       : React.null}
    <VSpacing size={`px(4)} />
    <Text value={fee->Format.iPretty ++ " BAND"} color=Colors.grayHeader weight=Text.Semibold />
  </div>;
};

let renderHeight = height => {
  <div className=Styles.textContainer>
    <Text value="#" size=Text.Md weight=Text.Semibold color=Colors.purple />
    <HSpacing size=Spacing.xs />
    <Text block=true value={height->Format.iPretty} size=Text.Md weight=Text.Semibold />
  </div>;
};

let renderName = name => {
  <div className=Styles.hashContainer>
    <Text block=true code=true value=name size=Text.Lg weight=Text.Bold ellipsis=true />
  </div>;
};

let renderTime = time => {
  <div className=Styles.timeContainer> <TimeAgos time size=Text.Md /> </div>;
};

let msgIcon =
  fun
  | TxHook.Msg.Store(_) => Images.newScript
  | Send(_) => Images.sendCoin
  | Request(_) => Images.dataRequest
  | Report(_) => Images.report
  | _ => Images.checkIcon;

type t =
  | Icon(TxHook.Msg.t)
  | Height(int)
  | Name(string)
  | Timestamp(MomentRe.Moment.t)
  | TxHash(Hash.t, MomentRe.Moment.t)
  | TxType(list(TxHook.Msg.t))
  | Fee(int, bool)
  | Hash(Hash.t)
  | Address(Address.t);

[@react.component]
let make = (~elementType) => {
  switch (elementType) {
  | Icon(msg) => <img src={msg->msgIcon} className=Styles.msgIcon />
  | Height(height) => renderHeight(height)
  | Name(name) => renderName(name)
  | Timestamp(time) => renderTime(time)
  | TxHash(hash, timestamp) => renderTxHash(hash, timestamp)
  | TxType(msg) => renderTxType(msg)
  | Fee(fee, hasUsd) => renderFee(fee, hasUsd)
  | Hash(hash) => renderHash(hash)
  | Address(address) => renderAddress(address)
  };
};
