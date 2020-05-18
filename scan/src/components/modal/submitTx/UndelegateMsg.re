module Styles = {
  open Css;

  let warning =
    style([
      padding(`px(10)),
      color(Colors.yellow6),
      backgroundColor(Colors.yellow1),
      border(`px(1), `solid, Colors.yellow6),
      borderRadius(`px(4)),
    ]);
};

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
          TxCreator.Undelegate(
            validatorValue,
            {amount: amountValue *. 1e6 |> Js.Float.toString, denom: "uband"},
          ),
        |]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    (validator, amount),
  );

  <>
    <div className=Styles.warning>
      <Text value="Note: Undelegated balance are locked for 21 days" />
    </div>
    <EnhanceTxInput
      width=360
      inputData=validator
      setInputData=setValidator
      parse=Address.fromBech32Opt
      msg="Undelegate from"
      errMsg="Invalid Address"
    />
    <VSpacing size=Spacing.md />
    <EnhanceTxInput
      width=115
      inputData=amount
      setInputData=setAmount
      parse=float_of_string_opt
      msg="Amount (BAND)"
      errMsg="Invalid amount"
    />
  </>;
};
