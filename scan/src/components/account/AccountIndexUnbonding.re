module Styles = {
  open Css;

  let tableWrapper = style([Media.mobile([padding2(~v=`px(16), ~h=`zero)])]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

let renderBody =
    (
      reserveIndex,
      unbondingListSub: ApolloHooks.Subscription.variant(UnbondingSub.unbonding_list_t),
    ) => {
  <TBody.Grid
    key={
      switch (unbondingListSub) {
      | Data({validator: {operatorAddress}, amount, completionTime}) =>
        (operatorAddress |> Address.toBech32)
        ++ (completionTime |> MomentRe.Moment.toISOString)
        ++ (amount |> Coin.getBandAmountFromCoin |> Js.Float.toString)
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row alignItems=Row.Center>
      <Col.Grid col=Col.Six>
        {switch (unbondingListSub) {
         | Data({validator: {operatorAddress, moniker, identity}}) =>
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
          {switch (unbondingListSub) {
           | Data({amount}) =>
             <Text value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty} />
           | _ => <LoadingCensorBar width=200 height=20 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Three>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (unbondingListSub) {
           | Data({completionTime}) =>
             <Timestamp.Grid
               time=completionTime
               size=Text.Md
               weight=Text.Regular
               textAlign=Text.Right
             />
           | _ => <LoadingCensorBar width=200 height=20 />
           }}
        </div>
      </Col.Grid>
    </Row>
  </TBody.Grid>;
};

let renderBodyMobile =
    (
      reserveIndex,
      unbondingListSub: ApolloHooks.Subscription.variant(UnbondingSub.unbonding_list_t),
    ) => {
  switch (unbondingListSub) {
  | Data({validator: {operatorAddress, moniker, identity}, amount, completionTime}) =>
    let key_ =
      (operatorAddress |> Address.toBech32)
      ++ (completionTime |> MomentRe.Moment.toISOString)
      ++ (reserveIndex |> string_of_int);
    <MobileCard
      values=InfoMobileCard.[
        ("Validator", Validator(operatorAddress, moniker, identity)),
        ("Amount\n(BAND)", Coin({value: [amount], hasDenom: false})),
        ("Unbonded At", Timestamp(completionTime)),
      ]
      key=key_
      idx=key_
    />;
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Validator", Loading(230)),
        ("Amount\n(BAND)", Loading(230)),
        ("Unbonded At", Loading(230)),
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

  <div className=Styles.tableWrapper>
    {isMobile
       ? <Row marginBottom=16>
           <Col.Grid>
             {switch (unbondingCountSub) {
              | Data(unbondingCount) =>
                <div className={CssHelper.flexBox()}>
                  <Text
                    block=true
                    value={unbondingCount |> string_of_int}
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                  <HSpacing size=Spacing.xs />
                  <Text
                    block=true
                    value="Unbonding Entries"
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
               {switch (unbondingCountSub) {
                | Data(unbondingCount) =>
                  <div className={CssHelper.flexBox()}>
                    <Text
                      block=true
                      value={unbondingCount |> string_of_int}
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                    <HSpacing size=Spacing.xs />
                    <Text
                      block=true
                      value="Unbonding Entries"
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
                 value="Unbonded At"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Right
               />
             </Col.Grid>
           </Row>
         </THead.Grid>}
    {switch (unbondingListSub) {
     | Data(unbondingList) =>
       unbondingList->Belt.Array.size > 0
         ? unbondingList
           ->Belt_Array.mapWithIndex((i, e) =>
               isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
             )
           ->React.array
         : <EmptyContainer>
             <img src=Images.noBlock className=Styles.noDataImage />
             <Heading
               size=Heading.H4
               value="No Unbonding"
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
    {switch (unbondingCountSub) {
     | Data(unbondingCount) =>
       let pageCount = Page.getPageCount(unbondingCount, pageSize);
       <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />;
     | _ => React.null
     }}
  </div>;
};
