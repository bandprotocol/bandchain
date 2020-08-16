type sort_by_t =
  | MostRequested
  | LatestUpdate;

let getName =
  fun
  | MostRequested => "Most Requested"
  | LatestUpdate => "Latest Update";

let defaultCompare = (a: DataSourceSub.t, b: DataSourceSub.t) =>
  if (a.timestamp != b.timestamp) {
    let ID.DataSource.ID(a_) = a.id;
    let ID.DataSource.ID(b_) = b.id;
    compare(b_, a_);
  } else {
    compare(b.request, a.request);
  };

let sorting = (dataSources: array(DataSourceSub.t), sortedBy) => {
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

let renderBody =
    (reserveIndex, dataSourcesSub: ApolloHooks.Subscription.variant(DataSourceSub.t)) => {
  <TBody.Grid
    key={
      switch (dataSourcesSub) {
      | Data({id}) => id |> ID.DataSource.toString
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
      <Col.Grid col=Col.Five>
        {switch (dataSourcesSub) {
         | Data({id, name}) =>
           <div className={CssHelper.flexBox()}>
             <TypeID.DataSource id />
             <HSpacing size=Spacing.sm />
             <Text value=name ellipsis=true />
           </div>
         | _ => <LoadingCensorBar width=300 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Four>
        {switch (dataSourcesSub) {
         | Data({description}) => <Text value=description weight=Text.Medium block=true />
         | _ => <LoadingCensorBar width=270 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.One>
        {switch (dataSourcesSub) {
         | Data({request}) =>
           <div>
             <Text value={request |> Format.iPretty} weight=Text.Medium block=true ellipsis=true />
           </div>
         | _ => <LoadingCensorBar width=70 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (dataSourcesSub) {
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
    (reserveIndex, dataSourcesSub: ApolloHooks.Subscription.variant(DataSourceSub.t)) => {
  switch (dataSourcesSub) {
  | Data({id, timestamp, description, name, request}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Data Sourse", DataSource(id, name)),
        ("Description", Text(description)),
        ("Requests", Count(request)),
        ("Timestamp", Timestamp(timestamp)),
      ]
      key={id |> ID.DataSource.toString}
      idx={id |> ID.DataSource.toString}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Data Sources", Loading(70)),
        ("Description", Loading(136)),
        ("Requests", Loading(20)),
        ("Timestamp", Loading(166)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = () => {
  let (page, setPage) = React.useState(_ => 1);
  let (searchTerm, setSearchTerm) = React.useState(_ => "");
  let (sortedBy, setSortedBy) = React.useState(_ => LatestUpdate);
  let pageSize = 10;

  let dataSourcesCountSub = DataSourceSub.count(~searchTerm, ());
  let dataSourcesSub = DataSourceSub.getList(~pageSize, ~page, ~searchTerm, ());

  let allSub = Sub.all2(dataSourcesSub, dataSourcesCountSub);
  let isMobile = Media.isMobile();

  <div className=CssHelper.mobileSpacing>
    <Row.Grid alignItems=Row.Center marginBottom=40 marginBottomSm=24>
      <Col.Grid col=Col.Twelve>
        <Heading value="All Data Sources" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
        {switch (allSub) {
         | Data((_, dataSourcesCount)) =>
           <Heading value={dataSourcesCount->string_of_int ++ " In total"} size=Heading.H3 />
         | _ => <LoadingCensorBar width=65 height=21 />
         }}
      </Col.Grid>
    </Row.Grid>
    <>
      <Row.Grid alignItems=Row.Center marginBottom=16>
        <Col.Grid col=Col.Six colSm=Col.Eight>
          <SearchInput placeholder="Search Data Source" onChange=setSearchTerm />
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
                     value="Data Sources"
                     size=Text.Md
                     weight=Text.Semibold
                     color=Colors.gray7
                   />
                 </div>
               </Col.Grid>
               <Col.Grid col=Col.Four>
                 <Text
                   block=true
                   value="Descriptions"
                   size=Text.Md
                   weight=Text.Semibold
                   color=Colors.gray7
                 />
               </Col.Grid>
               <Col.Grid col=Col.One>
                 <Text
                   block=true
                   value="Requests"
                   size=Text.Md
                   weight=Text.Semibold
                   color=Colors.gray7
                 />
               </Col.Grid>
               <Col.Grid col=Col.Two>
                 <Text
                   block=true
                   value="Timestamp"
                   size=Text.Md
                   weight=Text.Semibold
                   color=Colors.gray7
                   align=Text.Right
                 />
               </Col.Grid>
             </Row.Grid>
           </THead.Grid>}
    </>
    {switch (allSub) {
     | Data((dataSources, dataSourcesCount)) =>
       let pageCount = Page.getPageCount(dataSourcesCount, pageSize);
       <>
         {dataSources
          ->sorting(sortedBy)
          ->Belt_Array.mapWithIndex((i, e) =>
              isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
            )
          ->React.array}
         {isMobile
            ? React.null
            : <>
                <VSpacing size=Spacing.lg />
                <Pagination
                  currentPage=page
                  pageCount
                  onPageChange={newPage => setPage(_ => newPage)}
                />
                <VSpacing size=Spacing.lg />
              </>}
       </>;
     | _ =>
       Belt_Array.make(10, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
         )
       ->React.array
     }}
  </div>;
};
