module Styles = {
  open Css;

  let info = style([display(`flex), justifyContent(`spaceBetween), alignItems(`center)]);

  let validator =
    style([display(`flex), flexDirection(`column), alignItems(`flexEnd), width(`px(330))]);

  let warning =
    style([
      display(`flex),
      flexDirection(`column),
      padding(`px(10)),
      color(Colors.yellow6),
      backgroundColor(Colors.yellow1),
      border(`px(1), `solid, Colors.yellow6),
      borderRadius(`px(4)),
    ]);

  let select = style([width(`px(1000)), height(`px(1))]);

  let distValidatorContainer =
    style([display(`flex), flexDirection(`column), alignItems(`flexEnd)]);
};

module DstValidatorSelection = {
  type control_t = {
    display: string,
    height: string,
    width: string,
    fontSize: string,
    backgroundColor: string,
    borderRadius: string,
    boxShadow: string,
  };

  type option_t = {
    display: string,
    alignItems: string,
    height: string,
    fontSize: string,
    paddingLeft: string,
    cursor: string,
  };

  [@react.component]
  let make = (~filteredValidators: array(BandScan.ValidatorSub.t), ~setDstValidatorOpt) => {
    let (selectedValidator, setSelectedValidator) =
      React.useState(_ =>
        ReactSelect.{value: "N/A", label: "Enter or select validator to delegate to"}
      );
    let validatorList =
      filteredValidators->Belt_Array.map(({operatorAddress, moniker}) =>
        ReactSelect.{value: operatorAddress |> Address.toOperatorBech32, label: moniker}
      );

    // TODO: Hack styles for react-select
    <div className=Styles.distValidatorContainer>
      <ReactSelect
        options=validatorList
        onChange={newOption => {
          let newVal = newOption;
          setSelectedValidator(_ => newVal);
          setDstValidatorOpt(_ => Some(newVal.value |> Address.fromBech32));
        }}
        value=selectedValidator
        styles={
          ReactSelect.control: _ => {
            display: "flex",
            height: "30px",
            width: "300px",
            fontSize: "11px",
            backgroundColor: "white",
            borderRadius: "4px",
            boxShadow: "0 1px 4px 0 rgba(11,29,142,0.1) inset",
          },
          ReactSelect.option: _ => {
            fontSize: "11px",
            height: "30px",
            display: "flex",
            alignItems: "center",
            paddingLeft: "10px",
            cursor: "pointer",
          },
        }
      />
      <VSpacing size=Spacing.xs />
      <Text
        value={"(" ++ selectedValidator.value ++ ")"}
        size=Text.Sm
        color=Colors.blueGray5
        code=true
      />
    </div>;
  };
};

[@react.component]
let make = (~address, ~validator, ~setMsgsOpt) => {
  let validatorInfoSub = ValidatorSub.get(validator);
  let validatorsSub = ValidatorSub.getList(~isActive=true, ());
  let delegationSub = DelegationSub.getStakeByValidator(address, validator);

  let allSub = Sub.all3(validatorInfoSub, validatorsSub, delegationSub);

  let (dstValidatorOpt, setDstValidatorOpt) = React.useState(_ => None);

  let (amount, setAmount) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect2(
    _ => {
      let msgsOpt = {
        let%Opt dstValidator = dstValidatorOpt;
        let%Opt amountValue = amount.value;
        Some([|
          TxCreator.Redelegate(
            validator,
            dstValidator,
            {amount: amountValue |> Js.Float.toString, denom: "uband"},
          ),
        |]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    (dstValidatorOpt, amount),
  );
  <>
    <VSpacing size=Spacing.sm />
    <div className=Styles.warning>
      <Text value="Please read before proceeding:" />
      <VSpacing size=Spacing.xs />
      <Text
        value="You can only redelegate a maximum of 7 times to/from the same validator pairs during any 21 day period."
      />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.info>
      <Text value="Current Stake" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      {switch (allSub) {
       | Data((_, _, {amount: stakedAmount})) =>
         <div>
           <Text
             value={stakedAmount |> Coin.getBandAmountFromCoin |> Format.fPretty(~digits=6)}
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
    <div className=Styles.info>
      <Text value="Redelegate From" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      {switch (allSub) {
       | Data(({moniker}, _, _)) =>
         <div className=Styles.validator>
           <Text value=moniker code=true ellipsis=true align=Text.Right />
           <Text
             value={"(" ++ validator->Address.toOperatorBech32 ++ ")"}
             size=Text.Sm
             color=Colors.blueGray5
             code=true
           />
         </div>
       | _ => <LoadingCensorBar width=300 height=26 />
       }}
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.info>
      <Text value="Redelegate To" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      {switch (allSub) {
       | Data(({operatorAddress}, validators, _)) =>
         let filteredValidators =
           validators->Belt_Array.keep(validator =>
             validator.operatorAddress != operatorAddress && validator.commission != 100.
           );
         <DstValidatorSelection filteredValidators setDstValidatorOpt />;
       | _ => <LoadingCensorBar width=300 height=43 />
       }}
    </div>
    <VSpacing size=Spacing.lg />
    {switch (allSub) {
     | Data((_, _, {amount: stakedAmount})) =>
       let maxValInUband = stakedAmount |> Coin.getUBandAmountFromCoin;
       <EnhanceTxInput
         width=300
         inputData=amount
         setInputData=setAmount
         parse={Parse.getBandAmount(maxValInUband)}
         maxValue={maxValInUband /. 1e6 |> Js.Float.toString}
         msg="Amount (BAND)"
         placeholder="Insert amount"
         inputType="number"
         code=true
         autoFocus=true
         id="redelegateAmountInput"
       />;
     | _ => <EnhanceTxInput.Loading msg="Amount (BAND)" width=300 />
     }}
    <VSpacing size=Spacing.lg />
  </>;
};
