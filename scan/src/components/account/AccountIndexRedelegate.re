module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);

  let alignRight = style([display(`flex), justifyContent(`flexEnd)]);
  let alignLeft = style([display(`flex), justifyContent(`flexStart)]);
  let emptyContainer =
    style([
      height(`px(300)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
    ]);
  let noTransactionLogo = style([width(`px(160))]);
};

let renderBody =
    (
      reserveIndex,
      redelegateListSub: ApolloHooks.Subscription.variant(RedelegateSub.redelegate_list_t),
    ) => {
  <TBody
    key={
      switch (redelegateListSub) {
      | Data({
          srcValidator: {operatorAddress: srcAddress},
          dstValidator: {operatorAddress: dstAddress},
          completionTime,
          amount,
        }) =>
        (srcAddress |> Address.toBech32)
        ++ (dstAddress |> Address.toBech32)
        ++ (completionTime |> MomentRe.Moment.toISOString)
        ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
      | _ => reserveIndex |> string_of_int
      }
    }
    minHeight=50>
    <Row>
      <Col> <HSpacing size=Spacing.lg /> </Col>
      <Col size=1.>
        {switch (redelegateListSub) {
         | Data({
             srcValidator: {
               operatorAddress: srcAddress,
               moniker: srcMoniker,
               identity: srcIdentity,
             },
           }) =>
           <ValidatorMonikerLink
             validatorAddress=srcAddress
             moniker=srcMoniker
             identity=srcIdentity
             width={`px(200)}
           />
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col>
      <Col size=1.>
        <div className=Styles.alignLeft>
          {switch (redelegateListSub) {
           | Data({
               dstValidator: {
                 operatorAddress: dstAddress,
                 moniker: dstMoniker,
                 identity: dstIdentity,
               },
             }) =>
             <ValidatorMonikerLink
               validatorAddress=dstAddress
               moniker=dstMoniker
               identity=dstIdentity
               width={`px(200)}
             />

           | _ => <LoadingCensorBar width=200 height=20 />
           }}
        </div>
      </Col>
      <Col size=0.6>
        {switch (redelegateListSub) {
         | Data({amount}) =>
           <div className=Styles.alignRight>
             <Text value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty} code=true />
           </div>
         | _ => <LoadingCensorBar width=145 height=20 />
         }}
      </Col>
      <Col size=1.>
        <div className=Styles.alignRight>
          {switch (redelegateListSub) {
           | Data({completionTime}) =>
             <Text
               value={
                 completionTime
                 |> MomentRe.Moment.format(Config.timestampDisplayFormat)
                 |> String.uppercase_ascii
               }
               code=true
             />

           | _ => <LoadingCensorBar width=200 height=20 />
           }}
        </div>
      </Col>
      <Col> <HSpacing size=Spacing.lg /> </Col>
    </Row>
  </TBody>;
};

