module Styles = {
  open Css;

  let msgContainer = style([selector("> * + *", [marginLeft(`px(5))])]);
};

module CreateDataSourceMsg = {
  [@react.component]
  let make = (~id, ~name) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <TypeID.DataSource id />
      <Text value=name color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module EditDataSourceMsg = {
  [@react.component]
  let make = (~id, ~name) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <TypeID.DataSource id />
      {name == Config.doNotModify
         ? React.null : <Text value=name color=Colors.gray7 nowrap=true block=true ellipsis=true />}
    </div>;
  };
};

module CreateOracleScriptMsg = {
  [@react.component]
  let make = (~id, ~name) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <TypeID.OracleScript id />
      <Text value=name color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module EditOracleScriptMsg = {
  [@react.component]
  let make = (~id, ~name) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <TypeID.OracleScript id />
      {name == Config.doNotModify
         ? React.null : <Text value=name color=Colors.gray7 nowrap=true block=true ellipsis=true />}
    </div>;
  };
};

module RequestMsg = {
  [@react.component]
  let make = (~id, ~oracleScriptID, ~oracleScriptName) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <TypeID.Request id />
      <Text value={j| to |j} size=Text.Md nowrap=true block=true />
      <TypeID.OracleScript id=oracleScriptID />
      <Text value=oracleScriptName color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ReportMsg = {
  [@react.component]
  let make = (~requestID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <TypeID.Request id=requestID />
    </div>;
  };
};
