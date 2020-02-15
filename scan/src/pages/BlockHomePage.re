module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50))]);

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
  let blocksOpt = BlockHook.latest(~limit, ~pollInterval=100000, ());
  let infoOpt = React.useContext(GlobalContext.context);

  let blocks = blocksOpt->Belt.Option.getWithDefault([]);

  let latestBlockOpt = blocks->Belt_List.get(0);
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
    <Row>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL BLOCKS"
            weight=Text.Bold
            size=Text.Xl
            nowrap=true
            color=Colors.grayHeader
          />
          <div className=Styles.seperatedLine />
          <Text
            value={
              switch (latestBlockOpt) {
              | Some(latestBlock) => latestBlock.height->Format.iPretty ++ " in total"
              | None => ""
              }
            }
          />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <THead>
      <Row>
        <Col> <div className=Styles.icon /> </Col>
        {[
           ("BLOCK", 0.6),
           ("PROPOSER", 2.0),
           ("TXN", 0.7),
           ("TOTAL FEE", 0.7),
           ("BLOCK REWARD", 0.8),
         ]
         ->Belt.List.map(((title, size)) => {
             <Col size key=title>
               <Text block=true value=title size=Text.Sm weight=Text.Bold color=Colors.grayText />
             </Col>
           })
         ->Array.of_list
         ->React.array}
      </Row>
    </THead>
    {blocksWithMonikers->Belt_List.toArray->Belt_Array.map(renderBody)->React.array}
    <VSpacing size=Spacing.lg />
    <LoadMore onClick={_ => setLimit(oldLimit => oldLimit + 10)} />
  </div>;
};
