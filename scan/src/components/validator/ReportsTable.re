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
  let emptyContainer =
    style([
      height(`px(300)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      flexDirection(`column),
      backgroundColor(white),
    ]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);

  // DataSource Table

  let dataSourceTable = show => {
    style([
      padding2(~v=show ? `px(16) : `zero, ~h=`px(24)),
      marginTop(show ? `px(24) : `zero),
      backgroundColor(Colors.profileBG),
      transition(~duration=200, "all"),
      height(show ? `auto : `zero),
      opacity(show ? 1. : 0.),
      selector("> div + div", [paddingTop(`px(16))]),
    ]);
  };
  let toggle = style([cursor(`pointer)]);
};

module DataSourceItem = {
  [@react.component]
  let make = (~dataSource: ReportSub.ValidatorReport.report_details_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Two>
        <Text block=true value={dataSource.externalID} color=Colors.gray7 />
      </Col.Grid>
      <Col.Grid col=Col.Three>
        <div className={CssHelper.flexBox(~wrap=`nowrap, ())}>
          <TypeID.DataSource
            id={
                 let rawRequest = dataSource.rawRequest |> Belt_Option.getExn;
                 rawRequest.dataSource.dataSourceID;
               }
          />
          <HSpacing size=Spacing.sm />
          <Text
            value={
                    let rawRequest = dataSource.rawRequest |> Belt_Option.getExn;
                    rawRequest.dataSource.dataSourceName;
                  }
            ellipsis=true
          />
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <Text
          block=true
          value={
                  let rawRequest = dataSource.rawRequest |> Belt_Option.getExn;
                  rawRequest.calldata |> JsBuffer.toUTF8;
                }
          color=Colors.gray7
        />
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <Text block=true value={dataSource.exitCode} color=Colors.gray7 />
      </Col.Grid>
      <Col.Grid col=Col.Three>
        <Text
          block=true
          value={dataSource.data |> JsBuffer.toUTF8}
          align=Text.Right
          color=Colors.gray7
          ellipsis=true
        />
      </Col.Grid>
    </Row.Grid>;
  };
};

module RenderBody = {
  [@react.component]
  let make =
      (~reserveIndex, ~reportsSub: ApolloHooks.Subscription.variant(ReportSub.ValidatorReport.t)) => {
    let (show, setShow) = React.useState(_ => false);
    <TBody.Grid
      key={
        switch (reportsSub) {
        | Data({txHash}) => txHash |> Hash.toHex
        | _ => reserveIndex |> string_of_int
        }
      }
      paddingH={`px(24)}>
      <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
        <Col.Grid col=Col.Three>
          {switch (reportsSub) {
           | Data({request: {id}}) => <TypeID.Request id />
           | _ => <LoadingCensorBar width=135 height=15 />
           }}
        </Col.Grid>
        <Col.Grid col=Col.Four>
          {switch (reportsSub) {
           | Data({request: {oracleScript: {oracleScriptID, name}}}) =>
             <div className={CssHelper.flexBox()}>
               <TypeID.OracleScript id=oracleScriptID />
               <HSpacing size=Spacing.sm />
               <Text value=name ellipsis=true />
             </div>
           | _ => <LoadingCensorBar width=270 height=15 />
           }}
        </Col.Grid>
        <Col.Grid col=Col.Three>
          {switch (reportsSub) {
           | Data({txHash}) => <TxLink txHash width=140 />
           | _ => <LoadingCensorBar width=170 height=15 />
           }}
        </Col.Grid>
        <Col.Grid col=Col.Two>
          <div
            onClick={_ => setShow(prev => !prev)}
            className={Css.merge([CssHelper.flexBox(~justify=`flexEnd, ()), Styles.toggle])}>
            {switch (reportsSub) {
             | Data(_) =>
               <>
                 <Text
                   block=true
                   value={show ? "Hide Report" : "Show Report"}
                   weight=Text.Semibold
                   color=Colors.bandBlue
                 />
                 <HSpacing size=Spacing.xs />
                 <Icon
                   name={show ? "fas fa-caret-up" : "fas fa-caret-down"}
                   color=Colors.bandBlue
                 />
               </>
             | _ => <LoadingCensorBar width=100 height=15 />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      <div className={Styles.dataSourceTable(show)}>
        <Row.Grid>
          <Col.Grid col=Col.Two>
            <Text block=true value="External ID" weight=Text.Semibold color=Colors.gray7 />
          </Col.Grid>
          <Col.Grid col=Col.Three>
            <Text block=true value="Data Source" weight=Text.Semibold color=Colors.gray7 />
          </Col.Grid>
          <Col.Grid col=Col.Two>
            <Text block=true value="Param" weight=Text.Semibold color=Colors.gray7 />
          </Col.Grid>
          <Col.Grid col=Col.Two>
            <Text block=true value="Exit Code" weight=Text.Semibold color=Colors.gray7 />
          </Col.Grid>
          <Col.Grid col=Col.Three>
            <Text
              block=true
              value="Value"
              weight=Text.Semibold
              align=Text.Right
              color=Colors.gray7
            />
          </Col.Grid>
        </Row.Grid>
        {switch (reportsSub) {
         | Data({reportDetails}) =>
           reportDetails
           ->Belt_Array.mapWithIndex((i, reportDetail) =>
               <DataSourceItem
                 key={(i |> string_of_int) ++ reportDetail.externalID}
                 dataSource=reportDetail
               />
             )
           ->React.array
         | _ => <LoadingCensorBar width=170 height=50 />
         }}
      </div>
    </TBody.Grid>;
  };
};

