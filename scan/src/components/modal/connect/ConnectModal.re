module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      justifyContent(`center),
      position(`relative),
      width(`px(800)),
      height(`px(520)),
    ]);

  let innerContainer = style([display(`flex), flexDirection(`column), width(`percent(100.))]);

  let loginSelectionContainer =
    style([padding2(~v=`zero, ~h=`px(24)), height(`percent(100.))]);

  let modalTitle =
    style([
      display(`flex),
      justifyContent(`center),
      flexDirection(`column),
      alignItems(`center),
      paddingTop(`px(30)),
      borderBottom(`px(1), `solid, Colors.gray9),
    ]);

  let row = style([height(`percent(100.))]);
  let rowContainer = style([margin2(~v=`zero, ~h=`px(12)), height(`percent(100.))]);
  let warning = style([display(`flex), flexDirection(`row)]);

  let header = active =>
    style([
      display(`flex),
      flexDirection(`row),
      alignSelf(`center),
      alignItems(`center),
      padding2(~v=`zero, ~h=`px(20)),
      color(active ? Colors.gray8 : Colors.gray6),
      backgroundColor(Colors.white),
    ]);

  let loginList = active =>
    style([
      display(`flex),
      width(`percent(100.)),
      height(`px(50)),
      borderRadius(`px(4)),
      border(`px(1), `solid, active ? Colors.bandBlue : Colors.white),
      backgroundColor(Colors.white),
      cursor(`pointer),
      overflow(`hidden),
    ]);

  let loginSelectionBackground = style([background(Colors.profileBG)]);

  let ledgerIcon = style([height(`px(28)), width(`px(28)), transform(translateY(`px(3)))]);
  let ledgerImageContainer = active =>
    style([opacity(active ? 1.0 : 0.5), marginRight(`px(15))]);
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
      <div className={Styles.header(active)}>
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
        <Text value={name |> toLoginMethodString} weight=Text.Medium size=Text.Lg />
      </div>
    </div>;
  };
};

[@react.component]
let make = (~chainID) => {
  let (loginMethod, setLoginMethod) = React.useState(_ => Mnemonic);
  <div className=Styles.container>
    <div className=Styles.innerContainer>
      <div className=Styles.modalTitle>
        <Text value="Connect with your wallet" weight=Text.Medium size=Text.Xl />
        {chainID == "band-guanyu-mainnet"
           ? <>
               <VSpacing size=Spacing.md />
               <div className=Styles.warning>
                 <Text value="Please check that you are visiting" size=Text.Lg weight=Text.Thin />
                 <HSpacing size=Spacing.sm />
                 <Text
                   value="https://www.cosmoscan.io"
                   size=Text.Lg
                   weight=Text.Medium
                   color=Colors.bandBlue
                 />
               </div>
             </>
           : <VSpacing size=Spacing.sm />}
        <VSpacing size=Spacing.xl />
      </div>
      <div className=Styles.rowContainer>
        <Row.Grid style=Styles.row>
          <Col.Grid col=Col.Five style=Styles.loginSelectionBackground>
            <div className=Styles.loginSelectionContainer>
              <VSpacing size=Spacing.xl />
              <Text
                value="Select your connection method"
                size=Text.Lg
                weight=Text.Thin
                color=Colors.gray7
              />
              <VSpacing size=Spacing.md />
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
          </Col.Grid>
          <Col.Grid col=Col.Seven>
            {switch (loginMethod) {
             | Mnemonic => <ConnectWithMnemonic chainID />
             | LedgerWithCosmos => <ConnectWithLedger chainID ledgerApp=Ledger.Cosmos />
             | LedgerWithBandChain => <ConnectWithLedger chainID ledgerApp=Ledger.BandChain />
             }}
          </Col.Grid>
        </Row.Grid>
      </div>
    </div>
  </div>;
};
