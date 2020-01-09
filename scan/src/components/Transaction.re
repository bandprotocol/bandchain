type txType =
  | DataRequest(string)
  | NewScript(string);

let txColor =
  fun
  | DataRequest(_) => Colors.purple
  | NewScript(_) => Colors.pink;

let txLabel =
  fun
  | DataRequest(_) => "DATA REQUEST"
  | NewScript(_) => "NEW SCRIPT";

let txSource =
  fun
  | DataRequest(source) => source
  | NewScript(source) => source;

let txIcon =
  fun
  | DataRequest(_) => Images.dataRequest
  | NewScript(_) => Images.newScript;

module Styles = {
  open Css;

  let typeContainer = style([marginRight(`px(20))]);

  let txTypeOval = txType =>
    style([
      marginLeft(`px(-2)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(15)),
      padding2(~v=Spacing.xs, ~h=Spacing.sm),
      color(txType->txColor),
      backgroundColor(
        switch (txType) {
        | DataRequest(_) => Colors.purpleLight
        | NewScript(_) => Colors.pinkLight
        },
      ),
    ]);

  let txIcon =
    style([
      width(`px(30)),
      height(`px(30)),
      marginTop(`px(5)),
      marginLeft(Spacing.xl),
      marginRight(Spacing.xl),
    ]);

  let hashCol = style([maxWidth(`px(250))]);
  let feeCol = style([maxWidth(`px(80))]);
};

let renderDataType = txType =>
  <div className=Styles.typeContainer>
    <div className={Styles.txTypeOval(txType)}>
      <Text value={txType->txLabel} size=Text.Xs block=true />
    </div>
    <VSpacing size=Spacing.xs />
    <Text value={txType->txSource} size=Text.Lg weight=Text.Semibold block=true />
  </div>;

let renderTxHash = (hash, time) => {
  <div className=Styles.hashCol>
    <VSpacing size={`px(9)} />
    <TimeAgos time />
    <VSpacing size={`px(6)} />
    <Text block=true code=true value=hash size=Text.Lg weight=Text.Bold ellipsis=true />
  </div>;
};

let renderFee = fee => {
  <div className=Styles.feeCol>
    <VSpacing size={`px(4)} />
    <Text size=Text.Sm block=true value="$0.002" color=Colors.grayText />
    <VSpacing size={`px(4)} />
    <Text value={fee->Js.Float.toString ++ " BAND"} color=Colors.grayHeader weight=Text.Semibold />
  </div>;
};
