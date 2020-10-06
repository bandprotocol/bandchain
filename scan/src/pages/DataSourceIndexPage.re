module Styles = {
  open Css;
  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding(`px(24)),
      Media.mobile([padding(`px(16))]),
    ]);
  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
};

module Content = {
  [@react.component]
  let make =
      (
        ~dataSourceSub: ApolloHooks.Subscription.variant(BandScan.DataSourceSub.t),
        ~dataSourceID,
        ~hashtag: Route.data_source_tab_t,
      ) => {
    <Section pbSm=0>
      <div className=CssHelper.container>
        <Row.Grid marginBottom=40 marginBottomSm=16>
          <Col.Grid>
            <Heading value="Data Source" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
            {switch (dataSourceSub) {
             | Data({id, name}) =>
               <div className={CssHelper.flexBox()}>
                 <TypeID.DataSource id position=TypeID.Title />
                 <HSpacing size=Spacing.sm />
                 <Heading size=Heading.H3 value=name />
               </div>
             | _ => <LoadingCensorBar width=270 height=15 />
             }}
          </Col.Grid>
        </Row.Grid>
        <Row.Grid marginBottom=24>
          <Col.Grid>
            <div className=Styles.infoContainer>
              <Heading
                value="Information"
                size=Heading.H4
                style=Styles.infoHeader
                marginBottom=24
              />
              <div className={CssHelper.flexBox()}>
                <Heading value="Owner" size=Heading.H5 />
                <HSpacing size=Spacing.xs />
                <CTooltip tooltipText="The owner of the data source">
                  <Icon name="fal fa-info-circle" size=10 />
                </CTooltip>
              </div>
              <VSpacing size=Spacing.sm />
              {switch (dataSourceSub) {
               | Data({owner}) => <AddressRender address=owner position=AddressRender.Subtitle />
               | _ => <LoadingCensorBar width=284 height=15 />
               }}
              <VSpacing size=Spacing.lg />
              <Heading value="Description" size=Heading.H5 marginBottom=16 />
              {switch (dataSourceSub) {
               | Data({description}) =>
                 <p>
                   <Text
                     value=description
                     weight=Text.Regular
                     size=Text.Lg
                     color=Colors.gray7
                     block=true
                   />
                 </p>
               | _ => <LoadingCensorBar width=284 height=15 />
               }}
            </div>
          </Col.Grid>
        </Row.Grid>
        <Tab
          tabs=[|
            {
              name: "Requests",
              route: dataSourceID |> ID.DataSource.getRouteWithTab(_, Route.DataSourceRequests),
            },
            {
              name: "Code",
              route: dataSourceID |> ID.DataSource.getRouteWithTab(_, Route.DataSourceCode),
            },
            {
              name: "Test Execution",
              route: dataSourceID |> ID.DataSource.getRouteWithTab(_, Route.DataSourceExecute),
            },
            // {
            //   name: "REVISIONS",
            //   route: dataSourceID |> ID.DataSource.getRouteWithTab(_, Route.DataSourceRevisions),
            // },
          |]
          currentRoute={dataSourceID |> ID.DataSource.getRouteWithTab(_, hashtag)}>
          {switch (hashtag) {
           | DataSourceExecute =>
             switch (dataSourceSub) {
             | Data({executable}) => <DataSourceExecute executable />
             | _ => <LoadingCensorBar fullWidth=true height=400 />
             }
           | DataSourceCode =>
             switch (dataSourceSub) {
             | Data({executable}) => <DataSourceCode executable />
             | _ => <LoadingCensorBar fullWidth=true height=300 />
             }
           | DataSourceRequests => <DataSourceRequestTable dataSourceID />
           | DataSourceRevisions => <DataSourceRevisionTable id=dataSourceID />
           }}
        </Tab>
      </div>
    </Section>;
  };
};

[@react.component]
let make = (~dataSourceID, ~hashtag: Route.data_source_tab_t) => {
  let dataSourceSub = DataSourceSub.get(dataSourceID);

  switch (dataSourceSub) {
  | NoData => <NotFound />
  | _ => <Content dataSourceSub dataSourceID hashtag />
  };
};
