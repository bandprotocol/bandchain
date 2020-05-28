module Styles = {
  open Css;

  let info = style([display(`flex), justifyContent(`spaceBetween)]);
};
[@react.component]
let make = (~address, ~receiver, ~setMsgsOpt) => {
  let accountSub = AccountSub.get(address);
  let (toAddress, setToAddress) =
    React.useState(_ => {
      switch (receiver) {
      | Some(receiver') =>
        EnhanceTxInput.{text: receiver' |> Address.toBech32, value: Some(receiver')}
      | None => EnhanceTxInput.empty
      }
    });
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
    <VSpacing size=Spacing.lg />
    <div className=Styles.info>
      <Text
        value="Available Balance"
        size=Text.Lg
        spacing={Text.Em(0.03)}
        nowrap=true
        block=true
      />
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.md />
      {switch (accountSub) {
       | Data({balance}) =>
         <div>
           <Text
             value={balance |> Coin.getBandAmountFromCoins |> Format.fPretty(~digits=6)}
             code=true
             size=Text.Lg
             weight=Text.Semibold
           />
           <Text value=" BAND" code=true />
         </div>
       | _ => <LoadingCensorBar width=300 height=18 />
       }}
    </div>
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.md />
    <EnhanceTxInput
      width=302
      inputData=toAddress
      setInputData=setToAddress
      parse=Parse.address
      msg="Recipient Address"
      code=true
      placeholder="Insert recipient address"
    />
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.md />
    {switch (accountSub) {
     | Data({balance}) =>
       <EnhanceTxInput
         width=236
         inputData=amount
         setInputData=setAmount
         //  TODO: hard-coded tx fee
         parse={Parse.getBandAmount((balance |> Coin.getUBandAmountFromCoins) -. 5000.)}
         msg="Send Amount (BAND)"
         inputType="number"
         code=true
         placeholder="Insert send amount"
       />
     | _ => <LoadingCensorBar width=300 height=18 isRight=true />
     }}
    <VSpacing size=Spacing.lg />
  </>;
};
