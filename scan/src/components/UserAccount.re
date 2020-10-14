module Styles = {
  open Css;

  let container = style([position(`relative)]);

  let loading = style([width(`px(65)), height(`px(20))]);

  let oval =
    style([
      display(`flex),
      width(`px(24)),
      height(`px(24)),
      justifyContent(`center),
      alignItems(`center),
      padding(`px(5)),
      backgroundColor(Colors.bandBlue),
      borderRadius(`percent(50.)),
    ]);

  let logo = style([width(`px(12))]);

  let profileCard = show =>
    style([
      position(`absolute),
      backgroundColor(Colors.white),
      top(`px(30)),
      right(`px(-10)),
      borderRadius(`px(4)),
      padding(`px(16)),
      boxShadow(Shadow.box(~x=`zero, ~y=`zero, ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08)))),
      transition(~duration=200, "all"),
      opacity(show ? 1. : 0.),
      pointerEvents(show ? `auto : `none),
      zIndex(5),
    ]);

  let innerProfileCard = style([padding(`px(16)), backgroundColor(Colors.profileBG)]);

  let connect = style([padding2(~v=`px(10), ~h=`zero)]);
  let disconnect = style([paddingTop(`px(16))]);
};

module ConnectBtn = {
  [@react.component]
  let make = (~connect) => {
    <div
      id="connectButton"
      className={Css.merge([
        CssHelper.flexBox(~justify=`center, ~align=`center, ()),
        CssHelper.clickable,
        Styles.connect,
      ])}
      onClick={_ => connect()}>
      <Text
        value="Connect"
        weight=Text.Medium
        color=Colors.bandBlue
        nowrap=true
        spacing={Text.Em(0.03)}
        block=true
      />
      <HSpacing size=Spacing.sm />
      <div className=Styles.oval> <Icon name="fal fa-link" color=Colors.white /> </div>
    </div>;
  };
};

module DisconnectBtn = {
  [@react.component]
  let make = (~disconnect) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~justify=`center, ~align=`center, ()),
        CssHelper.clickable,
        Styles.disconnect,
      ])}
      onClick={_ => disconnect()}>
      <Text value="Disconnect" weight=Text.Medium color=Colors.bandBlue nowrap=true block=true />
    </div>;
  };
};

module FaucetBtn = {
  [@react.component]
  let make = (~address) => {
    let (isRequest, setIsRequest) = React.useState(_ => false);
    isRequest
      ? <img src=Images.loadingCircles className=Styles.loading />
      : <div id="getFreeButton">
          <Button
            px=20
            py=5
            variant=Button.Outline
            onClick={_ => {
              setIsRequest(_ => true);
              let _ =
                AxiosFaucet.request({address: address |> Address.toBech32, amount: 10_000_000})
                |> Js.Promise.then_(_ => {
                     setIsRequest(_ => false);
                     Js.Promise.resolve();
                   });
              ();
            }}>
            <Text value="Get 10 Testnet BAND" weight=Text.Medium nowrap=true />
          </Button>
        </div>;
  };
};

module SendBtn = {
  [@react.component]
  let make = (~send) => {
    <div id="sendToken">
      <Button px=20 py=5 onClick={_ => {send()}}>
        <Text value="Send" weight=Text.Medium nowrap=true block=true />
      </Button>
    </div>;
  };
};

module Balance = {
  [@react.component]
  let make = (~address) => {
    let accountSub = AccountSub.get(address);

    <div className={CssHelper.flexBox(~justify=`spaceBetween, ())}>
      <Text value="Balance" weight=Text.Medium />
      <div className={CssHelper.flexBox()} id="bandBalance">
        <Text
          value={
            switch (accountSub) {
            | Data(account) =>
              account.balance |> Coin.getBandAmountFromCoins |> Format.fPretty(~digits=6)
            | _ => "0"
            }
          }
          code=true
        />
        <HSpacing size=Spacing.sm />
        <Text value="BAND" weight=Text.Thin />
      </div>
    </div>;
  };
};

[@react.component]
let make = () => {
  let trackingSub = TrackingSub.use();
  let (accountOpt, dispatchAccount) = React.useContext(AccountContext.context);
  let (_, dispatchModal) = React.useContext(ModalContext.context);
  let (show, setShow) = React.useState(_ => false);

  let connect = chainID => dispatchModal(OpenModal(Connect(chainID)));
  let disconnect = () => {
    dispatchAccount(Disconnect);
    setShow(_ => false);
  };
  let send = () => dispatchModal(OpenModal(SubmitTx(SubmitMsg.Send(None))));

  switch (accountOpt) {
  | Some({address}) =>
    <div className={Css.merge([CssHelper.flexBox(~justify=`flexEnd, ()), Styles.container])}>
      <div
        id="userInfoButton"
        className={Css.merge([CssHelper.flexBox(), CssHelper.clickable])}
        onClick={_ => setShow(prev => !prev)}>
        <div className=Styles.oval> <Icon name="fal fa-user" color=Colors.white /> </div>
        <HSpacing size=Spacing.sm />
        <Icon name="fas fa-caret-down" color=Colors.bandBlue />
      </div>
      <div className={Styles.profileCard(show)} id="addressWrapper">
        <AddressRender address position=AddressRender.Text />
        <VSpacing size={`px(16)} />
        <div className=Styles.innerProfileCard>
          <Balance address />
          <VSpacing size={`px(16)} />
          <div className={CssHelper.flexBox(~direction=`row, ~justify=`flexEnd, ())}>
            // <FaucetBtn address />
             <SendBtn send /> </div>
        </div>
        <DisconnectBtn disconnect />
      </div>
    </div>
  | None =>
    <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
      {switch (trackingSub) {
       | Data({chainID}) => <ConnectBtn connect={_ => connect(chainID)} />
       | Error(err) =>
         // log for err details
         Js.Console.log(err);
         <Text value="chain id not found" />;
       | _ => <LoadingCensorBar width=80 height=18 />
       }}
    </div>
  };
};
