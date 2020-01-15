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

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
};

[@react.component]
let make = () => {
  let step = 10;
  let (limit, setLimit) = React.useState(_ => step);
  let txsOpt = TxHook.latest(~limit, ~pollInterval=3000, ());
  let txs = txsOpt->Belt.Option.getWithDefault([]);

  let infoOpt = React.useContext(GlobalContext.context);

  <div className=Styles.pageContainer>
    <Row>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL TRANSACTIONS"
            weight=Text.Bold
            size=Text.Xl
            nowrap=true
            color=Colors.grayHeader
          />
          <div className=Styles.seperatedLine />
          {switch (infoOpt) {
           | Some(info) =>
             <Text value={(info.latestBlock.totalTxs |> Format.iPretty) ++ " in total"} />
           | None => React.null
           }}
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <TxsTable txs />
    <VSpacing size=Spacing.lg />
    {switch (infoOpt) {
     | Some(info) =>
       txs->Belt_List.size == 0
       || txs->Belt_List.size
       mod step != 0
       || txs->Belt_List.size == info.latestBlock.totalTxs
         ? React.null : <LoadMore onClick={_ => {setLimit(oldLimit => oldLimit + step)}} />
     | None => React.null
     }}
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.xl />
  </div>;
};
