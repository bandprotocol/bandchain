module Styles = {
  open Css;

  let pie = (size, color) =>
    style([
      overflow(`hidden),
      display(`flex),
      position(`relative),
      width(`px(size)),
      height(`px(size)),
      borderRadius(`percent(100.)),
      backgroundColor(color),
    ]);

  let segment = (offset, angle, color, isHidden) =>
    style([
      visibility(isHidden ? `hidden : `visible),
      overflow(`hidden),
      width(`percent(100.)),
      height(`percent(100.)),
      position(`absolute),
      transforms([
        `translate((`zero, `percent(-50.))),
        `rotate(`deg(90.)),
        `rotate(`deg(offset)),
      ]),
      transformOrigin(`percent(50.), `percent(100.)),
      before([
        contentRule(`text("")),
        width(`percent(100.)),
        height(`percent(100.)),
        position(`absolute),
        transforms([`translate((`zero, `percent(100.))), `rotate(`deg(angle))]),
        transformOrigin(`percent(50.), `percent(0.)),
        backgroundColor(color),
      ]),
    ]);
};

let renderSegment = (offset, angle, color) =>
  <>
    <div className={Styles.segment(offset, angle <= 180.0 ? angle : 180.0, color, false)} />
    <div
      className={Styles.segment(
        offset +. 180.0,
        angle <= 180.0 ? 0.0 : angle -. 180.0,
        color,
        angle <= 180.0,
      )}
    />
  </>;

[@react.component]
let make = (~size, ~availableBalance, ~balanceAtStake, ~reward, ~unbonding, ~commission) => {
  let totalBalance = availableBalance +. balanceAtStake +. unbonding +. reward +. commission;
  let balanceAtStakeAngle = totalBalance == 0. ? 0. : 360. *. balanceAtStake /. totalBalance;
  let unbondingAngle = totalBalance == 0. ? 0. : 360. *. unbonding /. totalBalance;
  let rewardAngle = totalBalance == 0. ? 0. : 360. *. reward /. totalBalance;
  let commissionAngle = totalBalance == 0. ? 0. : 360. *. commission /. totalBalance;

  <div className={Styles.pie(size, Colors.bandBlue)}>
    {renderSegment(0., balanceAtStakeAngle, Colors.chartBalanceAtStake)}
    {renderSegment(balanceAtStakeAngle, unbondingAngle, Colors.blue4)}
    {renderSegment(balanceAtStakeAngle +. unbondingAngle, rewardAngle, Colors.chartReward)}
    {renderSegment(
       balanceAtStakeAngle +. unbondingAngle +. rewardAngle,
       commissionAngle,
       Colors.gray6,
     )}
  </div>;
};
