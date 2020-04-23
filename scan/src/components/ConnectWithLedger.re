module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      flexDirection(`column),
      width(`percent(100.)),
      padding4(~top=`px(18), ~left=`px(18), ~right=`px(28), ~bottom=`px(0)),
    ]);

  let instructionCard =
    style([
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      width(`percent(100.)),
      justifyContent(`spaceBetween),
      backgroundColor(Css.rgba(197, 199, 211, 0.3)),
      borderRadius(`px(8)),
      padding2(~v=`px(25), ~h=`px(18)),
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
      height(`px(60)),
    ]);

  let ledgerGuide = style([width(`px(106)), height(`px(25))]);

  let loading = style([width(`px(100))]);

  let connectBtn =
    style([
      width(`px(140)),
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
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(82, 105, 255, 0.25))),
      borderRadius(`px(4)),
      cursor(`pointer),
      alignSelf(`flexEnd),
    ]);
};

module InstructionCard = {
  [@react.component]
  let make = (~idx, ~title, ~url) => {
    <div className=Styles.instructionCard>
      <div className=Styles.rFlex>
        <div className=Styles.oval>
          <Text
            value={idx |> string_of_int}
            color=Colors.white
            size=Text.Xxl
            spacing={Text.Em(0.03)}
            weight=Text.Bold
          />
        </div>
        <HSpacing size=Spacing.md />
        <Text value=title weight=Text.Semibold spacing={Text.Em(0.03)} />
      </div>
      <img src=url className=Styles.ledgerGuide />
    </div>;
  };
};

type result_t =
  | Nothing
  | Loading
  | Error(string);

[@react.component]
let make = () => {
  let (_, dispatchAccount) = React.useContext(AccountContext.context);
  let (_, dispatchModal) = React.useContext(ModalContext.context);
  let (result, setResult) = React.useState(_ => Nothing);

  let createLedger = () => {
    setResult(_ => Loading);
    let _ =
      Wallet.createFromLedger()
      |> Js.Promise.then_(wallet => {
           let%Promise (address, pubKey) = wallet->Wallet.getAddressAndPubKey;
           dispatchAccount(Connect(wallet, address, pubKey));
           dispatchModal(CloseModal);
           Promise.ret();
         })
      |> Js.Promise.catch(err => {
           Js.Console.log(err);
           setResult(_ => Error("An error occured"));
           Promise.ret();
         });
    ();
  };

  <div className=Styles.container>
    <InstructionCard idx=1 title="Enter Pin Code" url=Images.ledgerStep1 />
    <VSpacing size=Spacing.md />
    <InstructionCard idx=2 title="Open Cosmos" url=Images.ledgerStep2 />
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
           <img src=Images.loadingCircles className=Styles.loading />
         </>
       | Error(err) => <Text value=err color=Colors.red6 />
       | Nothing => React.null
       }}
    </div>
    <div className=Styles.connectBtn onClick={_ => createLedger()}>
      <Text value="Connect To Ledger" weight=Text.Bold size=Text.Md color=Colors.white />
    </div>
  </div>;
};
