type pos =
  | Msg
  | TxIndex
  | Fee;

module Styles = {
  open Css;

  let container = style([display(`flex), alignItems(`center)]);
};

[@react.component]
let make = (~coins, ~pos=Msg) => {
  <div className=Styles.container>
    <Text
      value={coins |> Coin.getBandAmountFromCoins |> Format.fPretty}
      weight=Text.Semibold
      code=true
      block=true
      nowrap=true
    />
    <HSpacing size=Spacing.sm />
    {switch (pos) {
     | Msg => <Text value="BAND" weight=Text.Regular code=true nowrap=true block=true />
     | TxIndex => <Text value="BAND" weight=Text.Thin code=true nowrap=true block=true />
     | Fee => React.null
     }}
  </div>;
};
