module Styles = {
  open Css;

  let container = style([paddingBottom(`px(24))]);
};

[@react.component]
let make = (~validator, ~amount, ~setMsgsOpt) => {
  let validatorInfoSub = ValidatorSub.get(validator);

  React.useEffect0(_ => {
    let msgsOpt = {
      let%Opt amountValue = Some(amount);
      Some([|
        TxCreator.WithdrawReward(validator),
        TxCreator.Delegate(
          validator,
          {amount: amountValue |> Js.Float.toFixedWithPrecision(~digits=0), denom: "uband"},
        ),
      |]);
    };
    setMsgsOpt(_ => msgsOpt);
    None;
  });

  <>
    <div className=Styles.container>
      <Heading
        value="Reinvest to"
        size=Heading.H5
        marginBottom=8
        align=Heading.Left
        weight=Heading.Medium
      />
      {switch (validatorInfoSub) {
       | Data({moniker}) =>
         <div>
           <Text value=moniker size=Text.Lg ellipsis=true align=Text.Right />
           <Text
             value={"(" ++ validator->Address.toOperatorBech32 ++ ")"}
             size=Text.Md
             color=Colors.gray6
             code=true
             block=true
           />
         </div>
       | _ => <LoadingCensorBar width=300 height=34 />
       }}
    </div>
    <div className=Styles.container>
      <Heading
        value="Current Reward"
        size=Heading.H5
        marginBottom=8
        align=Heading.Left
        weight=Heading.Medium
      />
      <div>
        <Text
          value={
            amount
            |> Coin.newUBANDFromAmount
            |> Coin.getBandAmountFromCoin
            |> Format.fPretty(~digits=6)
          }
          code=true
          size=Text.Lg
        />
        <Text value=" BAND" size=Text.Lg code=true />
      </div>
    </div>
  </>;
};
