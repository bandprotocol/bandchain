module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let header =
    style([display(`flex), flexDirection(`row), alignItems(`center), height(`px(50))]);

  let blockHash = style([height(`px(18)), display(`flex), alignItems(`center)]);

  let blockLogo = style([minWidth(`px(50)), marginRight(`px(10))]);

  let infoContainerFullwidth =
    style([
      Media.mobile([
        selector("> div", [flexBasis(`percent(100.))]),
        selector("> div + div", [marginTop(`px(15))]),
        selector("> div > div > div", [display(`block)]),
      ]),
    ]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let proposerContainer = style([width(`px(300))]);
};

[@react.component]
let make = (~height) => {
  let isMobile = Media.isMobile();
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let blockSub = BlockSub.get(height);
  let txsSub = TxSub.getListByBlockHeight(height, ~pageSize, ~page, ());

  <Section>
    <div className=CssHelper.container>
      <Row justify=Row.Between>
        <div className=Styles.header>
          <img src=Images.blockLogo className=Styles.blockLogo />
          <Text
            value="BLOCK"
            weight=Text.Medium
            size=Text.Md
            nowrap=true
            color=Colors.gray7
            block=true
            spacing={Text.Em(0.06)}
          />
          {switch (blockSub) {
           | Data({height}) =>
             <>
               <div className=Styles.seperatedLine />
               <Text
                 value={height |> ID.Block.toString}
                 weight=Text.Thin
                 spacing={Text.Em(0.06)}
               />
             </>
           | _ => React.null
           }}
        </div>
      </Row>
      <VSpacing size=Spacing.lg />
      <div className=Styles.blockHash>
        {switch (blockSub) {
         | Data({hash}) =>
           isMobile
             ? <Text
                 value={hash |> Hash.toHex(~upper=true)}
                 size=Text.Lg
                 weight=Text.Bold
                 nowrap=false
                 breakAll=true
                 code=true
                 color=Colors.gray7
               />
             : <Text
                 value={hash |> Hash.toHex(~upper=true)}
                 size=Text.Xxl
                 nowrap=true
                 ellipsis=true
                 code=true
                 weight=Text.Bold
               />
         | _ => <LoadingCensorBar width=700 height=15 />
         }}
      </div>
      <VSpacing size=Spacing.lg />
      <Row wrap=true style=Styles.infoContainerFullwidth>
        <Col size=1.8>
          {switch (blockSub) {
           | Data({txn}) => <InfoHL info={InfoHL.Count(txn)} header="TRANSACTIONS" />
           | _ => <InfoHL info={InfoHL.Loading(75)} header="TRANSACTIONS" />
           }}
        </Col>
        <Col size=4.6>
          {switch (blockSub) {
           | Data({timestamp}) =>
             <InfoHL info={InfoHL.Timestamp(timestamp)} header="TIMESTAMP" />
           | _ => <InfoHL info={InfoHL.Loading(isMobile ? 240 : 370)} header="TIMESTAMP" />
           }}
        </Col>
        <Col size=3.2>
          {switch (blockSub) {
           | Data({validator}) =>
             <div className=Styles.proposerContainer>
               <InfoHL
                 info={
                   InfoHL.Validator(
                     validator.operatorAddress,
                     validator.moniker,
                     validator.identity,
                   )
                 }
                 header="PROPOSED BY"
               />
             </div>
           | _ => <InfoHL info={InfoHL.Loading(80)} header="PROPOSED BY" />
           }}
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <BlockIndexTxsTable txsSub />
      {switch (blockSub) {
       | Data({txn}) =>
         let pageCount = Page.getPageCount(txn, pageSize);
         <>
           <VSpacing size=Spacing.lg />
           <Pagination
             currentPage=page
             pageCount
             onPageChange={newPage => setPage(_ => newPage)}
           />
           <VSpacing size=Spacing.lg />
         </>;
       | _ => React.null
       }}
    </div>
  </Section>;
};
