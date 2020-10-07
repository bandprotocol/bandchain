module Styles = {
  open Css;

  let msgContainer = style([selector("> * + *", [marginLeft(`px(5))])]);
};

module SendMsg = {
  [@react.component]
  let make = (~toAddress, ~amount) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <AmountRender coins=amount />
      <Text value={j| to |j} size=Text.Md nowrap=true block=true />
      <AddressRender address=toAddress />
    </div>;
  };
};

module ReceiveMsg = {
  [@react.component]
  let make = (~fromAddress, ~amount) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <AmountRender coins=amount />
      <Text value={j| from |j} size=Text.Md nowrap=true block=true />
      <AddressRender address=fromAddress />
    </div>;
  };
};

module MultisendMsg = {
  [@react.component]
  let make = (~inputs, ~outputs) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value={inputs |> Belt_List.length |> string_of_int} weight=Text.Semibold />
      <Text value="Inputs" />
      <Text value={j| to |j} size=Text.Md nowrap=true block=true />
      <Text value={outputs |> Belt_List.length |> string_of_int} weight=Text.Semibold />
      <Text value="Outputs" />
    </div>;
  };
};

module DelegateMsg = {
  [@react.component]
  let make = (~amount) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <AmountRender coins=[amount] />
    </div>;
  };
};

module UndelegateMsg = {
  [@react.component]
  let make = (~amount) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <AmountRender coins=[amount] />
    </div>;
  };
};

module RedelegateMsg = {
  [@react.component]
  let make = (~amount) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <AmountRender coins=[amount] />
    </div>;
  };
};

module WithdrawRewardMsg = {
  [@react.component]
  let make = (~amount) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <AmountRender coins=amount />
    </div>;
  };
};

module WithdrawCommissionMsg = {
  [@react.component]
  let make = (~amount) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <AmountRender coins=amount />
    </div>;
  };
};
