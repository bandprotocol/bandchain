module Styles = {
  open Css;
  let link = style([color(Colors.bandBlue), textDecoration(`none)]);
  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      padding(`px(24)),
      Media.mobile([padding(`px(16))]),
    ]);
  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
  let loadingBox = style([width(`percent(100.))]);
  let idContainer = style([marginBottom(`px(16)), Media.mobile([marginBottom(`px(8))])]);
  let containerSpacingSm = style([Media.mobile([marginTop(`px(16))])]);

  // Avatar Box

  let avatarContainer =
    style([
      position(`relative),
      marginRight(`px(40)),
      Media.mobile([marginRight(`zero), marginBottom(`px(16))]),
    ]);
  let rankContainer =
    style([
      backgroundColor(Colors.bandBlue),
      borderRadius(`percent(50.)),
      position(`absolute),
      right(`zero),
      bottom(`zero),
      width(`px(26)),
      height(`px(26)),
    ]);

  // Active Status
  let statusBox = isActive_ => {
    style([
      backgroundColor(isActive_ ? Colors.green4 : Colors.gray7),
      position(`relative),
      width(`px(20)),
      height(`px(20)),
      borderRadius(`percent(50.)),
    ]);
  };

  // Oracle Status
  let oracleStatusBox = isActive_ => {
    style([
      backgroundColor(isActive_ ? Colors.green4 : Colors.red4),
      borderRadius(`px(50)),
      padding2(~v=`px(2), ~h=`px(10)),
    ]);
  };

  //Mockup
  let mockup = {
    style([minHeight(`px(300))]);
  };
};

module UptimePercentage = {
  [@react.component]
  let make = (~consensusAddress) => {
    let uptimeSub = ValidatorSub.getUptime(consensusAddress);
    <>
      {switch (uptimeSub) {
       | Data(uptime) =>
         switch (uptime) {
         | Some(uptime) =>
           <Text
             value={uptime |> Format.fPercent(~digits=2)}
             size=Text.Xxxl
             align=Text.Center
             block=true
           />
         | None => <Text value="N/A" size=Text.Xxxl align=Text.Center block=true />
         }
       | _ => <LoadingCensorBar width=100 height=24 />
       }}
    </>;
  };
};

[@react.component]
let make = (~address, ~hashtag: Route.validator_tab_t) => {
  let validatorSub = ValidatorSub.get(address);
  let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();
  let oracleReportsCountSub =
    ReportSub.ValidatorReport.count(address |> Address.toOperatorBech32);
  let allSub = Sub.all3(validatorSub, bondedTokenCountSub, oracleReportsCountSub);
  let isMobile = Media.isMobile();

  <Section pbSm=0>
    <div className=CssHelper.container>
      <Heading value="Validator Details" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
      <Row.Grid marginBottom=40 marginBottomSm=16 alignItems=Row.Center>
        <Col.Grid col=Col.Nine>
          <div
            className={Css.merge([
              CssHelper.flexBox(),
              CssHelper.flexBoxSm(~direction=`column, ()),
              Styles.idContainer,
            ])}>
            <div className=Styles.avatarContainer>
              //TODO: Will get rank later

                {switch (allSub) {
                 | Data(({identity, rank, moniker}, _, _)) =>
                   <>
                     <Avatar moniker identity width=100 widthSm=80 />
                     //  <div
                     //    className={Css.merge([
                     //      Styles.rankContainer,
                     //      CssHelper.flexBox(~justify=`center, ()),
                     //    ])}>
                     //     <Text value={rank |> string_of_int} color=Colors.white />
                     //     </div>
                   </>
                 | _ => <LoadingCensorBar width=100 height=100 radius=100 />
                 }}
              </div>
            {switch (allSub) {
             | Data(({moniker}, _, _)) => <Heading size=Heading.H3 value=moniker />
             | _ => <LoadingCensorBar width=270 height=20 />
             }}
          </div>
        </Col.Grid>
        <Col.Grid col=Col.Three>
          <div
            className={Css.merge([
              CssHelper.flexBox(~justify=`flexEnd, ()),
              CssHelper.flexBoxSm(~justify=`center, ()),
            ])}>
            {switch (allSub) {
             | Data(({isActive}, _, _)) =>
               <div className={CssHelper.flexBox()}>
                 <div
                   className={Css.merge([
                     CssHelper.flexBox(~justify=`center, ()),
                     Styles.statusBox(isActive),
                   ])}>
                   <Icon name={isActive ? "fas fa-bolt" : "fas fa-moon"} color=Colors.white />
                 </div>
                 <HSpacing size=Spacing.sm />
                 <Text value={isActive ? "Active" : "Inactive"} />
               </div>
             | _ => <LoadingCensorBar width=60 height=14 />
             }}
            <HSpacing size=Spacing.md />
            {switch (allSub) {
             | Data(({oracleStatus}, _, _)) =>
               <div
                 className={Css.merge([
                   CssHelper.flexBox(~justify=`center, ()),
                   Styles.oracleStatusBox(oracleStatus),
                 ])}>
                 <Text value="Oracle" color=Colors.white />
                 <HSpacing size=Spacing.sm />
                 <Icon
                   name={oracleStatus ? "fas fa-check" : "fal fa-times"}
                   color=Colors.white
                   size=10
                 />
               </div>
             | _ => <LoadingCensorBar width=75 height=14 />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      <Row.Grid marginBottom=24 marginBottomSm=16>
        <Col.Grid>
          <div className=Styles.infoContainer>
            <Row.Grid>
              <Col.Grid col=Col.Three colSm=Col.Six mbSm=48>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading
                    value="Voting power"
                    size=Heading.H4
                    marginBottom=16
                    align=Heading.Center
                  />
                  {switch (allSub) {
                   | Data(({votingPower}, {amount}, _)) =>
                     <>
                       <Text
                         value={votingPower *. 100. /. amount |> Format.fPercent(~digits=2)}
                         size=Text.Xxxl
                         align=Text.Center
                         block=true
                       />
                     </>
                   | _ => <LoadingCensorBar width=100 height=24 />
                   }}
                  <VSpacing size=Spacing.xs />
                  {switch (allSub) {
                   | Data(({votingPower}, _, _)) =>
                     <>
                       <Text
                         value={(votingPower /. 1e6 |> Format.fPretty(~digits=0)) ++ " Band"}
                         size=Text.Lg
                         color=Colors.gray6
                         align=Text.Center
                         block=true
                       />
                     </>
                   | _ => <LoadingCensorBar width=80 height=14 />
                   }}
                </div>
              </Col.Grid>
              <Col.Grid col=Col.Three colSm=Col.Six mbSm=48>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading
                    value="Commission"
                    size=Heading.H4
                    marginBottom=27
                    align=Heading.Center
                  />
                  {switch (allSub) {
                   | Data(({commission}, _, _)) =>
                     <>
                       <Text
                         value={commission |> Format.fPercent(~digits=2)}
                         size=Text.Xxxl
                         align=Text.Center
                         block=true
                       />
                     </>
                   | _ => <LoadingCensorBar width=100 height=24 />
                   }}
                </div>
              </Col.Grid>
              <Col.Grid col=Col.Three colSm=Col.Six>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading value="Uptime" size=Heading.H4 marginBottom=27 align=Heading.Center />
                  {switch (allSub) {
                   | Data(({consensusAddress}, _, _)) => <UptimePercentage consensusAddress />
                   | _ => <LoadingCensorBar width=100 height=24 />
                   }}
                </div>
              </Col.Grid>
              <Col.Grid col=Col.Three colSm=Col.Six>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <div className={Css.merge([CssHelper.flexBox(), CssHelper.mb(~size=27, ())])}>
                    <Heading value="Oracle Reports" size=Heading.H4 align=Heading.Center />
                    <HSpacing size=Spacing.xs />
                    //TODO: remove mock message later
                    <CTooltip
                      width=100
                      tooltipPlacementSm=CTooltip.BottomRight
                      tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                      <Icon name="fal fa-info-circle" size=12 />
                    </CTooltip>
                  </div>
                  {switch (allSub) {
                   | Data((_, _, oracleReportsCount)) =>
                     <Text
                       value={oracleReportsCount |> Format.iPretty}
                       size=Text.Xxxl
                       color=Colors.gray7
                       align=Text.Center
                       block=true
                     />
                   | _ => <LoadingCensorBar width=100 height=24 />
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
              <Col.Grid col=Col.Six mbSm=24>
                <div className={CssHelper.flexBox()}>
                  <Heading value="Operator Address" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  //TODO: remove mock message later
                  <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                <VSpacing size=Spacing.sm />
                {switch (allSub) {
                 | Data(({operatorAddress}, _, _)) =>
                   <AddressRender
                     address=operatorAddress
                     position=AddressRender.Subtitle
                     accountType=`validator
                     clickable=false
                     wordBreak=true
                   />
                 | _ => <LoadingCensorBar width=284 height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six>
                <div className={CssHelper.flexBox()}>
                  <Heading value="Address" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  //TODO: remove mock message later
                  <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                <VSpacing size=Spacing.sm />
                {switch (allSub) {
                 | Data(({operatorAddress}, _, _)) =>
                   <AddressRender address=operatorAddress position=AddressRender.Subtitle />
                 | _ => <LoadingCensorBar width=284 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
            <Row.Grid marginBottom=24>
              <Col.Grid col=Col.Six mbSm=24>
                <div className={CssHelper.flexBox()}>
                  <Heading value="Commission Max Change" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  //TODO: remove mock message later
                  <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                <VSpacing size=Spacing.sm />
                {switch (allSub) {
                 | Data(({commissionMaxChange}, _, _)) =>
                   <Text
                     value={commissionMaxChange |> Format.fPercent(~digits=2)}
                     size=Text.Lg
                     block=true
                   />
                 | _ => <LoadingCensorBar width=284 height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six>
                <div className={CssHelper.flexBox()}>
                  <Heading value="Commission Max Rate" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  //TODO: remove mock message later
                  <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                <VSpacing size=Spacing.sm />
                {switch (allSub) {
                 | Data(({commissionMaxRate}, _, _)) =>
                   <Text
                     value={commissionMaxRate |> Format.fPercent(~digits=2)}
                     size=Text.Lg
                     block=true
                   />
                 | _ => <LoadingCensorBar width=284 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
            <Row.Grid marginBottom=24>
              <Col.Grid>
                <Heading value="Website" size=Heading.H5 marginBottom=16 />
                {switch (allSub) {
                 | Data(({website}, _, _)) =>
                   <a href=website target="_blank" className=Styles.link>
                     <Text value=website size=Text.Lg color=Colors.bandBlue block=true />
                   </a>
                 | _ => <LoadingCensorBar width=284 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
            <Row.Grid marginBottom=24>
              <Col.Grid>
                <Heading value="Description" size=Heading.H5 marginBottom=16 />
                {switch (allSub) {
                 | Data(({details}, _, _)) =>
                   <p> <Text value=details size=Text.Lg color=Colors.gray7 block=true /> </p>
                 | _ => <LoadingCensorBar width=284 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
          </div>
        </Col.Grid>
      </Row.Grid>
      <Row.Grid marginBottom=24>
        <Col.Grid col=Col.Four mbSm=24>
          <div className={Css.merge([Styles.infoContainer, CssHelper.px(~size=12, ())])}>
            <div
              className={Css.merge([
                CssHelper.flexBox(),
                Styles.infoHeader,
                CssHelper.px(~size=12, ()),
                CssHelper.mb(~size=14, ()),
              ])}>
              <Heading value="Bonded Token" size=Heading.H5 />
              <HSpacing size=Spacing.xs />
              //TODO: remove mock message later
              <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                <Icon name="fal fa-info-circle" size=10 />
              </CTooltip>
            </div>
            <div className={CssHelper.flexBox()}>
              {switch (allSub) {
               | Data(({operatorAddress}, _, _)) => <HistoricalBondedGraph operatorAddress />
               | _ => <LoadingCensorBar width=100 height=180 style=Styles.loadingBox />
               }}
            </div>
          </div>
        </Col.Grid>
        {isMobile
           ? React.null
           : <Col.Grid col=Col.Eight> <ValidatorStakingInfo validatorAddress=address /> </Col.Grid>}
      </Row.Grid>
      <Row.Grid marginBottom=24>
        <Col.Grid col=Col.Six mbSm=24>
          <div className=Styles.infoContainer>
            <div
              className={Css.merge([
                CssHelper.flexBox(),
                CssHelper.mb(~size=24, ()),
                Styles.infoHeader,
              ])}>
              <Heading value="Block Uptime" size=Heading.H4 />
              <HSpacing size=Spacing.xs />
              //TODO: remove mock message later
              <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
                <Icon name="fal fa-info-circle" size=10 />
              </CTooltip>
            </div>
            {switch (allSub) {
             | Data(({consensusAddress}, _, _)) => <BlockUptimeChart consensusAddress />
             | _ => <LoadingCensorBar width=400 height=90 style=Styles.loadingBox />
             }}
          </div>
        </Col.Grid>
        <Col.Grid col=Col.Six>
          <div className=Styles.infoContainer>
            <div
              className={Css.merge([
                CssHelper.flexBox(),
                CssHelper.mb(~size=24, ()),
                Styles.infoHeader,
              ])}>
              <Heading value="Oracle Data Report" size=Heading.H5 />
              <HSpacing size=Spacing.xs />
              <CTooltip
                tooltipText="Last 90 days of Report" align=`center tooltipPlacement=CTooltip.Top>
                <Icon name="fal fa-info-circle" size=10 />
              </CTooltip>
            </div>
            {switch (allSub) {
             | Data(({oracleStatus}, _, _)) =>
               <OracleDataReportChart oracleStatus operatorAddress=address />
             | _ => <LoadingCensorBar width=400 height=90 style=Styles.loadingBox />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      <Tab
        tabs=[|
          {name: "Oracle Reports", route: Route.ValidatorIndexPage(address, Route.Reports)},
          {name: "Delegators", route: Route.ValidatorIndexPage(address, Route.Delegators)},
          {
            name: "Proposed Blocks",
            route: Route.ValidatorIndexPage(address, Route.ProposedBlocks),
          },
          {name: "Reporters", route: Route.ValidatorIndexPage(address, Route.Reporters)},
        |]
        currentRoute={Route.ValidatorIndexPage(address, hashtag)}>
        {switch (hashtag) {
         | ProposedBlocks =>
           switch (validatorSub) {
           | Data(validator) =>
             <ProposedBlocksTable consensusAddress={validator.consensusAddress} />
           | _ => <ProposedBlocksTable.LoadingWithHeader />
           }
         | Delegators => <DelegatorsTable address />
         | Reports => <ReportsTable address />
         | Reporters => <ReportersTable address />
         }}
      </Tab>
    </div>
  </Section>;
};
