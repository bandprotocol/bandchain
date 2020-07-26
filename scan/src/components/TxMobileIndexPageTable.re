let addressWidth = 160;
let renderMuitisendList = (tx: TxSub.Msg.MultiSend.t) =>
  InfoMobileCard.[("INPUTS", Nothing)]
  ->Belt.List.concat(
      {
        let%IterList {address, coins} = tx.inputs;
        [
          ("FROM", InfoMobileCard.Address(address, addressWidth, false)),
          ("AMOUNT", Coin({value: coins, hasDenom: false})),
        ];
      },
    )
  ->Belt.List.concat([("OUTPUT", Nothing)])
  ->Belt.List.concat(
      {
        let%IterList {address, coins} = tx.outputs;
        [
          ("TO", InfoMobileCard.Address(address, addressWidth, false)),
          ("AMOUNT", Coin({value: coins, hasDenom: false})),
        ];
      },
    );

let renderDetailMobile =
  //TODO: implement Guan Yu's message later
  fun
  | TxSub.Msg.SendMsg({fromAddress, toAddress, amount}) =>
    InfoMobileCard.[
      ("FROM", Address(fromAddress, addressWidth, false)),
      ("TO", Address(toAddress, addressWidth, false)),
      ("AMOUNT", Coin({value: amount, hasDenom: true})),
    ]
  | DelegateMsg({validatorAddress, delegatorAddress, amount}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, false)),
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, true)),
      ("AMOUNT", Coin({value: [amount], hasDenom: true})),
    ]
  | UndelegateMsg({validatorAddress, delegatorAddress, amount}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, false)),
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, true)),
      ("AMOUNT", Coin({value: [amount], hasDenom: true})),
    ]
  | MultiSendMsg(tx) => renderMuitisendList(tx)
  | WithdrawRewardMsg({validatorAddress, delegatorAddress, amount}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, false)),
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, true)),
      ("AMOUNT", Coin({value: amount, hasDenom: true})),
    ]
  | RedelegateMsg({validatorSourceAddress, validatorDestinationAddress, delegatorAddress, amount}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, false)),
      ("SOURCE ADDRESS", Address(validatorSourceAddress, addressWidth, true)),
      ("DESTINATION ADDRESS", Address(validatorDestinationAddress, addressWidth, true)),
      ("AMOUNT", Coin({value: [amount], hasDenom: true})),
    ]
  | SetWithdrawAddressMsg({delegatorAddress, withdrawAddress}) => [
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, false)),
      ("WITHDRAW ADDRESS", Address(withdrawAddress, addressWidth, false)),
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
      ("DELEGATOR ADDRESS", Address(delegatorAddress, addressWidth, false)),
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, true)),
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
      ("VALIDATOR ADDRESS", Address(sender, addressWidth, true)),
      (
        "MIN SELF DELEGATION",
        switch (minSelfDelegation) {
        | Some(amount) => Coin({value: [amount], hasDenom: true})
        | None => Text("Unchanged")
        },
      ),
    ]
  | WithdrawCommissionMsg({validatorAddress, amount}) => [
      ("VALIDATOR ADDRESS", Address(validatorAddress, addressWidth, true)),
      ("AMOUNT", Coin({value: amount, hasDenom: true})),
    ]
  | UnjailMsg({address}) => [("VALIDATOR ADDRESS", Address(address, addressWidth, true))]
  | _ => [];

[@react.component]
let make = (~messages: list(TxSub.Msg.t)) => {
  <>
    //TODO: Change index to be uniqe something
    {messages
     ->Belt.List.mapWithIndex((index, msg) => {
         let renderList = msg |> renderDetailMobile;
         let theme = msg |> TxSub.Msg.getBadgeTheme;
         let creator = msg |> TxSub.Msg.getCreator;
         <MobileCard
           values={
             InfoMobileCard.[
               ("MESSAGE\nTYPE", Badge(theme)),
               ("CREATOR", Address(creator, addressWidth, false)),
             ]
             ->Belt.List.concat(renderList)
           }
           key={index |> string_of_int}
           idx={index |> string_of_int}
         />;
       })
     ->Array.of_list
     ->React.array}
  </>;
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
