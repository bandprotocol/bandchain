module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      justifyContent(`center),
      position(`relative),
      width(`px(640)),
      height(`px(480)),
    ]);

  let bg =
    style([
      position(`absolute),
      width(`percent(100.)),
      height(`percent(100.)),
      backgroundColor(Css.rgb(249, 249, 251)),
      backgroundImage(`url(Images.modalBg)),
      backgroundRepeat(`noRepeat),
      borderRadius(`px(8)),
      zIndex(-1),
    ]);

  let innerContainer = style([display(`flex), flexDirection(`column), width(`percent(100.))]);

  let loginSelectionContainer = style([margin2(~v=`zero, ~h=`px(24))]);

  let modalTitle =
    style([
      display(`flex),
      justifyContent(`center),
      flexDirection(`column),
      alignItems(`center),
    ]);

  let warning =
    style([
      display(`flex),
      flexDirection(`row),
      padding2(~v=`px(5), ~h=`px(8)),
      color(Colors.yellow6),
      backgroundColor(Colors.yellow1),
      border(`px(1), `solid, Colors.yellow6),
      borderRadius(`px(4)),
    ]);

  let header = active =>
    style([
      display(`flex),
      flexDirection(`row),
      alignSelf(`center),
      alignItems(`center),
      justifyContent(`spaceBetween),
      width(`px(480)),
      padding2(~v=`zero, ~h=`px(20)),
      color(active ? Colors.gray8 : Colors.gray6),
      backgroundColor(Colors.white),
    ]);

  let loginList = active =>
    style([
      display(`flex),
      width(`px(226)),
      height(`px(50)),
      borderRadius(`px(8)),
      backgroundColor(Colors.white),
      boxShadow(
        active
          ? Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(11, 29, 142, `num(0.1)))
          : Shadow.box(~x=`zero, ~y=`px(0), ~blur=`px(0), Css.rgba(0, 0, 0, `num(0.))),
      ),
      cursor(`pointer),
      overflow(`hidden),
    ]);

  let seperatedLongLine =
    style([height(`px(275)), width(`px(2)), backgroundColor(Colors.gray4)]);

  let ledgerIcon = style([height(`px(40)), width(`px(40))]);
  let ledgerImageContainer = active => style([opacity(active ? 1.0 : 0.5)]);

  let activeBar = active =>
    style([backgroundColor(active ? Colors.bandBlue : Colors.white), width(`px(15))]);
};

type login_method_t =
  | Mnemonic
  | LedgerWithCosmos
  | LedgerWithBandChain;

let toLoginMethodString = method => {
  switch (method) {
  | Mnemonic => "Mnemonic Phrase"
  | LedgerWithCosmos => "Ledger - Cosmos"
  | LedgerWithBandChain => "Ledger - Band (beta)"
  };
};

module LoginMethod = {
  [@react.component]
  let make = (~name, ~active, ~onClick) => {
    <div className={Styles.loginList(active)} onClick>
      <div className={Styles.activeBar(active)} />
      <div className={Styles.header(active)}>
        <Text value={name |> toLoginMethodString} weight=Text.Medium size=Text.Md />
        {switch (name) {
         | LedgerWithCosmos =>
           <div className={Styles.ledgerImageContainer(active)}>
             <img src=Images.ledgerCosmosIcon className=Styles.ledgerIcon />
           </div>
         | LedgerWithBandChain =>
           <div className={Styles.ledgerImageContainer(active)}>
             <img src=Images.ledgerBandChainIcon className=Styles.ledgerIcon />
           </div>
         | _ => <div />
         }}
      </div>
    </div>;
  };
};

[@react.component]
let make = (~chainID) => {
  let (loginMethod, setLoginMethod) = React.useState(_ => Mnemonic);
  <div className=Styles.container>
    // <div className=Styles.bg />

      <div className=Styles.innerContainer>
        <VSpacing size=Spacing.xxl />
        <div className=Styles.modalTitle>
          <Text value="Connect with your wallet" weight=Text.Medium size=Text.Xl />
          {chainID == "band-wenchang-mainnet"
             ? <>
                 <VSpacing size=Spacing.lg />
                 <div className=Styles.warning>
                   <Text value="Please check that you are visiting" />
                   <HSpacing size=Spacing.sm />
                   <Text value="https://www.cosmoscan.io" weight=Text.Bold />
                 </div>
               </>
             : <VSpacing size=Spacing.xxl />}
        </div>
        <VSpacing size=Spacing.xl />
        <Row alignItems=`flexStart>
          <Col>
            <div className=Styles.loginSelectionContainer>
              <Text
                value="Select your connection method"
                size=Text.Lg
                weight=Text.Medium
                color=Colors.gray7
              />
              {[|Mnemonic, LedgerWithCosmos, LedgerWithBandChain|]
               ->Belt_Array.map(method =>
                   <React.Fragment key={method |> toLoginMethodString}>
                     <VSpacing size=Spacing.lg />
                     <LoginMethod
                       name=method
                       active={loginMethod == method}
                       onClick={_ => setLoginMethod(_ => method)}
                     />
                   </React.Fragment>
                 )
               ->React.array}
            </div>
          </Col>
          <Col> <div className=Styles.seperatedLongLine /> </Col>
          <Col size=1.>
            {switch (loginMethod) {
             | Mnemonic => <ConnectWithMnemonic chainID />
             | LedgerWithCosmos => <ConnectWithLedger chainID ledgerApp=Ledger.Cosmos />
             | LedgerWithBandChain => <ConnectWithLedger chainID ledgerApp=Ledger.BandChain />
             }}
          </Col>
        </Row>
      </div>
    </div>;
};
