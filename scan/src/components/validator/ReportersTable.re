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

let renderBody = (reserveIndex, reporterSub: ApolloHooks.Subscription.variant(Address.t)) => {
  <TBody
    key={
      switch (reporterSub) {
      | Data(address) => address |> Address.toBech32
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row alignItems=Row.Center minHeight={`px(30)}>
      <Col.Grid>
        {switch (reporterSub) {
         | Data(address) => <AddressRender address />
         | _ => <LoadingCensorBar width=300 height=15 />
         }}
      </Col.Grid>
    </Row>
  </TBody>;
};

let renderBodyMobile = (reserveIndex, reporterSub: ApolloHooks.Subscription.variant(Address.t)) => {
  switch (reporterSub) {
  | Data(address) =>
    <MobileCard
      values=InfoMobileCard.[("Reporter", Address(address, 200, `account))]
      key={address |> Address.toBech32}
      idx={address |> Address.toBech32}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[("Reporter", Loading(150))]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = (~address) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 5;
  let isMobile = Media.isMobile();

  let reportersSub = ReporterSub.getList(~operatorAddress=address, ~pageSize, ~page, ());
  let reporterCountSub = ReporterSub.count(address);
  let allSub = Sub.all2(reportersSub, reporterCountSub);

  <div className=Styles.tableWrapper>
    {isMobile
       ? <Row marginBottom=16>
           <Col.Grid>
             {switch (allSub) {
              | Data((_, reporterCount)) =>
                <div className={CssHelper.flexBox()}>
                  <Text
                    block=true
                    value={reporterCount |> string_of_int}
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                  <HSpacing size=Spacing.xs />
                  <Text block=true value="Reporters" weight=Text.Semibold color=Colors.gray7 />
                </div>
              | _ => <LoadingCensorBar width=100 height=15 />
              }}
           </Col.Grid>
         </Row>
       : <THead>
           <Row alignItems=Row.Center>
             <Col.Grid>
               {switch (allSub) {
                | Data((_, reporterCount)) =>
                  <div className={CssHelper.flexBox()}>
                    <Text
                      block=true
                      value={reporterCount |> string_of_int}
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                    <HSpacing size=Spacing.xs />
                    <Text block=true value="Reporters" weight=Text.Semibold color=Colors.gray7 />
                  </div>
                | _ => <LoadingCensorBar width=100 height=15 />
                }}
             </Col.Grid>
           </Row>
         </THead>}
    {switch (allSub) {
     | Data((reporters, reporterCount)) =>
       let pageCount = Page.getPageCount(reporterCount, pageSize);
       <>
         {reporterCount > 0
            ? reporters
              ->Belt_Array.mapWithIndex((i, e) =>
                  isMobile
                    ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
                )
              ->React.array
            : <EmptyContainer>
                <img src=Images.noBlock className=Styles.noDataImage />
                <Heading
                  size=Heading.H4
                  value="No Reporter"
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
