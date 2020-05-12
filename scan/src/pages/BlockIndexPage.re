module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let header =
    style([display(`flex), flexDirection(`row), alignItems(`center), height(`px(50))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let addressContainer = style([marginTop(`px(15))]);
  let checkLogo = style([marginRight(`px(10))]);
  let blockLogo = style([width(`px(50)), marginRight(`px(10))]);
  let proposerContainer = style([maxWidth(`px(180))]);
};

[@react.component]
let make = (~height) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let blockSub = BlockSub.get(height);
  let txsSub = TxSub.getListByBlockHeight(height, ~pageSize, ~page, ());

  <>
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
             <Text value={height |> ID.Block.toString} weight=Text.Thin spacing={Text.Em(0.06)} />
           </>
         | _ => React.null
         }}
      </div>
    </Row>
    <VSpacing size=Spacing.lg />
    {switch (blockSub) {
     | Data({hash}) =>
       <Text
         value={hash |> Hash.toHex(~upper=true)}
         size=Text.Xxl
         weight=Text.Semibold
         code=true
         nowrap=true
         ellipsis=true
       />
     | _ => <LoadingCensorBar width=700 height=15 />
     }}
    <VSpacing size=Spacing.lg />
    <Row>
      <Col size=1.8>
        {switch (blockSub) {
         | Data({txn}) => <InfoHL info={InfoHL.Count(txn)} header="TRANSACTIONS" />
         | _ => <InfoHL info={InfoHL.Loading(75)} header="TRANSACTIONS" />
         }}
      </Col>
      <Col size=4.6>
        {switch (blockSub) {
         | Data({timestamp}) => <InfoHL info={InfoHL.Timestamp(timestamp)} header="TIMESTAMP" />
         | _ => <InfoHL info={InfoHL.Loading(370)} header="TIMESTAMP" />
         }}
      </Col>
      <Col size=3.2>
        {switch (blockSub) {
         | Data({validator}) =>
           <InfoHL info={InfoHL.Text(validator.moniker)} header="PROPOSED BY" />
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
         <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
         <VSpacing size=Spacing.lg />
       </>;
     | _ => React.null
     }}
  </>;
};
