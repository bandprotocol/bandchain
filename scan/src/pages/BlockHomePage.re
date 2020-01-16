module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50)), minHeight(`px(500))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let fullWidth = style([width(`percent(100.0)), display(`flex)]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
};

let renderBody = ((block, moniker): (BlockHook.Block.t, string)) => {
  let height = block.height;
  let timestamp = block.timestamp;
  let proposer = block.proposer->Address.toOperatorBech32;
  let totalTx = block.numTxs;

  <TBody key={height |> string_of_int}>
    <div className=Styles.fullWidth onClick={_ => Route.BlockIndexPage(height) |> Route.redirect}>
      <Row>
        <Col size=0.6>
          <div className=Styles.textContainer>
            <Text value="#" size=Text.Md weight=Text.Bold color=Colors.purple />
            <HSpacing size=Spacing.xs />
            <Text block=true value={height |> Format.iPretty} size=Text.Md weight=Text.Bold />
          </div>
        </Col>
        <Col size=0.8>
          <div className=Styles.textContainer>
            <TimeAgos time=timestamp size=Text.Md weight=Text.Semibold />
          </div>
        </Col>
        <Col size=2.0>
          <div className={Css.merge([Styles.textContainer, Styles.proposerBox])}>
            <Text
              block=true
              value=moniker
              size=Text.Sm
              weight=Text.Regular
              color=Colors.grayHeader
            />
            <VSpacing size=Spacing.sm />
            <Text
              block=true
              value=proposer
              size=Text.Md
              weight=Text.Bold
              code=true
              ellipsis=true
              color=Colors.black
            />
          </div>
        </Col>
        <Col size=0.7>
          <div className=Styles.textContainer>
            <Text block=true value={totalTx |> Format.iPretty} size=Text.Md weight=Text.Semibold />
          </div>
        </Col>
        <Col size=0.7>
          <div className=Styles.textContainer>
            <Text block=true value="FREE" size=Text.Md weight=Text.Semibold />
          </div>
        </Col>
        <Col size=0.8>
          <div className=Styles.textContainer>
            <Text block=true value="N/A" size=Text.Md weight=Text.Semibold />
          </div>
        </Col>
      </Row>
    </div>
  </TBody>;
};

[@react.component]
let make = () => {
  let (limit, setLimit) = React.useState(_ => 10);
  let blocksOpt = BlockHook.latest(~limit, ~pollInterval=3000, ());
  let infoOpt = React.useContext(GlobalContext.context);

  let blocks = blocksOpt->Belt.Option.getWithDefault([]);

  let latestBlockOpt = blocks->Belt_List.get(0);
  let validators =
    switch (infoOpt) {
    | Some(info) => info.validators
    | None => []
    };

  let blocksWithMonikers =
    blocks
    ->Belt_List.map(block => {
        Js.Console.log2(block, validators |> Belt_List.toArray);
        BlockHook.Block.getProposerMoniker(block, validators);
      })
    ->Belt_List.zip(blocks, _);

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
        {[
           ("BLOCK", 0.6),
           ("AGE", 0.8),
           ("PROPOSER", 2.0),
           ("TXN", 0.7),
           ("TOTAL FEE", 0.7),
           ("BLOCK REWARD", 0.8),
         ]
         ->Belt.List.map(((title, size)) => {
             <Col size key=title>
               <div className=Styles.textContainer>
                 <Text
                   block=true
                   value=title
                   size=Text.Sm
                   weight=Text.Bold
                   color=Colors.grayText
                 />
               </div>
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
