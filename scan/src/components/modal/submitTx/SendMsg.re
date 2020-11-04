module Styles = {
  open Css;
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
    <Heading size=Heading.H5 value="Available Balance" marginBottom=8 />
    <div className={CssHelper.mb(~size=24, ())}>
      {switch (accountSub) {
       | Data({balance}) =>
         <div>
           <Text
             value={balance |> Coin.getBandAmountFromCoins |> Format.fPretty(~digits=6)}
             code=true
             size=Text.Lg
           />
           <Text value=" BAND" code=true />
         </div>
       | _ => <LoadingCensorBar width=150 height=18 />
       }}
    </div>
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
         placeholder="0.000000"
         autoFocus={
           switch (toAddress.text) {
           | "" => false
           | _ => true
           }
         }
         id="sendAmountInput"
         maxWarningMsg=true
       />;
     | _ =>
       <EnhanceTxInput.Loading
         msg="Send Amount (BAND)"
         code=true
         useMax=true
         placeholder="0.000000"
       />
     }}
  </>;
};
