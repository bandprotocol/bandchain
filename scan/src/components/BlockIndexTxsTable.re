module Styles = {
  open Css;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let container = style([width(`px(68))]);
  let hashContainer = style([maxWidth(`px(140)), cursor(`pointer)]);
  let paddingTopContainer = style([paddingTop(`px(5))]);
  let statusContainer =
    style([maxWidth(`px(95)), display(`flex), flexDirection(`row), alignItems(`center)]);
  let logo = style([width(`px(20)), marginLeft(`auto), marginRight(`px(15))]);
  let noTransactionLogo = style([width(`px(160))]);
  let emptyContainer =
    style([
      height(`px(300)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
    ]);
};

let renderBody = (reserveIndex, txSub: ApolloHooks.Subscription.variant(TxSub.t)) => {
  <TBody
    key={
      switch (txSub) {
      | Data({txHash}) => txHash |> Hash.toHex
      | _ => reserveIndex |> string_of_int
      }
    }>
    <Row>
      <HSpacing size={`px(20)} />
      <Col size=1.67 alignSelf=Col.Start>
        {switch (txSub) {
         | Data({txHash}) => <TxLink txHash width=140 />
         | _ => <LoadingCensorBar width=170 height=15 />
         }}
      </Col>
      <Col size=1.05 alignSelf=Col.Start>
        {switch (txSub) {
         | Data({gasFee}) =>
           <div className={Css.merge([Styles.fullWidth, Styles.paddingTopContainer])}>
             <AutoSpacing dir="left" />
             <Text
               block=true
               code=true
               spacing={Text.Em(0.02)}
               value={gasFee->Coin.getBandAmountFromCoins->Format.fPretty}
               weight=Text.Medium
             />
             <HSpacing size={`px(20)} />
           </div>
         | _ => <LoadingCensorBar width=30 height=15 isRight=true />
         }}
      </Col>
      <Col> <div className=Styles.container /> </Col>
      <Col size=5. alignSelf=Col.Start>
        {switch (txSub) {
         | Data({messages}) =>
           messages
           ->Belt.List.map(msg => {<> <Msg msg width=530 /> <VSpacing size=Spacing.md /> </>})
           ->Belt.List.toArray
           ->React.array
         | _ => <LoadingCensorBar width=530 height=15 />
         }}
      </Col>
      <HSpacing size={`px(20)} />
    </Row>
  </TBody>;
};

[@react.component]
let make = (~txsSub: ApolloHooks.Subscription.variant(array(TxSub.t))) => {
  <>
    <THead>
      <Row>
        <HSpacing size={`px(20)} />
        <Col size=1.67>
          <div className=Styles.fullWidth>
            <Text value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <Col size=1.05>
          <div className=Styles.fullWidth>
            <AutoSpacing dir="left" />
            <Text value="GAS FEE (BAND)" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
            <HSpacing size={`px(20)} />
          </div>
        </Col>
        <Col> <div className=Styles.container /> </Col>
        <Col size=5.>
          <div className=Styles.fullWidth>
            <Text value="ACTIONS" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <HSpacing size={`px(20)} />
      </Row>
    </THead>
    {switch (txsSub) {
     | Data(txs) =>
       txs->Belt.Array.size > 0
         ? txs->Belt_Array.mapWithIndex((i, e) => renderBody(i, Sub.resolve(e)))->React.array
         : <div className=Styles.emptyContainer>
             <img src=Images.noTransaction className=Styles.noTransactionLogo />
           </div>
     | _ =>
       Belt_Array.make(10, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) => renderBody(i, noData))
       ->React.array
     }}
  </>;
};
