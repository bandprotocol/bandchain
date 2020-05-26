[@react.component]
let make = (~validator, ~setMsgsOpt) => {
  let (validator, setValidator) =
    React.useState(_ =>
      EnhanceTxInput.{text: validator |> Address.toOperatorBech32, value: Some(validator)}
    );
  let (amount, setAmount) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect2(
    _ => {
      let msgsOpt = {
        let%Opt validatorValue = validator.value;
        let%Opt amountValue = amount.value;
        Some([|
          TxCreator.Delegate(
            validatorValue,
            {amount: amountValue |> Js.Float.toString, denom: "uband"},
          ),
        |]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    (validator, amount),
  );

  <>
    <EnhanceTxInput
      width=360
      inputData=validator
      setInputData=setValidator
      parse=Address.fromBech32Opt
      msg="Delegate to"
      errMsg="Invalid Address"
      code=true
    />
    <VSpacing size=Spacing.md />
    <EnhanceTxInput
      width=115
      inputData=amount
      setInputData=setAmount
      parse=Parse.getBandAmount
      msg="Amount (BAND)"
      errMsg="Invalid amount"
      code=true
    />
  </>;
};
