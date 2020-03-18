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
      backgroundColor(Colors.gray7),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
};

[@react.component]
let make = () => {
  let (page, setPage) = React.useState(_ => 1);
  let limit = 10;

  let txsOpt = TxHook.latest(~limit, ~page, ());
  let txs = txsOpt->Belt.Option.mapWithDefault([], ({txs}) => txs);

  let totalTxsOpt = txsOpt->Belt.Option.map(info => info.totalCount);
  // TODO: add loading state later.
  let pageCount = txsOpt->Belt.Option.mapWithDefault(1, info => info.pageCount);

  <div className=Styles.pageContainer>
    <Row>
      <Col> <img src=Images.txLogo className=Styles.logo /> </Col>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL TRANSACTIONS"
            weight=Text.Semibold
            color=Colors.gray7
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
    <VSpacing size=Spacing.lg />
    <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
    <VSpacing size=Spacing.lg />
  </div>;
};
