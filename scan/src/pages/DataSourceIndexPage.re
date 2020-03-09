module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(40))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.mediumGray),
    ]);
};

[@react.component]
let make = (~dataSourceID, ~hashtag: Route.data_source_tab_t) => {
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.dataSourceLogo className=Styles.logo />
          <Text
            value="DATA SOURCE"
            weight=Text.Medium
            size=Text.Md
            spacing={Text.Em(0.06)}
            height={Text.Px(15)}
            nowrap=true
            color=Colors.mediumGray
            block=true
          />
          <div className=Styles.seperatedLine />
          <Text
            value="Last updated 4 hours ago"
            size=Text.Md
            weight=Text.Thin
            spacing={Text.Em(0.06)}
            color=Colors.mediumGray
            nowrap=true
          />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.md />
    <div className=Styles.vFlex>
      <TypeID.DataSource id={ID.DataSource.ID(34)} position=TypeID.Title />
      <HSpacing size=Spacing.md />
      <Text
        value="CoinGecko V.2"
        size=Text.Xxl
        height={Text.Px(22)}
        weight=Text.Bold
        nowrap=true
      />
    </div>
    <VSpacing size=Spacing.xl />
    <Row>
      <Col size=1.>
        <InfoHL
          header="OWNER"
          info={
            InfoHL.Address(
              "band1gfskuezzv9hxgsnpdejyyctwv3pxzmnywps0q9" |> Address.fromBech32,
              Colors.mediumGray,
            )
          }
        />
      </Col>
      <Col size=0.8> <InfoHL info={InfoHL.Fee(1000.)} header="REQUEST FEE" /> </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <Tab
      tabs=[|
        {
          name: "EXECUTION",
          route: Route.DataSourceIndexPage(dataSourceID, Route.DataSourceExecute),
        },
        {name: "CODE", route: Route.DataSourceIndexPage(dataSourceID, Route.DataSourceCode)},
        {
          name: "REQUESTS",
          route: Route.DataSourceIndexPage(dataSourceID, Route.DataSourceRequests),
        },
        {
          name: "REVISIONS",
          route: Route.DataSourceIndexPage(dataSourceID, Route.DataSourceRevisions),
        },
      |]
      currentRoute={Route.DataSourceIndexPage(dataSourceID, hashtag)}>
      {switch (hashtag) {
       | DataSourceExecute => <DataSourceExecute />
       | DataSourceCode => <DataSourceCode />
       | DataSourceRequests => <DataSourceRequestTable />
       | DataSourceRevisions => <RevisionTable />
       }}
    </Tab>
  </div>;
};
