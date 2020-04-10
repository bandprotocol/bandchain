module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let validatorsLogo = style([marginRight(`px(10))]);
  let highlight = style([margin2(~v=`px(28), ~h=`zero)]);
  let valueContainer = style([display(`flex), justifyContent(`flexStart)]);
  let monikerContainer = style([maxWidth(`px(180))]);

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
          backgroundColor(Colors.purple1),
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

let renderBody = (rank, validator: ValidatorSub.t, bondedTokenCount) => {
  let votingPower = validator.votingPower;
  let token = validator.tokens;
  let commission = validator.commission;
  let uptime = validator.nodeStatus.uptime;
  let allRequestCount =
    validator.completedRequestCount + validator.missedRequestCount |> float_of_int;
  let reportRate = (validator.completedRequestCount |> float_of_int) /. allRequestCount *. 100.;

  <TBody key={rank |> string_of_int}>
    <div className=Styles.fullWidth>
      <Row>
        <Col size=0.8 alignSelf=Col.Start>
          <Col size=1.6 alignSelf=Col.Start>
            <Text
              value={rank |> string_of_int}
              color=Colors.gray7
              code=true
              weight=Text.Regular
              spacing={Text.Em(0.02)}
              block=true
              size=Text.Md
            />
          </Col>
        </Col>
        <Col size=1.9 alignSelf=Col.Start>
          <div className=Styles.monikerContainer>
            <ValidatorMonikerLink
              validatorAddress={validator.operatorAddress}
              moniker={validator.moniker}
            />
          </div>
        </Col>
        <Col size=1.4 alignSelf=Col.Start>
          <div>
            <Text
              value={token |> Format.fPretty}
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
                ++ (votingPower /. bondedTokenCount *. 100.)
                   ->Js.Float.toFixedWithPrecision(~digits=2)
                ++ "%)"
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
        </Col>
        <Col size=1.2 alignSelf=Col.Start>
          <Text
            value={commission->Js.Float.toFixedWithPrecision(~digits=2)}
            color=Colors.gray7
            code=true
            weight=Text.Regular
            spacing={Text.Em(0.02)}
            block=true
            align=Text.Right
            size=Text.Md
          />
        </Col>
        <Col size=1.1 alignSelf=Col.Start>
          <Text
            value={uptime->Js.Float.toFixedWithPrecision(~digits=2)}
            color=Colors.gray7
            code=true
            weight=Text.Regular
            spacing={Text.Em(0.02)}
            block=true
            align=Text.Right
            size=Text.Md
          />
        </Col>
      </Row>
    </div>
  </TBody>;
  // <Col size=1.2 alignSelf=Col.Start>
  //   <Text
  //     value={reportRate->Js.Float.toFixedWithPrecision(~digits=2)}
  //     color=Colors.gray7
  //     code=true
  //     weight=Text.Regular
  //     spacing={Text.Em(0.02)}
  //     block=true
  //     align=Text.Right
  //     size=Text.Md
  //   />
  // </Col>
};

