module Styles = {
  open Css;

  let hFlex = style([display(`flex), alignItems(`center)]);

  let topicContainer =
    style([
      display(`flex),
      justifyContent(`spaceBetween),
      width(`percent(100.)),
      height(`px(16)),
      alignItems(`center),
    ]);

  let statusContainer = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(20))]);
};

[@react.component]
let make = (~packet: IBCSub.packet_t) => {
  switch (packet) {
  | IBCSub.Request(request) =>
    <>
      <div className=Styles.topicContainer>
        <Text value="REQUEST ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex> <TypeID.Request id={request.id} /> </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="ORACLE SCRIPT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={request.oracleScriptID} />
          <HSpacing size=Spacing.sm />
          <Text value={request.oracleScriptName} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
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
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text
          value="SUFFICIENT VALIDATOR COUNT"
          size=Text.Sm
          weight=Text.Thin
          spacing={Text.Em(0.06)}
        />
        <Text value={request.sufficientValidatorCount |> string_of_int} weight=Text.Bold />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="REPORT PERIOD" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value={request.expiration |> string_of_int} weight=Text.Bold code=true />
          <HSpacing size=Spacing.sm />
          <Text value="Blocks" code=true />
        </div>
      </div>
    </>
  | IBCSub.Response(response) =>
    <>
      <div className=Styles.topicContainer>
        <Text value="REQUEST ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex> <TypeID.Request id={response.requestID} /> </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="ORACLE SCRIPT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={response.oracleScriptID} />
          <HSpacing size=Spacing.sm />
          <Text value={response.oracleScriptName} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="STATUS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <div className=Styles.statusContainer>
            <Text
              block=true
              code=true
              spacing={Text.Em(0.02)}
              value={response.status == IBCSub.Response.Success ? "success" : "fail"}
              weight=Text.Medium
              ellipsis=true
            />
            <HSpacing size=Spacing.md />
            <img
              src={response.status == IBCSub.Response.Success ? Images.success : Images.fail}
              className=Styles.logo
            />
          </div>
        </div>
      </div>
      {switch (response.status, response.result) {
       | (IBCSub.Response.Success, Some(result)) =>
         <>
           <VSpacing size=Spacing.lg />
           <div className=Styles.hFlex>
             <Text value="CALLDATA" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
             <HSpacing size=Spacing.md />
             <CopyButton data=result />
           </div>
           <VSpacing size=Spacing.md />
           // TODO: Mock calldata
           <KVTable
             tableWidth=470
             rows=[
               [KVTable.Value("px"), KVTable.Value("70000")],
               [KVTable.Value("timestamp"), KVTable.Value("1586319310")],
             ]
           />
           <VSpacing size=Spacing.md />
         </>
       | _ => React.null
       }}
    </>
  | IBCSub.Unknown => React.null
  };
};
