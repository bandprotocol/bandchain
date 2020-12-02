module Styles = {
  open Css;

  let tableWrapper = style([Media.mobile([padding2(~v=`px(16), ~h=`zero)])]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
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
    paddingH={`px(24)}>
    <Row alignItems=Row.Center>
      <Col.Grid col=Col.Six>
        {switch (delegationsSub) {
         | Data({moniker, operatorAddress, identity}) =>
           <div className={CssHelper.flexBox()}>
             <ValidatorMonikerLink
               validatorAddress=operatorAddress
               moniker
               identity
               width={`px(300)}
             />
           </div>
         | _ => <LoadingCensorBar width=200 height=20 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Three>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (delegationsSub) {
           | Data({amount}) =>
             <Text value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty} />

           | _ => <LoadingCensorBar width=200 height=20 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Three>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (delegationsSub) {
           | Data({reward}) =>
             <Text value={reward |> Coin.getBandAmountFromCoin |> Format.fPretty} />
           | _ => <LoadingCensorBar width=200 height=20 />
           }}
        </div>
      </Col.Grid>
    </Row>
  </TBody>;
};

let renderBodyMobile =
    (reserveIndex, delegationsSub: ApolloHooks.Subscription.variant(DelegationSub.stake_t)) => {
  switch (delegationsSub) {
  | Data({amount, moniker, operatorAddress, reward, identity}) =>
    let key_ =
      (operatorAddress |> Address.toHex)
      ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
      ++ (reward |> Coin.getBandAmountFromCoin |> Js.Float.toString)
      ++ (reserveIndex |> string_of_int);
    <MobileCard
      values=InfoMobileCard.[
        ("Validator", Validator(operatorAddress, moniker, identity)),
        ("Amount\n(BAND)", Coin({value: [amount], hasDenom: false})),
        ("Reward\n(BAND)", Coin({value: [reward], hasDenom: false})),
      ]
      key=key_
      idx=key_
    />;
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Validator", Loading(230)),
        ("Amount\n(BAND)", Loading(100)),
        ("Reward\n(BAND)", Loading(100)),
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

  <div className=Styles.tableWrapper>
    {isMobile
       ? <Row marginBottom=16>
           <Col.Grid>
             {switch (delegationsCountSub) {
              | Data(delegationsCount) =>
                <div className={CssHelper.flexBox()}>
                  <Text
                    block=true
                    value={delegationsCount |> string_of_int}
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                  <HSpacing size=Spacing.xs />
                  <Text
                    block=true
                    value="Validators Delegated"
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                </div>
              | _ => <LoadingCensorBar width=100 height=15 />
              }}
           </Col.Grid>
         </Row>
       : <THead.Grid>
           <Row alignItems=Row.Center>
             <Col.Grid col=Col.Six>
               {switch (delegationsCountSub) {
                | Data(delegationsCount) =>
                  <div className={CssHelper.flexBox()}>
                    <Text
                      block=true
                      value={delegationsCount |> string_of_int}
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                    <HSpacing size=Spacing.xs />
                    <Text
                      block=true
                      value="Validators Delegated"
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                  </div>
                | _ => <LoadingCensorBar width=100 height=15 />
                }}
             </Col.Grid>
             <Col.Grid col=Col.Three>
               <Text
                 block=true
                 value="Amount (BAND)"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Right
               />
             </Col.Grid>
             <Col.Grid col=Col.Three>
               <Text
                 block=true
                 value="Reward (BAND)"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Right
               />
             </Col.Grid>
           </Row>
         </THead.Grid>}
    {switch (delegationsSub) {
     | Data(delegations) =>
       delegations->Belt.Array.size > 0
         ? delegations
           ->Belt_Array.mapWithIndex((i, e) =>
               isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
             )
           ->React.array
         : <EmptyContainer>
             <img src=Images.noBlock className=Styles.noDataImage />
             <Heading
               size=Heading.H4
               value="No Delegation"
               align=Heading.Center
               weight=Heading.Regular
               color=Colors.bandBlue
             />
           </EmptyContainer>
     | _ =>
       Belt_Array.make(pageSize, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
         )
       ->React.array
     }}
    {switch (delegationsCountSub) {
     | Data(delegationsCount) =>
       let pageCount = Page.getPageCount(delegationsCount, pageSize);
       <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />;
     | _ => React.null
     }}
  </div>;
};
