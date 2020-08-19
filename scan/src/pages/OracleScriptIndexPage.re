module Styles = {
  open Css;

  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      padding(`px(24)),
      Media.mobile([padding(`px(16))]),
    ]);
  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
  let titleSpacing = style([marginBottom(`px(8))]);
  let loadingBox = style([width(`percent(100.))]);
  let idCointainer = style([marginBottom(`px(16))]);
  let containerSpacingSm = style([Media.mobile([marginTop(`px(16))])]);
};

[@react.component]
let make = (~oracleScriptID, ~hashtag: Route.oracle_script_tab_t) => {
  let oracleScriptSub = OracleScriptSub.get(oracleScriptID);

  <Section pbSm=0>
    <div className=CssHelper.container>
      <Heading value="Oracle Script" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
      <Row.Grid marginBottom=40 marginBottomSm=16 alignItems=Row.Center>
        <Col.Grid col=Col.Eight>
          <div className={Css.merge([CssHelper.flexBox(), Styles.idCointainer])}>
            {switch (oracleScriptSub) {
             | Data({id, name}) =>
               <>
                 <TypeID.OracleScript id position=TypeID.Title />
                 <HSpacing size=Spacing.sm />
                 <Heading size=Heading.H3 value=name />
               </>
             | _ => <LoadingCensorBar width=270 height=15 />
             }}
          </div>
        </Col.Grid>
        <Col.Grid col=Col.Four>
          <div className=Styles.infoContainer>
            <Row.Grid>
              <Col.Grid col=Col.Six colSm=Col.Six>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading value="Requests" size=Heading.H4 marginBottom=8 align=Heading.Center />
                  {switch (oracleScriptSub) {
                   | Data({request}) =>
                     <Text
                       value={request |> Format.iPretty}
                       size=Text.Xxl
                       align=Text.Center
                       block=true
                     />
                   | _ => <LoadingCensorBar width=100 height=15 />
                   }}
                </div>
              </Col.Grid>
              <Col.Grid col=Col.Six colSm=Col.Six>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <div
                    className={Css.merge([
                      CssHelper.flexBox(~justify=`center, ()),
                      Styles.titleSpacing,
                    ])}>
                    <Heading value="Response time" size=Heading.H4 align=Heading.Center />
                    <HSpacing size=Spacing.xs />
                    //TODO: remove mock message later
                    <CTooltip
                      tooltipPlacementSm=CTooltip.BottomRight
                      tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                      <Icon name="fal fa-info-circle" size=12 />
                    </CTooltip>
                  </div>
                  {switch (oracleScriptSub) {
                   | Data({responseTime}) =>
                     <Text
                       value={responseTime |> Format.iPretty}
                       size=Text.Xxl
                       align=Text.Center
                       block=true
                     />
                   | _ => <LoadingCensorBar width=100 height=15 />
                   }}
                </div>
              </Col.Grid>
            </Row.Grid>
          </div>
        </Col.Grid>
      </Row.Grid>
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div className=Styles.infoContainer>
            <Heading value="Information" size=Heading.H4 style=Styles.infoHeader marginBottom=24 />
            <Row.Grid marginBottom=24>
              <Col.Grid col=Col.Six>
                <div className={CssHelper.flexBox()}>
                  <Heading value="Owner" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  //TODO: remove mock message later
                  <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                <VSpacing size=Spacing.sm />
                {switch (oracleScriptSub) {
                 | Data({owner}) =>
                   <AddressRender address=owner position=AddressRender.Subtitle />
                 | _ => <LoadingCensorBar width=284 height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six>
                <div className={Css.merge([CssHelper.flexBox(), Styles.containerSpacingSm])}>
                  <Heading value="Data Sources" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  //TODO: remove mock message later
                  <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                <VSpacing size=Spacing.sm />
                <div className={CssHelper.flexBox()}>
                  {switch (oracleScriptSub) {
                   | Data({relatedDataSources}) =>
                     //TODO: it will be correct after we get the actual data
                     relatedDataSources
                     ->Belt.List.mapWithIndex((i, id) =>
                         <>
                           <div key={i |> string_of_int}>
                             <CTooltip
                               mobile=false
                               tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                               <TypeID.DataSource id position=TypeID.Subtitle />
                             </CTooltip>
                           </div>
                           <HSpacing size=Spacing.xs />
                         </>
                       )
                     ->Belt.List.toArray
                     ->React.array
                   | _ => <LoadingCensorBar width=284 height=15 />
                   }}
                </div>
              </Col.Grid>
            </Row.Grid>
            <Heading value="Description" size=Heading.H5 marginBottom=16 />
            {switch (oracleScriptSub) {
             | Data({description}) =>
               <p> <Text value=description size=Text.Lg color=Colors.gray7 block=true /> </p>
             | _ => <LoadingCensorBar width=284 height=15 />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      <Tab
        tabs=[|
          {
            name: "Requests",
            route:
              oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptRequests),
          },
          {
            name: "OWASM Code",
            route: oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptCode),
          },
          {
            name: "Bridge Code",
            route:
              oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptBridgeCode),
          },
          {
            name: "Make New Request",
            route:
              oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptExecute),
          },
          // {
          //   name: "REVISIONS",
          //   route:
          //     oracleScriptID |> ID.OracleScript.getRouteWithTab(_, Route.OracleScriptRevisions),
          // },
        |]
        currentRoute={oracleScriptID |> ID.OracleScript.getRouteWithTab(_, hashtag)}>
        {switch (hashtag) {
         | OracleScriptExecute =>
           switch (oracleScriptSub) {
           | Data({schema}) => <OracleScriptExecute id=oracleScriptID schema />
           | _ => <LoadingCensorBar width=100 height=400 style=Styles.loadingBox />
           }

         | OracleScriptCode =>
           switch (oracleScriptSub) {
           | Data({sourceCodeURL}) => <OracleScriptCode url=sourceCodeURL />
           | _ => <LoadingCensorBar width=100 height=400 style=Styles.loadingBox />
           }

         | OracleScriptBridgeCode =>
           switch (oracleScriptSub) {
           | Data({schema}) => <OracleScriptBridgeCode schema />
           | _ => <LoadingCensorBar width=100 height=400 style=Styles.loadingBox />
           }
         | OracleScriptRequests => <OracleScriptRequestTable oracleScriptID />
         | OracleScriptRevisions => <OracleScriptRevisionTable id=oracleScriptID />
         }}
      </Tab>
    </div>
  </Section>;
};
