module Styles = {
  open Css;
  let badge = color => {
    style([backgroundColor(color), padding2(~v=`px(3), ~h=`px(10)), borderRadius(`px(50))]);
  };
};

let getBadgeText =
  fun
  | ProposalSub.Deposit => "Deposit Period"
  | Voting => "Voting Period"
  | Passed => "Passed"
  | Rejected => "Rejected"
  | Failed => "Failed";

let getBadgeColor =
  fun
  | ProposalSub.Deposit
  | Voting => Colors.bandBlue
  | Passed => Colors.green4
  | Rejected
  | Failed => Colors.red4;

[@react.component]
let make = (~status) => {
  <div
    className={Css.merge([
      Styles.badge(getBadgeColor(status)),
      CssHelper.flexBox(~justify=`center, ()),
    ])}>
    <Text value={getBadgeText(status)} color=Colors.white />
  </div>;
};
