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
      backgroundColor(Colors.gray7),
    ]);
};

[@react.component]
let make = (~dataSourceID, ~hashtag: Route.data_source_tab_t) =>
  {
    let%Sub dataSource = DataSourceSub.get(dataSourceID);

    <div className=Styles.pageContainer>
      <div className=Styles.vFlex>
        <img src=Images.dataSourceLogo className=Styles.logo />
        <Text
          value="DATA SOURCE"
          weight=Text.Medium
          size=Text.Md
          spacing={Text.Em(0.06)}
          height={Text.Px(15)}
          nowrap=true
          color=Colors.gray7
          block=true
        />
        <div className=Styles.seperatedLine />
        <TimeAgos
          time={dataSource.timestamp}
          prefix="Last updated "
          size=Text.Md
          weight=Text.Thin
          spacing={Text.Em(0.06)}
          height={Text.Px(18)}
          upper=true
        />
      </div>
      <>
        <VSpacing size=Spacing.md />
        <VSpacing size=Spacing.sm />
        <div className=Styles.vFlex>
          <TypeID.DataSource id={dataSource.id} position=TypeID.Title />
          <HSpacing size=Spacing.md />
          <Text
            value={dataSource.name}
            size=Text.Xxl
            height={Text.Px(22)}
            weight=Text.Bold
            nowrap=true
          />
        </div>
        <VSpacing size=Spacing.xl />
        <Row>
          <Col size=1.>
            <InfoHL header="OWNER" info={InfoHL.Address(dataSource.owner, 380)} />
          </Col>
          <Col size=0.8>
            <InfoHL
              info={InfoHL.Fee(dataSource.fee->Coin.getBandAmountFromCoins)}
              header="REQUEST FEE"
            />
          </Col>
        </Row>
        <VSpacing size=Spacing.xl />
        <Tab
          tabs=[|
            {
              name: "EXECUTION",
              route: dataSourceID |> ID.DataSource.getRouteWithTab(_, Route.DataSourceExecute),
            },
            {
              name: "CODE",
              route: dataSourceID |> ID.DataSource.getRouteWithTab(_, Route.DataSourceCode),
            },
            {
              name: "REQUESTS",
              route: dataSourceID |> ID.DataSource.getRouteWithTab(_, Route.DataSourceRequests),
            },
            {
              name: "REVISIONS",
              route: dataSourceID |> ID.DataSource.getRouteWithTab(_, Route.DataSourceRevisions),
            },
          |]
          currentRoute={dataSourceID |> ID.DataSource.getRouteWithTab(_, hashtag)}>
          {switch (hashtag) {
           | DataSourceExecute => <DataSourceExecute executable={dataSource.executable} />
           | DataSourceCode => <DataSourceCode executable={dataSource.executable} />
           | DataSourceRequests => <DataSourceRequestTable dataSourceID />
           | DataSourceRevisions => <DataSourceRevisionTable id=dataSourceID />
           }}
        </Tab>
      </>
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