[@react.component]
let make = () =>
  {
    let (page, setPage) = React.useState(_ => 1);

    let (prevDay, setPrevDay) =
      React.useState(_ =>
        (
          MomentRe.momentNow()
          |> MomentRe.Moment.subtract(~duration=MomentRe.duration(1., `days))
          |> MomentRe.Moment.toUnix
          |> float_of_int
        )
        *. 1000.
      );
    React.useEffect0(() => {
      let newPrevDay =
        (
          MomentRe.momentNow()
          |> MomentRe.Moment.subtract(~duration=MomentRe.duration(1., `days))
          |> MomentRe.Moment.toUnix
          |> float_of_int
        )
        *. 1000.;
      let timeOutId = Js.Global.setTimeout(() => setPrevDay(_ => newPrevDay), 60_000);
      Some(() => Js.Global.clearTimeout(timeOutId));
    });

    let pageSize = 10;

    let validatorsCountSub = ValidatorSub.count();
    let validatorsSub = ValidatorSub.getList(~page, ~pageSize, ());
    // TODO: Update once bonding status is available
    let bondedValidatorCountSub = ValidatorSub.count();
    let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();
    let pastDayBlockCountSub = BlockSub.pastDayCount(lastDay);
    let metadataSub = MetadataSub.use();

    let%Sub validators = validatorsSub;
    let%Sub validatorCount = validatorsCountSub;
    let%Sub bondedValidatorCount = bondedValidatorCountSub;
    let%Sub bondedTokenCount = bondedTokenCountSub;
    let%Sub pastDayBlockCount = pastDayBlockCountSub;
    let%Sub metadata = metadataSub;

    let pageCount = Page.getPageCount(validatorCount, pageSize);
    let globalInfo = ValidatorSub.GlobalInfo.getGlobalInfo();
    let unbondedValidatorCount = 0;
    let unbondingValidatorCount = 0;
    let allValidatorCount =
      bondedValidatorCount + unbondedValidatorCount + unbondingValidatorCount;

    //TODO: Replace with real value
    let allBondedAmount = 400;

    let pastDayAvgBlockTime = (pastDayBlockCount |> float_of_int) /. 86_400_000.00;

    <>
      <Row justify=Row.Between>
        <Col>
          <div className=Styles.vFlex>
            <img src=Images.validators className=Styles.validatorsLogo />
            <Text
              value="ALL VALIDATORS"
              weight=Text.Medium
              size=Text.Md
              nowrap=true
              color=Colors.gray7
              spacing={Text.Em(0.06)}
            />
            <div className=Styles.seperatedLine />
            <Text value={(allValidatorCount |> string_of_int) ++ " In total"} />
          </div>
        </Col>
      </Row>
      // <Col> <ToggleButton isActive setIsActive /> </Col>
      <div className=Styles.highlight>
        <Row>
          <Col size=0.7>
            <InfoHL
              info={InfoHL.Fraction(validators->Belt.Array.size, allValidatorCount, false)}
              header="VALIDATORS"
            />
          </Col>
          <Col size=1.1>
            <InfoHL
              info={
                InfoHL.Fraction(bondedTokenCount |> int_of_float, globalInfo.totalSupply, true)
              }
              header="BONDED TOKENS"
            />
          </Col>
          <Col size=0.9>
            <InfoHL
              info={InfoHL.FloatWithSuffix(metadata.inflationRate, "  %", 2)}
              header="INFLATION RATE"
            />
          </Col>
          <Col size=0.51>
            <InfoHL
              info={InfoHL.FloatWithSuffix(pastDayAvgBlockTime, "  secs", 6)}
              header="24 HOUR AVG BLOCK TIME"
            />
          </Col>
        </Row>
      </div>
      <THead>
        <div className=Styles.fullWidth>
          <Row>
            {[
               ("RANK", 0.8),
               ("VALIDATOR", 1.9),
               ("VOTING POWER (BAND)", 1.4),
               ("COMMISSION (%)", 1.2),
               ("UPTIME (%)", 1.1),
               //  ("REPORT RATE (%)", 1.2),
             ]
             ->Belt.List.mapWithIndex((idx, (title, size)) => {
                 <Col size key=title>
                   <Text
                     block=true
                     value=title
                     size=Text.Sm
                     weight=Text.Semibold
                     align=?{idx > 1 ? Some(Text.Right) : None}
                     color=Colors.gray6
                     spacing={Text.Em(0.1)}
                   />
                 </Col>
               })
             ->Array.of_list
             ->React.array}
          </Row>
        </div>
      </THead>
      {if (validators->Belt_Array.size > 0) {
         validators
         ->Belt_Array.mapWithIndex((idx, validator) =>
             renderBody(idx + 1 + (page - 1) * pageSize, validator, bondedTokenCount)
           )
         ->React.array;
       } else {
         <div className=Styles.emptyContainer> <Text value="No Validators" size=Text.Xxl /> </div>;
       }}
      <VSpacing size=Spacing.lg />
      <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
      <VSpacing size=Spacing.lg />
    </>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
