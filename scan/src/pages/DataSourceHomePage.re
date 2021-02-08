module Styles = {
  open Css;

  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

type sort_by_t =
  | MostRequested
  | LatestUpdate;

let getName =
  fun
  | MostRequested => "Most Requested"
  | LatestUpdate => "Latest Update";

let defaultCompare = (a: DataSourceSub.t, b: DataSourceSub.t) =>
  if (a.timestamp != b.timestamp) {
    compare(b.id |> ID.DataSource.toInt, a.id |> ID.DataSource.toInt);
  } else {
    compare(b.requestCount, a.requestCount);
  };

let sorting = (dataSources: array(DataSourceSub.t), sortedBy) => {
  dataSources
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
    <Row.Grid alignItems=Row.Center>
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
      <Col.Grid col=Col.Five>
        {switch (dataSourcesSub) {
         | Data({description}) => <Text value=description block=true />
         | _ => <LoadingCensorBar width=270 height=15 />
         }}
      </Col.Grid>
      // <Col.Grid col=Col.One>
      //   {switch (dataSourcesSub) {
      //    | Data({requestCount}) =>
      //      <div>
      //        <Text
      //          value={requestCount |> Format.iPretty}
      //          weight=Text.Medium
      //          block=true
      //          ellipsis=true
      //        />
      //      </div>
      //    | _ => <LoadingCensorBar width=70 height=15 />
      //    }}
      // </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (dataSourcesSub) {
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
    </Row.Grid>
  </TBody.Grid>;
};

let renderBodyMobile =
    (reserveIndex, dataSourcesSub: ApolloHooks.Subscription.variant(DataSourceSub.t)) => {
  switch (dataSourcesSub) {
  | Data({id, timestamp: timestampOpt, description, name, requestCount}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Data Source", DataSource(id, name)),
        ("Description", Text(description)),
        // ("Requests", Count(requestCount)),
        (
          "Timestamp",
          switch (timestampOpt) {
          | Some(timestamp') => Timestamp(timestamp')
          | None => Text("Genesis")
          },
        ),
      ]
      key={id |> ID.DataSource.toString}
      idx={id |> ID.DataSource.toString}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Data Source", Loading(70)),
        ("Description", Loading(136)),
        // ("Requests", Loading(20)),
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

  React.useEffect1(
    () => {
      if (searchTerm != "") {
        setPage(_ => 1);
      };
      None;
    },
    [|searchTerm|],
  );

  <Section>
    <div className=CssHelper.container id="datasourcesSection">
      <div className=CssHelper.mobileSpacing>
        <Row.Grid alignItems=Row.Center marginBottom=40 marginBottomSm=24>
          <Col.Grid col=Col.Twelve>
            <Heading value="All Data Sources" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
            {switch (allSub) {
             | Data((_, dataSourcesCount)) =>
               <Heading
                 value={(dataSourcesCount |> Format.iPretty) ++ " In total"}
                 size=Heading.H3
               />
             | _ => <LoadingCensorBar width=65 height=21 />
             }}
          </Col.Grid>
        </Row.Grid>
        <Row.Grid alignItems=Row.Center marginBottom=16>
          <Col.Grid col=Col.Six colSm=Col.Eight>
            <SearchInput placeholder="Search Data Source" onChange=setSearchTerm />
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
        </Row.Grid>
        {isMobile
           ? React.null
           : <THead.Grid>
               <Row.Grid alignItems=Row.Center>
                 <Col.Grid col=Col.Five>
                   <div className=TElement.Styles.hashContainer>
                     <Text
                       block=true
                       value="Data Source"
                       size=Text.Md
                       weight=Text.Semibold
                       color=Colors.gray7
                     />
                   </div>
                 </Col.Grid>
                 <Col.Grid col=Col.Five>
                   <Text
                     block=true
                     value="Description"
                     size=Text.Md
                     weight=Text.Semibold
                     color=Colors.gray7
                   />
                 </Col.Grid>
                 //  <Col.Grid col=Col.One>
                 //    <Text
                 //      block=true
                 //      value="Requests"
                 //      size=Text.Md
                 //      weight=Text.Semibold
                 //      color=Colors.gray7
                 //    />
                 //  </Col.Grid>
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
        {switch (allSub) {
         | Data((dataSources, dataSourcesCount)) =>
           let pageCount = Page.getPageCount(dataSourcesCount, pageSize);
           <>
             {dataSources->Belt.Array.length > 0
                ? dataSources
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
                      value="No Data Source"
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
           Belt_Array.make(10, ApolloHooks.Subscription.NoData)
           ->Belt_Array.mapWithIndex((i, noData) =>
               isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
             )
           ->React.array
         }}
      </div>
    </div>
  </Section>;
};
