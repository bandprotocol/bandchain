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
    (reserveIndex, delegationsSub: ApolloHooks.Subscription.variant(DelegationSub.stake_t)) => {
  <TBody
    key={
      switch (delegationsSub) {
      | Data({amount, operatorAddress, reward}) =>
        (operatorAddress |> Address.toHex)
        ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
        ++ (reward |> Coin.getBandAmountFromCoin |> Js.Float.toString)
        ++ (reserveIndex |> string_of_int)
      | _ => reserveIndex |> string_of_int
      }
    }
    minHeight=50>
    <Row>
      <Col> <HSpacing size=Spacing.lg /> </Col>
      <Col size=0.9>
        {switch (delegationsSub) {
         | Data({moniker, operatorAddress, identity}) =>
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
        {switch (delegationsSub) {
         | Data({amount}) =>
           <div className=Styles.alignRight>
             <Text value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty} code=true />
           </div>
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col>
      <Col size=0.6>
        {switch (delegationsSub) {
         | Data({reward}) =>
           <div className=Styles.alignRight>
             <Text value={reward |> Coin.getBandAmountFromCoin |> Format.fPretty} code=true />
           </div>
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col>
      <Col> <HSpacing size=Spacing.lg /> </Col>
    </Row>
  </TBody>;
};

let renderBodyMobile =
    (reserveIndex, delegationsSub: ApolloHooks.Subscription.variant(DelegationSub.stake_t)) => {
  switch (delegationsSub) {
  | Data({amount, moniker, operatorAddress, reward, identity}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("VALIDATOR", Validator(operatorAddress, moniker, identity)),
        ("AMOUNT\n(BAND)", Coin({value: [amount], hasDenom: false})),
        ("REWARD\n(BAND)", Coin({value: [reward], hasDenom: false})),
      ]
      key={
        (operatorAddress |> Address.toHex)
        ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
        ++ (reward |> Coin.getBandAmountFromCoin |> Js.Float.toString)
        ++ (reserveIndex |> string_of_int)
      }
      idx={
        (operatorAddress |> Address.toHex)
        ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
        ++ (reward |> Coin.getBandAmountFromCoin |> Js.Float.toString)
        ++ (reserveIndex |> string_of_int)
      }
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("VALIDATOR", Loading(230)),
        ("AMOUNT\n(BAND)", Loading(100)),
        ("REWARD\n(BAND)", Loading(100)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = (~address) => {
  let isMobile = Media.isMobile();
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;
  let delegationsCountSub = DelegationSub.getStakeCountByDelegator(address);
  let delegationsSub = DelegationSub.getStakeList(address, ~pageSize, ~page, ());

  <div className=Styles.tableLowerContainer>
    <VSpacing size=Spacing.md />
    <div className=Styles.hFlex>
      {switch (delegationsCountSub) {
       | Data(delegationsCount) =>
         <>
           <HSpacing size=Spacing.lg />
           <Text value={delegationsCount |> string_of_int} weight=Text.Semibold />
           <HSpacing size=Spacing.xs />
           <Text value="Validators Delegated" />
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
               <Col size=0.9>
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
               <Col size=0.6>
                 <div className=Styles.alignRight>
                   <Text
                     block=true
                     value="REWARD (BAND)"
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
      {switch (delegationsSub) {
       | Data(delegations) =>
         delegations->Belt.Array.size > 0
           ? delegations
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
      <VSpacing size=Spacing.lg />
      {switch (delegationsCountSub) {
       | Data(delegationsCount) =>
         let pageCount = Page.getPageCount(delegationsCount, pageSize);
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
