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
      parse=Address.fromBech32Opt
      msg="Recipient Address"
      errMsg="Invalid Address"
      code=true
      placeholder="Insert recipient address"
    />
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.md />
    <EnhanceTxInput
      width=236
      inputData=amount
      setInputData=setAmount
      parse=Parse.getBandAmount
      msg="Send Amount (BAND)"
      errMsg="Invalid amount"
      code=true
      placeholder="Insert send amount"
    />
    <VSpacing size=Spacing.lg />
  </>;
};
