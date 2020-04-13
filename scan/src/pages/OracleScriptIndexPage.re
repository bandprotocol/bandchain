module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);
  let headerContainer = style([lineHeight(`px(25))]);

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
let make = (~oracleScriptID, ~hashtag: Route.oracle_script_tab_t) =>
  {
    let%Sub oracleScript = OracleScriptSub.get(oracleScriptID);

    <>
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
          </div>
        </Col>
      </Row>
      // <div className=Styles.seperatedLine />
      // <TimeAgos
      //   time={oracleScript.timestamp}
      //   prefix="Last updated "
      //   size=Text.Md
      //   weight=Text.Thin
      //   spacing={Text.Em(0.06)}
      //   height={Text.Px(18)}
      //   upper=true
      // />
      <VSpacing size=Spacing.xl />
      <div className=Styles.vFlex>
        <TypeID.OracleScript id={oracleScript.id} position=TypeID.Title />
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
          <InfoHL header="OWNER" info={InfoHL.Address(oracleScript.owner, 430)} />
        </Col>
        <Col size=0.95>
          <InfoHL
            info={InfoHL.DataSources(oracleScript.relatedDataSources)}
            header="RELATED DATA SOURCES"
          />
        </Col>
      </Row>
      <VSpacing size=Spacing.sm />
      <Row>
        <Col size=1.>
          <InfoHL header="DESCRIPTION" info={InfoHL.Description(oracleScript.description)} />
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <Tab
        tabs=[|
          {
            name: "EXECUTION",
            route:
              oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptExecute),
          },
          {
            name: "OWASM CODE",
            route: oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptCode),
          },
          {
            name: "BRIDGE CODE",
            route:
              oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptBridgeCode),
          },
          {
            name: "REQUESTS",
            route:
              oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptRequests),
          },
          {
            name: "REVISIONS",
            route:
              oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptRevisions),
          },
        |]
        currentRoute={oracleScriptID |> ID.OracleScript.getRouteWithTab(_, hashtag)}>
        {switch (hashtag) {
         | OracleScriptExecute => <OracleScriptExecute code={oracleScript.codeHash} />
         | OracleScriptCode => <OracleScriptCode code={oracleScript.codeHash} />
         | OracleScriptBridgeCode => <OracleScriptBridgeCode />
         | OracleScriptRequests => <OracleScriptRequestTable oracleScriptID />
         | OracleScriptRevisions => <OracleScriptRevisionTable id=oracleScriptID />
         }}
      </Tab>
    </>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
