module Styles = {
  open Css;

  let info = style([display(`flex), justifyContent(`spaceBetween)]);

  let validator =
    style([display(`flex), flexDirection(`column), alignItems(`flexEnd), width(`px(330))]);

  let warning =
    style([
      padding(`px(10)),
      color(Colors.blue5),
      backgroundColor(Colors.blue1),
      border(`px(1), `solid, Colors.blue6),
      borderRadius(`px(4)),
    ]);
};

[@react.component]
let make = (~address, ~validator, ~setMsgsOpt) => {
  let accountSub = AccountSub.get(address);
  let validatorInfoSub = ValidatorSub.get(validator);

  let allSub = Sub.all2(accountSub, validatorInfoSub);

  let (amount, setAmount) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect1(
    _ => {
      let msgsOpt = {
        let%Opt amountValue = amount.value;
        Some([|
          TxCreator.Delegate(
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
    <VSpacing size=Spacing.lg />
    <div className=Styles.info>
      <Text value="Account Balance" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      {switch (allSub) {
       | Data(({balance}, _)) =>
         <div>
           <Text
             value={balance |> Coin.getBandAmountFromCoins |> Format.fPretty(~digits=6)}
             code=true
             size=Text.Lg
             weight=Text.Semibold
           />
           <Text value=" BAND" code=true />
         </div>
       | _ => <LoadingCensorBar width=150 height=18 />
       }}
    </div>
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.md />
    <div className=Styles.info>
      <Text value="Delegate To" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      {switch (allSub) {
       | Data((_, {moniker})) =>
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
    {switch (allSub) {
     | Data(({balance}, _)) =>
       //  TODO: hard-coded tx fee
       let maxValInUband = (balance |> Coin.getUBandAmountFromCoins) -. 5000.;
       <EnhanceTxInput
         width=300
         inputData=amount
         setInputData=setAmount
         parse={Parse.getBandAmount(maxValInUband)}
         maxValue={maxValInUband /. 1e6 |> Js.Float.toString}
         msg="Delegate Amount (BAND)"
         placeholder="Insert delegation amount"
         inputType="number"
         code=true
         autoFocus=true
         id="delegateAmountInput"
       />;
     | _ => <EnhanceTxInput.Loading msg="Delegate Amount (BAND)" width=300 />
     }}
    <VSpacing size=Spacing.lg />
  </>;
};
