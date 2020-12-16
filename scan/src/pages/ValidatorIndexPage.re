module Styles = {
  open Css;
  let link = style([color(Colors.bandBlue), textDecoration(`none)]);
  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding(`px(24)),
      height(`percent(100.)),
      Media.mobile([padding(`px(16))]),
    ]);
  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
  let idContainer = style([marginBottom(`px(16)), Media.mobile([marginBottom(`px(8))])]);
  let containerSpacingSm = style([Media.mobile([marginTop(`px(16))])]);

  let infoHeaderCustom = style([margin3(~top=`px(8), ~h=`px(12), ~bottom=`px(20))]);
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
  let isMobile = Media.isMobile();
  let validatorSub = ValidatorSub.get(address);
  let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();
  let oracleReportsCountSub =
    ReportSub.ValidatorReport.count(address |> Address.toOperatorBech32);
  // for finding validator rank
  let validatorsSub = ValidatorSub.getList(~isActive=true, ());

  let allSub = Sub.all4(validatorSub, validatorsSub, bondedTokenCountSub, oracleReportsCountSub);

  <Section pbSm=0>
    <div className=CssHelper.container>
      <Heading value="Validator Details" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
      <Row marginBottom=40 marginBottomSm=16 alignItems=Row.Center>
        <Col col=Col.Nine>
          <div
            className={Css.merge([
              CssHelper.flexBox(),
              CssHelper.flexBoxSm(~direction=`column, ()),
              Styles.idContainer,
            ])}>
            <div className=Styles.avatarContainer>
              {switch (allSub) {
               | Data(({identity, moniker}, validators, _, _)) =>
                 let rankOpt =
                   validators
                   ->Belt.Array.keepMap(({moniker: m, rank}) =>
                       moniker === m ? Some(rank) : None
                     )
                   ->Belt.Array.get(0);

                 <>
                   <Avatar moniker identity width=100 widthSm=80 />
                   {switch (rankOpt) {
                    | Some(rank) =>
                      <div
                        className={Css.merge([
                          Styles.rankContainer,
                          CssHelper.flexBox(~justify=`center, ()),
                        ])}>
                        <Text value={rank |> string_of_int} color=Colors.white />
                      </div>
                    | None => React.null
                    }}
                 </>;
               | _ => <LoadingCensorBar width=100 height=100 radius=100 />
               }}
            </div>
            {switch (allSub) {
             | Data(({moniker}, _, _, _)) => <Heading size=Heading.H3 value=moniker />
             | _ => <LoadingCensorBar width=270 height=20 />
             }}
          </div>
        </Col>
        <Col col=Col.Three>
          <div
            className={Css.merge([
              CssHelper.flexBox(~justify=`flexEnd, ()),
              CssHelper.flexBoxSm(~justify=`center, ()),
            ])}>
            {switch (allSub) {
             | Data(({isActive}, _, _, _)) =>
               <div className={CssHelper.flexBox()}>
                 <div className={CssHelper.flexBox(~justify=`center, ())}>
                   <img
                     src={isActive ? Images.activeValidatorLogo : Images.inactiveValidatorLogo}
                   />
                 </div>
                 <HSpacing size=Spacing.sm />
                 <Text value={isActive ? "Active" : "Inactive"} />
               </div>
             | _ => <LoadingCensorBar width=60 height=14 />
             }}
            <HSpacing size=Spacing.md />
            {switch (allSub) {
             | Data(({oracleStatus}, _, _, _)) =>
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
        </Col>
      </Row>
      <Row marginBottom=24 marginBottomSm=16>
        <Col>
          <div className=Styles.infoContainer>
            <Row>
              <Col col=Col.Three colSm=Col.Six mbSm=48>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading
                    value="Voting power"
                    size=Heading.H4
                    marginBottom=16
                    align=Heading.Center
                  />
                  {switch (allSub) {
                   | Data(({votingPower}, _, {amount}, _)) =>
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
                   | Data(({votingPower}, _, _, _)) =>
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
              </Col>
              <Col col=Col.Three colSm=Col.Six mbSm=48>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading
                    value="Commission"
                    size=Heading.H4
                    marginBottom=27
                    align=Heading.Center
                  />
                  {switch (allSub) {
                   | Data(({commission}, _, _, _)) =>
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
              </Col>
              <Col col=Col.Three colSm=Col.Six>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading value="Uptime" size=Heading.H4 marginBottom=27 align=Heading.Center />
                  {switch (allSub) {
                   | Data(({consensusAddress}, _, _, _)) => <UptimePercentage consensusAddress />
                   | _ => <LoadingCensorBar width=100 height=24 />
                   }}
                </div>
              </Col>
              <Col col=Col.Three colSm=Col.Six>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <div className={Css.merge([CssHelper.flexBox(), CssHelper.mb(~size=27, ())])}>
                    <Heading value="Oracle Reports" size=Heading.H4 align=Heading.Center />
                    <HSpacing size=Spacing.xs />
                    <CTooltip
                      tooltipPlacementSm=CTooltip.BottomRight
                      tooltipText="The number of reports this validator has submitted to date">
                      <Icon name="fal fa-info-circle" size=12 />
                    </CTooltip>
                  </div>
                  {switch (allSub) {
                   | Data((_, _, _, oracleReportsCount)) =>
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
              </Col>
            </Row>
          </div>
        </Col>
      </Row>
      <Row marginBottom=24>
        <Col>
          <div className=Styles.infoContainer>
            <Heading value="Information" size=Heading.H4 style=Styles.infoHeader marginBottom=24 />
            <Row marginBottom=24>
              <Col col=Col.Six mbSm=24>
                <div className={Css.merge([CssHelper.flexBox(), CssHelper.mb(~size=8, ())])}>
                  <Heading value="Operator Address" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  <CTooltip tooltipText="The address used to show the validator's entity status">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                {switch (allSub) {
                 | Data(({operatorAddress}, _, _, _)) =>
                   <AddressRender
                     address=operatorAddress
                     position=AddressRender.Subtitle
                     accountType=`validator
                     clickable=false
                     wordBreak=true
                   />
                 | _ => <LoadingCensorBar width=260 height=15 />
                 }}
              </Col>
              <Col col=Col.Six>
                <div className={Css.merge([CssHelper.flexBox(), CssHelper.mb(~size=8, ())])}>
                  <Heading value="Address" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  <CTooltip tooltipText="The validator entity's unique address">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                {switch (allSub) {
                 | Data(({operatorAddress}, _, _, _)) =>
                   <AddressRender address=operatorAddress position=AddressRender.Subtitle />
                 | _ => <LoadingCensorBar width=260 height=15 />
                 }}
              </Col>
            </Row>
            <Row marginBottom=24>
              <Col col=Col.Six mbSm=24>
                <div className={Css.merge([CssHelper.flexBox(), CssHelper.mb(~size=8, ())])}>
                  <Heading value="Commission Max Change" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  <CTooltip
                    tooltipText="The maximum increment by which the validator can increase their commission rate">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                {switch (allSub) {
                 | Data(({commissionMaxChange}, _, _, _)) =>
                   <Text
                     value={commissionMaxChange |> Format.fPercent(~digits=2)}
                     size=Text.Lg
                     block=true
                   />
                 | _ => <LoadingCensorBar width=260 height=15 />
                 }}
              </Col>
              <Col col=Col.Six>
                <div className={Css.merge([CssHelper.flexBox(), CssHelper.mb(~size=8, ())])}>
                  <Heading value="Commission Max Rate" size=Heading.H5 />
                  <HSpacing size=Spacing.xs />
                  <CTooltip
                    tooltipText="Highest possible commission rate that the validator can set">
                    <Icon name="fal fa-info-circle" size=10 />
                  </CTooltip>
                </div>
                {switch (allSub) {
                 | Data(({commissionMaxRate}, _, _, _)) =>
                   <Text
                     value={commissionMaxRate |> Format.fPercent(~digits=2)}
                     size=Text.Lg
                     block=true
                   />
                 | _ => <LoadingCensorBar width=260 height=15 />
                 }}
              </Col>
            </Row>
            <Row marginBottom=24>
              <Col>
                <Heading value="Website" size=Heading.H5 marginBottom=8 />
                {switch (allSub) {
                 | Data(({website}, _, _, _)) =>
                   <a href=website target="_blank" className=Styles.link>
                     <Text value=website size=Text.Lg color=Colors.bandBlue block=true />
                   </a>
                 | _ => <LoadingCensorBar width=260 height=15 />
                 }}
              </Col>
            </Row>
            <Row>
              <Col>
                <Heading value="Description" size=Heading.H5 marginBottom=8 />
                {switch (allSub) {
                 | Data(({details}, _, _, _)) =>
                   <p> <Text value=details size=Text.Lg color=Colors.gray7 block=true /> </p>
                 | _ => <LoadingCensorBar width=260 height=15 />
                 }}
              </Col>
            </Row>
          </div>
        </Col>
      </Row>
      <Row marginBottom=24>
        <Col col=Col.Four mbSm=24>
          <div className={Css.merge([Styles.infoContainer, CssHelper.px(~size=12, ())])}>
            <div
              className={Css.merge([
                CssHelper.flexBox(),
                Styles.infoHeader,
                Styles.infoHeaderCustom,
              ])}>
              <Heading value="Bonded Token" size=Heading.H4 />
              <HSpacing size=Spacing.xs />
              <CTooltip tooltipText="Total amount of tokens current staked across all validators">
                <Icon name="fal fa-info-circle" size=10 />
              </CTooltip>
            </div>
            <div className={CssHelper.flexBox()}>
              {switch (allSub) {
               | Data(({operatorAddress}, _, _, _)) => <HistoricalBondedGraph operatorAddress />
               | _ => <LoadingCensorBar fullWidth=true height=180 />
               }}
            </div>
          </div>
        </Col>
        {isMobile
           ? React.null
           : <Col col=Col.Eight> <ValidatorStakingInfo validatorAddress=address /> </Col>}
      </Row>
      <Row marginBottom=24>
        <Col col=Col.Six mbSm=24>
          <div className=Styles.infoContainer>
            <div
              className={Css.merge([
                CssHelper.flexBox(),
                CssHelper.mb(~size=24, ()),
                Styles.infoHeader,
              ])}>
              <Heading value="Block Uptime" size=Heading.H4 />
              <HSpacing size=Spacing.xs />
              <CTooltip tooltipText="The validator's non-oracle performance">
                <Icon name="fal fa-info-circle" size=10 />
              </CTooltip>
            </div>
            {switch (allSub) {
             | Data(({consensusAddress}, _, _, _)) => <BlockUptimeChart consensusAddress />
             | _ => <LoadingCensorBar fullWidth=true height=90 />
             }}
          </div>
        </Col>
        <Col col=Col.Six>
          <div className=Styles.infoContainer>
            <div
              className={Css.merge([
                CssHelper.flexBox(),
                CssHelper.mb(~size=24, ()),
                Styles.infoHeader,
              ])}>
              <Heading value="Oracle Data Report" size=Heading.H4 />
              <HSpacing size=Spacing.xs />
              <CTooltip tooltipText="Last 90 days of Report" align=`center>
                <Icon name="fal fa-info-circle" size=10 />
              </CTooltip>
            </div>
            {switch (allSub) {
             | Data(({oracleStatus}, _, _, _)) =>
               <OracleDataReportChart oracleStatus operatorAddress=address />
             | _ => <LoadingCensorBar fullWidth=true height=90 />
             }}
          </div>
        </Col>
      </Row>
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
