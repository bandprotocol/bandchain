module Styles = {
  open Css;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let hashContainer = style([maxWidth(`px(140))]);
  let statusContainer =
    style([maxWidth(`px(95)), display(`flex), flexDirection(`row), alignItems(`center)]);
  let logo = style([width(`px(20)), marginLeft(`auto), marginRight(`px(15))]);
};

let renderBody =
    (
      reserveIndex,
      txSub: ApolloHooks.Subscription.variant(TxSub.t),
      msgTransform: TxSub.Msg.t => TxSub.Msg.t,
    ) => {
  <TBody
    key={
      switch (txSub) {
      | Data({txHash}) => txHash |> Hash.toHex
      | _ => reserveIndex |> string_of_int
      }
    }>
    <Row minHeight={`px(30)} alignItems=`flexStart>
      <HSpacing size={`px(20)} />
      <Col size=1.6>
        <VSpacing size=Spacing.sm />
        {switch (txSub) {
         | Data({txHash}) => <TxLink txHash width=140 />
         | _ => <LoadingCensorBar width=140 height=15 />
         }}
      </Col>
      <Col size=0.88>
        <VSpacing size=Spacing.sm />
        {switch (txSub) {
         | Data({blockHeight}) => <TypeID.Block id=blockHeight />
         | _ => <LoadingCensorBar width=65 height=15 />
         }}
      </Col>
      <Col size=1.>
        <VSpacing size=Spacing.xs />
        {switch (txSub) {
         | Data({success}) =>
           <div className=Styles.statusContainer>
             <Text
               block=true
               code=true
               spacing={Text.Em(0.02)}
               value={success ? "success" : "fail"}
               weight=Text.Medium
               ellipsis=true
             />
             <img src={success ? Images.success : Images.fail} className=Styles.logo />
           </div>
         | _ =>
           <div className=Styles.statusContainer>
             <LoadingCensorBar width=55 height=15 />
             <HSpacing size=Spacing.sm />
             <LoadingCensorBar width=20 height=20 radius=20 />
           </div>
         }}
      </Col>
      <Col size=1.25>
        <VSpacing size=Spacing.sm />
        <div className=Styles.fullWidth>
          <AutoSpacing dir="left" />
          {switch (txSub) {
           | Data({gasFee}) =>
             <Text
               block=true
               code=true
               spacing={Text.Em(0.02)}
               value={gasFee->Coin.getBandAmountFromCoins->Format.fPretty}
               weight=Text.Medium
             />
           | _ => <LoadingCensorBar width=65 height=15 isRight=true />
           }}
          <HSpacing size={`px(20)} />
        </div>
      </Col>
      <Col size=5.>
        {switch (txSub) {
         | Data({messages, txHash, success, errMsg}) =>
           <div>
             <TxMessages
               txHash
               messages={messages->Belt_List.map(msgTransform)}
               success
               errMsg
               width=450
             />
           </div>
         | _ => <> <VSpacing size=Spacing.sm /> <LoadingCensorBar width=450 height=15 /> </>
         }}
      </Col>
      <HSpacing size={`px(20)} />
    </Row>
  </TBody>;
};

let renderBodyMobile =
    (
      reserveIndex,
      txSub: ApolloHooks.Subscription.variant(TxSub.t),
      msgTransform: TxSub.Msg.t => TxSub.Msg.t,
    ) => {
  switch (txSub) {
  | Data({txHash, blockHeight, gasFee, success, messages, errMsg}) =>
    let msgTransform = messages->Belt_List.map(msgTransform);
    <MobileCard
      values=InfoMobileCard.[
        ("TX HASH", TxHash(txHash, Media.isSmallMobile() ? 170 : 200)),
        ("BLOCK", Height(blockHeight)),
        ("GAS FEE\n(BAND)", Coin({value: gasFee, hasDenom: false})),
        ("ACTIONS", Messages(txHash, msgTransform, success, errMsg)),
      ]
      key={txHash |> Hash.toHex}
      idx={txHash |> Hash.toHex}
      status=success
    />;
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("TX HASH", Loading(Media.isSmallMobile() ? 170 : 200)),
        ("BLOCK", Loading(70)),
        (
          "ACTIONS",
          Loading(
            {
              Media.isSmallMobile() ? 160 : 230;
            },
          ),
        ),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make =
    (
      ~txsSub: ApolloHooks.Subscription.variant(array(TxSub.t)),
      ~msgTransform: TxSub.Msg.t => TxSub.Msg.t=x => x,
    ) => {
  let isMobile = Media.isMobile();
  <>
    {isMobile
       ? React.null
       : <THead>
           <Row>
             <HSpacing size={`px(20)} />
             <Col size=1.6>
               <div className=Styles.fullWidth>
                 <Text
                   value="TX HASH"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </div>
             </Col>
             <Col size=0.88>
               <div className=Styles.fullWidth>
                 <Text
                   value="BLOCK"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </div>
             </Col>
             <Col size=1.>
               <div className=Styles.fullWidth>
                 <Text
                   value="STATUS"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </div>
             </Col>
             <Col size=1.25>
               <div className=Styles.fullWidth>
                 <AutoSpacing dir="left" />
                 <Text
                   value="GAS FEE (BAND)"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
                 <HSpacing size={`px(20)} />
               </div>
             </Col>
             <Col size=5.>
               <div className=Styles.fullWidth>
                 <Text
                   value="ACTIONS"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </div>
             </Col>
             <HSpacing size={`px(20)} />
           </Row>
         </THead>}
    {switch (txsSub) {
     | Data(txs) =>
       txs
       ->Belt_Array.mapWithIndex((i, e) =>
           isMobile
             ? renderBodyMobile(i, Sub.resolve(e), msgTransform)
             : renderBody(i, Sub.resolve(e), msgTransform)
         )
       ->React.array
     | _ =>
       Belt_Array.make(10, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile
             ? renderBodyMobile(i, noData, msgTransform) : renderBody(i, noData, msgTransform)
         )
       ->React.array
     }}
  </>;
};
