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
  let oracleScriptSub = OracleScriptSub.get(oracleScriptID);
  switch (packet) {
  | IBCSub.Request(request) =>
    // TODO: support loading state, no data later
    let outputKVsOpt =
      switch (oracleScriptSub) {
      | Data(oracleScript) => Obi.decode(oracleScript.schema, "input", request.calldata)
      | _ => None
      };
    <>
      {switch (request.idOpt) {
       | Some(id) =>
         <>
           <div className=Styles.topicContainer>
             <Text value="REQUEST ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
             <div className=Styles.hFlex> <TypeID.Request id /> </div>
           </div>
           <VSpacing size=Spacing.md />
         </>
       | None => React.null
       }}
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
           <CopyButton
             data={request.calldata |> JsBuffer.toHex(~with0x=false)}
             title="Copy as bytes"
             width=125
           />
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
             rows={
               outputKVs
               ->Belt_Array.map(({fieldName, fieldValue}) =>
                   [KVTable.Value(fieldName), KVTable.Value(fieldValue)]
                 )
               ->Belt_List.fromArray
             }
           />
           <VSpacing size=Spacing.xl />
         </>
       | None => <VSpacing size=Spacing.lg />
       }}
      <div className=Styles.topicContainer>
        <Text value="ASK COUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={request.askCount |> string_of_int} weight=Text.Bold />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="MIN COUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={request.minCount |> string_of_int} weight=Text.Bold />
      </div>
    </>;
  | IBCSub.Response(response) =>
    // TODO: support loading state, no data later
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
           switch (oracleScriptSub) {
           | Data(oracleScript) => Obi.decode(oracleScript.schema, "output", result)
           | _ => None
           };
         <>
           <VSpacing size=Spacing.lg />
           {let resultHeadRender =
              <div className=Styles.hFlex>
                <Text value="RESULT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
                <HSpacing size=Spacing.md />
                <CopyButton
                  data={result |> JsBuffer.toHex(~with0x=false)}
                  title="Copy as bytes"
                  width=125
                />
              </div>;
            switch (outputKVsOpt) {
            | Some(_) => resultHeadRender
            | None =>
              <div className=Styles.topicContainer>
                resultHeadRender
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
                  rows={
                    outputKVs
                    ->Belt_Array.map(({fieldName, fieldValue}) =>
                        [KVTable.Value(fieldName), KVTable.Value(fieldValue)]
                      )
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
