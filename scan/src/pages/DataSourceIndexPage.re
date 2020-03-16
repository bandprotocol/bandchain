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
let make = (~dataSourceID, ~hashtag: Route.data_source_tab_t) => {
  let dataSourceOpt = DataSourceHook.get(dataSourceID);

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
            color=Colors.gray7
            block=true
          />
          <div className=Styles.seperatedLine />
          {switch (dataSourceOpt) {
           | Some(dataSource) =>
             dataSource.revisions
             ->Belt_List.get(0)
             ->Belt_Option.mapWithDefault(React.null, ({timestamp}) =>
                 <TimeAgos
                   time=timestamp
                   prefix="Last updated "
                   size=Text.Md
                   weight=Text.Thin
                   spacing={Text.Em(0.06)}
                   height={Text.Px(18)}
                   upper=true
                 />
               )
           | None =>
             <Text
               value="???"
               size=Text.Md
               weight=Text.Thin
               spacing={Text.Em(0.06)}
               height={Text.Px(18)}
             />
           }}
        </div>
      </Col>
    </Row>
    {switch (dataSourceOpt) {
     | Some(dataSource) =>
       <>
         <VSpacing size=Spacing.md />
         <div className=Styles.vFlex>
           <TypeID.DataSource id={ID.DataSource.ID(dataSource.id)} position=TypeID.Title />
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
           <Col size=1.> <InfoHL header="OWNER" info={InfoHL.Address(dataSource.owner)} /> </Col>
           <Col size=0.8>
             <InfoHL
               info={
                 InfoHL.Fee(
                   dataSource.fee
                   ->Belt_List.get(0)
                   ->Belt_Option.mapWithDefault(0., coin => coin.amount),
                 )
               }
               header="REQUEST FEE"
             />
           </Col>
         </Row>
         <VSpacing size=Spacing.xl />
         <Tab
           tabs=[|
             {
               name: "EXECUTION",
               route: Route.DataSourceIndexPage(dataSourceID, Route.DataSourceExecute),
             },
             {
               name: "CODE",
               route: Route.DataSourceIndexPage(dataSourceID, Route.DataSourceCode),
             },
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
            | DataSourceExecute => <DataSourceExecute executable={dataSource.executable} />
            | DataSourceCode => <DataSourceCode />
            | DataSourceRequests => <DataSourceRequestTable />
            | DataSourceRevisions => <RevisionTable />
            }}
         </Tab>
       </>
     | None => React.null
     }}
  </div>;
};
