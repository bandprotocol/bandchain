module Styles = {
  open Css;

  let addressContainer = width_ => style([width(`px(width_))]);

  let badgeContainer = style([display(`block), marginTop(`px(-4))]);

  let badge = color =>
    style([
      display(`inlineFlex),
      padding2(~v=`px(5), ~h=`px(8)),
      backgroundColor(color),
      borderRadius(`px(50)),
    ]);

  let hFlex = style([display(`flex), alignItems(`center)]);

  let topicContainer =
    style([
      display(`flex),
      justifyContent(`spaceBetween),
      width(`percent(100.)),
      lineHeight(`px(16)),
      alignItems(`center),
    ]);

  let detailContainer = style([display(`flex), maxWidth(`px(360)), justifyContent(`flexEnd)]);

  let hashContainer =
    style([
      display(`flex),
      maxWidth(`px(350)),
      justifyContent(`flexEnd),
      wordBreak(`breakAll),
    ]);

  let firstCol = 0.45;
  let secondCol = 0.50;
  let thirdCol = 1.20;

  let failIcon = style([width(`px(16)), height(`px(16))]);

  let failedMessageDetails =
    style([
      display(`flex),
      width(`px(120)),
      alignItems(`center),
      justifyContent(`spaceBetween),
    ]);

  let separatorLine =
    style([
      borderStyle(`none),
      backgroundColor(Colors.gray9),
      height(`px(1)),
      margin2(~v=`px(10), ~h=`auto),
    ]);

  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
};

let renderAddReporter = (address: TxSub.Msg.AddReporter.success_t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={address.validatorMoniker} code=true />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="REPORTER ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AddressRender address={address.reporter} />
    </div>
  </Col>;
};

let renderRemoveReporter = (address: TxSub.Msg.RemoveReporter.success_t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={address.validatorMoniker} code=true />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="REPORTER ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AddressRender address={address.reporter} />
    </div>
  </Col>;
};

let renderCreateValidator = (validator: TxSub.Msg.CreateValidator.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="MONIKER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={validator.moniker} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="IDENTITY" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={validator.identity} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="WEBSITE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={validator.website} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="DETAILS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className=Styles.detailContainer>
        <Text value={validator.details} code=true height={Text.Px(16)} align=Text.Right />
      </div>
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="COMMISSION RATE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.commissionRate->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="COMMISSION MAX RATE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.commissionMaxRate->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="COMMISSION MAX CHANGE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.commissionMaxChange->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="DELEGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AddressRender address={validator.delegatorAddress} />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AddressRender address={validator.validatorAddress} accountType=`validator />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="PUBLIC KEY" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <PubKeyRender pubKey={validator.publicKey} />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="MIN SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins=[validator.minSelfDelegation] />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins=[validator.selfDelegation] />
    </div>
  </Col>;
};

let renderEditValidator = (validator: TxSub.Msg.EditValidator.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="MONIKER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.moniker == Config.doNotModify ? "Unchanged" : validator.moniker}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="IDENTITY" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.identity == Config.doNotModify ? "Unchanged" : validator.identity}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="WEBSITE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.website == Config.doNotModify ? "Unchanged" : validator.website}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="DETAILS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className=Styles.detailContainer>
        <Text
          value={validator.details == Config.doNotModify ? "Unchanged" : validator.details}
          code=true
          height={Text.Px(16)}
          align=Text.Right
        />
      </div>
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="COMMISSION RATE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={
          switch (validator.commissionRate) {
          | Some(rate) => rate->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"
          | None => "Unchanged"
          }
        }
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AddressRender address={validator.sender} accountType=`validator />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="MIN SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      {switch (validator.minSelfDelegation) {
       | Some(minSelfDelegation') => <AmountRender coins=[minSelfDelegation'] />
       | None => <Text value="Unchanged" code=true />
       }}
    </div>
  </Col>;
};

let renderCreateClient = (info: TxSub.Msg.CreateClient.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.clientID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="TRUSTING PERIOD" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.trustingPeriod |> MomentRe.Duration.toISOString} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="UNBOUNDING PERIOD" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.unbondingPeriod |> MomentRe.Duration.toISOString} code=true />
    </div>
  </Col>;
};

let renderUpdateClient = (info: TxSub.Msg.UpdateClient.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.clientID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text
        value="VALIDATOR HASH"
        size=Text.Sm
        weight=Text.Thin
        height={Text.Px(16)}
        spacing={Text.Em(0.06)}
      />
      <div className=Styles.hashContainer>
        <Text
          value={info.validatorHash |> Hash.toHex}
          code=true
          height={Text.Px(16)}
          align=Text.Right
        />
      </div>
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text
        value="PREVIOUS VALIDATOR HASH"
        size=Text.Sm
        weight=Text.Thin
        height={Text.Px(16)}
        spacing={Text.Em(0.06)}
      />
      <div className=Styles.hashContainer>
        <Text
          value={info.prevValidatorHash |> Hash.toHex}
          code=true
          height={Text.Px(16)}
          align=Text.Right
        />
      </div>
    </div>
  </Col>;
};

let renderSubmitClientMisbehaviour = (info: TxSub.Msg.SubmitClientMisbehaviour.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.clientID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text
        value="VALIDATOR HASH"
        size=Text.Sm
        weight=Text.Thin
        height={Text.Px(16)}
        spacing={Text.Em(0.06)}
      />
      <div className=Styles.hashContainer>
        <Text
          value={info.validatorHash |> Hash.toHex}
          code=true
          height={Text.Px(16)}
          align=Text.Right
        />
      </div>
    </div>
  </Col>;
};

let renderPacketVariant = (msg: TxSub.Msg.t, common: TxSub.Msg.Packet.common_t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="SEQUENCE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.sequence |> string_of_int} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="SOURCE PORT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.sourcePort} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="SOURCE CHANNEL" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.sourceChannel} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="DESTINATION PORT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.destinationPort} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="DESTINATION CHANNEL" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.destinationChannel} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="TIMEOUT HEIGHT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.timeoutHeight |> string_of_int} code=true />
    </div>
    {switch (msg) {
     | AcknowledgementMsg({acknowledgement}) =>
       <>
         <VSpacing size=Spacing.md />
         <div className=Styles.topicContainer>
           <Text value="ACKNOWLEDGEMENT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=acknowledgement code=true />
         </div>
       </>
     | TimeoutMsg({nextSequenceReceive}) =>
       <>
         <VSpacing size=Spacing.md />
         <div className=Styles.topicContainer>
           <Text value="ACKNOWLEDGEMENT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value={nextSequenceReceive |> string_of_int} code=true />
         </div>
       </>
     | _ => React.null
     }}
  </Col>;
};

let renderChannelVariant = (common: TxSub.Msg.ChannelCommon.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="PORT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.portID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CHANNEL ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.channelID} code=true />
    </div>
  </Col>;
};

let renderConnectionVariant = (msg: TxSub.Msg.t, common: TxSub.Msg.ConnectionCommon.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CONNECTION ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.connectionID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    {switch (msg) {
     | ConnectionOpenInitMsg({clientID}) =>
       <>
         <div className=Styles.topicContainer>
           <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=clientID code=true />
         </div>
       </>
     | ConnectionOpenTryMsg({clientID}) =>
       <>
         <div className=Styles.topicContainer>
           <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=clientID code=true />
         </div>
       </>
     | _ => React.null
     }}
  </Col>;
};

let renderUnjail = (unjail: TxSub.Msg.Unjail.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={unjail.address} accountType=`validator />
      </div>
    </div>
  </Col>;
};
let renderSubmitProposal = (proposal: TxSub.Msg.SubmitProposal.success_t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="TITLE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={proposal.title} code=true />
    </div>
    // <VSpacing size=Spacing.lg />
    //TODO: Will re-visit
    // <div className=Styles.topicContainer>
    //   <Text value="DESCRIPTION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
    //   <Text value={proposal.description} code=true />
    // </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="PROPOSER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={proposal.proposer} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins={proposal.initialDeposit} pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderDeposit = (deposit: TxSub.Msg.Deposit.success_t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="DEPOSITOR" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={deposit.depositor} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="PROPOSAL ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={deposit.proposalID |> ID.Proposal.toString} code=true />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins={deposit.amount} pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderSetWithdrawAddress = (set: TxSub.Msg.SetWithdrawAddress.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="DELEGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={set.delegatorAddress} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="WITHDRAW ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={set.withdrawAddress} />
      </div>
    </div>
  </Col>;
};

let renderVote = (vote: TxSub.Msg.Vote.success_t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VOTER ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={vote.voterAddress} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="PROPOSAL ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className=Styles.hFlex> <Text value={vote.proposalID |> ID.Proposal.toString} /> </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="OPTION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className=Styles.hFlex> <Text value={vote.option} /> </div>
    </div>
  </Col>;
};

let renderActivate = (activate: TxSub.Msg.Activate.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={activate.validatorAddress} accountType=`validator />
      </div>
    </div>
  </Col>;
};

let renderUnknownMessage = () => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="UNKNOWN MESSAGE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <img src=Images.fail className=Styles.failIcon />
    </div>
  </Col>;
};

let renderBody = (msg: TxSub.Msg.t) =>
  switch (msg) {
  | SendMsgSuccess(send)
  | SendMsgFail(send) => <IndexTokenMsg.SendMsg send />
  | DelegateMsgSuccess(delegation) => <IndexTokenMsg.DelegateMsg delegation />
  | DelegateMsgFail(delegation) => <IndexTokenMsg.DelegateFailMsg delegation />
  | UndelegateMsgSuccess(undelegation) => <IndexTokenMsg.UndelegateMsg undelegation />
  | UndelegateMsgFail(undelegation) => <IndexTokenMsg.UndelegateFailMsg undelegation />
  | RedelegateMsgSuccess(redelegation) => <IndexTokenMsg.RedelegateMsg redelegation />
  | RedelegateMsgFail(redelegation) => <IndexTokenMsg.RedelegateFailMsg redelegation />
  | WithdrawRewardMsgSuccess(withdrawal) => <IndexTokenMsg.WithdrawRewardMsg withdrawal />
  | WithdrawRewardMsgFail(withdrawal) => <IndexTokenMsg.WithdrawRewardFailMsg withdrawal />
  | WithdrawCommissionMsgSuccess(withdrawal) => <IndexTokenMsg.WithdrawComissionMsg withdrawal />
  | WithdrawCommissionMsgFail(withdrawal) => <IndexTokenMsg.WithdrawComissionFailMsg withdrawal />
  | MultiSendMsgSuccess(tx)
  | MultiSendMsgFail(tx) => <IndexTokenMsg.MultisendMsg tx />
  | CreateDataSourceMsgSuccess(dataSource) => <IndexDataMsg.CreateDataSourceMsg dataSource />
  | CreateDataSourceMsgFail(dataSource) => <IndexDataMsg.CreateDataSourceFailMsg dataSource />
  | EditDataSourceMsgSuccess(dataSource)
  | EditDataSourceMsgFail(dataSource) => <IndexDataMsg.EditDataSourceMsg dataSource />
  | CreateOracleScriptMsgSuccess(oracleScript) =>
    <IndexDataMsg.CreateOracleScriptMsg oracleScript />
  | CreateOracleScriptMsgFail(oracleScript) =>
    <IndexDataMsg.CreateOracleScriptFailMsg oracleScript />
  | EditOracleScriptMsgSuccess(oracleScript)
  | EditOracleScriptMsgFail(oracleScript) => <IndexDataMsg.EditOracleScriptMsg oracleScript />
  | RequestMsgSuccess(request) => <IndexDataMsg.RequestMsg request />
  | RequestMsgFail(request) => <IndexDataMsg.RequestFailMsg request />
  | ReportMsgSuccess(report)
  | ReportMsgFail(report) => <IndexDataMsg.ReportMsg report />
  | AddReporterMsgSuccess(address) => renderAddReporter(address)
  | AddReporterMsgFail(address) => React.null
  | RemoveReporterMsgSuccess(address) => renderRemoveReporter(address)
  | RemoveReporterMsgFail(address) => React.null
  | CreateValidatorMsgSuccess(validator) => <IndexValidatorMsg.CreateValidatorMsg validator />
  | CreateValidatorMsgFail(validator) => <IndexValidatorMsg.CreateValidatorMsg validator />
  | EditValidatorMsgSuccess(validator) => <IndexValidatorMsg.EditValidatorMsg validator />
  | EditValidatorMsgFail(validator) => React.null
  | UnjailMsgSuccess(unjail) => renderUnjail(unjail)
  | UnjailMsgFail(unjail) => React.null
  | SetWithdrawAddressMsgSuccess(set) => renderSetWithdrawAddress(set)
  | SetWithdrawAddressMsgFail(set) => React.null
  | SubmitProposalMsgSuccess(proposal) => <IndexProposalMsg.SubmitProposalMsg proposal />
  | SubmitProposalMsgFail(proposal) => React.null
  | DepositMsgSuccess(deposit) => <IndexProposalMsg.DepositMsg deposit />
  | DepositMsgFail(deposit) => <IndexProposalMsg.DepositFailMsg deposit />
  | VoteMsgSuccess(vote) => <IndexProposalMsg.VoteMsg vote />
  | VoteMsgFail(vote) => React.null
  | ActivateMsgSuccess(activate) => renderActivate(activate)
  | ActivateMsgFail(activate) => React.null
  | UnknownMsg => renderUnknownMessage()
  //TODO: Re-visit IBC Msg
  | CreateClientMsg(info) => renderCreateClient(info)
  | UpdateClientMsg(info) => renderUpdateClient(info)
  | SubmitClientMisbehaviourMsg(info) => renderSubmitClientMisbehaviour(info)
  | ConnectionOpenInitMsg(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenTryMsg(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenAckMsg(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenConfirmMsg(info) => renderConnectionVariant(msg, info.common)
  | ChannelOpenInitMsg(info) => renderChannelVariant(info.common)
  | ChannelOpenTryMsg(info) => renderChannelVariant(info.common)
  | ChannelOpenAckMsg(info) => renderChannelVariant(info.common)
  | ChannelOpenConfirmMsg(info) => renderChannelVariant(info.common)
  | ChannelCloseInitMsg(info) => renderChannelVariant(info.common)
  | ChannelCloseConfirmMsg(info) => renderChannelVariant(info.common)
  | PacketMsg(info) => renderPacketVariant(msg, info.common)
  | AcknowledgementMsg(info) => renderPacketVariant(msg, info.common)
  | TimeoutMsg(info) => renderPacketVariant(msg, info.common)
  | _ => React.null
  };

[@react.component]
let make = (~messages: list(TxSub.Msg.t)) => {
  <>
    {messages
     ->Belt.List.mapWithIndex((index, msg) => {
         let theme = msg |> TxSub.Msg.getBadgeTheme;
         <div className=CssHelper.infoContainer key={(index |> string_of_int) ++ theme.name}>
           <div
             className={Css.merge([
               CssHelper.flexBox(),
               Styles.infoHeader,
               CssHelper.mb(~size=21, ()),
               CssHelper.mbSm(~size=16, ()),
             ])}>
             <IndexMsgIcon category={theme.category} />
             <HSpacing size=Spacing.sm />
             <Heading value={theme.name} size=Heading.H4 />
           </div>
           {renderBody(msg)}
         </div>;
       })
     ->Array.of_list
     ->React.array}
  </>;
};

module Loading = {
  [@react.component]
  let make = () => {
    <div className=CssHelper.infoContainer>
      <div
        className={Css.merge([
          CssHelper.flexBox(),
          Styles.infoHeader,
          CssHelper.mb(~size=21, ()),
          CssHelper.mbSm(~size=16, ()),
        ])}>
        <LoadingCensorBar width=24 height=24 radius=24 />
        <HSpacing size=Spacing.sm />
        <LoadingCensorBar width=75 height=15 />
      </div>
      <Row.Grid>
        <Col.Grid col=Col.Six mb=24>
          <LoadingCensorBar width=75 height=15 mb=8 />
          <LoadingCensorBar width=150 height=15 />
        </Col.Grid>
        <Col.Grid col=Col.Six mb=24>
          <LoadingCensorBar width=75 height=15 mb=8 />
          <LoadingCensorBar width=150 height=15 />
        </Col.Grid>
        <Col.Grid col=Col.Six>
          <LoadingCensorBar width=75 height=15 mb=8 />
          <LoadingCensorBar width=150 height=15 />
        </Col.Grid>
      </Row.Grid>
    </div>;
  };
};
