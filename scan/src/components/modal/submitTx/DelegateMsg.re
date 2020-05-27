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
let make = (~address, ~validator, ~setMsgsOpt) =>
  {
    let accountSub = AccountSub.get(address);
    let validatorInfoSub = ValidatorSub.get(validator);
    let delegationSub = DelegationSub.getStakeByValiator(address, validator);

    let allSub = Sub.all3(accountSub, validatorInfoSub, delegationSub);

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
        <Text
          value="Account Balance"
          size=Text.Lg
          spacing={Text.Em(0.03)}
          nowrap=true
          block=true
        />
        {switch (allSub) {
         | Data(({balance}, _, _)) =>
           <div>
             <Text
               value={balance |> Coin.getBandAmountFromCoins |> Format.fPretty(~digits=6)}
               code=true
             />
             <Text value=" BAND" code=true />
           </div>
         | _ => <LoadingCensorBar width=300 height=18 />
         }}
      </div>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.md />
      <div className=Styles.info>
        <Text value="Delegate To" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
        {switch (allSub) {
         | Data((_, {moniker}, _)) =>
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
      <EnhanceTxInput
        width=226
        inputData=amount
        setInputData=setAmount
        parse=Parse.getBandAmount
        msg="Delegate Amount (BAND)"
        errMsg="Invalid amount"
        placeholder="Insert delegation amount"
        code=true
      />
      <VSpacing size=Spacing.lg />
    </>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
