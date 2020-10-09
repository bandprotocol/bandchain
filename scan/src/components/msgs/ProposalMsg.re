module Styles = {
  open Css;

  let msgContainer = style([selector("> * + *", [marginLeft(`px(5))])]);
};

module SubmitProposal = {
  [@react.component]
  let make = (~proposer, ~title) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
       <Text value=title size=Text.Md nowrap=true block=true /> </div>;
      // TODO: Proposal ID
  };
};

module Deposit = {
  [@react.component]
  let make = (~depositor, ~amount, ~proposalID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>

        <AmountRender coins=amount />
        <Text value={j| to |j} size=Text.Md nowrap=true block=true />
        <TypeID.Proposal id=proposalID />
      </div>;
      // TODO: Proposal Name
  };
};

module Vote = {
  [@react.component]
  let make = (~voterAddress, ~proposalID, ~option) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
       <TypeID.Proposal id=proposalID /> </div>;
      // TODO: Proposal Name
  };
};