module RenderBodyMobile = {
  [@react.component]
  let make =
      (~reserveIndex, ~reportsSub: ApolloHooks.Subscription.variant(ReportSub.ValidatorReport.t)) => {
    switch (reportsSub) {
    | Data({txHash, request: {id, oracleScript: {oracleScriptID, name}}, reportDetails}) =>
      <MobileCard
        values=InfoMobileCard.[
          ("Request ID", RequestID(id)),
          ("Oracle Script", OracleScript(oracleScriptID, name)),
          ("TX Hash", TxHash(txHash, Media.isSmallMobile() ? 170 : 200)),
        ]
        key={id |> ID.Request.toString}
        idx={id |> ID.Request.toString}
        panels={
          reportDetails
          ->Belt_Array.map(({externalID, exitCode, data, rawRequest: rawRequestOpt}) => {
              let ReportSub.ValidatorReport.{
                    dataSource: {dataSourceID, dataSourceName},
                    calldata,
                  } =
                rawRequestOpt->Belt.Option.getExn;
              InfoMobileCard.[
                ("External ID", Text(externalID)),
                ("Data Source", DataSource(dataSourceID, dataSourceName)),
                ("Param", Text(calldata |> JsBuffer.toUTF8)),
                ("Exit Code", Text(exitCode)),
                ("Value", Text(data |> JsBuffer.toUTF8)),
              ];
            })
          ->Belt.List.fromArray
        }
      />
    | _ =>
      <MobileCard
        values=InfoMobileCard.[
          ("Request ID", Loading(70)),
          ("Oracle Script", Loading(136)),
          ("TX Hash", Loading(Media.isSmallMobile() ? 170 : 200)),
        ]
        key={reserveIndex |> string_of_int}
        idx={reserveIndex |> string_of_int}
      />
    };
  };
};

[@react.component]
let make = (~address) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 5;

  let reportsSub =
    ReportSub.ValidatorReport.getListByValidator(
      ~page,
      ~pageSize,
      ~validator={
        address |> Address.toOperatorBech32;
      },
    );
  let reportsCountSub = ReportSub.ValidatorReport.count(address |> Address.toOperatorBech32);

  let allSub = Sub.all2(reportsSub, reportsCountSub);

  let isMobile = Media.isMobile();

  <div className=Styles.tableWrapper>
    {isMobile
       ? <Row.Grid marginBottom=16>
           <Col.Grid>
             {switch (allSub) {
              | Data((_, reportsCount)) =>
                <div className={CssHelper.flexBox()}>
                  <Text
                    block=true
                    value={reportsCount |> string_of_int}
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
             <Col.Grid col=Col.Three>
               {switch (allSub) {
                | Data((_, reportsCount)) =>
                  <div className={CssHelper.flexBox()}>
                    <Text
                      block=true
                      value={reportsCount |> string_of_int}
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                    <HSpacing size=Spacing.xs />
                    <Text
                      block=true
                      value="Oracle Reports"
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                  </div>
                | _ => <LoadingCensorBar width=100 height=15 />
                }}
             </Col.Grid>
             <Col.Grid col=Col.Four>
               <Text block=true value="Oracle Script" weight=Text.Semibold color=Colors.gray7 />
             </Col.Grid>
             <Col.Grid col=Col.Five>
               <Text block=true value="TX Hash" weight=Text.Semibold color=Colors.gray7 />
             </Col.Grid>
           </Row.Grid>
         </THead.Grid>}
    {switch (allSub) {
     | Data((reports, reportsCount)) =>
       let pageCount = Page.getPageCount(reportsCount, pageSize);
       <>
         {reportsCount > 0
            ? reports
              ->Belt_Array.mapWithIndex((i, e) =>
                  isMobile
                    ? <RenderBodyMobile
                        key={(i |> string_of_int) ++ (e.txHash |> Hash.toHex)}
                        reserveIndex=i
                        reportsSub={Sub.resolve(e)}
                      />
                    : <RenderBody
                        key={(i |> string_of_int) ++ (e.txHash |> Hash.toHex)}
                        reserveIndex=i
                        reportsSub={Sub.resolve(e)}
                      />
                )
              ->React.array
            : <div className=Styles.emptyContainer>
                <img src=Images.noSource className=Styles.noDataImage />
                <Heading
                  size=Heading.H4
                  value="No Report"
                  align=Heading.Center
                  weight=Heading.Regular
                  color=Colors.bandBlue
                />
              </div>}
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
           isMobile
             ? <RenderBodyMobile key={i |> string_of_int} reserveIndex=i reportsSub=noData />
             : <RenderBody key={i |> string_of_int} reserveIndex=i reportsSub=noData />
         )
       ->React.array
     }}
  </div>;
};
