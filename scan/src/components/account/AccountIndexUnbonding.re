module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);

  let alignRight = style([display(`flex), justifyContent(`flexEnd)]);
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
      unbondingListSub: ApolloHooks.Subscription.variant(UnbondingSub.unbonding_list_t),
    ) => {
  <TBody
    key={
      switch (unbondingListSub) {
      | Data({validator: {operatorAddress}, amount, completionTime}) =>
        (operatorAddress |> Address.toBech32)
        ++ (completionTime |> MomentRe.Moment.toISOString)
        ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
      | _ => reserveIndex |> string_of_int
      }
    }
    minHeight=50>
    <Row>
      <Col> <HSpacing size=Spacing.lg /> </Col>
      <Col size=1.>
        {switch (unbondingListSub) {
         | Data({validator: {operatorAddress, moniker, identity}}) =>
           <div className=Styles.hFlex>
             <ValidatorMonikerLink
               validatorAddress=operatorAddress
               moniker
               identity
               width={`px(300)}
             />
           </div>
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col>
      <Col size=0.6>
        {switch (unbondingListSub) {
         | Data({amount}) =>
           <div className=Styles.alignRight>
             <Text value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty} code=true />
           </div>
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col>
      <Col size=1.>
        {switch (unbondingListSub) {
         | Data({completionTime}) =>
           <div className=Styles.alignRight>
             <Text
               value={
                 completionTime
                 |> MomentRe.Moment.format(Config.timestampDisplayFormat)
                 |> String.uppercase_ascii
               }
               code=true
             />
           </div>
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col>
      <Col> <HSpacing size=Spacing.lg /> </Col>
    </Row>
  </TBody>;
};

let renderBodyMobile =
    (
      reserveIndex,
      unbondingListSub: ApolloHooks.Subscription.variant(UnbondingSub.unbonding_list_t),
    ) => {
  switch (unbondingListSub) {
  | Data({validator: {operatorAddress, moniker, identity}, amount, completionTime}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("VALIDATOR", Validator(operatorAddress, moniker, identity)),
        ("AMOUNT\n(BAND)", Coin({value: [amount], hasDenom: false})),
        ("UNBONDED AT", Timestamp(completionTime)),
      ]
      key={
        (operatorAddress |> Address.toBech32)
        ++ (completionTime |> MomentRe.Moment.toISOString)
        ++ (reserveIndex |> string_of_int)
      }
      idx={
        (operatorAddress |> Address.toBech32)
        ++ (completionTime |> MomentRe.Moment.toISOString)
        ++ (reserveIndex |> string_of_int)
      }
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("VALIDATOR", Loading(230)),
        ("AMOUNT\n(BAND)", Loading(230)),
        ("UNBONDED AT", Loading(230)),
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

  let unbondingListSub =
    UnbondingSub.getUnbondingByDelegator(address, currentTime, ~pageSize, ~page, ());
  let unbondingCountSub = UnbondingSub.getUnbondingCountByDelegator(address, currentTime);

  <div className=Styles.tableLowerContainer>
    <VSpacing size=Spacing.md />
    <div className=Styles.hFlex>
      {switch (unbondingCountSub) {
       | Data(unbondingCount) =>
         <>
           <HSpacing size=Spacing.lg />
           <Text value={unbondingCount |> string_of_int} weight=Text.Semibold />
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
                   value="VALIDATOR"
                   size=Text.Sm
                   weight=Text.Bold
                   spacing={Text.Em(0.05)}
                   color=Colors.gray6
                 />
               </Col>
               <Col size=0.6>
                 <div className=Styles.alignRight>
                   <Text
                     block=true
                     value="AMOUNT (BAND)"
                     size=Text.Sm
                     weight=Text.Bold
                     spacing={Text.Em(0.05)}
                     color=Colors.gray6
                   />
                 </div>
               </Col>
               <Col size=1.>
                 <div className=Styles.alignRight>
                   <Text
                     block=true
                     value="UNBONDED AT"
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
      {switch (unbondingListSub) {
       | Data(unbondingList) =>
         unbondingList->Belt.Array.size > 0
           ? unbondingList
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
      {switch (unbondingCountSub) {
       | Data(unbondingCount) =>
         let pageCount = Page.getPageCount(unbondingCount, pageSize);
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
