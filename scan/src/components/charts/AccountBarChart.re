module Styles = {
  open Css;
  let barChart =
    style([
      width(`percent(100.)),
      height(`px(10)),
      borderRadius(`px(50)),
      overflow(`hidden),
    ]);
  let barItem = (width_, color_) =>
    style([width(`percent(width_)), height(`percent(100.)), backgroundColor(color_)]);
};

[@react.component]
let make = (~availableBalance, ~balanceAtStake, ~reward, ~unbonding, ~commission) => {
  let totalBalance = availableBalance +. balanceAtStake +. unbonding +. reward +. commission;
  let availableBalancePercent = totalBalance == 0. ? 0. : 100. *. availableBalance /. totalBalance;
  let balanceAtStakePercent = totalBalance == 0. ? 0. : 100. *. balanceAtStake /. totalBalance;
  let unbondingPercent = totalBalance == 0. ? 0. : 100. *. unbonding /. totalBalance;
  let rewardPercent = totalBalance == 0. ? 0. : 100. *. reward /. totalBalance;
  let commissionPercent = totalBalance == 0. ? 0. : 100. *. commission /. totalBalance;

  <div className={Css.merge([Styles.barChart, CssHelper.flexBox()])}>
    <div className={Styles.barItem(availableBalancePercent, Colors.bandBlue)} />
    <div className={Styles.barItem(balanceAtStakePercent, Colors.chartBalanceAtStake)} />
    <div className={Styles.barItem(unbondingPercent, Colors.blue4)} />
    <div className={Styles.barItem(rewardPercent, Colors.chartReward)} />
    {commission == 0.
       ? React.null : <div className={Styles.barItem(commissionPercent, Colors.gray6)} />}
  </div>;
};
