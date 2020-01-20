module Styles = {
  open Css;

  let typeContainer = w => style([marginRight(`px(20)), width(w)]);

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

  let hashContainer = style([maxWidth(`px(220))]);
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
  <div className={Styles.typeContainer(`px(100))}>
    <div className={Styles.txTypeOval(textColor, bgColor)}>
      <Text value=typeName size=Text.Xs block=true />
    </div>
  </div>;
};

let renderText = (text, weight) =>
  <div className={Styles.typeContainer(`px(150))}>
    <Text value=text size=Text.Lg weight block=true ellipsis=true />
  </div>;

let renderSource = text =>
  <div className={Styles.typeContainer(`px(150))}>
    <Text value=text size=Text.Lg align=Text.Right block=true ellipsis=true />
  </div>;

let renderTxTypeWithDetail = (msg: TxHook.Msg.t) => {
  let (typeName, textColor, bgColor) = txTypeMapping(msg.action);
  <div className={Styles.typeContainer(`px(150))}>
    <div className={Styles.txTypeOval(textColor, bgColor)}>
      <Text value=typeName size=Text.Xs block=true />
    </div>
    <VSpacing size=Spacing.xs />
    <Text
      value={msg->TxHook.Msg.getDescription}
      size=Text.Lg
      weight=Text.Semibold
      block=true
      ellipsis=true
    />
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

let renderHashWithLink = hash => {
  <div className=Styles.hashContainer onClick={_ => hash->Route.TxIndexPage->Route.redirect}>
    <Text
      block=true
      code=true
      value={hash |> Hash.toHex}
      size=Text.Lg
      weight=Text.Bold
      ellipsis=true
      color=Colors.brightPurple
    />
  </div>;
};

let renderAddress = address => {
  <div className=Styles.hashContainer>
    <Text
      block=true
      code=true
      value={address |> Address.toBech32}
      size=Text.Lg
      weight=Text.Bold
      ellipsis=true
    />
  </div>;
};

let renderFee = fee => {
  <div className=Styles.feeContainer>
    {fee == 0.0 ? React.null : <VSpacing size={`px(4)} />}
    {fee == 0.0
       ? React.null : <Text size=Text.Sm block=true value="$0.002" color=Colors.grayText />}
    {fee == 0.0 ? React.null : <VSpacing size={`px(4)} />}
    <Text
      value={fee == 0.0 ? "FREE" : fee->Format.fPretty ++ " BAND"}
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
  | Fee(float)
  | Hash(Hash.t)
  | HashWithLink(Hash.t)
  | Address(Address.t)
  | Source(string)
  | Value(Js.Json.t);

[@react.component]
let make = (~elementType) => {
  switch (elementType) {
  | Icon({action, _}) => <img src={action->msgIcon} className=Styles.msgIcon />
  | Height(height) => renderHeight(height)
  | Name(name) => renderName(name)
  | Timestamp(time) => renderTime(time)
  | TxHash(hash, timestamp) => renderTxHash(hash, timestamp)
  | TxTypeWithDetail(msg) => renderTxTypeWithDetail(msg)
  | TxType({action, _}) => renderTxType(action)
  | Detail(detail) => renderText(detail, Text.Semibold)
  | Status(status) => renderText(status, Text.Semibold)
  | Fee(fee) => renderFee(fee)
  | Hash(hash) => renderHash(hash)
  | HashWithLink(hash) => renderHashWithLink(hash)
  | Address(address) => renderAddress(address)
  | Source(source) => renderSource(source)
  | Value(value) => renderText(value->Js.Json.stringify, Text.Regular)
  };
};
