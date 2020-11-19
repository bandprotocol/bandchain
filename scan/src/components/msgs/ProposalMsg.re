module Styles = {
  open Css;

  let msgContainer = style([selector("> * + *", [marginLeft(`px(5))])]);
};

module SubmitProposal = {
  [@react.component]
  let make = (~proposalID, ~title) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <TypeID.Proposal id=proposalID />
      <Text value=title size=Text.Md nowrap=true block=true />
    </div>;
  };
};

module Deposit = {
  [@react.component]
  let make = (~amount, ~proposalID, ~title) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <AmountRender coins=amount />
      <Text value={j| to |j} size=Text.Md nowrap=true block=true />
      <TypeID.Proposal id=proposalID />
      <Text value=title size=Text.Md nowrap=true block=true />
    </div>;
  };
};

module Vote = {
  [@react.component]
  let make = (~proposalID, ~title) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <TypeID.Proposal id=proposalID />
      <Text value=title size=Text.Md nowrap=true block=true />
    </div>;
  };
};
