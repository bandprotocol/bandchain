module Styles = {
  open Css;

  let withWH = (w, h) =>
    style([
      width(w),
      height(h),
      display(`flex),
      justifyContent(`flexEnd),
      alignItems(`center),
    ]);

  let connectBtn =
    style([
      backgroundColor(Colors.green1),
      padding2(~h=`px(8), ~v=`px(2)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(10)),
      cursor(`pointer),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(17, 85, 78, 0.1))),
    ]);

  let disconnectBtn =
    style([
      backgroundColor(Colors.yellowAccent1),
      padding2(~h=`px(8), ~v=`px(2)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(10)),
      cursor(`pointer),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(99, 81, 3, 0.1))),
    ]);

  let faucetBtn =
    style([
      backgroundColor(Colors.blue1),
      padding2(~h=`px(8), ~v=`px(2)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(10)),
      cursor(`pointer),
      height(`px(16)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(11, 29, 142, 0.1))),
    ]);

  let logo = style([width(`px(10))]);

  let balanceContainer = style([display(`flex), alignItems(`center)]);

  let modalLogin =
    style([
      width(`px(640)),
      height(`px(480)),
      backgroundColor(Css.rgb(249, 249, 251)),
      backgroundImage(`url(Images.modalBg)),
      borderRadius(`px(5)),
    ]);

  let modalTitle = style([display(`flex), justifyContent(`center)]);

  let modalSelectText = style([display(`flex), marginLeft(`px(34))]);

  let container = style([backgroundColor(Colors.transparent)]);

  let header = active =>
    style([
      display(`flex),
      flexDirection(`row),
      alignSelf(`center),
      alignItems(`center),
      justifyContent(`spaceBetween),
      cursor(`pointer),
      width(`px(480)),
      padding2(~v=`zero, ~h=`px(20)),
      color(active ? Colors.gray8 : Colors.gray6),
      backgroundColor(Colors.white),
    ]);

  let buttonContainer = active =>
    style([
      display(`flex),
      width(`px(226)),
      height(`px(50)),
      marginLeft(`px(34)),
      borderLeft(`px(6), `solid, active ? Colors.bandBlue : Colors.white),
      borderRadius(`px(8)),
      backgroundColor(Colors.white),
      boxShadow(
        active
          ? Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(11, 29, 142, 0.1))
          : Shadow.box(~x=`zero, ~y=`px(0), ~blur=`px(0), Css.rgba(0, 0, 0, 0.)),
      ),
    ]);
  let ledgerIcon = style([height(`px(30)), width(`px(30)), display(`flex)]);
  let seperatedLongLine =
    style([height(`px(275)), width(`px(2)), backgroundColor(Colors.gray4)]);

  let itemCol =
    style([height(`px(275)), display(`flex), flexDirection(`column), verticalAlign(`top)]);
};

module ConnectBtn = {
  [@react.component]
  let make = (~connect) => {
    <div className=Styles.connectBtn onClick={_ => connect()}>
      <Text
        value="connect"
        size=Text.Xs
        weight=Text.Medium
        color=Colors.green7
        nowrap=true
        height={Text.Px(10)}
        spacing={Text.Em(0.03)}
        block=true
      />
      <HSpacing size=Spacing.sm />
      <img src=Images.connectIcon className=Styles.logo />
    </div>;
  };
};

module DisconnectBtn = {
  [@react.component]
  let make = (~disconnect) => {
    <div className=Styles.disconnectBtn onClick={_ => disconnect()}>
      <Text
        value="disconnect"
        size=Text.Xs
        weight=Text.Medium
        color=Colors.yellowAccent7
        nowrap=true
        height={Text.Px(10)}
        spacing={Text.Em(0.03)}
        block=true
      />
      <HSpacing size=Spacing.sm />
      <img src=Images.disconnectIcon className=Styles.logo />
    </div>;
  };
};

module FaucetBtn = {
  let loadingRender = (wDiv, wImg, h) => {
    <div className={Styles.withWH(wDiv, h)}>
      <img src=Images.loadingCircles className={Styles.withWH(wImg, h)} />
    </div>;
  };

