module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50)), minHeight(`px(500))]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

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

let msg_actions = [
  TxHook.Msg.Send({
    fromAddress: "fACeB00cCC40D220de2A22aA206A43e7A49d54CB" |> Address.fromHex,
    toAddress: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
    amount: [{denom: "uband", amount: 100.0}],
  }),
];

let overwriteTxAt = (i: int, tx: TxHook.Tx.t) =>
  i < (msg_actions |> Belt_List.length)
    ? {
      ...tx,
      messages: [
        {...tx.messages->Belt_List.getExn(0), action: msg_actions->Belt_List.getExn(i)},
      ],
    }
    : tx;

[@react.component]
let make = () => {
  let step = 10;
  let (limit, setLimit) = React.useState(_ => step);
  let txsOpt = TxHook.latest(~limit, ());
  let txs =
    txsOpt
    ->Belt.Option.mapWithDefault([], ({txs}) => txs)
    ->Belt_List.mapWithIndex(overwriteTxAt);

  let infoOpt = React.useContext(GlobalContext.context);
  let totalTxsOpt = infoOpt->Belt.Option.map(info => info.latestBlock.totalTxs);

  <div className=Styles.pageContainer>
    <Row>
      <Col> <img src=Images.txLogo className=Styles.logo /> </Col>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL TRANSACTIONS"
            weight=Text.Semibold
            nowrap=true
            spacing={Text.Em(0.06)}
            color=Colors.grayHeader
          />
          <div className=Styles.seperatedLine />
          {switch (totalTxsOpt) {
           | Some(totalTxs) => <Text value={(totalTxs * 100 |> Format.iPretty) ++ " in total"} />
           | None => React.null
           }}
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <TxsTable txs />
    <VSpacing size={`px(70)} />
  </div>;
};
