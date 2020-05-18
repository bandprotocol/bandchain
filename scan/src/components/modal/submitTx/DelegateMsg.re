[@react.component]
let make = (~setMsgsOpt) => {
  let (validator, setValidator) = React.useState(_ => EnhanceTxInput.empty);
  let (amount, setAmount) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect2(
    _ => {
      let msgsOpt = {
        let%Opt validatorValue = validator.value;
        let%Opt amountValue = amount.value;
        Some([|
          TxCreator.Delegate(
            validatorValue,
            {
              amount: amountValue *. 1e6 |> Js.Math.floor_float |> Js.Float.toString,
              denom: "uband",
            },
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
      parse=float_of_string_opt
      msg="Amount (BAND)"
      errMsg="Invalid amount"
      code=true
    />
  </>;
};
