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
     | AddReporterMsgSuccess({validator, reporter}) => <ValidatorMsg.AddReporter reporter />
     | RemoveReporterMsgSuccess({reporter}) => <ValidatorMsg.RemoveReporter reporter />
     | CreateValidatorMsgSuccess({moniker}) => <ValidatorMsg.CreateValidator moniker />
     | EditValidatorMsgSuccess({moniker}) => <ValidatorMsg.EditValidator moniker />
     | CreateClientMsg({address, clientID, chainID}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="CREATE CLIENT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=clientID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | UpdateClientMsg({address, clientID, chainID}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="UPDATE CLIENT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=clientID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | SubmitClientMisbehaviourMsg({address, clientID, chainID}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="SUBMIT CLIENT MISBEHAVIOUR"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=clientID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ConnectionOpenInitMsg({signer, common: {connectionID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 140)}>
             <Text
               value="CONNECTION OPEN INIT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=connectionID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ConnectionOpenTryMsg({signer, common: {connectionID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 120)}>
             <Text
               value="CONNECTION OPEN TRY"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=connectionID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ConnectionOpenAckMsg({signer, common: {connectionID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 130)}>
             <Text
               value="CONNECTION OPEN ACK"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=connectionID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ConnectionOpenConfirmMsg({signer, common: {connectionID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 140)}>
             <Text
               value="CONNECTION OPEN CONFIRM"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=connectionID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelOpenInitMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="CHANNEL OPEN INIT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelOpenTryMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="CHANNEL OPEN TRY"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelOpenAckMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="CHANNEL OPEN ACK"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelOpenConfirmMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 130)}>
             <Text
               value="CHANNEL OPEN CONFIRM"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelCloseInitMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="CHANNEL CLOSE INIT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelCloseConfirmMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="CHANNEL CLOSE CONFIRM"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | PacketMsg({sender, data, common: {chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 50)}>
             <Text
               value="PACKET"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value="data"
           color=Colors.gray7
           weight=Text.Semibold
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <HSpacing size=Spacing.sm />
         <div className={Styles.withWidth(110)}>
           <Text value=data color=Colors.gray7 code=true nowrap=true block=true ellipsis=true />
         </div>
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | TimeoutMsg({sender, common: {chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="TIMEOUT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value="at"
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | UnjailMsgSuccess(_) => React.null
     | SetWithdrawAddressMsgSuccess({withdrawAddress}) =>
       <ValidatorMsg.SetWithdrawAddress withdrawAddress />
     | SubmitProposalMsgSuccess({proposalID, title}) =>
       <ProposalMsg.SubmitProposal proposalID title />
     | DepositMsgSuccess({amount, proposalID, title}) =>
       <ProposalMsg.Deposit amount proposalID title />
     | VoteMsgSuccess({voterAddress, proposalID, option, title}) =>
       <ProposalMsg.Vote voterAddress proposalID option title />
     | ActivateMsgSuccess(_) => React.null
     | _ => React.null
     }}
  </div>;
};
