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
let make = (~oracleScriptID, ~hashtag: Route.oracle_script_tab_t) => {
  let oracleScriptOpt = OracleScriptHook.get(oracleScriptID);

  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.oracleScriptLogo className=Styles.logo />
          <Text
            value="ORACLE SCRIPT"
            weight=Text.Medium
            size=Text.Md
            spacing={Text.Em(0.06)}
            height={Text.Px(15)}
            nowrap=true
            color=Colors.gray7
            block=true
          />
          <div className=Styles.seperatedLine />
          {switch (oracleScriptOpt) {
           | Some(oracleScript) =>
             oracleScript.revisions
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
    {switch (oracleScriptOpt) {
     | Some(oracleScript) =>
       <>
         <VSpacing size=Spacing.xl />
         <div className=Styles.vFlex>
           <TypeID.OracleScript id={ID.OracleScript.ID(oracleScript.id)} position=TypeID.Title />
           <HSpacing size=Spacing.md />
           <Text
             value={oracleScript.name}
             size=Text.Xxl
             height={Text.Px(22)}
             weight=Text.Bold
             nowrap=true
           />
         </div>
         <VSpacing size=Spacing.xl />
         <Row>
           <Col size=1.>
             <InfoHL header="OWNER" info={InfoHL.Address(oracleScript.owner)} />
           </Col>
           <Col size=0.8>
             <InfoHL
               info={InfoHL.DataSources(oracleScript.relatedDataSource)}
               header="RELATED DATA SOURCES"
             />
           </Col>
         </Row>
         <VSpacing size=Spacing.xl />
         <Tab
           tabs=[|
             {
               name: "EXECUTION",
               route: Route.OracleScriptIndexPage(oracleScriptID, Route.OracleScriptExecute),
             },
             {
               name: "CODE",
               route: Route.OracleScriptIndexPage(oracleScriptID, Route.OracleScriptCode),
             },
             {
               name: "REQUESTS",
               route: Route.OracleScriptIndexPage(oracleScriptID, Route.OracleScriptRequests),
             },
             {
               name: "REVISIONS",
               route: Route.OracleScriptIndexPage(oracleScriptID, Route.OracleScriptRevisions),
             },
           |]
           currentRoute={Route.OracleScriptIndexPage(oracleScriptID, hashtag)}>
           {switch (hashtag) {
            | OracleScriptExecute => <OracleScriptExecute code={oracleScript.code} />
            | OracleScriptCode => <OracleScriptCode />
            | OracleScriptRequests => <OracleScriptRequestTable />
            | OracleScriptRevisions => <RevisionTable />
            }}
         </Tab>
       </>
     | None => React.null
     }}
  </div>;
};
