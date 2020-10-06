module Styles = {
  open Css;

  let tableWrapper =
    style([
      Media.mobile([
        padding2(~v=`px(16), ~h=`px(12)),
        backgroundColor(Colors.profileBG),
        margin2(~v=`zero, ~h=`px(-12)),
      ]),
    ]);
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

let renderBody = (reserveIndex, depositSub: ApolloHooks.Subscription.variant(DepositSub.t)) => {
  <TBody.Grid
    key={
      switch (depositSub) {
      | Data({depositor}) => depositor |> Address.toBech32
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
      <Col.Grid col=Col.Five>
        {switch (depositSub) {
         | Data({depositor}) => <AddressRender address=depositor />
         | _ => <LoadingCensorBar width=300 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Five>
        {switch (depositSub) {
         | Data({txHash}) => <TxLink txHash width=240 />
         | _ => <LoadingCensorBar width=100 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (depositSub) {
           | Data({amount}) =>
             <Text
               block=true
               value={amount |> Coin.getBandAmountFromCoins |> Format.fPretty(~digits=6)}
               color=Colors.gray7
             />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
      </Col.Grid>
    </Row.Grid>
  </TBody.Grid>;
};

let renderBodyMobile = (reserveIndex, depositSub: ApolloHooks.Subscription.variant(DepositSub.t)) => {
  switch (depositSub) {
  | Data({depositor, txHash, amount}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Depositor", Address(depositor, 200, `account)),
        ("TX Hash", TxHash(txHash, 200)),
        ("Amount", Coin({value: amount, hasDenom: false})),
      ]
      key={depositor |> Address.toBech32}
      idx={depositor |> Address.toBech32}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Depositor", Loading(200)),
        ("TX Hash", Loading(200)),
        ("Amount", Loading(80)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = (~proposalID) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 5;
  let isMobile = Media.isMobile();

  let depositsSub = DepositSub.getList(proposalID, ~pageSize, ~page, ());
  let depositCountSub = DepositSub.count(proposalID);
  let allSub = Sub.all2(depositsSub, depositCountSub);

  <div className=Styles.tableWrapper>
    {isMobile
       ? <Row.Grid marginBottom=16>
           <Col.Grid>
             {switch (allSub) {
              | Data((_, depositCount)) =>
                <div className={CssHelper.flexBox()}>
                  <Text
                    block=true
                    value={depositCount |> string_of_int}
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                  <HSpacing size=Spacing.xs />
                  <Text
                    block=true
                    value={depositCount > 1 ? "Depositors" : "Depositor"}
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                </div>
              | _ => <LoadingCensorBar width=100 height=15 />
              }}
           </Col.Grid>
         </Row.Grid>
       : <THead.Grid>
           <Row.Grid alignItems=Row.Center>
             <Col.Grid col=Col.Five>
               {switch (allSub) {
                | Data((_, depositCount)) =>
                  <div className={CssHelper.flexBox()}>
                    <Text
                      block=true
                      value={depositCount |> string_of_int}
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                    <HSpacing size=Spacing.xs />
                    <Text block=true value="Depositors" weight=Text.Semibold color=Colors.gray7 />
                  </div>
                | _ => <LoadingCensorBar width=100 height=15 />
                }}
             </Col.Grid>
             <Col.Grid col=Col.Five>
               <Text block=true value="TX Hash" weight=Text.Semibold color=Colors.gray7 />
             </Col.Grid>
             <Col.Grid col=Col.Two>
               <Text
                 block=true
                 value="Amount"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Right
               />
             </Col.Grid>
           </Row.Grid>
         </THead.Grid>}
    {switch (allSub) {
     | Data((delegators, depositCount)) =>
       let pageCount = Page.getPageCount(depositCount, pageSize);
       <>
         {depositCount > 0
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
                  value="No Depositors"
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
