module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(35))]);
  let validatorsLogo = style([marginRight(`px(10))]);
  let highlight = style([margin2(~v=`px(28), ~h=`zero)]);
  let tableSpace = style([width(`px(26))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let fullWidth = style([width(`percent(100.0)), display(`flex)]);

  let icon =
    style([
      width(`px(30)),
      height(`px(30)),
      marginTop(`px(5)),
      marginLeft(Spacing.xl),
      marginRight(Spacing.xl),
    ]);
};

let renderBody = ((block, moniker): (BlockHook.Block.t, string)) => {
  let height = block.height;
  let timestamp = block.timestamp;
  let proposer = block.proposer->Address.toOperatorBech32;
  let totalTx = block.numTxs;

  <TBody key={height |> string_of_int}>
    <div className=Styles.fullWidth onClick={_ => Route.BlockIndexPage(height) |> Route.redirect}>
      <Row>
        <Col> <img src=Images.blockLogo className=Styles.icon /> </Col>
        <Col size=0.6>
          <TElement elementType={TElement.HeightWithTime(height, timestamp)} />
        </Col>
        <Col size=2.0> <TElement elementType={TElement.Proposer(moniker, proposer)} /> </Col>
        <Col size=0.7> <TElement elementType={TElement.Count(totalTx)} /> </Col>
        <Col size=0.7> <TElement elementType={TElement.Fee(0.0)} /> </Col>
        <Col size=0.8> <Text block=true value="N/A" size=Text.Md weight=Text.Semibold /> </Col>
      </Row>
    </div>
  </TBody>;
};

[@react.component]
let make = () => {
  let (limit, setLimit) = React.useState(_ => 10);
  let blocksOpt = BlockHook.latest(~limit, ());
  let infoOpt = React.useContext(GlobalContext.context);

  let blocks = blocksOpt->Belt.Option.getWithDefault([]);

  let validators =
    switch (infoOpt) {
    | Some(info) => info.validators
    | None => []
    };

  let blocksWithMonikers =
    blocks->Belt_List.map(block =>
      (block, BlockHook.Block.getProposerMoniker(block, validators))
    );

  <div className=Styles.pageContainer>
    <div className=Styles.vFlex>
      <img src=Images.validators className=Styles.validatorsLogo />
      <Text
        value="ALL VALIDATORS"
        weight=Text.Medium
        size=Text.Md
        nowrap=true
        color=Colors.grayHeader
        spacing={Text.Em(0.06)}
      />
      <div className=Styles.seperatedLine />
      <Text value={20->Format.iPretty ++ " In total"} />
    </div>
    <div className=Styles.highlight>
      <Row>
        <Col size=0.7> <InfoHL info={InfoHL.Fraction(8, 20, false)} header="VALIDATOR" /> </Col>
        <Col size=1.1>
          <InfoHL info={InfoHL.Fraction(5352500, 10849023, true)} header="BONDED TOKENS" />
        </Col>
        <Col size=0.9>
          <InfoHL info={InfoHL.FloatWithSuffix(12.45, "  %")} header="INFLATION RATE" />
        </Col>
        <Col size=0.51>
          <InfoHL info={InfoHL.FloatWithSuffix(2.59, "  secs")} header="24 HOUR AVG BLOCK TIME" />
        </Col>
      </Row>
    </div>
    // TODO : Add toggle button
    <THead>
      <Row>
        <Col> <div className=Styles.tableSpace /> </Col>
        {[
           ("RANK", 0.8),
           ("VALIDATOR", 1.6),
           ("VOTING POWER (BAND)", 2.1),
           ("COMMISSION (%)", 1.9),
           ("UPTIME (%)", 1.3),
           ("REPORT RATE (%)", 1.5),
         ]
         ->Belt.List.map(((title, size)) => {
             <Col size key=title>
               <Text
                 block=true
                 value=title
                 size=Text.Sm
                 weight=Text.Semibold
                 color=Colors.graySubHeader
                 spacing={Text.Em(0.1)}
               />
             </Col>
           })
         ->Array.of_list
         ->React.array}
      </Row>
    </THead>
    // <Col> <div className=Styles.tableSpace /> </Col>
    {blocksWithMonikers->Belt_List.toArray->Belt_Array.map(renderBody)->React.array}
    <VSpacing size=Spacing.lg />
    <LoadMore onClick={_ => setLimit(oldLimit => oldLimit + 10)} />
  </div>;
};
