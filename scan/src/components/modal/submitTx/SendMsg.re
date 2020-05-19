[@react.component]
let make = (~setMsgsOpt) => {
  let (toAddress, setToAddress) = React.useState(_ => EnhanceTxInput.empty);
  let (amount, setAmount) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect2(
    _ => {
      let msgsOpt = {
        let%Opt toAddressValue = toAddress.value;
        let%Opt amountValue = amount.value;
        Some([|
          TxCreator.Send(
            toAddressValue,
            {amount: amountValue |> Js.Float.toString, denom: "uband"},
          ),
        |]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    (toAddress, amount),
  );

  <>
    <EnhanceTxInput
      width=360
      inputData=toAddress
      setInputData=setToAddress
      parse=Address.fromBech32Opt
      msg="To"
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
