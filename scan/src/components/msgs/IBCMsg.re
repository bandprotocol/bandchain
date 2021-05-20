module Styles = {
  open Css;

  let msgContainer = style([selector("> * + *", [marginLeft(`px(5))])]);
};

module CreateClient = {
  [@react.component]
  let make = (~chainID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=chainID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module UpdateClient = {
  [@react.component]
  let make = (~chainID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=chainID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module SubmitClientMisbehaviour = {
  [@react.component]
  let make = (~chainID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=chainID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ConnectionOpenInit = {
  [@react.component]
  let make = (~clientID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=clientID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ConnectionOpenTry = {
  [@react.component]
  let make = (~connectionID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=connectionID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ConnectionOpenAck = {
  [@react.component]
  let make = (~connectionID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=connectionID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ConnectionOpenConfirm = {
  [@react.component]
  let make = (~connectionID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=connectionID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ChannelOpenInit = {
  [@react.component]
  let make = (~channelID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=channelID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ChannelOpenTry = {
  [@react.component]
  let make = () => {
    // <div
    //   className={Css.merge([
    //     CssHelper.flexBox(~wrap=`nowrap, ()),
    //     CssHelper.overflowHidden,
    //     Styles.msgContainer,
    //   ])}>
    //   <Text value=channelID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    // </div>;
    React.null;
  };
};

module ChannelOpenAck = {
  [@react.component]
  let make = (~channelID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=channelID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ChannelOpenConfirm = {
  [@react.component]
  let make = () => {
    // <div
    //   className={Css.merge([
    //     CssHelper.flexBox(~wrap=`nowrap, ()),
    //     CssHelper.overflowHidden,
    //     Styles.msgContainer,
    //   ])}>
    //   <Text value=channelID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    // </div>;
    React.null;
  };
};

module ChannelCloseInit = {
  [@react.component]
  let make = (~channelID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=channelID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module ChannelCloseConfirm = {
  [@react.component]
  let make = (~channelID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=channelID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module Packet = {
  [@react.component]
  let make = (~data) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=data color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};

module Timeout = {
  [@react.component]
  let make = (~chainID) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~wrap=`nowrap, ()),
        CssHelper.overflowHidden,
        Styles.msgContainer,
      ])}>
      <Text value=chainID color=Colors.gray7 nowrap=true block=true ellipsis=true />
    </div>;
  };
};
