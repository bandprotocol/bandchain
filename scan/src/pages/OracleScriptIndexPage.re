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
  let titleSpacing = style([marginBottom(`px(8))]);
  let idCointainer = style([marginBottom(`px(16))]);
  let containerSpacingSm = style([Media.mobile([marginTop(`px(16))])]);
};

module Content = {
  [@react.component]
  let make =
      (
        ~oracleScriptSub: ApolloHooks.Subscription.variant(BandScan.OracleScriptSub.t),
        ~oracleScriptID,
        ~hashtag,
      ) => {
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
                    <Heading
                      value="Requests"
                      size=Heading.H4
                      marginBottom=8
                      align=Heading.Center
                    />
                    {switch (oracleScriptSub) {
                     | Data({requestCount}) =>
                       <Text
                         value={requestCount |> Format.iPretty}
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
                      <CTooltip
                        tooltipPlacementSm=CTooltip.BottomRight
                        tooltipText="The average time requests to this oracle script takes to resolve">
                        <Icon name="fal fa-info-circle" size=12 />
                      </CTooltip>
                    </div>
                    {switch (oracleScriptSub) {
                     | Data({responseTime: responseTimeOpt}) =>
                       <Text
                         value={
                           switch (responseTimeOpt) {
                           | Some(responseTime') => responseTime' |> Format.fPretty(~digits=2)
                           | None => "TBD"
                           }
                         }
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
              <Heading
                value="Information"
                size=Heading.H4
                style=Styles.infoHeader
                marginBottom=24
              />
              <Row.Grid marginBottom=24>
                <Col.Grid col=Col.Six>
                  <div className={CssHelper.flexBox()}>
                    <Heading value="Owner" size=Heading.H5 />
                    <HSpacing size=Spacing.xs />
                    <CTooltip tooltipText="The owner of the oracle script">
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
                    <CTooltip tooltipText="The data sources used in this oracle script">
                      <Icon name="fal fa-info-circle" size=10 />
                    </CTooltip>
                  </div>
                  <VSpacing size=Spacing.sm />
                  <div className={CssHelper.flexBox()}>
                    {switch (oracleScriptSub) {
                     | Data({relatedDataSources}) =>
                       relatedDataSources->Belt.List.size > 0
                         ? relatedDataSources
                           ->Belt.List.map(({dataSourceName, dataSourceID}) =>
                               <>
                                 <div key={dataSourceID |> ID.DataSource.toString}>
                                   <CTooltip
                                     mobile=false
                                     align=`center
                                     tooltipText={Ellipsis.format(
                                       ~text=dataSourceName,
                                       ~limit=32,
                                       (),
                                     )}>
                                     <TypeID.DataSource id=dataSourceID position=TypeID.Subtitle />
                                   </CTooltip>
                                 </div>
                                 <HSpacing size=Spacing.xs />
                               </>
                             )
                           ->Belt.List.toArray
                           ->React.array
                         : <Text value="TBD" />

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
             | _ => <LoadingCensorBar fullWidth=true height=400 />
             }

           | OracleScriptCode =>
             switch (oracleScriptSub) {
             | Data({sourceCodeURL}) => <OracleScriptCode url=sourceCodeURL />
             | _ => <LoadingCensorBar fullWidth=true height=400 />
             }

           | OracleScriptBridgeCode =>
             switch (oracleScriptSub) {
             | Data({schema}) => <OracleScriptBridgeCode schema />
             | _ => <LoadingCensorBar fullWidth=true height=400 />
             }
           | OracleScriptRequests => <OracleScriptRequestTable oracleScriptID />
           | OracleScriptRevisions => <OracleScriptRevisionTable id=oracleScriptID />
           }}
        </Tab>
      </div>
    </Section>;
  };
};

[@react.component]
let make = (~oracleScriptID, ~hashtag) => {
  let oracleScriptSub = OracleScriptSub.get(oracleScriptID);

  switch (oracleScriptSub) {
  | NoData => <NotFound />
  | _ => <Content oracleScriptSub oracleScriptID hashtag />
  };
};
