module Styles = {
  open Css;

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
};

let getPrevDay = _ => {
  MomentRe.momentNow()
  |> MomentRe.Moment.defaultUtc
  |> MomentRe.Moment.subtract(~duration=MomentRe.duration(1., `days))
  |> MomentRe.Moment.format(Config.timestampUseFormat);
};

let getCurrentDay = _ => {
  MomentRe.momentNow() |> MomentRe.Moment.format(Config.timestampUseFormat);
};

[@react.component]
let make = () => {
  let currentTime =
    React.useContext(TimeContext.context) |> MomentRe.Moment.format(Config.timestampUseFormat);

  let (prevDayTime, setPrevDayTime) = React.useState(getPrevDay);
  let (searchTerm, setSearchTerm) = React.useState(_ => "");
  let (sortedBy, setSortedBy) = React.useState(_ => ValidatorsTable.VotingPowerDesc);
  let (isActive, setIsActive) = React.useState(_ => true);

  React.useEffect0(() => {
    let timeOutID = Js.Global.setInterval(() => {setPrevDayTime(getPrevDay)}, 60_000);
    Some(() => {Js.Global.clearInterval(timeOutID)});
  });

  let validatorsSub = ValidatorSub.getList(~isActive, ());
  let validatorsCountSub = ValidatorSub.count();
  let activeValidatorCountSub = ValidatorSub.countByActive(true);
  let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();
  let avgBlockTimeSub = BlockSub.getAvgBlockTime(prevDayTime, currentTime);
  let latestBlock = BlockSub.getLatest();
  let votesBlockSub = ValidatorSub.getListVotesBlock();

  let topPartAllSub =
    Sub.all5(
      validatorsCountSub,
      activeValidatorCountSub,
      bondedTokenCountSub,
      avgBlockTimeSub,
      latestBlock,
    );

  let allSub = Sub.all3(topPartAllSub, validatorsSub, votesBlockSub);
  let isMobile = Media.isMobile();

  <Section>
    <div className=CssHelper.container id="validatorsSection">
      <Row.Grid alignItems=Row.Center marginBottom=40 marginBottomSm=24>
        <Col.Grid>
          <Heading value="All Validators" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
          {switch (topPartAllSub) {
           | Data((validatorCount, _, _, _, _)) =>
             <Heading value={validatorCount->string_of_int ++ " In total"} size=Heading.H3 />
           | _ => <LoadingCensorBar width=65 height=21 />
           }}
        </Col.Grid>
      </Row.Grid>
      <Row.Grid marginBottom=36 marginBottomSm=24>
        <Col.Grid>
          <div className=Styles.infoContainer>
            <Row.Grid>
              <Col.Grid col=Col.Three colSm=Col.Six mbSm=48>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading
                    value="Active Validators"
                    size=Heading.H4
                    marginBottom=28
                    align=Heading.Center
                  />
                  {switch (topPartAllSub) {
                   | Data((_, isActiveValidatorCount, _, _, _)) =>
                     <>
                       <Text
                         value={isActiveValidatorCount |> string_of_int}
                         size=Text.Xxxl
                         align=Text.Center
                         block=true
                       />
                     </>
                   | _ => <LoadingCensorBar width=100 height=24 />
                   }}
                </div>
              </Col.Grid>
              <Col.Grid col=Col.Three colSm=Col.Six mbSm=48>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading
                    value="Bonded Tokens"
                    size=Heading.H4
                    marginBottom=28
                    align=Heading.Center
                  />
                  {switch (topPartAllSub) {
                   | Data((_, _, bondedTokenCount, _, _)) =>
                     <Text
                       value={bondedTokenCount->Coin.getBandAmountFromCoin |> Format.fCurrency}
                       size=Text.Xxxl
                       align=Text.Center
                       block=true
                     />
                   | _ => <LoadingCensorBar width=100 height=24 />
                   }}
                </div>
              </Col.Grid>
              <Col.Grid col=Col.Three colSm=Col.Six>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading
                    value="Inflation Rate"
                    size=Heading.H4
                    marginBottom=28
                    align=Heading.Center
                  />
                  {switch (topPartAllSub) {
                   | Data((_, _, _, _, {inflation})) =>
                     <Text
                       value={(inflation *. 100. |> Format.fPretty(~digits=2)) ++ "%"}
                       size=Text.Xxxl
                       color=Colors.gray7
                       align=Text.Center
                       block=true
                     />

                   | _ => <LoadingCensorBar width=100 height=24 />
                   }}
                </div>
              </Col.Grid>
              <Col.Grid col=Col.Three colSm=Col.Six>
                <div className={CssHelper.flexBox(~direction=`column, ())}>
                  <Heading value="24 Hour AVG " size=Heading.H4 align=Heading.Center />
                  <Heading
                    value="Block Time"
                    size=Heading.H4
                    align=Heading.Center
                    marginBottom=10
                  />
                  {switch (topPartAllSub) {
                   | Data((_, _, _, avgBlockTime, _)) =>
                     <Text
                       value={(avgBlockTime |> Format.fPretty(~digits=2)) ++ " secs"}
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
      <Row.Grid marginBottom=16 marginBottomSm=24>
        <Col.Grid col=Col.Six colSm=Col.Eight mbSm=16>
          <SearchInput placeholder="Search Validator" onChange=setSearchTerm />
        </Col.Grid>
        {isMobile
           ? <Col.Grid col=Col.Six colSm=Col.Four mbSm=16>
               <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
                 <SortableDropdown
                   sortedBy
                   setSortedBy
                   sortList=ValidatorsTable.[
                     (NameAsc, getName(NameAsc)),
                     (NameDesc, getName(NameDesc)),
                     (VotingPowerAsc, getName(VotingPowerAsc)),
                     (VotingPowerDesc, getName(VotingPowerDesc)),
                     (CommissionAsc, getName(CommissionAsc)),
                     (CommissionDesc, getName(CommissionDesc)),
                     (UptimeAsc, getName(UptimeAsc)),
                     (UptimeDesc, getName(UptimeDesc)),
                   ]
                 />
               </div>
             </Col.Grid>
           : React.null}
        <Col.Grid col=Col.Six>
          <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
            <ToggleButton isActive setIsActive />
          </div>
        </Col.Grid>
      </Row.Grid>
      <ValidatorsTable allSub searchTerm sortedBy setSortedBy />
    </div>
  </Section>;
};
