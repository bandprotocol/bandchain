module Styles = {
  open Css;

  let container =
    style([
      flexDirection(`column),
      width(`px(640)),
      minHeight(`px(300)),
      height(`auto),
      padding2(~v=`px(50), ~h=`px(50)),
      backgroundColor(rgb(249, 249, 251)),
      borderRadius(`px(5)),
      justifyContent(`flexStart),
    ]);

  let disable = isActive => style([display(isActive ? `flex : `none)]);

  let modalTitle = style([display(`flex), justifyContent(`center)]);

  let rowContainer =
    style([display(`flex), alignItems(`center), justifyContent(`spaceBetween)]);

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
      background(rgba(255, 255, 255, `num(1.))),
      borderRadius(`px(100)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(0, 0, 0, `num(0.1)))),
      float(`left),
      fontSize(`px(14)),
    ]);

  let selectContent =
    style([
      background(rgba(255, 255, 255, `num(1.))),
      border(`px(0), `solid, hex("FFFFFF")),
      width(`px(135)),
      focus([outlineColor(Colors.white)]),
    ]);

  let nextBtn =
    style([
      marginTop(`px(30)),
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
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(82, 105, 255, `num(0.25))),
      ),
      borderRadius(`px(4)),
      border(`zero, `solid, Colors.white),
      alignSelf(`center),
      cursor(`pointer),
      color(Colors.white),
      transition(~duration=600, "all"),
      disabled([
        backgroundImage(
          `linearGradient((
            `deg(90.),
            [(`percent(0.), Colors.gray3), (`percent(100.), Colors.gray3)],
          )),
        ),
        color(Colors.gray6),
        boxShadow(
          Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(11, 29, 142, `num(0.1))),
        ),
        cursor(`default),
      ]),
    ]);

  let info = style([display(`flex), justifyContent(`spaceBetween)]);
};

module SubmitTxStep = {
  [@react.component]
  let make = (~account: AccountContext.t, ~setRawTx, ~isActive, ~msg) => {
    let (msgsOpt, setMsgsOpt) = React.useState(_ => None);

    let gas = SubmitMsg.gasLimit(msg);
    let fee = 5000.;
    let (memo, setMemo) = React.useState(_ => EnhanceTxInput.{text: "", value: Some("")});

    <div className={Css.merge([Styles.container, Styles.disable(isActive)])}>
      <div className=Styles.modalTitle>
        <Text value="Submit Transaction" weight=Text.Bold size=Text.Xxxl />
      </div>
      <VSpacing size=Spacing.xl />
      <div className=Styles.rowContainer>
        <Text value="Message Type" size=Text.Lg spacing={Text.Em(0.03)} />
        <Text
          value={SubmitMsg.toString(msg)}
          size=Text.Md
          spacing={Text.Em(0.03)}
          weight=Text.Semibold
        />
      </div>
      <VSpacing size=Spacing.md />
      {switch (msg) {
       | SubmitMsg.Send(receiver) => <SendMsg address={account.address} receiver setMsgsOpt />
       | Delegate(validator) => <DelegateMsg address={account.address} validator setMsgsOpt />
       | Undelegate(validator) => <UndelegateMsg address={account.address} validator setMsgsOpt />
       | Redelegate(validator) => <RedelegateMsg address={account.address} validator setMsgsOpt />
       | WithdrawReward(validator) =>
         <WithdrawRewardMsg validator setMsgsOpt address={account.address} />
       | Vote(proposalID, proposalName) => <VoteMsg proposalID proposalName setMsgsOpt />
       }}
      <VSpacing size=Spacing.sm />
      <EnhanceTxInput
        width=300
        inputData=memo
        setInputData=setMemo
        parse={newVal => {
          newVal->Js.String.length <= 32 ? Result.Ok(newVal) : Err("Exceed limit length")
        }}
        msg="Memo (optional)"
        placeholder="Insert memo"
        code=true
      />
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.md />
      <div className=Styles.info>
        <Text
          value="Transaction Fee"
          size=Text.Lg
          spacing={Text.Em(0.03)}
          nowrap=true
          block=true
        />
        <Text value="0.005 BAND" code=true />
      </div>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.md />
      <button
        className=Styles.nextBtn
        disabled={msgsOpt->Belt.Option.isNone}
        onClick={_ => {
          let rawTxOpt =
            {let%Opt memo' = memo.value;
             let%Opt msgs = msgsOpt;

             Some(
               TxCreator.createRawTx(
                 ~address=account.address,
                 ~msgs,
                 ~chainID=account.chainID,
                 ~feeAmount=fee |> Js.Float.toString,
                 ~gas=gas |> string_of_int,
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
              Webapi.Dom.(window |> Window.alert("invalid msgs"));
              Promise.ret();
            };
          ();
        }}>
        <Text value="Next" weight=Text.Bold size=Text.Md />
      </button>
    </div>;
  };
};

module CreateTxFlow = {
  [@react.component]
  let make = (~account, ~msg) => {
    let (rawTx, setRawTx) = React.useState(_ => None);

    <>
      <SubmitTxStep account setRawTx isActive={rawTx == None} msg />
      {switch (rawTx) {
       | None => React.null
       | Some(rawTx') =>
         <PreviewJsonStep rawTx=rawTx' onBack={_ => setRawTx(_ => None)} account />
       }}
    </>;
  };
};

[@react.component]
let make = (~msg) => {
  let (account, _) = React.useContext(AccountContext.context);

  switch (account) {
  | Some(account') => <CreateTxFlow account=account' msg />
  | None => <div className=Styles.container> <Text value="Please sign in" size=Text.Lg /> </div>
  };
};
