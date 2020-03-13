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
      backgroundColor(Colors.lightBlue),
      borderRadius(`px(15)),
    ]);

  let msgAmount =
    style([borderRadius(`percent(50.)), padding(`px(3)), backgroundColor(Colors.lightGray)]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.mediumGray),
    ]);

  let addressContainer = style([marginTop(`px(15))]);

  let checkLogo = style([marginRight(`px(10))]);
  let blockLogo = style([width(`px(50)), marginRight(`px(10))]);
  let proposerContainer = style([maxWidth(`px(180))]);

  let seperatorLine =
    style([
      width(`percent(100.)),
      height(`pxFloat(1.4)),
      backgroundColor(Colors.lightGray),
      display(`flex),
    ]);
};

[@react.component]
let make = (~height: int) => {
  let (limit, setLimit) = React.useState(_ => 10);
  let txsOpt = TxHook.atHeight(height, ~limit, ());
  let infoOpt = React.useContext(GlobalContext.context);
  let blockOpt = BlockHook.atHeight(height);
  let monikerOpt = {
    let%Opt info = infoOpt;
    let%Opt block = blockOpt;
    let validators = info.validators;
    Some(BlockHook.Block.getProposerMoniker(block, validators));
  };
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
            color=Colors.mediumGray
            block=true
            spacing={Text.Em(0.06)}
          />
          <div className=Styles.seperatedLine />
          {switch (blockOpt) {
           | Some(block) => <TypeID.Block id={ID.Block.ID(height)} />
           | None => <Text value="in the future" size=Text.Xl />
           }}
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
         | None => <InfoHL info={InfoHL.Text("?")} header="TRANSACTIONS" />
         }}
      </Col>
      <Col size=4.6>
        {switch (blockOpt) {
         | Some(block) => <InfoHL info={InfoHL.Timestamp(block.timestamp)} header="TIMESTAMP" />
         | None => <InfoHL info={InfoHL.Text("?")} header="TRANSACTIONS" />
         }}
      </Col>
      <Col size=3.2>
        <div className=Styles.proposerContainer>
          {switch (monikerOpt) {
           | Some(moniker) => <InfoHL info={InfoHL.Text(moniker)} header="PROPOSED BY" />

           | None => <InfoHL info={InfoHL.Text("?")} header="PROPOSED BY" />
           }}
        </div>
      </Col>
    </Row>
    {switch (blockOpt, txsOpt) {
     | (Some(_), Some({txs})) =>
       switch (txs->Belt_List.size) {
       | 0 => <VSpacing size={`px(280)} />
       | _ =>
         <>
           <VSpacing size=Spacing.xl />
           <BlockIndexTxsTable txs />
           <VSpacing size=Spacing.lg />
           {txs->Belt_List.size < limit
              ? React.null : <LoadMore onClick={_ => setLimit(oldLimit => oldLimit + 10)} />}
           <VSpacing size=Spacing.xl />
         </>
       }
     | _ => <VSpacing size={`px(280)} />
     }}
  </div>;
};
