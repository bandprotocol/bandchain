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
    amount: [{denom: "uband", amount: 1234.34}],
  }),
  TxHook.Msg.CreateDataSource({
    id: 408,
    owner: "fACeB00cCC40D220de2A22aA206A43e7A49d54CB" |> Address.fromHex,
    name: "CoinGecko V.1.1",
    fee: [{denom: "uband", amount: 10.0}],
    executable: JsBuffer.fromHex("aa"),
    sender: "fACeB00cCC40D220de2A22aA206A43e7A49d54CB" |> Address.fromHex,
  }),
  TxHook.Msg.EditDataSource({
    id: 2430,
    owner: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
    name: "Binance API Beta",
    fee: [{denom: "uband", amount: 10.0}],
    executable: JsBuffer.fromHex("aa"),
    sender: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
  }),
  TxHook.Msg.CreateOracleScript({
    id: 35,
    owner: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
    name: "Mean Fx Price With Timestamp",
    code: JsBuffer.fromHex("aa"),
    sender: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
  }),
  TxHook.Msg.EditOracleScript({
    id: 7235,
    owner: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
    name: "Deep Learning Price",
    code: JsBuffer.fromHex("aa"),
    sender: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
  }),
  TxHook.Msg.Request({
    id: 644,
    oracleScriptID: 22,
    calldata: JsBuffer.fromHex("aa"),
    requestedValidatorCount: 0,
    sufficientValidatorCount: 0,
    expiration: 0,
    prepareGas: 0,
    executeGas: 0,
    sender: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
  }),
  TxHook.Msg.Report({
    requestID: 40,
    dataSet: [],
    sender: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
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
