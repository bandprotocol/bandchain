module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      flexDirection(`column),
      width(`percent(100.)),
      padding4(~top=`px(45), ~left=`px(18), ~right=`px(20), ~bottom=`px(0)),
    ]);

  let inputBar =
    style([
      width(`percent(100.)),
      height(`px(30)),
      paddingLeft(`px(9)),
      borderRadius(`px(8)),
      boxShadow(
        Shadow.box(
          ~inset=true,
          ~x=`zero,
          ~y=`px(3),
          ~blur=`px(4),
          Css.rgba(11, 29, 142, `num(0.1)),
        ),
      ),
      focus([outline(`zero, `none, Colors.white)]),
    ]);

  let mnemonicHelper =
    style([
      width(`px(130)),
      height(`px(16)),
      display(`flex),
      justifyContent(`spaceBetween),
      alignContent(`center),
      color(Css.hex("5269FF")),
    ]);

  let connectBtn =
    style([
      width(`px(100)),
      height(`px(30)),
      display(`flex),
      justifySelf(`right),
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
      cursor(`pointer),
      alignSelf(`flexEnd),
    ]);
};

[@react.component]
let make = (~chainID) => {
  let (_, dispatchAccount) = React.useContext(AccountContext.context);
  let (_, dispatchModal) = React.useContext(ModalContext.context);
  let (mnemonic, setMnemonic) = React.useState(_ => "");
  let (errMsg, setErrMsg) = React.useState(_ => "");

  let createMnemonic = () =>
    if (mnemonic->Js.String.trim == "") {
      setErrMsg(_ => "Invalid mnemonic");
    } else {
      let wallet = Wallet.createFromMnemonic(mnemonic);
      let _ =
        wallet->Wallet.getAddressAndPubKey
        |> Js.Promise.then_(((address, pubKey)) => {
             dispatchAccount(Connect(wallet, address, pubKey, chainID));
             dispatchModal(CloseModal);
             Promise.ret();
           })
        |> Js.Promise.catch(err => {
             Js.Console.log(err);
             setErrMsg(_ => "An error occurred.");
             Promise.ret();
           });
      ();
    };

  <div className=Styles.container>
    <Text value="Enter Your Mnemonic" size=Text.Md weight=Text.Medium />
    <VSpacing size=Spacing.sm />
    <input
      id="mnemonicInput"
      autoFocus=true
      value=mnemonic
      className=Styles.inputBar
      onChange={event => setMnemonic(ReactEvent.Form.target(event)##value)}
      onKeyDown={event =>
        switch (ReactEvent.Keyboard.key(event)) {
        | "Enter" =>
          createMnemonic();
          ReactEvent.Keyboard.preventDefault(event);
        | _ => ()
        }
      }
    />
    <VSpacing size={`px(35)} />
    <div id="mnemonicConnectButton" className={CssHelper.flexBox(~justify=`flexEnd, ())}>
      <Button px=20 py=8 onClick={_ => createMnemonic()}>
        <Text value="Connect" weight=Text.Bold size=Text.Md color=Colors.white />
      </Button>
    </div>
    <VSpacing size=Spacing.lg />
    <Text value=errMsg color=Colors.red6 />
  </div>;
  // </Col>
  //   </div>
  //     <img src=Images.linkIcon />
  //     <Text value="What is Mnemonic" />
  //   <div className=Styles.mnemonicHelper>
  // <Col>
};
