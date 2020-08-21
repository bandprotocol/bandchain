module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let header =
    style([display(`flex), flexDirection(`row), alignItems(`center), height(`px(50))]);

  let logo = style([minWidth(`px(50)), marginRight(`px(10))]);

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

  <Section>
    <div className=CssHelper.container>
      <Row>
        <div className=Styles.header>
          <img src=Images.txLogo className=Styles.logo />
          <Text
            value="ALL TRANSACTIONS"
            weight=Text.Medium
            size=Text.Md
            spacing={Text.Em(0.06)}
            height={Text.Px(15)}
            nowrap=true
            block=true
            color=Colors.gray7
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
      </Row>
      <VSpacing size=Spacing.xl />
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
