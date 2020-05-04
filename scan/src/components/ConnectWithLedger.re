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
      height(`px(50)),
      width(`percent(100.)),
      justifyContent(`spaceBetween),
      backgroundColor(Css.rgb(255, 255, 255)),
      border(`px(1), `solid, Colors.blueGray3),
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
      height(`px(35)),
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

  let selectWrapper =
    style([
      display(`flex),
      padding2(~v=`px(3), ~h=`px(8)),
      justifyContent(`center),
      alignItems(`center),
      width(`percent(100.)),
      height(`px(30)),
      left(`zero),
      top(`px(32)),
      background(rgba(255, 255, 255, 1.)),
      border(`px(1), `solid, Colors.blueGray3),
      borderRadius(`px(6)),
      float(`left),
    ]);

  let selectContent =
    style([
      background(rgba(255, 255, 255, 1.)),
      border(`px(0), `solid, hex("FFFFFF")),
      width(`px(100)),
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
  let (deviation, setDeviation) = React.useState(_ => 0);

  let createLedger = deviation => {
    setResult(_ => Loading);
    let _ =
      Wallet.createFromLedger(deviation)
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
    <Text value="1. Select HD Derivation Path" weight=Text.Semibold />
    <VSpacing size=Spacing.sm />
    <div className=Styles.selectWrapper>
      <select
        className=Styles.selectContent
        onChange={event => {
          let newDeviation = ReactEvent.Form.target(event)##value |> int_of_string;
          setDeviation(_ => newDeviation);
        }}>
        {[|0, 1, 2, 3, 4, 5|]
         |> Belt.Array.map(_, deviation =>
              <option value={deviation |> string_of_int}>
                {{
                   "44/118/0/0/" ++ (deviation |> string_of_int);
                 }
                 |> React.string}
              </option>
            )
         |> React.array}
      </select>
    </div>
    <VSpacing size=Spacing.sm />
    <Text value="2. On Your Ledger" weight=Text.Semibold />
    <VSpacing size=Spacing.sm />
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
       | Error(err) =>
         <Text value=err color=Colors.red5 weight=Text.Medium spacing={Text.Em(0.03)} />
       | Nothing => React.null
       }}
    </div>
    <div className=Styles.connectBtn onClick={_ => createLedger(deviation)}>
      <Text value="Connect To Ledger" weight=Text.Bold size=Text.Md color=Colors.white />
    </div>
  </div>;
};
