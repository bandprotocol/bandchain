[@react.component]
let make = (~validator, ~setMsgsOpt) => {
  let (validator, setValidator) =
    React.useState(_ =>
      EnhanceTxInput.{text: validator |> Address.toOperatorBech32, value: Some(validator)}
    );

  React.useEffect1(
    _ => {
      let msgsOpt = {
        let%Opt validatorValue = validator.value;
        Some([|TxCreator.WithdrawReward(validatorValue)|]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    [|validator|],
  );

  <>
    <EnhanceTxInput
      width=360
      inputData=validator
      setInputData=setValidator
      parse=Address.fromBech32Opt
      msg="Withdraw from"
      errMsg="Invalid Address"
    />
  </>;
};
