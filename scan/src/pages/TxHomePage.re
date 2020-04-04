module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);
};

[@react.component]
let make = () => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  {
    let txsSub = TxSub.getList(~pageSize, ~page, ());
    let numTotalTxsSub = TxSub.count();

    let%Sub txs = txsSub;
    let%Sub numTotalTxs = numTotalTxsSub;

    // TODO: add loading state later.
    let pageCount = Page.getPageCount(numTotalTxs, pageSize);

    <>
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
            <Text value={(numTotalTxs |> Format.iPretty) ++ " in total"} />
          </div>
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <TxsTable txs />
      <VSpacing size=Spacing.lg />
      <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
      <VSpacing size=Spacing.lg />
    </>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
};
