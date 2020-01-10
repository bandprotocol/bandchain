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

  let hashCol = style([maxWidth(`px(250))]);
  let feeCol = style([maxWidth(`px(80))]);

  let textContainer = style([display(`flex)]);
};

let renderTxType = txType =>
  <div className=Styles.typeContainer>
    <div className=Styles.txTypeOval> <Text value="DATA REQUEST" size=Text.Xs block=true /> </div>
    <VSpacing size=Spacing.xs />
    <Text value="ETH/USD Price Feed" size=Text.Lg weight=Text.Semibold block=true ellipsis=true />
  </div>;

let renderTxHash = (hash, time) => {
  <div className=Styles.hashCol>
    <TimeAgos time />
    <VSpacing size={`px(6)} />
    <Text block=true code=true value=hash size=Text.Lg weight=Text.Bold ellipsis=true />
  </div>;
};

let renderHash = hash => {
  <div className=Styles.hashCol>
    <Text block=true code=true value=hash size=Text.Lg weight=Text.Bold ellipsis=true />
  </div>;
};

let renderFee = fee => {
  <div className=Styles.feeCol>
    <VSpacing size={`px(4)} />
    <Text size=Text.Sm block=true value="$0.002" color=Colors.grayText />
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
  | Timestamp(MomentRe.Moment.t)
  | TxHash(string, MomentRe.Moment.t)
  | TxType(list(TxHook.Msg.t))
  | Fee(int)
  | Hash(string);

[@react.component]
let make = (~elementType) => {
  switch (elementType) {
  | Icon(msg) => <img src={msg->msgIcon} className=Styles.msgIcon />
  | Height(height) => renderHeight(height)
  | TxHash(hash, timestamp) => renderTxHash(hash, timestamp)
  | TxType(msg) => renderTxType(msg)
  | Fee(fee) => renderFee(fee)
  | Hash(hash) => renderHash(hash)
  | _ => <div />
  };
};
