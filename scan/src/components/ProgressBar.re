module Styles = {
  open Css;

  let barContainer =
    style([
      position(`relative),
      paddingTop(`px(20)),
      Media.mobile([display(`flex), alignItems(`center), paddingTop(`zero)]),
    ]);
  let progressOuter =
    style([
      position(`relative),
      width(`percent(100.)),
      height(`px(12)),
      borderRadius(`px(7)),
      border(`px(1), `solid, Colors.gray9),
      padding(`px(1)),
      overflow(`hidden),
    ]);
  let progressInner = (p, success) =>
    style([
      width(`percent(p)),
      height(`percent(100.)),
      borderRadius(`px(7)),
      transition(~duration=200, "all"),
      background(success ? Colors.bandBlue : Colors.red4),
    ]);
  let leftText =
    style([
      position(`absolute),
      top(`zero),
      left(`zero),
      Media.mobile([
        position(`static),
        flexGrow(0.),
        flexShrink(0.),
        flexBasis(`px(50)),
        paddingRight(`px(10)),
      ]),
    ]);
  let rightText =
    style([
      position(`absolute),
      top(`zero),
      right(`zero),
      Media.mobile([
        position(`static),
        flexGrow(0.),
        flexShrink(0.),
        flexBasis(`px(70)),
        paddingLeft(`px(10)),
      ]),
    ]);

  // uptimeBar

  let progressUptimeInner = (p, color) =>
    style([
      width(`percent(p)),
      height(`percent(100.)),
      borderRadius(`px(7)),
      background(color),
      transition(~duration=200, "all"),
    ]);
};

[@react.component]
let make = (~reportedValidators, ~minimumValidators, ~requestValidators) => {
  let progressPercentage =
    (reportedValidators * 100 |> float_of_int) /. (requestValidators |> float_of_int);
  let success = reportedValidators >= minimumValidators;

  <div className=Styles.barContainer>
    <div className=Styles.leftText>
      <Text value={"Min " ++ (minimumValidators |> Format.iPretty)} color=Colors.gray7 />
    </div>
    <div className=Styles.progressOuter>
      <div className={Styles.progressInner(progressPercentage, success)} />
    </div>
    <div className=Styles.rightText>
      <Text
        value={
          (reportedValidators |> Format.iPretty) ++ " of " ++ (requestValidators |> Format.iPretty)
        }
        color=Colors.gray7
      />
    </div>
  </div>;
};

module Uptime = {
  [@react.component]
  let make = (~percent) => {
    let color =
      if (percent == 100.) {
        Colors.bandBlue;
      } else if (percent < 100. && percent >= 79.) {
        Colors.blue11;
      } else {
        Colors.red4;
      };

    <div className=Styles.progressOuter>
      <div className={Styles.progressUptimeInner(percent, color)} />
    </div>;
  };
};

module Deposit = {
  [@react.component]
  let make = (~totalDeposit) => {
    // TODO: remove hard-coded later.
    let minDeposit = 1000.;
    let totalDeposit_ = totalDeposit |> Coin.getBandAmountFromCoins;
    let percent = totalDeposit_ /. minDeposit *. 100.;
    let formatedMinDeposit = minDeposit |> Format.fPretty(~digits=0);
    let formatedTotalDeposit = totalDeposit_ |> Format.fPretty(~digits=0);

    <div>
      <div
        className={Css.merge([
          CssHelper.mb(~size=8, ()),
          CssHelper.flexBox(~justify=`spaceBetween, ()),
        ])}>
        <Text value={j|Min Deposit $formatedMinDeposit BAND|j} color=Colors.gray7 size=Text.Lg />
        <Text
          value={j|$formatedTotalDeposit / $formatedMinDeposit|j}
          color=Colors.gray7
          size=Text.Lg
        />
      </div>
      <div className=Styles.progressOuter>
        <div className={Styles.progressUptimeInner(percent, Colors.bandBlue)} />
      </div>
    </div>;
  };
};

module Voting = {
  [@react.component]
  let make = (~percent, ~label, ~amount) => {
    <div>
      <div
        className={Css.merge([
          CssHelper.flexBox(~justify=`spaceBetween, ()),
          CssHelper.mb(~size=8, ()),
        ])}>
        <Heading value={VoteSub.toString(label, ~withSpace=true)} size=Heading.H4 />
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          <Text value={percent |> Format.fPercent(~digits=2)} size=Text.Lg block=true />
          {Media.isMobile()
             ? React.null
             : <>
                 <HSpacing size=Spacing.sm />
                 <Text
                   value={"(" ++ (amount |> Format.fPretty(~digits=2)) ++ " BAND)"}
                   size=Text.Lg
                   block=true
                   color=Colors.gray6
                 />
               </>}
        </div>
      </div>
      <div className=Styles.progressOuter>
        <div className={Styles.progressInner(percent, true)} />
      </div>
    </div>;
  };
};
