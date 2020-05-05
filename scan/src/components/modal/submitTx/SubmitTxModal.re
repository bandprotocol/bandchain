module Styles = {
  open Css;

  let container =
    style([
      flexDirection(`column),
      width(`px(640)),
      height(`px(480)),
      padding2(~v=`px(50), ~h=`px(50)),
      backgroundColor(rgb(249, 249, 251)),
      borderRadius(`px(5)),
      justifyContent(`flexStart),
    ]);

  let disable = isActive => style([display(isActive ? `flex : `none)]);

  let modalTitle = style([display(`flex), justifyContent(`center)]);

  let rowContainer =
    style([display(`flex), alignItems(`center), justifyContent(`spaceBetween)]);

  let input = wid =>
    style([
      width(`px(wid)),
      height(`px(30)),
      paddingRight(`px(9)),
      borderRadius(`px(8)),
      fontSize(`px(12)),
      textAlign(`right),
      boxShadow(
        Shadow.box(
          ~inset=false,
          ~x=`zero,
          ~y=`px(3),
          ~blur=`px(4),
          Css.rgba(11, 29, 142, 0.1),
        ),
      ),
      focus([outlineColor(Colors.white)]),
    ]);

  let selectWrapper =
    style([
      display(`flex),
      flexDirection(`row),
      padding2(~v=`px(3), ~h=`px(8)),
      position(`static),
      width(`px(130)),
      height(`px(30)),
      left(`zero),
      top(`px(32)),
      background(rgba(255, 255, 255, 1.)),
      borderRadius(`px(100)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(0, 0, 0, 0.1))),
      float(`left),
      fontSize(`px(14)),
    ]);

  let selectContent =
    style([
      background(rgba(255, 255, 255, 1.)),
      border(`px(0), `solid, hex("FFFFFF")),
      width(`px(135)),
      focus([outlineColor(Colors.white)]),
    ]);

  let nextBtn =
    style([
      width(`px(100)),
      height(`px(30)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      backgroundImage(
        `linearGradient((
          `deg(90.),
          [(`percent(0.), Colors.blue7), (`percent(100.), Colors.bandBlue)],
        )),
      ),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(82, 105, 255, 0.25))),
      borderRadius(`px(4)),
      cursor(`pointer),
      alignSelf(`center),
    ]);
};

type selection_t =
  | Send
  | Delegate
  | Undelegate;

exception WrongMsgChoice(string);

let toVariant =
  fun
  | "Send" => Send
  | "Delegate" => Delegate
  | "Undelegate" => Undelegate
  | _ => raise(WrongMsgChoice("Chosen message does not exist"));

let toString =
  fun
  | Send => "Send"
  | Delegate => "Delegate"
  | Undelegate => "Undelegate";

module SubmitTxStep = {
  [@react.component]
  let make = (~account: AccountContext.t, ~setRawTx, ~isActive) => {
    let (msgType, setMsgType) = React.useState(_ => Send);
    let (msgsOpt, setMsgsOpt) = React.useState(_ => None);
    let (gas, setGas) =
      React.useState(_ => EnhanceTxInput.{text: "200000", value: Some(200000)});
    let (fee, setFee) = React.useState(_ => EnhanceTxInput.{text: "100", value: Some(100)});
    let (memo, setMemo) = React.useState(_ => EnhanceTxInput.{text: "", value: Some("")});

    <div className={Css.merge([Styles.container, Styles.disable(isActive)])}>
      <div className=Styles.modalTitle>
        <Text value="Submit Transaction" weight=Text.Bold size=Text.Xxxl />
      </div>
      <VSpacing size=Spacing.xl />
      <div className=Styles.rowContainer>
        <Text value="Message Type" size=Text.Lg spacing={Text.Em(0.03)} />
        <div className=Styles.selectWrapper>
          <select
            className=Styles.selectContent
            onChange={event => {
              let newMsg = ReactEvent.Form.target(event)##value |> toVariant;
              setMsgType(_ => newMsg);
            }}>
            {[|Send, Delegate, Undelegate|]
             ->Belt_Array.map(symbol =>
                 <option value={symbol |> toString}> {symbol |> toString |> React.string} </option>
               )
             |> React.array}
          </select>
        </div>
      </div>
      <VSpacing size=Spacing.md />
      {switch (msgType) {
       | Send => <SendMsg setMsgsOpt />
       | Delegate => <DelegateMsg setMsgsOpt />
       | Undelegate => <UndelegateMsg setMsgsOpt />
       }}
      <VSpacing size=Spacing.md />
      <EnhanceTxInput
        width=115
        inputData=gas
        setInputData=setGas
        parse=int_of_string_opt
        msg="Gas"
        errMsg="Invalid amount"
      />
      <VSpacing size=Spacing.md />
      <EnhanceTxInput
        width=115
        inputData=fee
        setInputData=setFee
        parse=int_of_string_opt
        msg="Fee"
        errMsg="Invalid amount"
      />
      <VSpacing size=Spacing.md />
      <EnhanceTxInput
        width=300
        inputData=memo
        setInputData=setMemo
        parse={newVal => {newVal->Js.String.length <= 32 ? Some(newVal) : None}}
        msg="Memo"
        errMsg="Exceed limit length"
      />
      <VSpacing size=Spacing.xxl />
      <div
        className=Styles.nextBtn
        onClick={_ => {
          let rawTxOpt =
            {let%Opt fee' = fee.value;
             let%Opt gas' = gas.value;
             let%Opt memo' = memo.value;
             let%Opt msgs = msgsOpt;

             Some(
               TxCreator.createRawTx(
                 ~address=account.address,
                 ~msgs,
                 ~chainID=account.chainID,
                 ~feeAmount=fee' |> string_of_int,
                 ~gas=gas' |> string_of_int,
                 ~memo=memo',
                 (),
               ),
             )};
          let _ =
            switch (rawTxOpt) {
            | Some(rawTxPromise) =>
              let%Promise rawTx = rawTxPromise;
              setRawTx(_ => Some(rawTx));
              Promise.ret();
            | None =>
              Window.alert("invalid msgs");
              Promise.ret();
            };
          ();
        }}>
        <Text value="Next" weight=Text.Bold size=Text.Md color=Colors.white />
      </div>
    </div>;
  };
};

module CreateTxFlow = {
  [@react.component]
  let make = (~account) => {
    let (rawTx, setRawTx) = React.useState(_ => None);

    <>
      <SubmitTxStep account setRawTx isActive={rawTx == None} />
      {switch (rawTx) {
       | None => React.null
       | Some(rawTx') =>
         <PreviewJsonStep rawTx=rawTx' onBack={_ => setRawTx(_ => None)} account />
       }}
    </>;
  };
};

[@react.component]
let make = () => {
  let (account, _) = React.useContext(AccountContext.context);

  switch (account) {
  | Some(account') => <CreateTxFlow account=account' />
  | None => <div className=Styles.container> <Text value="Please sign in" size=Text.Lg /> </div>
  };
};