let renderBodyMobile =
    (
      reserveIndex,
      redelegateListSub: ApolloHooks.Subscription.variant(RedelegateSub.redelegate_list_t),
    ) => {
  switch (redelegateListSub) {
  | Data({
      srcValidator: {operatorAddress: srcAddress, moniker: srcMoniker, identity: srcIdentity},
      dstValidator: {operatorAddress: dstAddress, moniker: dstMoniker, identity: dstIdentity},
      completionTime,
      amount,
    }) =>
    let key_ =
      (srcAddress |> Address.toBech32)
      ++ (dstAddress |> Address.toBech32)
      ++ (completionTime |> MomentRe.Moment.toISOString)
      ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
      ++ (reserveIndex |> string_of_int);
    <MobileCard
      values=InfoMobileCard.[
        ("SOURCE\nVALIDATOR", Validator(srcAddress, srcMoniker, srcIdentity)),
        ("DESTINATION\nVALIDATOR", Validator(dstAddress, dstMoniker, dstIdentity)),
        ("AMOUNT\n(BAND)", Coin({value: [amount], hasDenom: false})),
        ("REDELEGATE\nCOMPLETE AT", Timestamp(completionTime)),
      ]
      key=key_
      idx=key_
    />;
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("SOURCE\nVALIDATOR", Loading(230)),
        ("DESTINATION\nVALIDATOR", Loading(100)),
        ("AMOUNT\n(BAND)", Loading(100)),
        ("REDELEGATE\nCOMPLETE AT", Loading(230)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = (~address) => {
  let isMobile = Media.isMobile();
  let currentTime =
    React.useContext(TimeContext.context) |> MomentRe.Moment.format(Config.timestampUseFormat);

  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let redelegateCountSub = RedelegateSub.getRedelegateCountByDelegator(address, currentTime);
  let redelegateListSub =
    RedelegateSub.getRedelegationByDelegator(address, currentTime, ~pageSize, ~page, ());

  <div className=Styles.tableLowerContainer>
    <VSpacing size=Spacing.md />
    <div className=Styles.hFlex>
      {switch (redelegateCountSub) {
       | Data(redelegateCount) =>
         <>
           <HSpacing size=Spacing.lg />
           <Text value={redelegateCount |> string_of_int} weight=Text.Semibold />
           <HSpacing size=Spacing.xs />
           <Text value="Unbonding Entries" />
         </>
       | _ =>
         <div className=Styles.hFlex>
           <HSpacing size=Spacing.lg />
           <LoadingCensorBar width=130 height=15 />
         </div>
       }}
    </div>
    <VSpacing size=Spacing.lg />
    <>
      {isMobile
         ? React.null
         : <THead>
             <Row>
               <Col> <HSpacing size=Spacing.lg /> </Col>
               <Col size=1.>
                 <Text
                   block=true
                   value="SOURCE VALIDATOR"
                   size=Text.Sm
                   weight=Text.Bold
                   spacing={Text.Em(0.05)}
                   color=Colors.gray6
                 />
               </Col>
               <Col size=1.>
                 <div className=Styles.alignLeft>
                   <Text
                     block=true
                     value="DESTINATION VALIDATOR"
                     size=Text.Sm
                     weight=Text.Bold
                     spacing={Text.Em(0.05)}
                     color=Colors.gray6
                   />
                 </div>
               </Col>
               <Col size=0.6>
                 <div className=Styles.alignRight>
                   <Text
                     block=true
                     value="AMOUNT (BAND)"
                     size=Text.Sm
                     spacing={Text.Em(0.05)}
                     weight=Text.Bold
                     color=Colors.gray6
                   />
                 </div>
               </Col>
               <Col size=1.>
                 <div className=Styles.alignRight>
                   <Text
                     block=true
                     value="REDELEGATE COMPLETE AT"
                     size=Text.Sm
                     spacing={Text.Em(0.05)}
                     weight=Text.Bold
                     color=Colors.gray6
                   />
                 </div>
               </Col>
               <Col> <HSpacing size=Spacing.lg /> </Col>
             </Row>
           </THead>}
      {switch (redelegateListSub) {
       | Data(redelegateList) =>
         redelegateList->Belt.Array.size > 0
           ? redelegateList
             ->Belt_Array.mapWithIndex((i, e) =>
                 isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
               )
             ->React.array
           : <div className=Styles.emptyContainer>
               <img src=Images.noTransaction className=Styles.noTransactionLogo />
             </div>
       | _ =>
         Belt_Array.make(1, ApolloHooks.Subscription.NoData)
         ->Belt_Array.mapWithIndex((i, noData) =>
             isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
           )
         ->React.array
       }}
      {switch (redelegateCountSub) {
       | Data(redelegateCount) =>
         let pageCount = Page.getPageCount(redelegateCount, pageSize);
         <>
           <VSpacing size=Spacing.lg />
           <Pagination
             currentPage=page
             pageCount
             onPageChange={newPage => setPage(_ => newPage)}
           />
           <VSpacing size=Spacing.lg />
         </>;
       | _ => React.null
       }}
    </>
  </div>;
};
