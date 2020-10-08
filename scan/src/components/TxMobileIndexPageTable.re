module Styles = {
  open Css;
  let blockWrapper = style([paddingBottom(`px(20))]);
};

let addressWidth = 160;
let renderMuitisendList = (tx: TxSub.Msg.MultiSend.t) =>
  InfoMobileCard.[("Inputs", Nothing)]
  ->Belt.List.concat(
      {
        let%IterList {address, coins} = tx.inputs;
        [
          ("From", InfoMobileCard.Address(address, addressWidth, `account)),
          ("Amount", Coin({value: coins, hasDenom: false})),
        ];
      },
    )
  ->Belt.List.concat([("Output", Nothing)])
  ->Belt.List.concat(
      {
        let%IterList {address, coins} = tx.outputs;
        [
          ("To", InfoMobileCard.Address(address, addressWidth, `account)),
          ("Amount", Coin({value: coins, hasDenom: false})),
        ];
      },
    );

let renderDetailMobile =
  //TODO: implement Guan Yu's message later
  fun
  | TxSub.Msg.SendMsgSuccess({fromAddress, toAddress, amount}) =>
    InfoMobileCard.[
      ("From", Address(fromAddress, addressWidth, `account)),
      ("To", Address(toAddress, addressWidth, `account)),
      ("Amount", Coin({value: amount, hasDenom: true})),
    ]
  | DelegateMsgSuccess({validatorAddress, delegatorAddress, amount}) => [
      ("Delegator Address", Address(delegatorAddress, addressWidth, `account)),
      ("Validator Address", Address(validatorAddress, addressWidth, `validator)),
      ("Amount", Coin({value: [amount], hasDenom: true})),
    ]
  | UndelegateMsgSuccess({validatorAddress, delegatorAddress, amount}) => [
      ("Delegator Address", Address(delegatorAddress, addressWidth, `account)),
      ("Validator Address", Address(validatorAddress, addressWidth, `validator)),
      ("Amount", Coin({value: [amount], hasDenom: true})),
    ]
  | MultiSendMsgSuccess(tx) => renderMuitisendList(tx)
  | WithdrawRewardMsgSuccess({validatorAddress, delegatorAddress, amount}) => [
      ("Delegator Address", Address(delegatorAddress, addressWidth, `account)),
      ("Validator Address", Address(validatorAddress, addressWidth, `validator)),
      ("Amount", Coin({value: amount, hasDenom: true})),
    ]
  | RedelegateMsgSuccess({
      validatorSourceAddress,
      validatorDestinationAddress,
      delegatorAddress,
      amount,
    }) => [
      ("Delegator Address", Address(delegatorAddress, addressWidth, `account)),
      ("Source Address", Address(validatorSourceAddress, addressWidth, `validator)),
      ("Destination Address", Address(validatorDestinationAddress, addressWidth, `validator)),
      ("Amount", Coin({value: [amount], hasDenom: true})),
    ]
  | SetWithdrawAddressMsgSuccess({delegatorAddress, withdrawAddress}) => [
      ("Delegator Address", Address(delegatorAddress, addressWidth, `account)),
      ("Withdraw Address", Address(withdrawAddress, addressWidth, `account)),
    ]
  | CreateValidatorMsgSuccess({
      moniker,
      identity,
      website,
      details,
      commissionRate,
      commissionMaxRate,
      commissionMaxChange,
      delegatorAddress,
      validatorAddress,
      publicKey,
      minSelfDelegation,
      selfDelegation,
    }) => [
      ("Moniker", Text(moniker)),
      ("Identity", Text(identity)),
      ("Website", Text(website)),
      ("Detail", Text(details)),
      ("Commission Rate", Percentage(commissionRate, Some(4))),
      ("Commission Max Rate", Percentage(commissionMaxRate, Some(4))),
      ("Commission Max Change", Percentage(commissionMaxChange, Some(4))),
      ("Delegator Address", Address(delegatorAddress, addressWidth, `account)),
      ("Validator Address", Address(validatorAddress, addressWidth, `validator)),
      ("Public Key", PubKey(publicKey)),
      ("Min Self Delegation", Coin({value: [minSelfDelegation], hasDenom: true})),
      ("Self Delegation", Coin({value: [selfDelegation], hasDenom: true})),
    ]
  | EditValidatorMsgSuccess({
      moniker,
      identity,
      website,
      details,
      commissionRate,
      sender,
      minSelfDelegation,
    }) => [
      ("Moniker", moniker == Config.doNotModify ? Text("Unchanged") : Text(moniker)),
      ("Identity", identity == Config.doNotModify ? Text("Unchanged") : Text(identity)),
      ("Website", website == Config.doNotModify ? Text("Unchanged") : Text(website)),
      ("Detail", details == Config.doNotModify ? Text("Unchanged") : Text(details)),
      (
        "Commission Rate",
        switch (commissionRate) {
        | Some(rate) => Percentage(rate, Some(4))
        | None => Text("Unchanged")
        },
      ),
      ("Validator Address", Address(sender, addressWidth, `validator)),
      (
        "Min Self Delegation",
        switch (minSelfDelegation) {
        | Some(amount) => Coin({value: [amount], hasDenom: true})
        | None => Text("Unchanged")
        },
      ),
    ]
  | WithdrawCommissionMsgSuccess({validatorAddress, amount}) => [
      ("Validator Address", Address(validatorAddress, addressWidth, `validator)),
      ("Amount", Coin({value: amount, hasDenom: true})),
    ]
  | UnjailMsgSuccess({address}) => [
      ("Validator Address", Address(address, addressWidth, `validator)),
    ]
  | CreateDataSourceMsgSuccess({id, owner, name})
  | EditDataSourceMsgSuccess({id, owner, name}) => [
      ("Owner", Address(owner, addressWidth, `account)),
      ("Name", DataSource(id, name)),
    ]
  | CreateOracleScriptMsgSuccess({id, owner, name})
  | EditOracleScriptMsgSuccess({id, owner, name}) => [
      ("Owner", Address(owner, addressWidth, `account)),
      ("Name", OracleScript(id, name)),
    ]
  | RequestMsgSuccess({oracleScriptID, oracleScriptName, calldata, askCount, schema, minCount}) => {
      let calldataKVsOpt = Obi.decode(schema, "input", calldata);
      [
        ("Oracle Script", OracleScript(oracleScriptID, oracleScriptName)),
        ("Calldata", CopyButton(calldata)),
        ("", KVTableRequest(calldataKVsOpt)),
        ("Ask Count", Count(askCount)),
        ("Min Count", Count(minCount)),
      ];
    }
  | ReportMsgSuccess({requestID, rawReports}) => [
      ("Request ID", RequestID(requestID)),
      ("Raw Data Reports", KVTableReport(["EXTERNAL ID", "EXIT CODE", "VALUE"], rawReports)),
    ]
  | AddReporterMsgSuccess({reporter, validatorMoniker})
  | RemoveReporterMsgSuccess({reporter, validatorMoniker}) => [
      ("Validator", Text(validatorMoniker)),
      ("Reporter Address", Address(reporter, addressWidth, `account)),
    ]
  | ActivateMsgSuccess({validatorAddress}) => [
      ("Validator Address", Address(validatorAddress, addressWidth, `validator)),
    ]
  | SubmitProposalMsgSuccess({proposer, title, description, initialDeposit}) => [
      ("Title", Text(title)),
      //TODO: will re-visit
      // ("Description", Text(description)),
      ("Proposer", Address(proposer, addressWidth, `account)),
      ("Amount", Coin({value: initialDeposit, hasDenom: true})),
    ]
  | DepositMsgSuccess({depositor, proposalID, amount}) => [
      ("Dopositor", Address(depositor, addressWidth, `account)),
      ("Proposal ID", Count(proposalID)),
      ("Amount", Coin({value: amount, hasDenom: true})),
    ]
  | VoteMsgSuccess({voterAddress, proposalID, option}) => [
      ("Voter Address", Address(voterAddress, addressWidth, `account)),
      ("Proposal ID", Count(proposalID)),
      ("Option", Text(option)),
    ]
  | _ => [];

[@react.component]
let make = (~messages: list(TxSub.Msg.t)) => {
  <div className=Styles.blockWrapper>
    {messages
     ->Belt.List.mapWithIndex((index, msg) => {
         let renderList = msg |> renderDetailMobile;
         let theme = msg |> TxSub.Msg.getBadgeTheme;
         let creator = msg |> TxSub.Msg.getCreator;
         let key_ = (index |> string_of_int) ++ (creator |> Address.toBech32);
         <MobileCard
           values={
             InfoMobileCard.[
               ("Message\nType", Badge(theme)),
               ("Creator", Address(creator, addressWidth, `account)),
             ]
             ->Belt.List.concat(renderList)
           }
           key=key_
           idx=key_
         />;
       })
     ->Array.of_list
     ->React.array}
  </div>;
};

module Loading = {
  [@react.component]
  let make = () => {
    <MobileCard
      values=InfoMobileCard.[
        ("Message\nType", Loading(80)),
        ("Creator", Loading(80)),
        ("Detail", Loading(80)),
      ]
      idx="1"
    />;
  };
};
