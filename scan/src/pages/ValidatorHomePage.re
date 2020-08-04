module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row)]);

  let header =
    style([display(`flex), flexDirection(`row), alignItems(`center), height(`px(50))]);

  let validatorsLogo = style([minWidth(`px(50)), marginRight(`px(10))]);
  let highlight =
    style([
      margin2(~v=`px(28), ~h=`zero),
      Media.mobile([
        selector(
          "> div",
          [flexGrow(0.), flexShrink(0.), flexBasis(`calc((`sub, `percent(50.), `px(6))))],
        ),
        selector("> div + div + div", [marginTop(`px(24))]),
      ]),
    ]);
  let valueContainer = style([display(`flex), justifyContent(`flexStart)]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let controlContainer =
    style([
      display(`flex),
      justifyContent(`flexEnd),
      alignItems(`center),
      Media.mobile([width(`percent(100.)), flexDirection(`columnReverse)]),
    ]);
  let featureContainer =
    style([
      width(`percent(30.)),
      marginRight(`px(16)),
      Media.mobile([
        display(`flex),
        width(`percent(100.)),
        alignItems(`center),
        justifyContent(`spaceBetween),
        marginRight(`zero),
      ]),
    ]);
  let searchContainer =
    style([
      Media.mobile([
        marginRight(`zero),
        display(`flex),
        flexBasis(`percent(60.)),
        alignItems(`center),
        before([
          backgroundImage(`url(Images.searchGray)),
          contentRule(`text("")),
          width(`px(15)),
          height(`px(15)),
          backgroundRepeat(`noRepeat),
          display(`block),
          backgroundPositions([`center, `center]),
          position(`absolute),
        ]),
      ]),
    ]);

  let searchBar =
    style([
      display(`flex),
      width(`percent(100.)),
      height(`px(30)),
      paddingLeft(`px(9)),
      borderRadius(`px(4)),
      border(`px(1), `solid, Colors.blueGray3),
      marginRight(`px(6)),
      Media.mobile([
        backgroundColor(Colors.transparent),
        borderRadius(`zero),
        border(`zero, `none, Colors.white),
        borderBottom(`px(1), `solid, Colors.gray8),
        placeholder([color(Colors.blueGray3)]),
        paddingLeft(`px(20)),
        focus([outlineStyle(`none)]),
      ]),
    ]);

  let sortedContainer = style([display(`flex), flexDirection(`column), alignItems(`flexEnd)]);
  let sortDrowdownContainer =
    style([position(`relative), zIndex(2), flexBasis(`percent(40.))]);
  let sortDrowdownPanel = show => {
    style([
      display(
        {
          show ? `block : `none;
        },
      ),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      backgroundColor(Colors.white),
      position(`absolute),
      right(`zero),
      top(`percent(100.)),
      width(`px(165)),
    ]);
  };
  let sortDropdownItem = isActive => {
    style([
      backgroundColor(
        {
          isActive ? Colors.blue1 : Colors.white;
        },
      ),
      display(`flex),
      alignItems(`center),
      padding2(~v=`px(8), ~h=`px(10)),
      selector("> img", [marginRight(`px(5))]),
    ]);
  };
  let sortDropdownTextItem = {
    style([
      paddingRight(`px(15)),
      after([
        contentRule(`text("")),
        backgroundImage(`url(Images.sortDown)),
        width(`px(8)),
        height(`px(8)),
        backgroundRepeat(`noRepeat),
        backgroundSize(`contain),
        display(`block),
        position(`absolute),
        top(`percent(50.)),
        right(`zero),
        transform(translateY(`percent(-50.))),
      ]),
    ]);
  };
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

module SortableDropdown = {
  [@react.component]
  let make = (~sortedBy, ~setSortedBy) => {
    let (show, setShow) = React.useState(_ => false);
    let sortList = [
      ValidatorsTable.NameAsc,
      NameDesc,
      VotingPowerAsc,
      VotingPowerDesc,
      CommissionAsc,
      CommissionDesc,
      UptimeAsc,
      UptimeDesc,
    ];
    <div className=Styles.sortDrowdownContainer>
      <div className=Styles.sortDropdownTextItem onClick={_ => setShow(prev => !prev)}>
        <Text
          block=true
          value="Sort By"
          size=Text.Md
          weight=Text.Regular
          color=Colors.gray6
          align=Text.Right
        />
      </div>
      <div className={Styles.sortDrowdownPanel(show)}>
        {sortList
         ->Belt.List.mapWithIndex((i, value) => {
             let isActive = sortedBy == value;
             <div
               key={i |> string_of_int}
               className={Styles.sortDropdownItem(isActive)}
               onClick={_ => {
                 setSortedBy(_ => value);
                 setShow(_ => false);
               }}>
               <Text
                 block=true
                 value={ValidatorsTable.getName(value)}
                 size=Text.Md
                 weight=Text.Regular
                 color={isActive ? Colors.blue7 : Colors.gray6}
               />
             </div>;
           })
         ->Array.of_list
         ->React.array}
      </div>
    </div>;
  };
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
    <Row wrap=true style=Styles.highlight>
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
    <div className=Styles.controlContainer>
      <div className=Styles.featureContainer>
        <div className=Styles.searchContainer>
          <input
            type_="text"
            className=Styles.searchBar
            placeholder="Search Validator"
            onChange={event => {
              let newVal = ReactEvent.Form.target(event)##value |> String.lowercase_ascii;
              setSearchTerm(_ => newVal);
            }}
          />
        </div>
        {Media.isMobile() ? <SortableDropdown sortedBy setSortedBy /> : React.null}
      </div>
      <ToggleButton isActive setIsActive />
    </div>
    <VSpacing size=Spacing.md />
    <ValidatorsTable allSub searchTerm sortedBy setSortedBy />
    <VSpacing size=Spacing.lg />
  </>;
};
