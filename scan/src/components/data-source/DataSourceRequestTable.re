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

let renderBody = (reserveIndex, requestsSub: ApolloHooks.Subscription.variant(RequestSub.Mini.t)) => {
  <TBody.Grid
    key={
      switch (requestsSub) {
      | Data({id}) => id |> ID.Request.toString
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row.Grid alignItems=Row.Center>
      <Col.Grid col=Col.Two>
        {switch (requestsSub) {
         | Data({id}) => <TypeID.Request id />
         | _ => <LoadingCensorBar width=135 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Four>
        {switch (requestsSub) {
         | Data({oracleScriptID, oracleScriptName}) =>
           <div className={CssHelper.flexBox()}>
             <TypeID.OracleScript id=oracleScriptID />
             <HSpacing size=Spacing.sm />
             <Text value=oracleScriptName ellipsis=true />
           </div>
         | _ => <LoadingCensorBar width=270 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Three>
        {switch (requestsSub) {
         | Data({minCount, askCount, reportsCount}) =>
           <ProgressBar
             reportedValidators=reportsCount
             minimumValidators=minCount
             requestValidators=askCount
           />
         | _ => <LoadingCensorBar width=212 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.One>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (requestsSub) {
           | Data({resolveStatus}) => <RequestStatus resolveStatus />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (requestsSub) {
           | Data({txTimestamp}) =>
             <Timestamp.Grid
               time=txTimestamp
               size=Text.Md
               weight=Text.Regular
               textAlign=Text.Right
             />
           | _ =>
             <>
               <LoadingCensorBar width=70 height=15 />
               <LoadingCensorBar width=80 height=15 mt=5 />
             </>
           }}
        </div>
      </Col.Grid>
    </Row.Grid>
  </TBody.Grid>;
};

let renderBodyMobile =
    (reserveIndex, requestsSub: ApolloHooks.Subscription.variant(RequestSub.Mini.t)) => {
  switch (requestsSub) {
  | Data({
      id,
      txTimestamp,
      oracleScriptID,
      oracleScriptName,
      minCount,
      askCount,
      reportsCount,
      resolveStatus,
    }) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Request ID", RequestID(id)),
        ("Oracle Script", OracleScript(oracleScriptID, oracleScriptName)),
        (
          "Report Status",
          ProgressBar({
            reportedValidators: reportsCount,
            minimumValidators: minCount,
            requestValidators: askCount,
          }),
        ),
        ("Timestamp", Timestamp(txTimestamp)),
      ]
      key={id |> ID.Request.toString}
      idx={id |> ID.Request.toString}
      requestStatus=resolveStatus
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Request ID", Loading(70)),
        ("Oracle Script", Loading(136)),
        ("Report Status", Loading(20)),
        ("Timestamp", Loading(166)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = (~dataSourceID: ID.DataSource.t) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 5;

  let requestsSub = RequestSub.Mini.getListByDataSource(dataSourceID, ~pageSize, ~page, ());
  let totalRequestCountSub = RequestSub.countByDataSource(dataSourceID);

  let allSub = Sub.all2(requestsSub, totalRequestCountSub);

  let isMobile = Media.isMobile();

  <div className=Styles.tableWrapper>
    {isMobile
       ? <Row.Grid marginBottom=16>
           <Col.Grid>
             {switch (allSub) {
              | Data((_, totalRequestCount)) =>
                <div className={CssHelper.flexBox()}>
                  <Text
                    block=true
                    value={totalRequestCount |> Format.iPretty}
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                  <HSpacing size=Spacing.xs />
                  <Text block=true value="Requests" weight=Text.Semibold color=Colors.gray7 />
                </div>
              | _ => <LoadingCensorBar width=100 height=15 />
              }}
           </Col.Grid>
         </Row.Grid>
       : <THead.Grid>
           <Row.Grid alignItems=Row.Center>
             <Col.Grid col=Col.Two>
               {switch (allSub) {
                | Data((_, totalRequestCount)) =>
                  <div className={CssHelper.flexBox()}>
                    <Text
                      block=true
                      value={totalRequestCount |> Format.iPretty}
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                    <HSpacing size=Spacing.xs />
                    <Text block=true value="Requests" weight=Text.Semibold color=Colors.gray7 />
                  </div>
                | _ => <LoadingCensorBar width=100 height=15 />
                }}
             </Col.Grid>
             <Col.Grid col=Col.Four>
               <Text block=true value="Oracle Script" weight=Text.Semibold color=Colors.gray7 />
             </Col.Grid>
             <Col.Grid col=Col.Four>
               <Text
                 block=true
                 value="Report Status"
                 size=Text.Md
                 weight=Text.Semibold
                 color=Colors.gray7
               />
             </Col.Grid>
             <Col.Grid col=Col.Two>
               <Text
                 block=true
                 value="Timestamp"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Right
               />
             </Col.Grid>
           </Row.Grid>
         </THead.Grid>}
    {switch (allSub) {
     | Data((requests, requestsCount)) =>
       let pageCount = Page.getPageCount(requestsCount, pageSize);
       <>
         {requestsCount > 0
            ? requests
              ->Belt_Array.mapWithIndex((i, e) =>
                  isMobile
                    ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
                )
              ->React.array
            : <EmptyContainer>
                <img src=Images.noSource className=Styles.noDataImage />
                <Heading
                  size=Heading.H4
                  value="No Request"
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
