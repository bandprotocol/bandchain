[@react.component]
let make = (~setMsgsOpt) => {
  let (validator, setValidator) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect2(
    _ => {
      let msgsOpt = {
        let%Opt validatorValue = validator.value;
        Some([|TxCreator.WithdrawReward(validatorValue)|]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    (validator, 0),
  );

  <>
    <EnhanceTxInput
      width=360
      inputData=validator
      setInputData=setValidator
      parse=Address.fromBech32Opt
      msg="Undelegate from"
      errMsg="Invalid Address"
    />
  </>;
};
