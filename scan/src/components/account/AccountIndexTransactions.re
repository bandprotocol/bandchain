module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);
};

[@react.component]
let make = (~accountAddress: Address.t) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let txsSub = TxSub.getListBySender(accountAddress, ~pageSize, ~page, ());
  let txsCountSub = TxSub.countBySender(accountAddress);

  <div className=Styles.tableLowerContainer>
    <VSpacing size=Spacing.md />
    {switch (txsCountSub) {
     | Data(txsCount) =>
       <div className=Styles.hFlex>
         <HSpacing size=Spacing.lg />
         <Text value={txsCount |> string_of_int} weight=Text.Semibold />
         <HSpacing size=Spacing.xs />
         <Text value="Transactions In Total" />
       </div>
     | _ =>
       <div className=Styles.hFlex>
         <HSpacing size=Spacing.lg />
         <LoadingCensorBar width=130 height=15 />
       </div>
     }}
    <VSpacing size=Spacing.lg />
    <TxsTable txsSub />
    {switch (txsCountSub) {
     | Data(txsCount) =>
       let pageCount = Page.getPageCount(txsCount, pageSize);
       <>
         <VSpacing size=Spacing.lg />
         <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
         <VSpacing size=Spacing.lg />
       </>;
     | _ => React.null
     }}
  </div>;
};
