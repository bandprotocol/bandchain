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

  let maxWidth = w => style([maxWidth(`px(w))]);
};

[@react.component]
let make = (~packet: IBCSub.packet_t, ~oracleScriptID: ID.OracleScript.t) => {
  // TODO: If we can get the schema out of IBCSub directly then this sub is not necessary any more.
  let schemaSub = OracleScriptCodeSub.getSchemaByOracleScriptID(oracleScriptID);

  switch (packet) {
  | IBCSub.Request(request) =>
    let outputKVsOpt =
      switch (schemaSub) {
      | Data(schema) => Borsh.decode(schema, "Input", request.calldata)
      | _ => None
      };
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
      {let calldataHeadRender =
         <div className=Styles.hFlex>
           <Text value="CALLDATA" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <HSpacing size=Spacing.md />
           <CopyButton data={request.calldata} />
         </div>;
       switch (outputKVsOpt) {
       | Some(_) => calldataHeadRender
       | None =>
         <div className=Styles.topicContainer>
           calldataHeadRender
           <div className={Styles.maxWidth(250)}>
             <Text value={request.calldata |> JsBuffer.toHex} code=true ellipsis=true block=true />
           </div>
         </div>
       }}
      {switch (outputKVsOpt) {
       | Some(outputKVs) =>
         <>
           <VSpacing size=Spacing.md />
           <KVTable
             tableWidth=470
             rows={
               outputKVs
               ->Belt_Array.map(((k, v)) => [KVTable.Value(k), KVTable.Value(v)])
               ->Belt_List.fromArray
             }
           />
           <VSpacing size=Spacing.xl />
         </>
       | None => <VSpacing size=Spacing.lg />
       }}
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
    </>;
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
         let outputKVsOpt =
           switch (schemaSub) {
           | Data(schema) => Borsh.decode(schema, "Output", result)
           | _ => None
           };
         <>
           <VSpacing size=Spacing.lg />
           {let resultRender =
              <div className=Styles.hFlex>
                <Text value="RESULT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
                <HSpacing size=Spacing.md />
                <CopyButton data=result />
              </div>;
            switch (outputKVsOpt) {
            | Some(_) => resultRender
            | None =>
              <div className=Styles.topicContainer>
                resultRender
                <div className={Styles.maxWidth(250)}>
                  <Text value={result |> JsBuffer.toHex} code=true ellipsis=true block=true />
                </div>
              </div>
            }}
           <VSpacing size=Spacing.md />
           {switch (outputKVsOpt) {
            | Some(outputKVs) =>
              <>
                <KVTable
                  tableWidth=470
                  rows={
                    outputKVs
                    ->Belt_Array.map(((k, v)) => [KVTable.Value(k), KVTable.Value(v)])
                    ->Belt_List.fromArray
                  }
                />
                <VSpacing size=Spacing.md />
              </>
            | None => <VSpacing size=Spacing.md />
            }}
         </>;
       | _ => React.null
       }}
    </>
  | IBCSub.Unknown => React.null
  };
};
