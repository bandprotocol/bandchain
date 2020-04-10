module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);
};

[@react.component]
let make = (~accountAddress: Address.t) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  {
    let txsSub = TxSub.getListBySender(accountAddress, ~pageSize, ~page, ());
    let txsCountSub = TxSub.countBySender(accountAddress);

    let%Sub txs = txsSub;
    let%Sub txsCount = txsCountSub;

    let pageCount = Page.getPageCount(txsCount, pageSize);

    <div className=Styles.tableLowerContainer>
      <VSpacing size=Spacing.md />
      <div className=Styles.hFlex>
        <HSpacing size=Spacing.lg />
        <Text value={txsCount |> string_of_int} weight=Text.Semibold />
        <HSpacing size=Spacing.xs />
        <Text value="Transactions In Total" />
      </div>
      <VSpacing size=Spacing.lg />
      <TxsTable txs />
      <VSpacing size=Spacing.lg />
      <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
      <VSpacing size=Spacing.lg />
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
};
