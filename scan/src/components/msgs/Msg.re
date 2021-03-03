module Styles = {
  open Css;
  let rowWithWidth = (w: int) =>
    style([
      width(`px(w)),
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      Media.mobile([
        width(`auto),
        flexWrap(`wrap),
        selector("> div:nth-child(1)", [width(`px(90)), marginBottom(`px(10))]),
        selector(
          "> .labelContainer",
          [
            display(`flex),
            flexBasis(`calc((`sub, `percent(100.), `px(100)))),
            marginBottom(`px(10)),
          ],
        ),
      ]),
      Media.smallMobile([
        selector("> div:nth-child(1)", [width(`px(68)), marginBottom(`px(10))]),
      ]),
    ]);
  let withWidth = (w: int) => style([width(`px(w))]);
  let withBg = (color: Types.Color.t, mw: int) =>
    style([
      minWidth(`px(mw)),
      height(`px(16)),
      backgroundColor(color),
      borderRadius(`px(100)),
      margin2(~v=`zero, ~h=`px(5)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
    ]);

  let addressWrapper = style([width(`px(120))]);

  let msgContainer =
    style([
      Media.mobile([
        selector("> div", [width(`percent(100.))]),
        selector("> div + div", [marginTop(`px(8))]),
      ]),
    ]);
};

let makeBadge = (name, length, color1, color2) =>
  <div className="labelContainer">
    <div className={Styles.withBg(color1, length)}>
      <Text value=name size=Text.Xs spacing={Text.Em(0.07)} weight=Text.Medium color=color2 />
    </div>
  </div>;

[@react.component]
let make = (~msg: TxSub.Msg.t, ~width: int) => {
  let theme = msg |> TxSub.Msg.getBadgeTheme;
  <div
    className={Css.merge([
      CssHelper.flexBox(~wrap=`nowrap, ()),
      CssHelper.flexBoxSm(~wrap=`wrap, ()),
      CssHelper.overflowHidden,
      Styles.msgContainer,
    ])}>
    <MsgFront
      msgType={theme.category}
      name={theme.name}
      fromAddress={msg |> TxSub.Msg.getCreator}
    />
    {switch (msg) {
     | SendMsgSuccess({toAddress, amount}) => <TokenMsg.SendMsg toAddress amount />
     | ReceiveMsg({fromAddress, amount}) => <TokenMsg.ReceiveMsg fromAddress amount />
     | MultiSendMsgSuccess({inputs, outputs}) => <TokenMsg.MultisendMsg inputs outputs />
     | DelegateMsgSuccess({amount}) => <TokenMsg.DelegateMsg amount />
     | UndelegateMsgSuccess({amount}) => <TokenMsg.UndelegateMsg amount />
     | RedelegateMsgSuccess({amount}) => <TokenMsg.RedelegateMsg amount />
     | WithdrawRewardMsgSuccess({amount}) => <TokenMsg.WithdrawRewardMsg amount />
     | WithdrawCommissionMsgSuccess({amount}) => <TokenMsg.WithdrawCommissionMsg amount />
     | CreateDataSourceMsgSuccess({id, name}) => <DataMsg.CreateDataSourceMsg id name />
     | EditDataSourceMsgSuccess({id, name}) => <DataMsg.EditDataSourceMsg id name />
     | CreateOracleScriptMsgSuccess({id, name}) => <DataMsg.CreateOracleScriptMsg id name />
     | EditOracleScriptMsgSuccess({id, name}) => <DataMsg.EditOracleScriptMsg id name />
     | RequestMsgSuccess({id, oracleScriptID, oracleScriptName}) =>
       <DataMsg.RequestMsg id oracleScriptID oracleScriptName />
     | ReportMsgSuccess({requestID}) => <DataMsg.ReportMsg requestID />
     | AddReporterMsgSuccess({reporter}) => <ValidatorMsg.AddReporter reporter />
     | RemoveReporterMsgSuccess({reporter}) => <ValidatorMsg.RemoveReporter reporter />
     | CreateValidatorMsgSuccess({moniker}) => <ValidatorMsg.CreateValidator moniker />
     | EditValidatorMsgSuccess({moniker}) => <ValidatorMsg.EditValidator moniker />
     | UnjailMsgSuccess(_) => React.null
     | SetWithdrawAddressMsgSuccess({withdrawAddress}) =>
       <ValidatorMsg.SetWithdrawAddress withdrawAddress />
     | SubmitProposalMsgSuccess({proposalID, title}) =>
       <ProposalMsg.SubmitProposal proposalID title />
     | DepositMsgSuccess({amount, proposalID, title}) =>
       <ProposalMsg.Deposit amount proposalID title />
     | VoteMsgSuccess({proposalID, title}) => <ProposalMsg.Vote proposalID title />
     | ActivateMsgSuccess(_) => React.null
     | CreateClientMsg({chainID}) => <IBCMsg.CreateClient chainID />
     | UpdateClientMsg({chainID}) => <IBCMsg.UpdateClient chainID />
     | SubmitClientMisbehaviourMsg({chainID}) => <IBCMsg.SubmitClientMisbehaviour chainID />
     | ConnectionOpenInitMsg({clientID}) => <IBCMsg.ConnectionOpenInit clientID />
     | ConnectionOpenTryMsg({common: {connectionID}}) =>
       <IBCMsg.ConnectionOpenTry connectionID />
     | ConnectionOpenAckMsg({common: {connectionID}}) =>
       <IBCMsg.ConnectionOpenAck connectionID />
     | ConnectionOpenConfirmMsg({common: {connectionID}}) =>
       <IBCMsg.ConnectionOpenConfirm connectionID />
     | ChannelOpenInitMsg({common: {channelID}}) => <IBCMsg.ChannelOpenInit channelID />
     | ChannelOpenTryMsg({common: {channelID}}) => <IBCMsg.ChannelOpenTry channelID />
     | ChannelOpenAckMsg({common: {channelID}}) => <IBCMsg.ChannelOpenAck channelID />
     | ChannelOpenConfirmMsg({common: {channelID}}) => <IBCMsg.ChannelOpenConfirm channelID />
     | ChannelCloseInitMsg({common: {channelID}}) => <IBCMsg.ChannelCloseInit channelID />
     | ChannelCloseConfirmMsg({common: {channelID}}) => <IBCMsg.ChannelCloseConfirm channelID />
     | PacketMsg({data}) => <IBCMsg.Packet data />
     | TimeoutMsg({common: {chainID}}) => <IBCMsg.Timeout chainID />
     | _ => React.null
     }}
  </div>;
};
