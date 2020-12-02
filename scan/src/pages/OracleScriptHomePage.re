module Styles = {
  open Css;
  let mostRequestCard =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding3(~top=`px(24), ~h=`px(24), ~bottom=`px(16)),
      height(`calc((`sub, `percent(100.), `px(24)))),
      marginBottom(`px(24)),
    ]);
  let requestResponseBox = style([flexShrink(0.), flexGrow(0.), flexBasis(`percent(50.))]);
  let descriptionBox = style([minHeight(`px(36)), margin2(~v=`px(16), ~h=`zero)]);
  let idBox = style([marginBottom(`px(4))]);
  let tbodyContainer = style([minHeight(`px(600))]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

type sort_by_t =
  | MostRequested
  | LatestUpdate;

let getName =
  fun
  | MostRequested => "Most Requested"
  | LatestUpdate => "Latest Update";

let defaultCompare = (a: OracleScriptSub.t, b: OracleScriptSub.t) =>
  if (a.timestamp != b.timestamp) {
    compare(b.id |> ID.OracleScript.toInt, a.id |> ID.OracleScript.toInt);
  } else {
    compare(b.requestCount, a.requestCount);
  };

let sorting = (oracleSctipts: array(OracleScriptSub.t), sortedBy) => {
  oracleSctipts
  ->Belt.List.fromArray
  ->Belt.List.sort((a, b) => {
      let result = {
        switch (sortedBy) {
        | MostRequested => compare(b.requestCount, a.requestCount)
        | LatestUpdate => compare(b.timestamp, a.timestamp)
        };
      };
      if (result != 0) {
        result;
      } else {
        defaultCompare(a, b);
      };
    })
  ->Belt.List.toArray;
};

let renderMostRequestedCard =
    (reserveIndex, oracleScriptSub: ApolloHooks.Subscription.variant(OracleScriptSub.t)) => {
  <Col.Grid
    key={
      switch (oracleScriptSub) {
      | Data({id}) => id |> ID.OracleScript.toString
      | _ => reserveIndex |> string_of_int
      }
    }
    col=Col.Four>
    <div
      className={Css.merge([
        Styles.mostRequestCard,
        CssHelper.flexBox(~direction=`column, ~justify=`spaceBetween, ~align=`stretch, ()),
      ])}>
      <div className=Styles.idBox>
        {switch (oracleScriptSub) {
         | Data({id}) => <TypeID.OracleScript id position=TypeID.Subtitle />
         | _ => <LoadingCensorBar width=40 height=15 />
         }}
      </div>
      {switch (oracleScriptSub) {
       | Data({name}) => <Heading size=Heading.H4 value=name />
       | _ => <LoadingCensorBar width=200 height=15 />
       }}
      <div className=Styles.descriptionBox>
        {switch (oracleScriptSub) {
         | Data({description}) =>
           let text = Ellipsis.format(~text=description, ~limit=70, ());
           <Text size=Text.Lg value=text weight=Text.Regular block=true />;
         | _ => <LoadingCensorBar width=250 height=15 />
         }}
      </div>
      <div className={CssHelper.flexBox()}>
        <div className=Styles.requestResponseBox>
          <Heading size=Heading.H5 value="Requests" marginBottom=8 />
          {switch (oracleScriptSub) {
           | Data({requestCount}) =>
             <Text
               size=Text.Lg
               value={requestCount |> Format.iPretty}
               weight=Text.Regular
               block=true
             />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
        <div className=Styles.requestResponseBox>
          <Heading size=Heading.H5 value="Response time" marginBottom=8 />
          {switch (oracleScriptSub) {
           | Data({responseTime: responseTimeOpt}) =>
             switch (responseTimeOpt) {
             | Some(responseTime') =>
               <Text
                 size=Text.Lg
                 value={(responseTime' |> Format.fPretty(~digits=2)) ++ " s"}
                 weight=Text.Regular
                 block=true
               />
             | None => <Text value="TBD" />
             }

           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
      </div>
    </div>
  </Col.Grid>;
};

let renderBody =
    (reserveIndex, oracleScriptSub: ApolloHooks.Subscription.variant(OracleScriptSub.t)) => {
  <TBody
    key={
      switch (oracleScriptSub) {
      | Data({id}) => id |> ID.OracleScript.toString
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row alignItems=Row.Center>
      <Col.Grid col=Col.Five>
        {switch (oracleScriptSub) {
         | Data({id, name}) =>
           <div className={CssHelper.flexBox()}>
             <TypeID.OracleScript id />
             <HSpacing size=Spacing.sm />
             <Text value=name ellipsis=true />
           </div>
         | _ => <LoadingCensorBar width=300 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Three>
        {switch (oracleScriptSub) {
         | Data({description}) =>
           let text = Ellipsis.format(~text=description, ~limit=70, ());
           <Text value=text block=true />;
         | _ => <LoadingCensorBar width=270 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div
          className={CssHelper.flexBox(
            ~justify=`flexStart,
            ~align=`flexEnd,
            ~direction=`column,
            (),
          )}>
          {switch (oracleScriptSub) {
           | Data({requestCount, responseTime: responseTimeOpt}) =>
             <>
               <div>
                 <Text
                   value={requestCount |> Format.iPretty}
                   weight=Text.Medium
                   block=true
                   ellipsis=true
                   align=Text.Right
                 />
               </div>
               <div>
                 <Text
                   value={
                     switch (responseTimeOpt) {
                     | Some(responseTime') =>
                       "(" ++ (responseTime' |> Format.fPretty(~digits=2)) ++ " s)"
                     | None => "(TBD)"
                     }
                   }
                   weight=Text.Medium
                   block=true
                   color=Colors.gray6
                   align=Text.Right
                 />
               </div>
             </>
           | _ => <LoadingCensorBar width=70 height=15 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (oracleScriptSub) {
           | Data({timestamp: timestampOpt}) =>
             switch (timestampOpt) {
             | Some(timestamp') =>
               <Timestamp.Grid
                 time=timestamp'
                 size=Text.Md
                 weight=Text.Regular
                 textAlign=Text.Right
               />
             | None => <Text value="Genesis" />
             }
           | _ =>
             <>
               <LoadingCensorBar width=70 height=15 />
               <LoadingCensorBar width=80 height=15 mt=5 />
             </>
           }}
        </div>
      </Col.Grid>
    </Row>
  </TBody>;
};

let renderBodyMobile =
    (reserveIndex, oracleScriptSub: ApolloHooks.Subscription.variant(OracleScriptSub.t)) => {
  switch (oracleScriptSub) {
  | Data({id, timestamp: timestampOpt, description, name, requestCount, responseTime}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Oracle Script", OracleScript(id, name)),
        ("Description", Text(description)),
        ("Request&\nResponse time", RequestResponse({requestCount, responseTime})),
        (
          "Timestamp",
          switch (timestampOpt) {
          | Some(timestamp') => Timestamp(timestamp')
          | None => Text("Genesis")
          },
        ),
      ]
      key={id |> ID.OracleScript.toString}
      idx={id |> ID.OracleScript.toString}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Oracle Script", Loading(200)),
        ("Description", Loading(200)),
        ("Request&\nResponse time", Loading(80)),
        ("Timestamp", Loading(180)),
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
  let (searchTerm, setSearchTerm) = React.useState(_ => "");
  let (sortedBy, setSortedBy) = React.useState(_ => LatestUpdate);

  React.useEffect1(
    () => {
      if (searchTerm != "") {
        setPage(_ => 1);
      };
      None;
    },
    [|searchTerm|],
  );

  let pageSize = 10;
  let mostRequestedPageSize = isMobile ? 3 : 6;
  let oracleScriptsCountSub = OracleScriptSub.count(~searchTerm, ());
  let oracleScriptsSub = OracleScriptSub.getList(~pageSize, ~page, ~searchTerm, ());

  let mostRequestedOracleScriptSub =
    OracleScriptSub.getList(~pageSize=mostRequestedPageSize, ~page=1, ~searchTerm="", ());

  let allSub = Sub.all2(oracleScriptsSub, oracleScriptsCountSub);

  <Section>
    <div className=CssHelper.container id="oraclescriptsSection">
      <div className=CssHelper.mobileSpacing>
        <Heading value="All Oracle Scripts" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
        {switch (mostRequestedOracleScriptSub) {
         | Data(oracleScripts) =>
           oracleScripts->Belt.Array.length > 0
             ? <>
                 <Heading value="Most Requested" size=Heading.H4 marginBottom=16 />
                 <Row>
                   {oracleScripts
                    ->Belt_Array.mapWithIndex((i, e) =>
                        renderMostRequestedCard(i, Sub.resolve(e))
                      )
                    ->React.array}
                 </Row>
               </>
             : React.null
         | _ =>
           <>
             <Heading value="Most Requested" size=Heading.H4 marginBottom=16 />
             <Row>
               {Belt_Array.make(mostRequestedPageSize, ApolloHooks.Subscription.NoData)
                ->Belt_Array.mapWithIndex((i, noData) => renderMostRequestedCard(i, noData))
                ->React.array}
             </Row>
           </>
         }}
        <Row alignItems=Row.Center marginBottom=40 marginBottomSm=24>
          <Col.Grid>
            {switch (allSub) {
             | Data((_, oracleScriptsCount)) =>
               <Heading
                 value={(oracleScriptsCount |> Format.iPretty) ++ " In total"}
                 size=Heading.H3
               />
             | _ => <LoadingCensorBar width=65 height=21 />
             }}
          </Col.Grid>
        </Row>
        <Row alignItems=Row.Center marginBottom=16>
          <Col.Grid col=Col.Six colSm=Col.Eight>
            <SearchInput placeholder="Search Oracle Script" onChange=setSearchTerm />
          </Col.Grid>
          <Col.Grid col=Col.Six colSm=Col.Four>
            <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
              <SortableDropdown
                sortedBy
                setSortedBy
                sortList=[
                  (MostRequested, getName(MostRequested)),
                  (LatestUpdate, getName(LatestUpdate)),
                ]
              />
            </div>
          </Col.Grid>
        </Row>
        {isMobile
           ? React.null
           : <THead.Grid>
               <Row alignItems=Row.Center>
                 <Col.Grid col=Col.Five>
                   <div className=TElement.Styles.hashContainer>
                     <Text
                       block=true
                       value="Oracle Script"
                       weight=Text.Semibold
                       color=Colors.gray7
                     />
                   </div>
                 </Col.Grid>
                 <Col.Grid col=Col.Three>
                   <Text block=true value="Description" weight=Text.Semibold color=Colors.gray7 />
                 </Col.Grid>
                 <Col.Grid col=Col.Two>
                   <Text
                     block=true
                     value="Request"
                     weight=Text.Semibold
                     color=Colors.gray7
                     align=Text.Right
                   />
                   <Text
                     block=true
                     value="& Response time"
                     weight=Text.Semibold
                     color=Colors.gray7
                     align=Text.Right
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
               </Row>
             </THead.Grid>}
        {switch (allSub) {
         | Data((oracleScripts, oracleScriptsCount)) =>
           let pageCount = Page.getPageCount(oracleScriptsCount, pageSize);
           <div className=Styles.tbodyContainer>
             {oracleScripts->Belt.Array.length > 0
                ? oracleScripts
                  ->sorting(sortedBy)
                  ->Belt_Array.mapWithIndex((i, e) =>
                      isMobile
                        ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
                    )
                  ->React.array
                : <EmptyContainer>
                    <img src=Images.noSource className=Styles.noDataImage />
                    <Heading
                      size=Heading.H4
                      value="No Oracle Script"
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
           </div>;
         | _ =>
           <div className=Styles.tbodyContainer>
             {Belt_Array.make(10, ApolloHooks.Subscription.NoData)
              ->Belt_Array.mapWithIndex((i, noData) =>
                  isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
                )
              ->React.array}
           </div>
         }}
      </div>
    </div>
  </Section>;
};
