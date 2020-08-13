module Styles = {
  open Css;
  let mostRequestCard =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      padding(`px(24)),
      height(`calc((`sub, `percent(100.), `px(24)))),
      marginBottom(`px(24)),
    ]);
  let requestResponseBox = style([flexShrink(0.), flexGrow(0.), flexBasis(`percent(50.))]);
  let descriptionBox = style([minHeight(`px(36)), margin2(~v=`px(16), ~h=`zero)]);
  let idBox = style([marginBottom(`px(4))]);
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
    let ID.OracleScript.ID(a_) = a.id;
    let ID.OracleScript.ID(b_) = b.id;
    compare(b_, a_);
  } else {
    compare(b.request, a.request);
  };

let sorting = (dataSources: array(OracleScriptSub.t), sortedBy) => {
  dataSources
  ->Belt.List.fromArray
  ->Belt.List.sort((a, b) => {
      let result = {
        switch (sortedBy) {
        | MostRequested => compare(b.request, a.request)
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
           <Text size=Text.Lg value=description weight=Text.Regular block=true />
         | _ => <LoadingCensorBar width=250 height=15 />
         }}
      </div>
      <div className={CssHelper.flexBox()}>
        <div className=Styles.requestResponseBox>
          <Heading size=Heading.H5 value="Requests" marginBottom=8 />
          {switch (oracleScriptSub) {
           | Data({request}) =>
             <Text size=Text.Lg value={request |> Format.iPretty} weight=Text.Regular block=true />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
        <div className=Styles.requestResponseBox>
          <Heading size=Heading.H5 value="Response time" marginBottom=8 />
          {switch (oracleScriptSub) {
           | Data({responseTime}) =>
             <Text
               size=Text.Lg
               value={(responseTime |> Format.iPretty) ++ "ms"}
               weight=Text.Regular
               block=true
             />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
      </div>
    </div>
  </Col.Grid>;
};

let renderBody =
    (reserveIndex, oracleScriptSub: ApolloHooks.Subscription.variant(OracleScriptSub.t)) => {
  <TBody.Grid
    key={
      switch (oracleScriptSub) {
      | Data({id}) => id |> ID.OracleScript.toString
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
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
      <Col.Grid col=Col.Four>
        {switch (oracleScriptSub) {
         | Data({description}) => <Text value=description weight=Text.Medium block=true />
         | _ => <LoadingCensorBar width=270 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.One>
        {switch (oracleScriptSub) {
         | Data({request, responseTime}) =>
           <>
             <div>
               <Text
                 value={request |> Format.iPretty}
                 weight=Text.Medium
                 block=true
                 ellipsis=true
                 align=Text.Right
               />
             </div>
             <div>
               <Text
                 value={"(" ++ (responseTime |> Format.iPretty) ++ "ms)"}
                 weight=Text.Medium
                 block=true
                 color=Colors.gray6
                 align=Text.Right
               />
             </div>
           </>
         | _ => <LoadingCensorBar width=70 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (oracleScriptSub) {
           | Data({timestamp}) =>
             <Timestamp.Grid
               time=timestamp
               size=Text.Md
               weight=Text.Regular
               textAlign=Text.Right
             />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </div>
      </Col.Grid>
    </Row.Grid>
  </TBody.Grid>;
};

let renderBodyMobile =
    (reserveIndex, oracleScriptSub: ApolloHooks.Subscription.variant(OracleScriptSub.t)) => {
  switch (oracleScriptSub) {
  | Data({id, timestamp, description, name, request, responseTime}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Oracle Script", OracleScript(id, name)),
        ("Description", Text(description)),
        ("Requests&\nResponse time", RequestResponse({request, responseTime})),
        ("Timestamp", Timestamp(timestamp)),
      ]
      key={id |> ID.OracleScript.toString}
      idx={id |> ID.OracleScript.toString}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Oracle Script", Loading(200)),
        ("Description", Loading(200)),
        ("Requests&\nResponse time", Loading(80)),
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

  let pageSize = 10;
  let oracleScriptsCountSub = OracleScriptSub.count(~searchTerm, ());
  let oracleScriptsSub = OracleScriptSub.getList(~pageSize, ~page, ~searchTerm, ());
  let oracleScriptTopPart = OracleScriptSub.getList(~pageSize=6, ~page=1, ~searchTerm="", ());

  let allSub = Sub.all2(oracleScriptsSub, oracleScriptsCountSub);

  <div className=CssHelper.mobileSpacing>
    <Heading value="All Oracle Scripts" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
    <Heading value="Most Requested" size=Heading.H4 marginBottom=16 />
    <Row.Grid>
      {switch (oracleScriptTopPart) {
       | Data(oracleScripts) =>
         <>
           {oracleScripts
            ->sorting(MostRequested)
            ->Belt_Array.mapWithIndex((i, e) => renderMostRequestedCard(i, Sub.resolve(e)))
            ->React.array}
         </>
       | _ =>
         Belt_Array.make(6, ApolloHooks.Subscription.NoData)
         ->Belt_Array.mapWithIndex((i, noData) => renderMostRequestedCard(i, noData))
         ->React.array
       }}
    </Row.Grid>
    <Row.Grid alignItems=Row.Center marginBottom=40 marginBottomSm=24>
      <Col.Grid>
        {switch (allSub) {
         | Data((_, oracleScriptsCount)) =>
           <Heading value={oracleScriptsCount->string_of_int ++ " In total"} size=Heading.H3 />
         | _ => <LoadingCensorBar width=65 height=21 />
         }}
      </Col.Grid>
    </Row.Grid>
    <>
      <Row.Grid alignItems=Row.Center marginBottom=16>
        <Col.Grid col=Col.Six colSm=Col.Eight>
          <SearchInput placeholder="Search Oracle Script" onChange=setSearchTerm />
        </Col.Grid>
        <Col.Grid col=Col.Six colSm=Col.Four>
          <SortableDropdown
            sortedBy
            setSortedBy
            sortList=[
              (MostRequested, getName(MostRequested)),
              (LatestUpdate, getName(LatestUpdate)),
            ]
          />
        </Col.Grid>
      </Row.Grid>
      {isMobile
         ? React.null
         : <THead.Grid>
             <Row.Grid alignItems=Row.Center>
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
                 <Text block=true value="Descriptions" weight=Text.Semibold color=Colors.gray7 />
               </Col.Grid>
               <Col.Grid col=Col.Two>
                 <Text
                   block=true
                   value="Requests"
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
             </Row.Grid>
           </THead.Grid>}
      {switch (allSub) {
       | Data((oracleScripts, oracleScriptsCount)) =>
         let pageCount = Page.getPageCount(oracleScriptsCount, pageSize);
         <>
           {oracleScripts
            ->sorting(sortedBy)
            ->Belt_Array.mapWithIndex((i, e) =>
                isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
              )
            ->React.array}
           {isMobile
              ? React.null
              : <Pagination
                  currentPage=page
                  pageCount
                  onPageChange={newPage => setPage(_ => newPage)}
                />}
         </>;
       | _ =>
         Belt_Array.make(10, ApolloHooks.Subscription.NoData)
         ->Belt_Array.mapWithIndex((i, noData) =>
             isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
           )
         ->React.array
       }}
    </>
  </div>;
};
