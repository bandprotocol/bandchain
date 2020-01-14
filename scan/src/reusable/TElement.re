module Styles = {
  open Css;

  let typeContainer = style([marginRight(`px(20)), maxWidth(`px(210))]);

  let txTypeOval = (textColor, bgColor) =>
    style([
      marginLeft(`px(-2)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(15)),
      padding2(~v=Spacing.xs, ~h=Spacing.sm),
      color(textColor),
      backgroundColor(bgColor),
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

let txTypeMapping = msg => {
  switch (msg) {
  | TxHook.Msg.Request(_) => ("DATA REQUEST", Colors.darkBlue, Colors.lightBlue)
  | TxHook.Msg.Store(_) => ("NEW SCRIPT", Colors.darkGreen, Colors.lightGreen)
  | TxHook.Msg.Send(_) => ("SEND TOKEN", Colors.purple, Colors.lightPurple)
  | TxHook.Msg.Report(_) => ("DATA REPORT", Colors.darkIndigo, Colors.lightIndigo)
  | Unknown => ("Unknown", Colors.darkGrayText, Colors.grayHeader)
  };
};

let renderTxType = txType => {
  let (typeName, textColor, bgColor) = txTypeMapping(txType);
  <div className=Styles.typeContainer>
    <div className={Styles.txTypeOval(textColor, bgColor)}>
      <Text value=typeName size=Text.Xs block=true />
    </div>
  </div>;
};

let renderText = text =>
  <div className=Styles.typeContainer>
    <Text value=text size=Text.Lg weight=Text.Semibold block=true ellipsis=true />
  </div>;

let renderTxTypeWithDetail = txType => {
  let (typeName, textColor, bgColor) = txTypeMapping(txType);
  <div className=Styles.typeContainer>
    <div className={Styles.txTypeOval(textColor, bgColor)}>
      <Text value=typeName size=Text.Xs block=true />
    </div>
    <VSpacing size=Spacing.xs />
    <Text value="ETH/USD Price Feed" size=Text.Lg weight=Text.Semibold block=true ellipsis=true />
  </div>;
};

let renderTxHash = (hash, time) => {
  <div className=Styles.hashContainer>
    <TimeAgos time />
    <VSpacing size={`px(6)} />
    <Text
      block=true
      code=true
      value={hash |> Hash.toHex(~with0x=true)}
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

let renderFee = fee => {
  <div className=Styles.feeContainer>
    {fee == 0 ? React.null : <VSpacing size={`px(4)} />}
    {fee == 0 ? React.null : <Text size=Text.Sm block=true value="$0.002" color=Colors.grayText />}
    {fee == 0 ? React.null : <VSpacing size={`px(4)} />}
    <Text
      value={fee == 0 ? "FREE" : fee->Format.iPretty ++ " BAND"}
      color=Colors.grayHeader
      weight=Text.Semibold
    />
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
  | Unknown => Images.checkIcon;

type t =
  | Icon(TxHook.Msg.t)
  | Height(int)
  | Name(string)
  | Timestamp(MomentRe.Moment.t)
  | TxHash(Hash.t, MomentRe.Moment.t)
  | TxTypeWithDetail(TxHook.Msg.t)
  | TxType(TxHook.Msg.t)
  | Detail(string)
  | Status(string)
  | Fee(int)
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
  | TxTypeWithDetail(msg) => renderTxTypeWithDetail(msg)
  | TxType(msg) => renderTxType(msg)
  | Detail(detail) => renderText(detail)
  | Status(status) => renderText(status)
  | Fee(fee) => renderFee(fee)
  | Hash(hash) => renderHash(hash)
  | Address(address) => renderAddress(address)
  };
};
