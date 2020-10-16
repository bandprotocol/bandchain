module Styles = {
  open Css;

  let container = style([paddingBottom(`px(24))]);

  let validator =
    style([display(`flex), flexDirection(`column), alignItems(`flexEnd), width(`px(330))]);
};

[@react.component]
let make = (~address, ~validator, ~setMsgsOpt) => {
  let validatorInfoSub = ValidatorSub.get(validator);
  let delegationSub = DelegationSub.getStakeByValidator(address, validator);

  let allSub = Sub.all2(validatorInfoSub, delegationSub);

  React.useEffect1(
    _ => {
      let msgsOpt = {
        Some([|TxCreator.WithdrawReward(validator)|]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    [|validator|],
  );

  <>
    <div className=Styles.container>
      <Text
        value="Withdraw Delegation Rewards"
        size=Text.Md
        weight=Text.Medium
        nowrap=true
        block=true
      />
      <VSpacing size=Spacing.sm />
      {switch (allSub) {
       | Data(({moniker}, _)) =>
         <div>
           <Text value=moniker size=Text.Lg weight=Text.Thin ellipsis=true align=Text.Right />
           <Text
             value={"(" ++ validator->Address.toOperatorBech32 ++ ")"}
             size=Text.Md
             weight=Text.Thin
             color=Colors.gray6
             code=true
             block=true
           />
         </div>
       | _ => <LoadingCensorBar width=300 height=34 />
       }}
    </div>
    <div className=Styles.container>
      <Text value="Current Reward" size=Text.Md weight=Text.Medium nowrap=true block=true />
      <VSpacing size=Spacing.sm />
      {switch (allSub) {
       | Data((_, {reward})) =>
         <div>
           <NumberCountup
             value={reward |> Coin.getBandAmountFromCoin}
             size=Text.Lg
             weight=Text.Thin
             spacing={Text.Em(0.0)}
           />
           <Text value=" BAND" size=Text.Lg weight=Text.Thin code=true />
         </div>
       | _ => <LoadingCensorBar width=150 height=18 />
       }}
    </div>
  </>;
};
