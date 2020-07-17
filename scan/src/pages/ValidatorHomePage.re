module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let header =
    style([display(`flex), flexDirection(`row), alignItems(`center), height(`px(50))]);

  let validatorsLogo = style([minWidth(`px(50)), marginRight(`px(10))]);
  let highlight = style([margin2(~v=`px(28), ~h=`zero)]);
  let valueContainer = style([display(`flex), justifyContent(`flexStart)]);

  let emptyContainer =
    style([
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      height(`px(300)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
    ]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let fullWidth =
    style([
      width(`percent(100.0)),
      display(`flex),
      paddingLeft(`px(26)),
      paddingRight(`px(46)),
    ]);

  let sortableTHead = isRight =>
    style([
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      cursor(`pointer),
      justifyContent(isRight ? `flexEnd : `flexStart),
    ]);

  let sort = style([width(`px(10))]);
  let downIcon = down =>
    style([
      width(`px(8)),
      marginLeft(`pxFloat(1.6)),
      transform(`rotate(`deg(down ? 0. : 180.))),
    ]);
  let searchContainer = style([display(`flex), justifyContent(`flexEnd), alignItems(`center)]);
  let searchBar =
    style([
      display(`flex),
      width(`percent(30.)),
      height(`px(30)),
      paddingLeft(`px(9)),
      borderRadius(`px(4)),
      border(`px(1), `solid, Colors.blueGray3),
    ]);

  let oracleStatus = style([display(`flex), justifyContent(`center)]);
  let logo = style([width(`px(20))]);
};

module ToggleButton = {
  open Css;

