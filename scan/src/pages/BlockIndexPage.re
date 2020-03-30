module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(50))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let txTypeBadge =
    style([
      paddingLeft(`px(12)),
      paddingRight(`px(12)),
      paddingTop(`px(5)),
      paddingBottom(`px(5)),
      backgroundColor(Colors.blue1),
      borderRadius(`px(15)),
    ]);

  let msgAmount =
    style([borderRadius(`percent(50.)), padding(`px(3)), backgroundColor(Colors.gray4)]);

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
  let noTransactionLogo = style([width(`px(160))]);
  let emptyContainer =
    style([
      height(`px(300)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
    ]);
  let proposerContainer = style([maxWidth(`px(180))]);

  let seperatorLine =
    style([
      width(`percent(100.)),
      height(`pxFloat(1.4)),
      backgroundColor(Colors.gray4),
      display(`flex),
    ]);
};

[@react.component]
let make = (~height) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  {
    let blockSub = BlockSub.get(height);
    let txsSub = TxSub.getListByBlockHeight(height, ~pageSize, ~page, ());

<<<<<<< HEAD
    let%Sub block = blockSub;
    let%Sub txs = txsSub;

    let pageCount = Page.getPageCount(block.txn, pageSize);

    <div className=Styles.pageContainer>
      <Row justify=Row.Between>
        <Col>
          <div className=Styles.vFlex>
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
            <div className=Styles.seperatedLine />
            <Text value={height |> ID.Block.toString} weight=Text.Thin spacing={Text.Em(0.06)} />
          </div>
        </Col>
      </Row>
      <VSpacing size=Spacing.lg />
      <div className=Styles.vFlex> <HSpacing size=Spacing.xs /> </div>
      <Text
        value={block.hash |> Hash.toHex(~upper=true)}
        size=Text.Xxl
        weight=Text.Semibold
        code=true
        nowrap=true
        ellipsis=true
      />
      <VSpacing size=Spacing.lg />
      <Row>
        <Col size=1.8> <InfoHL info={InfoHL.Count(block.txn)} header="TRANSACTIONS" /> </Col>
        <Col size=4.6>
          <div className=Styles.vFlex>
            <InfoHL info={InfoHL.Timestamp(block.timestamp)} header="TIME STAMP" />
          </div>
        </Col>
        <Col size=3.2>
          <div className=Styles.proposerContainer>
            <InfoHL info={InfoHL.Text(block.validator.moniker)} header="PROPOSED BY" />
          </div>
        </Col>
      </Row>
      {txs->Belt_Array.size == 0
         ? <>
             <VSpacing size=Spacing.xl />
             <BlockIndexTxsTable txs />
             <div className=Styles.emptyContainer>
               <img src=Images.noTransaction className=Styles.noTransactionLogo />
             </div>
           </>
         : <>
             <VSpacing size=Spacing.xl />
             <BlockIndexTxsTable txs />
             <VSpacing size=Spacing.lg />
             <Pagination
               currentPage=page
               pageCount
               onPageChange={newPage => setPage(_ => newPage)}
             />
             <VSpacing size=Spacing.xl />
           </>}
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
=======
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
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
          <div className=Styles.seperatedLine />
          <Text
            value={"#B" ++ (height |> Format.iPretty)}
            weight=Text.Thin
            spacing={Text.Em(0.06)}
          />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.lg />
    <div className=Styles.vFlex> <HSpacing size=Spacing.xs /> </div>
    {switch (blockOpt) {
     | Some(block) =>
       <Text
         value={block.hash |> Hash.toHex(~upper=true)}
         size=Text.Xxl
         weight=Text.Semibold
         code=true
         nowrap=true
         ellipsis=true
       />
     | None => <Text value="in the future" size=Text.Xxl />
     }}
    <VSpacing size=Spacing.lg />
    <Row>
      <Col size=1.8>
        {switch (blockOpt) {
         | Some(block) => <InfoHL info={InfoHL.Count(block.numTxs)} header="TRANSACTIONS" />
         | None => <InfoHL info={InfoHL.Text("?", InfoHL.THEME_1)} header="TRANSACTIONS" />
         }}
      </Col>
      <Col size=4.6>
        {switch (blockOpt) {
         | Some(block) =>
           <>
             <div className=Styles.vFlex>
               <InfoHL info={InfoHL.Timestamp(block.timestamp)} header="TIME STAMP" />
             </div>
           </>
         | None => <InfoHL info={InfoHL.Text("?", InfoHL.THEME_1)} header="TIME STAMP" />
         }}
      </Col>
      <Col size=3.2>
        <div className=Styles.proposerContainer>
          {switch (monikerOpt) {
           | Some(moniker) =>
             <InfoHL info={InfoHL.Text(moniker, InfoHL.THEME_1)} header="PROPOSED BY" />
           | None => <InfoHL info={InfoHL.Text("?", InfoHL.THEME_1)} header="PROPOSED BY" />
           }}
        </div>
      </Col>
    </Row>
    {switch (blockOpt, txsOpt) {
     | (Some(_), Some({txs})) =>
       switch (txs->Belt_List.size) {
       | 0 =>
         <>
           <VSpacing size=Spacing.xl />
           <BlockIndexTxsTable txs />
           <div className=Styles.emptyContainer>
             <img src=Images.noTransaction className=Styles.noTransactionLogo />
           </div>
         </>
       | _ =>
         <>
           <VSpacing size=Spacing.xl />
           <BlockIndexTxsTable txs />
           <VSpacing size=Spacing.lg />
           <Pagination
             currentPage=page
             pageCount
             onPageChange={newPage => setPage(_ => newPage)}
           />
           <VSpacing size=Spacing.xl />
         </>
       }
     | _ => <VSpacing size={`px(280)} />
     }}
  </div>;
>>>>>>> update InfoHL for Block
};
