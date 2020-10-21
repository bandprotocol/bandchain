module Styles = {
  open Css;

  let container = style([paddingBottom(`px(24))]);

  let warning =
    style([
      display(`flex),
      flexDirection(`column),
      padding2(~v=`px(16), ~h=`px(24)),
      backgroundColor(Colors.profileBG),
      borderRadius(`px(4)),
      marginBottom(`px(24)),
    ]);

  let select = style([width(`px(1000)), height(`px(1))]);
};

module DstValidatorSelection = {
  type control_t = {
    display: string,
    height: string,
    width: string,
    fontSize: string,
    backgroundColor: string,
    borderRadius: string,
    border: string,
  };

  type option_t = {
    display: string,
    alignItems: string,
    height: string,
    fontSize: string,
    paddingLeft: string,
    cursor: string,
  };

  type container_t = {
    width: string,
    position: string,
    boxSizing: string,
  };

  type singleValue_t = {
    margin: string,
    maxWidth: string,
    overflow: string,
    position: string,
    textOverflow: string,
    whiteSpace: string,
    top: string,
    transform: string,
    boxSizing: string,
    fontWeight: string,
    lineHeight: string,
  };

  type indicatorSeparator_t = {display: string};

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
    <div
      className={CssHelper.flexBox(~align=`flexStart, ~direction=`column, ())}
      id="redelegateContainer">
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
            height: "37px",
            width: "100%",
            fontSize: "14px",
            backgroundColor: "white",
            borderRadius: "4px",
            border: "#EAEAEA solid 1px",
          },
          ReactSelect.option: _ => {
            fontSize: "14px",
            height: "37px",
            display: "flex",
            alignItems: "center",
            paddingLeft: "10px",
            cursor: "pointer",
          },
          ReactSelect.container: _ => {
            width: "100%",
            position: "relative",
            boxSizing: "border-box",
          },
          ReactSelect.singleValue: _ => {
            margin: "0px 2px",
            maxWidth: "calc(100% - 8px)",
            overflow: "hidden",
            position: "absolute",
            textOverflow: "ellipsis",
            whiteSpace: "nowrap",
            top: "50%",
            transform: "translateY(-50%)",
            boxSizing: "border-box",
            fontWeight: "300",
            lineHeight: "1.3em",
          },
          ReactSelect.indicatorSeparator: _ => {display: "none"},
        }
      />
      <VSpacing size=Spacing.sm />
      <Text
        value={"(" ++ selectedValidator.value ++ ")"}
        size=Text.Md
        color=Colors.gray6
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
    <div className=Styles.warning>
      <Heading
        value="Please read before proceeding:"
        size=Heading.H5
        marginBottom=4
        align=Heading.Left
      />
      <Text
        value="You can only redelegate a maximum of 7 times to/from the same validator pairs during any 21 day period."
      />
    </div>
    <div className=Styles.container>
      <Heading
        value="Redelegate from"
        size=Heading.H5
        marginBottom=8
        align=Heading.Left
        weight=Heading.Medium
      />
      {switch (allSub) {
       | Data(({moniker}, _, _)) =>
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
        value="Redelegate to"
        size=Heading.H5
        marginBottom=8
        align=Heading.Left
        weight=Heading.Medium
      />
      {switch (allSub) {
       | Data(({operatorAddress}, validators, _)) =>
         let filteredValidators =
           validators->Belt_Array.keep(validator =>
             validator.operatorAddress != operatorAddress && validator.commission != 100.
           );
         <DstValidatorSelection filteredValidators setDstValidatorOpt />;
       | _ => <LoadingCensorBar width=300 height=59 />
       }}
    </div>
    <div className=Styles.container>
      <Heading
        value="Current Stake"
        size=Heading.H5
        marginBottom=8
        align=Heading.Left
        weight=Heading.Medium
      />
      {switch (allSub) {
       | Data((_, _, {amount: stakedAmount})) =>
         <div>
           <Text
             value={stakedAmount |> Coin.getBandAmountFromCoin |> Format.fPretty(~digits=6)}
             code=true
             size=Text.Lg
           />
           <Text value=" BAND" size=Text.Lg code=true />
         </div>
       | _ => <LoadingCensorBar width=150 height=18 />
       }}
    </div>
    {switch (allSub) {
     | Data((_, _, {amount: stakedAmount})) =>
       let maxValInUband = stakedAmount |> Coin.getUBandAmountFromCoin;
       <EnhanceTxInput
         width=300
         inputData=amount
         setInputData=setAmount
         parse={Parse.getBandAmount(maxValInUband)}
         maxValue={maxValInUband /. 1e6 |> Js.Float.toString}
         msg="Amount"
         placeholder="0.000000"
         inputType="number"
         code=true
         autoFocus=true
         id="redelegateAmountInput"
       />;
     | _ => <EnhanceTxInput.Loading msg="Amount" code=true useMax=true placeholder="0.000000" />
     }}
  </>;
};
