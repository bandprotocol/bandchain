module Styles = {
  open Css;

  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

let renderBody = (reserveIndex, requestsSub: ApolloHooks.Subscription.variant(RequestSub.t)) => {
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
         | Data({oracleScript: {oracleScriptID, name}}) =>
           <div className={CssHelper.flexBox()}>
             <TypeID.OracleScript id=oracleScriptID />
             <HSpacing size=Spacing.sm />
             <Text value=name ellipsis=true />
           </div>
         | _ => <LoadingCensorBar width=270 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Three>
        {switch (requestsSub) {
         | Data({requestedValidators, minCount, reports}) =>
           <ProgressBar
             reportedValidators={reports |> Belt.Array.size}
             minimumValidators=minCount
             requestValidators={requestedValidators |> Belt.Array.size}
           />
         | _ => <LoadingCensorBar width=212 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.One>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (requestsSub) {
           | Data({resolveStatus}) => <RequestStatus.Sub resolveStatus />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (requestsSub) {
           | Data({transaction}) =>
             <Timestamp.Grid
               time={transaction.block.timestamp}
               size=Text.Md
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
    (reserveIndex, requestsSub: ApolloHooks.Subscription.variant(RequestSub.t)) => {
  switch (requestsSub) {
  | Data({
      id,
      transaction,
      oracleScript: {oracleScriptID, name},
      requestedValidators,
      minCount,
      reports,
      resolveStatus,
    }) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Request ID", RequestID(id)),
        ("Oracle Script", OracleScript(oracleScriptID, name)),
        (
          "Report Status",
          ProgressBar({
            reportedValidators: {
              reports |> Belt_Array.size;
            },
            minimumValidators: minCount,
            requestValidators: {
              requestedValidators |> Belt_Array.size;
            },
          }),
        ),
        ("Timestamp", Timestamp(transaction.block.timestamp)),
      ]
      key={id |> ID.Request.toString}
      idx={id |> ID.Request.toString}
      requestStatusSub=resolveStatus
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
let make = () => {
  let isMobile = Media.isMobile();

  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let latestRequestSub = RequestSub.getList(~pageSize=1, ~page=1, ());
  let requestsSub = RequestSub.getList(~pageSize, ~page, ());

  let allSub = Sub.all2(requestsSub, latestRequestSub);

  <Section>
    <div className=CssHelper.container id="requestsSection">
      <Row.Grid alignItems=Row.Center marginBottom=40 marginBottomSm=24>
        <Col.Grid col=Col.Twelve>
          <Heading value="All Requests" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
          {switch (latestRequestSub) {
           | Data(latestRequest) =>
             <Heading
               value={
                 latestRequest
                 ->Belt.Array.get(0)
                 ->Belt.Option.mapWithDefault(0, ({id}) => id |> ID.Request.toInt)
                 ->Format.iPretty
                 ++ " In total"
               }
               size=Heading.H3
             />
           | _ => <LoadingCensorBar width=65 height=21 />
           }}
        </Col.Grid>
      </Row.Grid>
      {isMobile
         ? React.null
         : <THead.Grid>
             <Row.Grid alignItems=Row.Center>
               <Col.Grid col=Col.Two>
                 <Text
                   block=true
                   value="Request ID"
                   size=Text.Md
                   weight=Text.Semibold
                   color=Colors.gray7
                 />
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
       | Data((requests, latestRequest)) =>
         let requestsCount =
           latestRequest
           ->Belt.Array.get(0)
           ->Belt.Option.mapWithDefault(0, ({id}) => id |> ID.Request.toInt);
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
    </div>
  </Section>;
};
