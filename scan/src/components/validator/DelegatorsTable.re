module Styles = {
  open Css;

  let tableWrapper = style([Media.mobile([padding2(~v=`px(16), ~h=`zero)])]);
  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

let renderBody =
    (reserveIndex, delegatorSub: ApolloHooks.Subscription.variant(DelegationSub.stake_t)) => {
  <TBody
    key={
      switch (delegatorSub) {
      | Data({delegatorAddress}) => delegatorAddress |> Address.toBech32
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row alignItems=Row.Center minHeight={`px(30)}>
      <Col col=Col.Six>
        {switch (delegatorSub) {
         | Data({delegatorAddress}) => <AddressRender address=delegatorAddress />
         | _ => <LoadingCensorBar width=300 height=15 />
         }}
      </Col>
      <Col col=Col.Four>
        {switch (delegatorSub) {
         | Data({sharePercentage}) =>
           <Text block=true value={sharePercentage |> Format.fPretty} color=Colors.gray7 />
         | _ => <LoadingCensorBar width=100 height=15 />
         }}
      </Col>
      <Col col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (delegatorSub) {
           | Data({amount}) =>
             <Text
               block=true
               value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty}
               color=Colors.gray7
             />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
      </Col>
    </Row>
  </TBody>;
};

let renderBodyMobile =
    (reserveIndex, delegatorSub: ApolloHooks.Subscription.variant(DelegationSub.stake_t)) => {
  switch (delegatorSub) {
  | Data({amount, sharePercentage, delegatorAddress}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Delegator", Address(delegatorAddress, 149, `account)),
        ("Shares (%)", Float(sharePercentage, Some(4))),
        ("Amount\n(BAND)", Coin({value: [amount], hasDenom: false})),
      ]
      key={delegatorAddress |> Address.toBech32}
      idx={delegatorAddress |> Address.toBech32}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Delegator", Loading(150)),
        ("Shares (%)", Loading(60)),
        ("Amount\n(BAND)", Loading(80)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = (~address) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let delegatorsSub = DelegationSub.getDelegatorsByValidator(address, ~pageSize, ~page, ());
  let delegatorCountSub = DelegationSub.getDelegatorCountByValidator(address);

  let allSub = Sub.all2(delegatorsSub, delegatorCountSub);

  let isMobile = Media.isMobile();

  <div className=Styles.tableWrapper>
    {isMobile
       ? <Row marginBottom=16>
           <Col>
             {switch (allSub) {
              | Data((_, delegatorCount)) =>
                <div className={CssHelper.flexBox()}>
                  <Text
                    block=true
                    value={delegatorCount |> Format.iPretty}
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                  <HSpacing size=Spacing.xs />
                  <Text block=true value="Delegators" weight=Text.Semibold color=Colors.gray7 />
                </div>
              | _ => <LoadingCensorBar width=100 height=15 />
              }}
           </Col>
         </Row>
       : <THead>
           <Row alignItems=Row.Center>
             <Col col=Col.Six>
               {switch (allSub) {
                | Data((_, delegatorCount)) =>
                  <div className={CssHelper.flexBox()}>
                    <Text
                      block=true
                      value={delegatorCount |> Format.iPretty}
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                    <HSpacing size=Spacing.xs />
                    <Text block=true value="Delegators" weight=Text.Semibold color=Colors.gray7 />
                  </div>
                | _ => <LoadingCensorBar width=100 height=15 />
                }}
             </Col>
             <Col col=Col.Four>
               <Text block=true value="Share(%)" weight=Text.Semibold color=Colors.gray7 />
             </Col>
             <Col col=Col.Two>
               <Text
                 block=true
                 value="Amount"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Right
               />
             </Col>
           </Row>
         </THead>}
    {switch (allSub) {
     | Data((delegators, delegatorCount)) =>
       let pageCount = Page.getPageCount(delegatorCount, pageSize);
       <>
         {delegatorCount > 0
            ? delegators
              ->Belt_Array.mapWithIndex((i, e) =>
                  isMobile
                    ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
                )
              ->React.array
            : <EmptyContainer>
                <img src=Images.noBlock className=Styles.noDataImage />
                <Heading
                  size=Heading.H4
                  value="No Delegators"
                  align=Heading.Center
                  weight=Heading.Regular
                  color=Colors.bandBlue
                />
              </EmptyContainer>}
         {isMobile
            ? React.null
            : <Pagination
                currentPage=page
                pageCount
                onPageChange={newPage => setPage(_ => newPage)}
              />}
       </>;
     | _ =>
       Belt_Array.make(pageSize, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
         )
       ->React.array
     }}
  </div>;
};