  [@react.component]
  let make = (~address) => {
    let (isRequest, setIsRequest) = React.useState(_ => false);
    isRequest
      ? loadingRender(`pxFloat(98.5), `px(65), `px(16))
      : <div
          className=Styles.faucetBtn
          onClick={_ => {
            setIsRequest(_ => true);
            let _ =
              AxiosFaucet.request({address, amount: 10_000_000})
              |> Js.Promise.then_(_ => {
                   setIsRequest(_ => false);
                   Js.Promise.resolve();
                 });
            ();
          }}>
          <Text
            value="get 10 testnet BAND"
            size=Text.Xs
            weight=Text.Medium
            color=Colors.blue7
            nowrap=true
            height={Text.Px(10)}
            spacing={Text.Em(0.03)}
            block=true
          />
        </div>;
  };
};

module Balance = {
  [@react.component]
  let make = (~address) =>
    {
      let accountSub = AccountSub.get(address);
      let%Sub account = accountSub;

      <div className=Styles.balanceContainer>
        <Text
          value={account.balance |> Coin.getBandAmountFromCoins |> Js.Float.toString}
          code=true
          size=Text.Sm
          height={Text.Px(13)}
        />
        <HSpacing size=Spacing.sm />
        <Text value="BAND" size=Text.Sm height={Text.Px(13)} weight=Text.Thin />
      </div>
      |> Sub.resolve;
    }
    |> Sub.default(
         _,
         <div className=Styles.balanceContainer>
           <Text value="0" code=true size=Text.Sm height={Text.Px(13)} />
           <HSpacing size=Spacing.sm />
           <Text value="BAND" size=Text.Sm height={Text.Px(13)} weight=Text.Thin />
         </div>,
       );
};

type login_method_t =
  | Mnemonic
  | Ledger;

let toLoginMethodString = method => {
  switch (method) {
  | Mnemonic => "Mnemonic Phrase"
  | Ledger => "Ledger"
  };
};

module LoginMethod = {
  [@react.component]
  let make = (~name, ~active) => {
    <div className={Styles.buttonContainer(active)}>
      <div className={Styles.header(active)}>
        <Text value={name |> toLoginMethodString} weight=Text.Medium size=Text.Md />
        {switch (name) {
         | Ledger =>
           active
             ? <img src=Images.ledgerIconActive className=Styles.ledgerIcon />
             : <img src=Images.ledgerIconInactive className=Styles.ledgerIcon />
         | _ => <div />
         }}
      </div>
    </div>;
  };
};

module ModalLogin = {
  [@react.component]
  let make = () => {
    let (loginMethod, setLoginMethod) = React.useState(_ => Mnemonic);
    <div className=Styles.modalLogin>
      <VSpacing size=Spacing.xxl />
      <div className=Styles.modalTitle>
        <Text value="Connect With Your Wallet" weight=Text.Bold size=Text.Xxxl />
      </div>
      <VSpacing size=Spacing.xxl />
      <VSpacing size=Spacing.sm />
      <div className=Styles.modalSelectText>
        <Text value="Select your connection method" size=Text.Lg weight=Text.Medium />
      </div>
      <Row>
        <div className=Styles.itemCol>
          {[|Mnemonic, Ledger|]
           ->Belt_Array.map(method =>
               <div>
                 <VSpacing size=Spacing.md />
                 <VSpacing size=Spacing.xs />
                 <div onClick={_ => setLoginMethod(_ => method)}>
                   <LoginMethod name=method active={loginMethod == method} />
                 </div>
               </div>
             )
           ->React.array}
        </div>
        <HSpacing size=Spacing.lg />
        <HSpacing size=Spacing.sm />
        <div className=Styles.itemCol> <div className=Styles.seperatedLongLine /> </div>
        <HSpacing size=Spacing.lg />
        <HSpacing size=Spacing.sm />
        <div className=Styles.itemCol>
          {switch (loginMethod) {
           | Mnemonic => "Mnemonic Phrase" |> React.string
           | Ledger => "Ledger" |> React.string
           }}
        </div>
      </Row>
    </div>;
  };
};

[@react.component]
let make = () => {
  let (addressOpt, dispatch) = React.useContext(AccountContext.context);
  let (_, dispatchModal) = React.useContext(ModalContext.context);

  let connect = () => {
    let mnemonicOpt = Window.prompt("Please enter your mnemonic.", "") |> Js.Nullable.toOption;

    switch (mnemonicOpt) {
    | Some(mnemonic) => dispatch(Connect(mnemonic))
    | None => ()
    };
  };

  let disconnect = () => dispatch(Disconnect);

  <>
    <Row justify=Row.Right>
      {switch (addressOpt) {
       | Some(address) =>
         <>
           <Col> <AddressRender address position=AddressRender.Nav /> </Col>
           <Col> <HSpacing size={`px(27)} /> </Col>
           <Col> <DisconnectBtn disconnect /> </Col>
         </>
       | None =>
         <Col>
           <ConnectBtn connect />
           // TODO: remove later
           <button
             onClick={_ => {
               dispatchModal(OpenModal(Connect(<ModalLogin />)));
               ();
             }}>
             {"modal" |> React.string}
           </button>
         </Col>
       }}
    </Row>
    {switch (addressOpt) {
     | Some(address) =>
       <>
         <VSpacing size=Spacing.md />
         <Row justify=Row.Right>
           <Col> <Balance address /> </Col>
           <Col> <HSpacing size={`px(5)} /> </Col>
           <Col> <FaucetBtn address={addressOpt->Belt_Option.getExn->Address.toBech32} /> </Col>
         </Row>
       </>
     | None => React.null
     }}
  </>;
};
