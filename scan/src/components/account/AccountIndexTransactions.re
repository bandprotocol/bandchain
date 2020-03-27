module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);
};

// TODO: Mock
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
    validator: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
    reporter: "F4F9994D5E59aEf6281739b046f0E28c33b3A847" |> Address.fromHex,
  }),
];

// TODO: Mock
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
  // TODO: Mock for now
  let txsOpt = TxHook.latest(~limit=10, ());
  let txs =
    txsOpt
    ->Belt.Option.mapWithDefault([], ({txs}) => txs)
    ->Belt_List.mapWithIndex(overwriteTxAt);

  <div className=Styles.tableLowerContainer>
    <VSpacing size=Spacing.md />
    <div className=Styles.hFlex>
      <HSpacing size=Spacing.lg />
      <Text value="28" weight=Text.Semibold />
      <HSpacing size=Spacing.xs />
      <Text value="Transactions In Total" />
    </div>
    <VSpacing size=Spacing.lg />
    <TxsTable txs />
  </div>;
};
