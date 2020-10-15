module Styles = {
  open Css;

  let container = style([paddingBottom(`px(24))]);

  let warning =
    style([
      display(`flex),
      flexDirection(`column),
      padding2(~v=`px(16), ~h=`px(24)),
      backgroundColor(Colors.profileBG),
      borderRadius(`px(4)),
      marginBottom(`px(24)),
    ]);
};

[@react.component]
let make = (~address, ~validator, ~setMsgsOpt) => {
  let validatorInfoSub = ValidatorSub.get(validator);
  let delegationSub = DelegationSub.getStakeByValidator(address, validator);

  let allSub = Sub.all2(validatorInfoSub, delegationSub);

  let (amount, setAmount) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect1(
    _ => {
      let msgsOpt = {
        let%Opt amountValue = amount.value;
        Some([|
          TxCreator.Undelegate(
            validator,
            {amount: amountValue |> Js.Float.toString, denom: "uband"},
          ),
        |]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    [|amount|],
  );

  <>
    <div className=Styles.warning>
      <Text weight=Text.Semibold value="Please read before proceeding:" />
      <VSpacing size=Spacing.xs />
      <Text
        weight=Text.Thin
        value="1. Undelegated balance are locked for 21 days. After the unbonding period, the balance will automatically be added to your account"
      />
      <VSpacing size=Spacing.xs />
      <Text
        weight=Text.Thin
        value="2. You can have a maximum of 7 pending unbonding transactions at any one time."
      />
    </div>
    <div className=Styles.container>
      <Text value="Undelegate From" size=Text.Md weight=Text.Medium nowrap=true block=true />
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
      <Text value="Current Stake" size=Text.Md weight=Text.Medium nowrap=true block=true />
      <VSpacing size=Spacing.sm />
      {switch (allSub) {
       | Data((_, {amount: stakedAmount})) =>
         <div>
           <Text
             value={stakedAmount |> Coin.getBandAmountFromCoin |> Format.fPretty(~digits=6)}
             code=true
             size=Text.Lg
             weight=Text.Thin
           />
           <Text value=" BAND" size=Text.Lg weight=Text.Thin code=true />
         </div>
       | _ => <LoadingCensorBar width=150 height=18 />
       }}
    </div>
    {switch (allSub) {
     | Data((_, {amount: stakedAmount})) =>
       let maxValInUband = stakedAmount |> Coin.getUBandAmountFromCoin;
       <EnhanceTxInput
         width=300
         inputData=amount
         setInputData=setAmount
         parse={Parse.getBandAmount(maxValInUband)}
         maxValue={maxValInUband /. 1e6 |> Js.Float.toString}
         msg="Undelegate Amount (BAND)"
         placeholder="Insert unbonding amount"
         inputType="number"
         code=true
         autoFocus=true
         id="undelegateAmountInput"
       />;
     | _ =>
       <EnhanceTxInput.Loading2
         msg="Delegate Amount"
         code=true
         useMax=true
         placeholder="0.000000"
       />
     }}
  </>;
};
