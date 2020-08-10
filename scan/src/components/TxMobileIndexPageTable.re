module Styles = {
  open Css;
  let blockWrapper = style([paddingBottom(`px(20))]);
};

let addressWidth = 160;
let renderMuitisendList = (tx: TxSub.Msg.MultiSend.t) =>
  InfoMobileCard.[("INPUTS", Nothing)]
  ->Belt.List.concat(
      {
        let%IterList {address, coins} = tx.inputs;
        [
          ("FROM", InfoMobileCard.Address(address, addressWidth, `account)),
          ("AMOUNT", Coin({value: coins, hasDenom: false})),
        ];
      },
    )
  ->Belt.List.concat([("OUTPUT", Nothing)])
  ->Belt.List.concat(
      {
        let%IterList {address, coins} = tx.outputs;
        [
          ("TO", InfoMobileCard.Address(address, addressWidth, `account)),
          ("AMOUNT", Coin({value: coins, hasDenom: false})),
        ];
      },
    );

let renderDetailMobile =
  //TODO: implement Guan Yu's message later
  fun
  | TxSub.Msg.SendMsg({fromAddress, toAddress, amount}) =>
    InfoMobileCard.[
      ("FROM", Address(fromAddress, addressWidth, `account)),
      ("TO", Address(toAddress, addressWidth, `account)),
      ("AMOUNT", Coin({value: amount, hasDenom: true})),
    ]
  | DelegateMsg({validatorAddress, delegatorAddress, amount}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, `account)),
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, `validator)),
      ("AMOUNT", Coin({value: [amount], hasDenom: true})),
    ]
  | UndelegateMsg({validatorAddress, delegatorAddress, amount}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, `account)),
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, `validator)),
      ("AMOUNT", Coin({value: [amount], hasDenom: true})),
    ]
  | MultiSendMsg(tx) => renderMuitisendList(tx)
  | WithdrawRewardMsg({validatorAddress, delegatorAddress, amount}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, `account)),
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, `validator)),
      ("AMOUNT", Coin({value: amount, hasDenom: true})),
    ]
  | RedelegateMsg({validatorSourceAddress, validatorDestinationAddress, delegatorAddress, amount}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, `account)),
      ("SOURCE ADDRESS", Address(validatorSourceAddress, addressWidth, `validator)),
      ("DESTINATION ADDRESS", Address(validatorDestinationAddress, addressWidth, `validator)),
      ("AMOUNT", Coin({value: [amount], hasDenom: true})),
    ]
  | SetWithdrawAddressMsg({delegatorAddress, withdrawAddress}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, `account)),
      ("WITHDRAW ADDRESS", Address(withdrawAddress, addressWidth, `account)),
    ]
  | CreateValidatorMsg({
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
      ("MONIKER", Text(moniker)),
      ("IDENTITY", Text(identity)),
      ("WEBSITE", Text(website)),
      ("DETAIL", Text(details)),
      ("COMMISSION RATE", Percentage(commissionRate, Some(4))),
      ("COMMISSION MAX RATE", Percentage(commissionMaxRate, Some(4))),
      ("COMMISSION MAX CHANGE", Percentage(commissionMaxChange, Some(4))),
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, `account)),
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, `validator)),
      ("PUBLIC KEY", PubKey(publicKey)),
      ("MIN SELF DELEGATION", Coin({value: [minSelfDelegation], hasDenom: true})),
      ("SELF DELEGATION", Coin({value: [selfDelegation], hasDenom: true})),
    ]
  | EditValidatorMsg({
      moniker,
      identity,
      website,
      details,
      commissionRate,
      sender,
      minSelfDelegation,
    }) => [
      ("MONIKER", moniker == Config.doNotModify ? Text("Unchanged") : Text(moniker)),
      ("IDENTITY", identity == Config.doNotModify ? Text("Unchanged") : Text(identity)),
      ("WEBSITE", website == Config.doNotModify ? Text("Unchanged") : Text(website)),
      ("DETAIL", details == Config.doNotModify ? Text("Unchanged") : Text(details)),
      (
        "COMMISSION RATE",
        switch (commissionRate) {
        | Some(rate) => Percentage(rate, Some(4))
        | None => Text("Unchanged")
        },
      ),
      ("VALIDATOR ADDRESS", Address(sender, addressWidth, `validator)),
      (
        "MIN SELF DELEGATION",
        switch (minSelfDelegation) {
        | Some(amount) => Coin({value: [amount], hasDenom: true})
        | None => Text("Unchanged")
        },
      ),
    ]
  | WithdrawCommissionMsg({validatorAddress, amount}) => [
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, `validator)),
      ("AMOUNT", Coin({value: amount, hasDenom: true})),
    ]
  | UnjailMsg({address}) => [("VALIDATOR ADDRESS", Address(address, addressWidth, `validator))]
  | CreateDataSourceMsg({id, owner, name})
  | EditDataSourceMsg({id, owner, name}) => [
      ("OWNER", Address(owner, addressWidth, `account)),
      ("NAME", DataSource(id, name)),
    ]
  | CreateOracleScriptMsg({id, owner, name})
  | EditOracleScriptMsg({id, owner, name}) => [
      ("OWNER", Address(owner, addressWidth, `account)),
      ("NAME", OracleScript(id, name)),
    ]
  | RequestMsg({oracleScriptID, oracleScriptName, calldata, askCount, schema, minCount}) => {
      let calldataKVsOpt = Obi.decode(schema, "input", calldata);
      [
        ("ORACLE SCRIPT", OracleScript(oracleScriptID, oracleScriptName)),
        ("CALLDATA", CopyButton(calldata)),
        ("", KVTableRequest(calldataKVsOpt)),
        ("ASK COUNT", Count(askCount)),
        ("MIN COUNT", Count(minCount)),
      ];
    }
  | ReportMsg({requestID, rawReports}) => [
      ("REQUEST ID", RequestID(requestID)),
      ("RAW DATA REPORTS", KVTableReport(["EXTERNAL ID", "EXIT CODE", "VALUE"], rawReports)),
    ]
  | AddReporterMsg({reporter, validatorMoniker})
  | RemoveReporterMsg({reporter, validatorMoniker}) => [
      ("VALIDATOR", Text(validatorMoniker)),
      ("REPORTER ADDRESS", Address(reporter, addressWidth, `account)),
    ]
  | ActivateMsg({validatorAddress}) => [
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, `validator)),
    ]
  | SubmitProposalMsg({proposer, title, description, initialDeposit}) => [
      ("TITLE", Text(title)),
      ("DESCRIPTION", Text(description)),
      ("PROPOSER", Address(proposer, addressWidth, `account)),
      ("AMOUNT", Coin({value: initialDeposit, hasDenom: true})),
    ]
  | DepositMsg({depositor, proposalID, amount}) => [
      ("DOPOSITER", Address(depositor, addressWidth, `account)),
      ("PROPOSAL ID", Count(proposalID)),
      ("AMOUNT", Coin({value: amount, hasDenom: true})),
    ]
  | VoteMsg({voterAddress, proposalID, option}) => [
      ("VOTER ADDRESS", Address(voterAddress, addressWidth, `account)),
      ("PROPOSAL ID", Count(proposalID)),
      ("OPTION", Text(option)),
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
               ("MESSAGE\nTYPE", Badge(theme)),
               ("CREATOR", Address(creator, addressWidth, `account)),
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
        ("MESSAGE\nTYPE", Loading(80)),
        ("CREATOR", Loading(80)),
        ("Detail", Loading(80)),
      ]
      idx="1"
    />;
  };
};
