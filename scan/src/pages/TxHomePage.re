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

  let txsSub = TxSub.getList(~pageSize, ~page, ());
  let txsCountSub = TxSub.count();

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
          {switch (txsCountSub) {
           | Data(txsCount) =>
             <>
               <div className=Styles.seperatedLine />
               <Text
                 value={txsCount->Format.iPretty ++ " in total"}
                 size=Text.Md
                 weight=Text.Thin
                 spacing={Text.Em(0.06)}
                 color=Colors.gray7
                 nowrap=true
               />
             </>
           | _ => React.null
           }}
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
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
  </>;
};
