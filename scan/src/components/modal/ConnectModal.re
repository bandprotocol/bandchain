module Styles = {
  open Css;

  let loginContainer =
    style([
      display(`flex),
      justifyContent(`center),
      width(`percent(100.)),
      height(`percent(100.)),
    ]);

  let modalLogin =
    style([
      width(`px(640)),
      height(`px(480)),
      backgroundColor(Css.rgb(249, 249, 251)),
      backgroundImage(`url(Images.modalBg)),
      borderRadius(`px(8)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(8), ~blur=`px(32), Css.rgba(0, 0, 0, 0.5))),
    ]);

  let modalTitle = style([display(`flex), justifyContent(`center)]);

  let modalSelectText = style([display(`flex), marginLeft(`px(34))]);

  let itemCol =
    style([height(`px(275)), display(`flex), flexDirection(`column), verticalAlign(`top)]);

  let container = style([backgroundColor(Colors.transparent)]);

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

  let buttonContainer = active =>
    style([
      display(`flex),
      width(`px(226)),
      height(`px(50)),
      marginLeft(`px(34)),
      borderRadius(`px(8)),
      backgroundColor(Colors.white),
      boxShadow(
        active
          ? Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(11, 29, 142, 0.1))
          : Shadow.box(~x=`zero, ~y=`px(0), ~blur=`px(0), Css.rgba(0, 0, 0, 0.)),
      ),
      cursor(`pointer),
      overflow(`hidden),
    ]);
  let seperatedLongLine =
    style([
      height(`px(275)),
      width(`px(2)),
      backgroundColor(Colors.gray4),
      marginLeft(`px(24)),
    ]);

  let ledgerIcon = style([height(`px(30)), width(`px(30)), display(`flex)]);
  let ledgerImageContainer = active => style([opacity(active ? 1.0 : 0.5)]);

  let activeBar = active =>
    style([backgroundColor(active ? Colors.bandBlue : Colors.white), width(`px(15))]);
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
  let make = (~name, ~active, ~onClick) => {
    <div className={Styles.buttonContainer(active)} onClick>
      <div className={Styles.activeBar(active)} />
      <div className={Styles.header(active)}>
        <Text value={name |> toLoginMethodString} weight=Text.Medium size=Text.Md />
        {switch (name) {
         | Ledger =>
           <div className={Styles.ledgerImageContainer(active)}>
             <img src=Images.ledgerIconActive className=Styles.ledgerIcon />
           </div>
         | _ => <div />
         }}
      </div>
    </div>;
  };
};

[@react.component]
let make = _ => {
  let (loginMethod, setLoginMethod) = React.useState(_ => Mnemonic);
  <div className=Styles.loginContainer>
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
                 <LoginMethod
                   name=method
                   active={loginMethod == method}
                   onClick={_ => setLoginMethod(_ => method)}
                 />
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
           | Mnemonic => <ConnectMnemonic />
           | Ledger => "Ledger" |> React.string
           }}
        </div>
      </Row>
    </div>
  </div>;
};
