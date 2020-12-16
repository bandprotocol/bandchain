module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      flexDirection(`column),
      width(`percent(100.)),
      padding2(~v=`px(0), ~h=`px(18)),
    ]);

  let instructionCard =
    style([
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      height(`px(50)),
      width(`percent(100.)),
      justifyContent(`spaceBetween),
      backgroundColor(Css.rgb(255, 255, 255)),
    ]);

  let oval =
    style([
      backgroundImage(
        `linearGradient((
          `deg(180.),
          [(`percent(0.), Colors.blue2), (`percent(100.), Colors.bandBlue)],
        )),
      ),
      width(`px(30)),
      height(`px(30)),
      borderRadius(`percent(50.)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
    ]);

  let rFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let resultContainer =
    style([
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      justifyContent(`spaceBetween),
      height(`px(35)),
    ]);

  let ledgerGuide = style([width(`px(248)), height(`px(38))]);

  let connectBtn = (~isLoading, ()) =>
    style([
      marginTop(`px(10)),
      width(`percent(100.)),
      height(`px(36)),
      display(`flex),
      justifySelf(`right),
      justifyContent(`center),
      alignItems(`center),
      backgroundColor(isLoading ? Colors.blueGray3 : Colors.bandBlue),
      boxShadow(
        isLoading
          ? `none
          : Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(82, 105, 255, `num(0.25))),
      ),
      borderRadius(`px(4)),
      cursor(isLoading ? `default : `pointer),
      pointerEvents(isLoading ? `none : `auto),
      alignSelf(`flexEnd),
    ]);

  let selectWrapper =
    style([
      display(`flex),
      padding2(~v=`px(3), ~h=`px(8)),
      justifyContent(`center),
      alignItems(`center),
      width(`percent(100.)),
      height(`px(37)),
      left(`zero),
      top(`px(32)),
      background(rgba(255, 255, 255, `num(1.))),
      border(`px(1), `solid, Colors.blueGray3),
      borderRadius(`px(6)),
      float(`left),
    ]);

  let selectContent =
    style([
      background(rgba(255, 255, 255, `num(1.))),
      border(`px(0), `solid, hex("FFFFFF")),
      width(`px(100)),
      lineHeight(`em(1.41)),
    ]);

  let connectingBtnContainer =
    style([width(`px(104)), display(`flex), justifyContent(`spaceBetween)]);
};

module InstructionCard = {
  [@react.component]
  let make = (~title, ~url) => {
    <div className=Styles.instructionCard>
      <div className=Styles.rFlex> <Text value=title /> </div>
      <img src=url className=Styles.ledgerGuide />
    </div>;
  };
};

type result_t =
  | Nothing
  | Loading
  | Error(string);

[@react.component]
let make = (~chainID, ~ledgerApp) => {
  let (_, dispatchAccount) = React.useContext(AccountContext.context);
  let (_, dispatchModal) = React.useContext(ModalContext.context);
  let (result, setResult) = React.useState(_ => Nothing);
  let (accountIndex, setAccountIndex) = React.useState(_ => 0);

  let createLedger = accountIndex => {
    dispatchModal(DisableExit);
    setResult(_ => Loading);
    let _ =
      Wallet.createFromLedger(ledgerApp, accountIndex)
      |> Js.Promise.then_(wallet => {
           let%Promise (address, pubKey) = wallet->Wallet.getAddressAndPubKey;
           dispatchAccount(Connect(wallet, address, pubKey, chainID));
           dispatchModal(EnableExit);
           dispatchModal(CloseModal);
           Promise.ret();
         })
      |> Js.Promise.catch(err => {
           Js.Console.log(err);
           setResult(_ => Error("An error occured"));
           dispatchModal(EnableExit);
           Promise.ret();
         });
    ();
  };

  <div className=Styles.container>
    <VSpacing size=Spacing.xl />
    <Text value="1. Select HD Derivation Path" weight=Text.Semibold size=Text.Lg />
    <VSpacing size=Spacing.md />
    <div className=Styles.selectWrapper>
      <div className={CssHelper.selectWrapper(~pRight=8, ~mW=100, ~size=10, ())}>
        <select
          className=Styles.selectContent
          onChange={event => {
            let newAccountIndex = ReactEvent.Form.target(event)##value |> int_of_string;
            setAccountIndex(_ => newAccountIndex);
          }}>
          {[|0, 1, 2, 3, 4, 5, 6, 7, 8, 9|]
           |> Belt.Array.map(_, index =>
                <option key={index |> string_of_int} value={index |> string_of_int}>
                  {let prefixPath =
                     switch (ledgerApp) {
                     | Ledger.Cosmos => "44/118/0/0/"
                     | BandChain => "44/494/0/0/"
                     };
                   prefixPath ++ (index |> string_of_int) |> React.string}
                </option>
              )
           |> React.array}
        </select>
      </div>
    </div>
    <VSpacing size=Spacing.xl />
    <Text value="2. On Your Ledger" weight=Text.Semibold size=Text.Lg />
    <VSpacing size=Spacing.xl />
    <InstructionCard title="1. Enter Pin Code" url=Images.ledgerStep1 />
    <VSpacing size=Spacing.lg />
    {switch (ledgerApp) {
     | Ledger.Cosmos => <InstructionCard title="2. Open Cosmos" url=Images.ledgerStep2Cosmos />
     | BandChain => <InstructionCard title="2. Open BandChain" url=Images.ledgerStep2BandChain />
     }}
    <div className=Styles.resultContainer>
      {switch (result) {
       | Loading =>
         <>
           <Text
             value="Please accept with ledger"
             color=Colors.bandBlue
             spacing={Text.Em(0.03)}
             weight=Text.Medium
           />
           <Loading width={`px(100)} />
         </>
       | Error(err) =>
         <Text
           value=err
           color=Colors.red5
           weight=Text.Medium
           size=Text.Lg
           spacing={Text.Em(0.03)}
         />
       | Nothing => React.null
       }}
    </div>
    {result == Loading
       ? <div className={Styles.connectBtn(~isLoading=true, ())}>
           <div className=Styles.connectingBtnContainer>
             <Icon name="fad fa-spinner-third fa-spin" size=16 />
             <Text value="Connecting..." weight=Text.Bold size=Text.Md color=Colors.blueGray7 />
           </div>
         </div>
       : <Button
           style={Styles.connectBtn(~isLoading=false, ())}
           onClick={_ => {
             switch (Os.isWindows(), Os.checkHID()) {
             | (true, false) =>
               let isConfirm =
                 Webapi.Dom.(
                   window
                   |> Window.confirm(
                        {j|To use Ledger Nano on Windows 10, please enable "Experimental Web Platform Features" by copy-paste "chrome://flags/#enable-experimental-web-platform-features". Click OK to copy.|j},
                      )
                 );
               isConfirm
                 ? Copy.copy("chrome://flags/#enable-experimental-web-platform-features") : ();
             | (_, _) => createLedger(accountIndex)
             }
           }}>
           <Text value="Connect to Ledger" weight=Text.Bold size=Text.Lg color=Colors.white />
         </Button>}
  </div>;
};
