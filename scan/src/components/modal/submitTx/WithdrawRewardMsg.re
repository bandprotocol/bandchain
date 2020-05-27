module Styles = {
  open Css;

  let validator =
    style([display(`flex), flexDirection(`column), alignItems(`flexEnd), width(`px(330))]);

  let info = style([display(`flex), justifyContent(`spaceBetween)]);
};

[@react.component]
let make = (~address, ~validator, ~setMsgsOpt) => {
  let validatorInfoSub = ValidatorSub.get(validator);
  let delegationSub = DelegationSub.getStakeByValiator(address, validator);

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
    <VSpacing size=Spacing.lg />
    <div className=Styles.info>
      <Text
        value="Withdraw Reward From"
        size=Text.Lg
        spacing={Text.Em(0.03)}
        nowrap=true
        block=true
      />
      {switch (allSub) {
       | Data(({moniker}, _)) =>
         <div className=Styles.validator>
           <Text value=moniker code=true ellipsis=true align=Text.Right />
           <Text
             value={"(" ++ validator->Address.toOperatorBech32 ++ ")"}
             size=Text.Sm
             color=Colors.blueGray5
             code=true
           />
         </div>
       | _ => <LoadingCensorBar width=300 height=26 />
       }}
    </div>
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.md />
    <div className=Styles.info>
      <Text value="Current Reward" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      {switch (allSub) {
       | Data((_, {reward})) =>
         <div>
           <Text
             value={reward |> Coin.getBandAmountFromCoin |> Format.fPretty(~digits=6)}
             code=true
           />
           <Text value=" BAND" code=true />
         </div>
       | _ => <LoadingCensorBar width=300 height=18 />
       }}
    </div>
    <VSpacing size=Spacing.lg />
  </>;
};