  [@react.component]
  let make = (~isActive, ~setIsActive) => {
    <div className={style([display(`flex), alignItems(`center)])}>
      <div
        onClick={_ => setIsActive(_ => true)}
        className={style([display(`flex), cursor(`pointer)])}>
        <Text value="Active" color=Colors.purple8 />
      </div>
      <HSpacing size=Spacing.sm />
      <div
        className={style([
          display(`flex),
          justifyContent(isActive ? `flexStart : `flexEnd),
          backgroundColor(Colors.gray2),
          borderRadius(`px(15)),
          padding2(~v=`px(2), ~h=`px(3)),
          width(`px(45)),
          cursor(`pointer),
          boxShadow(
            Shadow.box(
              ~inset=true,
              ~x=`zero,
              ~y=`zero,
              ~blur=`px(4),
              isActive ? Colors.purple2 : Colors.gray7,
            ),
          ),
        ])}
        onClick={_ => setIsActive(oldVal => !oldVal)}>
        <img
          src={isActive ? Images.activeValidatorLogo : Images.inactiveValidatorLogo}
          className={style([width(`px(15))])}
        />
      </div>
      <HSpacing size=Spacing.sm />
      <div
        onClick={_ => setIsActive(_ => false)}
        className={style([display(`flex), cursor(`pointer)])}>
        <Text value="Inactive" />
      </div>
    </div>;
  };
};

let renderBody =
    (rank, validatorSub: ApolloHooks.Subscription.variant(ValidatorSub.t), bondedTokenCount) => {
  <TBody
    key={
      switch (validatorSub) {
      | Data({operatorAddress}) => operatorAddress |> Address.toOperatorBech32
      | _ => rank |> string_of_int
      }
    }
    minHeight=60>

      <div className=Styles.fullWidth>
        <Row alignItems=`center>
          <Col size=0.4>
            {switch (validatorSub) {
             | Data(_) =>
               <Text
                 value={rank |> string_of_int}
                 color=Colors.gray7
                 code=true
                 weight=Text.Regular
                 spacing={Text.Em(0.02)}
                 block=true
                 size=Text.Md
               />
             | _ => <LoadingCensorBar width=20 height=15 />
             }}
          </Col>
          <Col size=0.9>
            {switch (validatorSub) {
             | Data({operatorAddress, moniker, identity}) =>
               <ValidatorMonikerLink
                 validatorAddress=operatorAddress
                 moniker
                 identity
                 width={`px(180)}
               />
             | _ => <LoadingCensorBar width=150 height=15 />
             }}
          </Col>
          <Col size=0.7>
            {switch (validatorSub) {
             | Data({tokens, votingPower}) =>
               <div>
                 <Text
                   value={tokens |> Coin.getBandAmountFromCoin |> Format.fPretty(~digits=0)}
                   color=Colors.gray7
                   code=true
                   weight=Text.Regular
                   spacing={Text.Em(0.02)}
                   block=true
                   align=Text.Right
                   size=Text.Md
                 />
                 <VSpacing size=Spacing.sm />
                 <Text
                   value={
                     "("
                     ++ (votingPower /. bondedTokenCount *. 100. |> Format.fPercent(~digits=2))
                     ++ ")"
                   }
                   color=Colors.gray6
                   code=true
                   weight=Text.Thin
                   spacing={Text.Em(0.02)}
                   block=true
                   align=Text.Right
                   size=Text.Md
                 />
               </div>
             | _ =>
               <>
                 <LoadingCensorBar width=100 height=15 isRight=true />
                 <VSpacing size=Spacing.sm />
                 <LoadingCensorBar width=40 height=15 isRight=true />
               </>
             }}
          </Col>
          <Col size=0.8>
            {switch (validatorSub) {
             | Data({commission}) =>
               <Text
                 value={commission |> Format.fPercent(~digits=2)}
                 color=Colors.gray7
                 code=true
                 weight=Text.Regular
                 spacing={Text.Em(0.02)}
                 block=true
                 align=Text.Right
                 size=Text.Md
               />
             | _ => <LoadingCensorBar width=70 height=15 isRight=true />
             }}
          </Col>
          <Col size=0.3> <HSpacing size=Spacing.sm /> </Col>
          <Col size=1.1>
            {switch (validatorSub) {
             | Data({uptime}) =>
               switch (uptime) {
               | Some(uptime') =>
                 <>
                   <Text
                     value={uptime' |> Format.fPercent(~digits=2)}
                     color=Colors.gray7
                     code=true
                     weight=Text.Regular
                     spacing={Text.Em(0.02)}
                     block=true
                     size=Text.Md
                   />
                   <VSpacing size=Spacing.sm />
                   <UptimeBar percent=uptime' />
                 </>
               | None =>
                 <Text
                   value="N/A"
                   color=Colors.gray7
                   code=true
                   weight=Text.Regular
                   spacing={Text.Em(0.02)}
                   block=true
                   size=Text.Md
                 />
               }
             | _ =>
               <>
                 <LoadingCensorBar width=50 height=15 />
                 <VSpacing size=Spacing.sm />
                 <LoadingCensorBar width=220 height=15 />
               </>
             }}
          </Col>
        </Row>
      </div>
    </TBody>;
    // <Col size=0.1> <HSpacing size=Spacing.sm /> </Col>
    // <Col size=0.5>
    //   <div className=Styles.oracleStatus>
    //     {switch (validatorSub) {
    //      | Data({oracleStatus}) =>
    //        <img src={oracleStatus ? Images.success : Images.fail} className=Styles.logo />
    //      | _ => <LoadingCensorBar width=20 height=20 radius=50 />
    //      }}
    //   </div>
    // </Col>
};

let addUptimeOnValidators =
    (validators: array(ValidatorSub.t), votesBlock: array(ValidatorSub.validator_vote_t)) => {
  validators->Belt.Array.map(validator => {
    let signedBlock =
      votesBlock
      ->Belt.Array.keep(({consensusAddress, voted}) =>
          validator.consensusAddress == consensusAddress && voted == true
        )
      ->Belt.Array.get(0)
      ->Belt.Option.mapWithDefault(0, ({count}) => count)
      |> float_of_int;

    let missedBlock =
      votesBlock
      ->Belt.Array.keep(({consensusAddress, voted}) =>
          validator.consensusAddress == consensusAddress && voted == false
        )
      ->Belt.Array.get(0)
      ->Belt.Option.mapWithDefault(0, ({count}) => count)
      |> float_of_int;

    {
      ...validator,
      uptime:
        signedBlock == 0. && missedBlock == 0.
          ? None : Some(signedBlock /. (signedBlock +. missedBlock) *. 100.),
    };
  });
};

type sort_by_t =
  | NameAsc
  | NameDesc
  | VotingPowerAsc
  | VotingPowerDesc
  | CommissionAsc
  | CommissionDesc
  | UptimeAsc
  | UptimeDesc;

let compareString = (a, b) => Js.String.localeCompare(a, b) |> int_of_float;

let defaultCompare = (a: ValidatorSub.t, b: ValidatorSub.t) =>
  if (a.tokens != b.tokens) {
    compare(b.tokens, a.tokens);
  } else {
    compareString(b.moniker, a.moniker);
  };

let sorting = (validators: array(ValidatorSub.t), sortedBy) => {
  validators
  ->Belt.List.fromArray
  ->Belt.List.sort((a, b) => {
      let result = {
        switch (sortedBy) {
        | NameAsc => compareString(a.moniker, b.moniker)
        | NameDesc => compareString(b.moniker, a.moniker)
        | VotingPowerAsc => compare(a.tokens, b.tokens)
        | VotingPowerDesc => compare(b.tokens, a.tokens)
        | CommissionAsc => compare(a.commission, b.commission)
        | CommissionDesc => compare(b.commission, a.commission)
        | UptimeAsc =>
          compare(
            a.uptime->Belt.Option.getWithDefault(0.),
            b.uptime->Belt.Option.getWithDefault(0.),
          )
        | UptimeDesc =>
          compare(
            b.uptime->Belt.Option.getWithDefault(0.),
            a.uptime->Belt.Option.getWithDefault(0.),
          )
        };
      };
      if (result != 0) {
        result;
      } else {
        defaultCompare(a, b);
      };
    })
  ->Belt.List.toArray;
};

module SortableTHead = {
  [@react.component]
  let make =
      (
        ~title,
        ~asc,
        ~desc,
        ~toggle,
        ~sortedBy,
        ~isRight=true,
        ~tooltipItem=?,
        ~tooltipPlacement=Text.AlignBottomStart,
      ) => {
    <div className={Styles.sortableTHead(isRight)} onClick={_ => toggle(asc, desc)}>
      <Text
        block=true
        value=title
        size=Text.Sm
        weight=Text.Semibold
        color=Colors.gray6
        spacing={Text.Em(0.1)}
        tooltipItem={tooltipItem->Belt_Option.mapWithDefault(React.null, React.string)}
        tooltipPlacement
      />
      <HSpacing size=Spacing.xs />
      {if (sortedBy == asc) {
         <img src=Images.sortDown className={Styles.downIcon(false)} />;
       } else if (sortedBy == desc) {
         <img src=Images.sortDown className={Styles.downIcon(true)} />;
       } else {
         <img src=Images.sort className=Styles.sort />;
       }}
    </div>;
  };
};

module ValidatorList = {
  [@react.component]
  let make = (~allSub, ~searchTerm) => {
    let (sortedBy, setSortedBy) = React.useState(_ => VotingPowerDesc);

    let toggle = (sortedByAsc, sortedByDesc) =>
      if (sortedBy == sortedByDesc) {
        setSortedBy(_ => sortedByAsc);
      } else {
        setSortedBy(_ => sortedByDesc);
      };

    <>
      <THead>
        <div className=Styles.fullWidth>
          <Row>
            <Col size=0.4 key="RANK">
              <Text
                block=true
                value="RANK"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.gray6
                spacing={Text.Em(0.1)}
              />
            </Col>
            <Col size=0.9>
              <SortableTHead
                title="VALIDATOR"
                asc=NameAsc
                desc=NameDesc
                toggle
                sortedBy
                isRight=false
              />
            </Col>
            <Col size=0.7>
              <SortableTHead
                title="VOTING POWER"
                asc=VotingPowerAsc
                desc=VotingPowerDesc
                toggle
                sortedBy
                tooltipItem="Sum of self-bonded and delegated tokens"
              />
            </Col>
            <Col size=0.8>
              <SortableTHead
                title="COMMISSION"
                asc=CommissionAsc
                desc=CommissionDesc
                toggle
                sortedBy
                tooltipItem="Validator service fees charged to delegators"
              />
            </Col>
            <Col size=0.3> <HSpacing size=Spacing.sm /> </Col>
            <Col size=1.1>
              <SortableTHead
                title="UPTIME"
                asc=UptimeAsc
                desc=UptimeDesc
                toggle
                sortedBy
                isRight=false
                tooltipItem="Percentage of the blocks that the validator is active for out of the last 250"
              />
            </Col>
          </Row>
        </div>
      </THead>
      // <Col size=0.1> <HSpacing size=Spacing.sm /> </Col>
      // <Col size=0.5>
      //   <Text
      //     block=true
      //     value="ORACLE STATUS"
      //     size=Text.Sm
      //     weight=Text.Semibold
      //     color=Colors.gray6
      //     spacing={Text.Em(0.1)}
      //     tooltipItem={"Oracle status" |> React.string}
      //   />
      // </Col>
      {switch (allSub) {
       | ApolloHooks.Subscription.Data((
           (_, _, bondedTokenCount: Coin.t, _, _),
           rawValidators,
           votesBlock,
         )) =>
         let validators = addUptimeOnValidators(rawValidators, votesBlock);
         let filteredValidator =
           searchTerm |> Js.String.length == 0
             ? validators
             : validators->Belt_Array.keep(validator => {
                 Js.String.includes(searchTerm, validator.moniker |> Js.String.toLowerCase)
               });
         <>
           {filteredValidator->Belt_Array.size > 0
              ? filteredValidator
                ->sorting(sortedBy)
                ->Belt_Array.map(e =>
                    renderBody(e.rank, Sub.resolve(e), bondedTokenCount.amount)
                  )
                ->React.array
              : <div className=Styles.emptyContainer>
                  <Text value="No Validators" size=Text.Xxl />
                </div>}
           <VSpacing size=Spacing.lg />
         </>;
       | _ =>
         Belt_Array.make(10, ApolloHooks.Subscription.NoData)
         ->Belt_Array.mapWithIndex((i, noData) => renderBody(i, noData, 1.0))
         ->React.array
       }}
    </>;
  };
};

let getPrevDay = _ => {
  MomentRe.momentNow()
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

  React.useEffect0(() => {
    let timeOutID = Js.Global.setInterval(() => {setPrevDayTime(getPrevDay)}, 60_000);
    Some(() => {Js.Global.clearInterval(timeOutID)});
  });

  let (isActive, setIsActive) = React.useState(_ => true);

  let validatorsSub = ValidatorSub.getList(~isActive, ());
  let validatorsCountSub = ValidatorSub.count();
  let isActiveValidatorCountSub = ValidatorSub.countByActive(isActive);
  let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();
  let avgBlockTimeSub = BlockSub.getAvgBlockTime(prevDayTime, currentTime);
  let latestBlock = BlockSub.getLatest();
  let votesBlockSub = ValidatorSub.getListVotesBlock();

  let topPartAllSub =
    Sub.all5(
      validatorsCountSub,
      isActiveValidatorCountSub,
      bondedTokenCountSub,
      avgBlockTimeSub,
      latestBlock,
    );

  let allSub = Sub.all3(topPartAllSub, validatorsSub, votesBlockSub);

  <>
    <Row justify=Row.Between>
      <div className=Styles.header>
        <img src=Images.validators className=Styles.validatorsLogo />
        <Text
          value="ALL VALIDATORS"
          weight=Text.Medium
          size=Text.Md
          nowrap=true
          color=Colors.gray7
          spacing={Text.Em(0.06)}
        />
        {switch (topPartAllSub) {
         | Data((validatorCount, _, _, _, _)) =>
           <>
             <div className=Styles.seperatedLine />
             <Text value={(validatorCount |> string_of_int) ++ " In total"} />
           </>
         | _ => React.null
         }}
      </div>
    </Row>
    <div className=Styles.highlight>
      <Row>
        <Col size=0.7>
          {switch (topPartAllSub) {
           | Data((validatorCount, isActiveValidatorCount, _, _, _)) =>
             <InfoHL
               info={InfoHL.Fraction(isActiveValidatorCount, validatorCount, false)}
               header="VALIDATORS"
             />
           | _ =>
             <>
               <LoadingCensorBar width=105 height=15 />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=45 height=15 />
             </>
           }}
        </Col>
        <Col size=1.1>
          {switch (topPartAllSub) {
           | Data((_, _, bondedTokenCount, _, _)) =>
             <InfoHL
               info={InfoHL.Currency(bondedTokenCount->Coin.getBandAmountFromCoin)}
               header="BONDED TOKENS"
             />
           | _ =>
             <>
               <LoadingCensorBar width=105 height=15 />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=45 height=15 />
             </>
           }}
        </Col>
        <Col size=0.9>
          {switch (topPartAllSub) {
           | Data((_, _, _, _, {inflation})) =>
             <InfoHL
               info={InfoHL.FloatWithSuffix(inflation *. 100., "  %", 2)}
               header="INFLATION RATE"
             />
           | _ =>
             <>
               <LoadingCensorBar width=105 height=15 />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=45 height=15 />
             </>
           }}
        </Col>
        <Col size=0.51>
          {switch (topPartAllSub) {
           | Data((_, _, _, avgBlockTime, _)) =>
             <InfoHL
               info={InfoHL.FloatWithSuffix(avgBlockTime, "  secs", 2)}
               header="24 HOUR AVG BLOCK TIME"
             />
           | _ =>
             <>
               <LoadingCensorBar width=105 height=15 />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=45 height=15 />
             </>
           }}
        </Col>
      </Row>
    </div>
    <div className=Styles.searchContainer>
      <input
        type_="text"
        className=Styles.searchBar
        placeholder="Search Validator"
        onChange={event => {
          let newVal = ReactEvent.Form.target(event)##value;
          setSearchTerm(_ => newVal);
        }}
      />
      <HSpacing size=Spacing.sm />
      <Col> <ToggleButton isActive setIsActive /> </Col>
    </div>
    <VSpacing size=Spacing.md />
    <ValidatorList allSub searchTerm />
    <VSpacing size=Spacing.lg />
  </>;
};
