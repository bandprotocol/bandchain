module Styles = {
  open Css;

  let table =
    style([
      display(`flex),
      borderRadius(`px(3)),
      border(`px(1), `solid, `hex("dadada")),
      backgroundColor(`hex("ffffff")),
      flexDirection(`column),
      width(`percent(100.)),
    ]);

  let tableHead =
    style([
      display(`flex),
      padding2(~v=px(12), ~h=`px(30)),
      fontSize(`em(0.8)),
      color(Css_Colors.black),
      fontWeight(`num(600)),
      borderBottom(`px(1), `solid, `hex("dadada")),
    ]);

  let tableBody = style([display(`flex), width(`percent(100.))]);
};

[@react.component]
let make = (~header, ~body) => {
  <div className=Styles.table>
    <div className=Styles.tableHead> header </div>
    <div className=Styles.tableBody> body </div>
  </div>;
};
