module Styles = {
  open Css;

  let hFlex = style([display(`flex), alignItems(`center)]);

  let topicContainer =
    style([display(`flex), justifyContent(`spaceBetween), width(`percent(100.))]);
};

[@react.component]
let make = (~packet: IBCSub.packet_t) => {
  switch (packet) {
  | IBCSub.Request(request) =>
    <>
      <div className=Styles.topicContainer>
        <Text value="ORACLE SCRIPT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={request.oracleScriptID} />
          <HSpacing size=Spacing.sm />
          <Text value={request.oracleScriptName} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.md />
      <div className=Styles.hFlex>
        <Text value="CALLDATA" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <HSpacing size=Spacing.md />
        <CopyButton data={request.calldata} />
      </div>
      <VSpacing size=Spacing.md />
      // TODO: Mock calldata
      <KVTable
        tableWidth=470
        rows=[
          [KVTable.Value("crypto_symbol"), KVTable.Value("BTC")],
          [KVTable.Value("aggregation_method"), KVTable.Value("mean")],
          [
            KVTable.Value("data_sources"),
            KVTable.Value("Binance v1, coingecko v1, coinmarketcap v1, band-validator"),
          ],
        ]
      />
      <VSpacing size=Spacing.xl />
      <div className=Styles.topicContainer>
        <Text
          value="REQUEST VALIDATOR COUNT"
          size=Text.Sm
          weight=Text.Thin
          spacing={Text.Em(0.06)}
        />
        <Text value={request.requestedValidatorCount |> string_of_int} weight=Text.Bold />
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text
          value="SUFFICIENT VALIDATOR COUNT"
          size=Text.Sm
          weight=Text.Thin
          spacing={Text.Em(0.06)}
        />
        <Text value={request.sufficientValidatorCount |> string_of_int} weight=Text.Bold />
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="REPORT PERIOD" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value={request.expiration |> string_of_int} weight=Text.Bold code=true />
          <HSpacing size=Spacing.sm />
          <Text value="Blocks" code=true />
        </div>
      </div>
    </>
  | IBCSub.Response(response) => React.null
  | IBCSub.Unknown => React.null
  };
};
