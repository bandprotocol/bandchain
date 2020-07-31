module Styles = {
  open Css;

  let fullWidth =
    style([
      width(`percent(100.0)),
      display(`flex),
      paddingLeft(`px(26)),
      paddingRight(`px(46)),
    ]);

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

  let oracleStatus = style([display(`flex), justifyContent(`center)]);
  let logo = style([width(`px(20))]);
};

let renderBody =
    (rank, validatorSub: ApolloHooks.Subscription.variant(ValidatorSub.t), votingPower) => {
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
        <Col size=0.5>
          {switch (validatorSub) {
           | Data({tokens}) =>
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
                 value={"(" ++ (votingPower |> Format.fPercent(~digits=2)) ++ ")"}
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
        <Col size=0.6>
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
        <Col size=0.2> <HSpacing size=Spacing.sm /> </Col>
        <Col size=0.6>
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
               <LoadingCensorBar width=130 height=15 />
             </>
           }}
        </Col>
        <Col size=0.1> <HSpacing size=Spacing.sm /> </Col>
        <Col size=0.5>
          <div className=Styles.oracleStatus>
            {switch (validatorSub) {
             | Data({oracleStatus}) =>
               <img src={oracleStatus ? Images.success : Images.fail} className=Styles.logo />
             | _ => <LoadingCensorBar width=20 height=20 radius=50 />
             }}
          </div>
        </Col>
      </Row>
    </div>
  </TBody>;
};

let renderBodyMobile =
    (rank, validatorSub: ApolloHooks.Subscription.variant(ValidatorSub.t), votingPower) => {
  switch (validatorSub) {
  | Data({operatorAddress, moniker, identity, tokens, commission, uptime}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("RANK", Count(rank)),
        ("VALIDATOR", Validator(operatorAddress, moniker, identity)),
        ("VOTING\nPOWER", VotingPower(tokens, votingPower)),
        ("COMMISSION", Float(commission, Some(2))),
        ("UPTIME (%)", Uptime(uptime)),
      ]
      key={rank |> string_of_int}
      idx={rank |> string_of_int}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("RANK", Loading(70)),
        ("VALIDATOR", Loading(166)),
        ("VOTING\nPOWER", Loading(166)),
        ("COMMISSION", Loading(136)),
        ("UPTIME (%)", Loading(200)),
      ]
      key={rank |> string_of_int}
      idx={rank |> string_of_int}
    />
  };
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
type sort_direction_t =
  | ASC
  | DESC;

type sort_by_t =
  | NameAsc
  | NameDesc
  | VotingPowerAsc
  | VotingPowerDesc
  | CommissionAsc
  | CommissionDesc
  | UptimeAsc
  | UptimeDesc;

let getDirection =
  fun
  | NameAsc
  | VotingPowerAsc
  | CommissionAsc
  | UptimeAsc => ASC
  | NameDesc
  | VotingPowerDesc
  | CommissionDesc
  | UptimeDesc => DESC;

let getName =
  fun
  | NameAsc => "Validator Name (A-Z)"
  | NameDesc => "Validator name (Z-A)"
  | VotingPowerAsc => "Voting Power (Low-High)"
  | VotingPowerDesc => "Voting Power (High-Low)"
  | CommissionAsc => "Commission (Low-High)"
  | CommissionDesc => "Commission (High-Low)"
  | UptimeAsc => "Uptime (Low-High)"
  | UptimeDesc => "Uptime (High-Low)";

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
        | NameAsc => compareString(b.moniker, a.moniker)
        | NameDesc => compareString(a.moniker, b.moniker)
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

[@react.component]
let make = (~allSub, ~searchTerm, ~sortedBy, ~setSortedBy) => {
  let isMobile = Media.isMobile();

  let toggle = (sortedByAsc, sortedByDesc) =>
    if (sortedBy == sortedByDesc) {
      setSortedBy(_ => sortedByAsc);
    } else {
      setSortedBy(_ => sortedByDesc);
    };

  <>
    {isMobile
       ? React.null
       : <THead>
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
               <Col size=0.5>
                 <SortableTHead
                   title="VOTING POWER"
                   asc=VotingPowerAsc
                   desc=VotingPowerDesc
                   toggle
                   sortedBy
                   tooltipItem="Sum of self-bonded and delegated tokens"
                 />
               </Col>
               <Col size=0.6>
                 <SortableTHead
                   title="COMMISSION"
                   asc=CommissionAsc
                   desc=CommissionDesc
                   toggle
                   sortedBy
                   tooltipItem="Validator service fees charged to delegators"
                 />
               </Col>
               <Col size=0.2> <HSpacing size=Spacing.sm /> </Col>
               <Col size=0.6>
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
               <Col size=0.1> <HSpacing size=Spacing.sm /> </Col>
               <Col size=0.5>
                 <Text
                   block=true
                   value="ORACLE STATUS"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.1)}
                   tooltipItem={"Oracle status" |> React.string}
                 />
               </Col>
             </Row>
           </div>
         </THead>}
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
              ->Belt_Array.map(e => {
                  let votingPower = e.votingPower /. bondedTokenCount.amount *. 100.;
                  isMobile
                    ? renderBodyMobile(e.rank, Sub.resolve(e), votingPower)
                    : renderBody(e.rank, Sub.resolve(e), votingPower);
                })
              ->React.array
            : <div className=Styles.emptyContainer>
                <Text value="No Validators" size=Text.Xxl />
              </div>}
         <VSpacing size=Spacing.lg />
       </>;
     | _ =>
       Belt_Array.make(10, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile ? renderBodyMobile(i, noData, 1.0) : renderBody(i, noData, 1.0)
         )
       ->React.array
     }}
  </>;
};
