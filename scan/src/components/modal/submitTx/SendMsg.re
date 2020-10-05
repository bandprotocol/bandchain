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
       | _ => <LoadingCensorBar width=150 height=18 />
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
      id="recipientAddressInput"
      placeholder="Insert recipient address"
      autoFocus={
        switch (toAddress.text) {
        | "" => true
        | _ => false
        }
      }
    />
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.md />
    {switch (accountSub) {
     | Data({balance}) =>
       //  TODO: hard-coded tx fee
       let maxValInUband = (balance |> Coin.getUBandAmountFromCoins) -. 5000.;
       <EnhanceTxInput
         width=300
         inputData=amount
         setInputData=setAmount
         parse={Parse.getBandAmount(maxValInUband)}
         maxValue={maxValInUband /. 1e6 |> Js.Float.toString}
         msg="Send Amount (BAND)"
         inputType="number"
         code=true
         placeholder="Insert send amount"
         autoFocus={
           switch (toAddress.text) {
           | "" => false
           | _ => true
           }
         }
         id="sendAmountInput"
       />;
     | _ => <EnhanceTxInput.Loading msg="Send Amount (BAND)" width=300 />
     }}
    <VSpacing size=Spacing.lg />
  </>;
};
