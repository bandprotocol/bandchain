module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      background(Css_Colors.white),
      justifyContent(`center),
      height(`px(120)),
      boxShadow(Css.Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(6), `rgba((0, 0, 0, 0.05)))),
    ]);

  let innerContainer =
    style([
      display(`flex),
      flexDirection(`column),
      width(`percent(100.)),
      padding2(~v=`em(1.4), ~h=`zero),
      justifyContent(`spaceBetween),
    ]);

  let flex = style([display(`flex)]);
  let flexRow = style([display(`flex), flexDirection(`row), justifyContent(`spaceBetween)]);
  let flexRowInner = style([display(`flex), flexDirection(`row), justifyContent(`flexEnd)]);

  let testnetLabel =
    style([
      display(`flex),
      backgroundColor(`hex("ececec")),
      color(`hex("505050")),
      paddingLeft(`px(7)),
      paddingRight(`px(7)),
      fontSize(`px(13)),
      height(`px(25)),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(3)),
    ]);

  let explorerLabel = style([fontSize(`px(18)), color(`hex("5b5b5b"))]);
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <PageContainer>
      <div className=Styles.innerContainer>
        <div className=Styles.flexRow>
          <div className=Styles.flexRowInner>
            <div className=Styles.testnetLabel> {"TESTNET v.10" |> React.string} </div>
          </div>
        </div>
        <div className=Styles.flexRow>
          <div className=Styles.explorerLabel> {"DATA REQUEST EXPLORER" |> React.string} </div>
          <div className=Styles.flex> <NavLink text="Validators" to_="/validators" /> </div>
          <div className=Styles.flex> <NavLink text="Blocks" to_="/blocks" /> </div>
          <div className=Styles.flex> <NavLink text="Transactions" to_="/transactions" /> </div>
          <div className=Styles.flex> <NavLink text="Request Scripts" to_="/scripts" /> </div>
          <div className=Styles.flex> <NavLink text="Data Providers" to_="/providers" /> </div>
          <div className=Styles.flex> <NavLink text="OWASM Studio" to_="/studio" /> </div>
        </div>
      </div>
    </PageContainer>
  </div>;
};
