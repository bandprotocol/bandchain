module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50)), minHeight(`px(500))]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.mediumGray),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
};

[@react.component]
let make = () => {
  let step = 10;
  let (limit, setLimit) = React.useState(_ => step);
  let txsOpt = TxHook.latest(~limit, ());
  let txs = txsOpt->Belt.Option.mapWithDefault([], ({txs}) => txs);

  let infoOpt = React.useContext(GlobalContext.context);
  let totalTxsOpt = infoOpt->Belt.Option.map(info => info.latestBlock.totalTxs);

  <div className=Styles.pageContainer>
    <Row>
      <Col> <img src=Images.txLogo className=Styles.logo /> </Col>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL TRANSACTIONS"
            weight=Text.Semibold
            color=Colors.mediumGray
            nowrap=true
            spacing={Text.Em(0.06)}
          />
          <div className=Styles.seperatedLine />
          {switch (totalTxsOpt) {
           | Some(totalTxs) => <Text value={(totalTxs |> Format.iPretty) ++ " in total"} />
           | None => React.null
           }}
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <TxsTable txs />
    <VSpacing size={`px(70)} />
  </div>;
};
