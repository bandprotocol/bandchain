module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      width(`percent(100.)),
      alignItems(`center),
      flexDirection(`column),
      paddingTop(`px(20)),
      paddingBottom(`px(20)),
    ]);

  let tableContainer = style([display(`flex), marginTop(`px(25)), width(`percent(100.))]);

  let flex = style([display(`flex), flex(`num(1.))]);

  let separator = style([display(`flex), width(`px(20))]);
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <Summary />
    <div className=Styles.tableContainer>
      <div className=Styles.flex> <LatestBlocks /> </div>
      <div className=Styles.separator />
      <div className=Styles.flex> <LatestTxs /> </div>
    </div>
  </div>;
};
