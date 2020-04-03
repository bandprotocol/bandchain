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

[@react.component]
let make = () => {
  let (addressOpt, dispatchAccount) = React.useContext(AccountContext.context);

  let connect = () => {
    let mnemonicOpt = Window.prompt("Please enter your mnemonic.", "") |> Js.Nullable.toOption;

    switch (mnemonicOpt) {
    | Some(mnemonic) => dispatchAccount(Connect(mnemonic))
    | None => ()
    };
  };

  let disconnect = () => dispatchAccount(Disconnect);

  <>
    <Row justify=Row.Right>
      {switch (addressOpt) {
       | Some(address) =>
         <>
           <Col> <AddressRender address position=AddressRender.Nav /> </Col>
           <Col> <HSpacing size={`px(27)} /> </Col>
           <Col> <DisconnectBtn disconnect /> </Col>
         </>
       | None => <Col> <ConnectBtn connect /> </Col>
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
