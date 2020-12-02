[@react.component]
let make = () => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let txsSub = TxSub.getList(~pageSize, ~page, ());
  let latestTxsSub = TxSub.getList(~pageSize=1, ~page=1, ());
  let txsCountSub = TxSub.count();

  let isMobile = Media.isMobile();

  <Section>
    <div className=CssHelper.container id="transactionsSection">
      <Row alignItems=Row.Center marginBottom=40 marginBottomSm=24>
        <Col.Grid col=Col.Twelve>
          <Heading value="All Transactions" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
          {switch (latestTxsSub) {
           | Data(txs) =>
             <Heading
               value={
                 txs
                 ->Belt.Array.get(0)
                 ->Belt.Option.mapWithDefault(0, ({id}) => id)
                 ->Format.iPretty
                 ++ " In total"
               }
               size=Heading.H3
             />
           | _ => <LoadingCensorBar width=65 height=21 />
           }}
        </Col.Grid>
      </Row>
      {isMobile
         ? React.null
         : <THead>
             <Row alignItems=Row.Center>
               <Col.Grid col=Col.Two>
                 <Text block=true value="TX Hash" weight=Text.Semibold color=Colors.gray7 />
               </Col.Grid>
               <Col.Grid col=Col.One>
                 <Text block=true value="Block" weight=Text.Semibold color=Colors.gray7 />
               </Col.Grid>
               <Col.Grid col=Col.One>
                 <Text
                   block=true
                   value="Status"
                   size=Text.Md
                   weight=Text.Semibold
                   color=Colors.gray7
                   align=Text.Center
                 />
               </Col.Grid>
               <Col.Grid col=Col.Two>
                 <Text
                   block=true
                   value="Gas Fee (BAND)"
                   weight=Text.Semibold
                   color=Colors.gray7
                   align=Text.Center
                 />
               </Col.Grid>
               <Col.Grid col=Col.Six>
                 <Text block=true value="Actions" weight=Text.Semibold color=Colors.gray7 />
               </Col.Grid>
             </Row>
           </THead>}
      <TxsTable txsSub />
      {switch (txsCountSub) {
       | Data(txsCount) =>
         let pageCount = Page.getPageCount(txsCount, pageSize);

         <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />;
       | _ => React.null
       }}
    </div>
  </Section>;
};
