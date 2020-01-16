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
      backgroundColor(Colors.grayHeader),
    ]);

  let addressContainer = style([marginTop(`px(15))]);

  let checkLogo = style([marginRight(`px(10))]);

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
  let blockOpt = BlockHook.atHeight(height, ());
  let txsOpt = TxHook.atHeight(height, ~limit, ());

  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="BLOCK"
            weight=Text.Semibold
            size=Text.Lg
            nowrap=true
            color=Colors.grayHeader
            block=true
          />
          <div className=Styles.seperatedLine />
          {switch (blockOpt) {
           | Some(block) => <TimeAgos time={block.timestamp} size=Text.Lg weight=Text.Regular />
           | None => <Text value="in the future" size=Text.Xl />
           }}
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.lg />
    <div className=Styles.vFlex>
      <Text value="#" size=Text.Xxl weight=Text.Semibold color=Colors.brightPurple />
      <HSpacing size=Spacing.xs />
      <Text value={height |> Format.iPretty} size=Text.Xxl weight=Text.Semibold />
    </div>
    <VSpacing size=Spacing.lg />
    <Row>
      <Col size=1.>
        {switch (blockOpt) {
         | Some(block) => <InfoHL info={InfoHL.Count(block.numTxs)} header="TRANSACTIONS" />
         | None => <InfoHL info={InfoHL.Text("?")} header="TRANSACTIONS" />
         }}
      </Col>
      <Col size=4.>
        <InfoHL
          info={
            InfoHL.Address(
              switch (blockOpt) {
              | Some(block) => block.proposer
              | None => "" |> Address.fromHex
              },
              Colors.grayHeader,
            )
          }
          header="PROPOSED BY"
        />
      </Col>
      <Col size=2.>
        {switch (blockOpt) {
         | Some(block) => <InfoHL info={InfoHL.Timestamp(block.timestamp)} header="TIMESTAMP" />
         | None => <InfoHL info={InfoHL.Text("?")} header="TRANSACTIONS" />
         }}
      </Col>
    </Row>
    {switch (blockOpt, txsOpt) {
     | (Some(_), Some({txs})) =>
       switch (txs->Belt_List.size) {
       | 0 => <VSpacing size={`px(280)} />
       | _ =>
         <>
           <VSpacing size=Spacing.xl />
           <TxsTable txs />
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
