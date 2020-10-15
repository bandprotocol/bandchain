module Styles = {
  open Css;

  let container = style([paddingBottom(`px(24))]);

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
    <div className=Styles.container>
      <Text value="Delegate To" size=Text.Md weight=Text.Medium nowrap=true block=true />
      <VSpacing size=Spacing.sm />
      {switch (allSub) {
       | Data((_, {moniker})) =>
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
      <Text value="Account Balance" size=Text.Md weight=Text.Medium nowrap=true block=true />
      <VSpacing size=Spacing.sm />
      {switch (allSub) {
       | Data(({balance}, _)) =>
         <div>
           <Text
             value={balance |> Coin.getBandAmountFromCoins |> Format.fPretty(~digits=6)}
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
     | Data(({balance}, _)) =>
       //  TODO: hard-coded tx fee
       let maxValInUband = (balance |> Coin.getUBandAmountFromCoins) -. 5000.;
       <EnhanceTxInput
         width=300
         inputData=amount
         setInputData=setAmount
         parse={Parse.getBandAmount(maxValInUband)}
         maxValue={maxValInUband /. 1e6 |> Js.Float.toString}
         msg="Delegate Amount"
         placeholder="0.000000"
         inputType="number"
         code=true
         autoFocus=true
         id="delegateAmountInput"
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
