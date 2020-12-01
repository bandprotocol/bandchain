module Styles = {
  open Css;

  let resultContainer =
    style([
      backgroundColor(Colors.white),
      margin2(~v=`px(20), ~h=`zero),
      selector("> div + div", [borderTop(`px(1), `solid, Colors.gray9)]),
    ]);
  let resultBox = style([padding(`px(20))]);
  let labelWrapper = style([flexShrink(0.), flexGrow(0.), flexBasis(`px(220))]);
  let resultWrapper =
    style([
      flexShrink(0.),
      flexGrow(0.),
      flexBasis(`calc((`sub, `percent(100.), `px(220)))),
    ]);
};

[@react.component]
let make = (~txResponse: TxCreator.tx_response_t, ~schema: string) =>
  {
    let requestsByTxHashSub = RequestSub.Mini.getListByTxHash(txResponse.txHash);
    let%Sub requestsByTxHash = requestsByTxHashSub;
    let requestOpt = requestsByTxHash->Belt_Array.get(0);

    <>
      <div className=Styles.resultContainer>
        <div className={Css.merge([CssHelper.flexBox(), Styles.resultBox])}>
          <div className=Styles.labelWrapper>
            <Text value="Exit Status" color=Colors.gray6 weight=Text.Regular />
          </div>
          <Text value={txResponse.success ? "0" : "1"} />
        </div>
        {switch (requestOpt) {
         | Some({id}) =>
           <div className={Css.merge([CssHelper.flexBox(), Styles.resultBox])}>
             <div className=Styles.labelWrapper>
               <Text value="Request ID" color=Colors.gray6 weight=Text.Regular />
             </div>
             <TypeID.Request id />
           </div>

         | None => React.null
         }}
        <div className={Css.merge([CssHelper.flexBox(), Styles.resultBox])}>
          <div className=Styles.labelWrapper>
            <Text value="Tx Hash" color=Colors.gray6 weight=Text.Regular />
          </div>
          <TxLink txHash={txResponse.txHash} width=500 />
        </div>
        {switch (requestOpt) {
         | Some({result: Some(result), id}) =>
           let outputKVsOpt = Obi.decode(schema, "output", result);
           <>
             <div className={Css.merge([CssHelper.flexBox(), Styles.resultBox])}>
               <div className=Styles.labelWrapper>
                 <Text
                   value="Output"
                   color=Colors.gray6
                   weight=Text.Regular
                   height={Text.Px(20)}
                 />
               </div>
               <div className=Styles.resultWrapper>
                 {switch (outputKVsOpt) {
                  | Some(outputKVs) =>
                    <KVTable
                      rows={
                        outputKVs
                        ->Belt_Array.map(({fieldName, fieldValue}) =>
                            [KVTable.Value(fieldName), KVTable.Value(fieldValue)]
                          )
                        ->Belt_List.fromArray
                      }
                    />
                  | None => React.null
                  }}
               </div>
             </div>
             <OracleScriptExecuteProof id />
           </>;
         | Some(request) =>
           <div className={Css.merge([CssHelper.flexBox(), Styles.resultBox])}>
             <div className=Styles.labelWrapper>
               <Text
                 value="Waiting for output and `proof`"
                 color=Colors.gray6
                 weight=Text.Regular
               />
             </div>
             <div className=Styles.resultWrapper>
               <ProgressBar
                 reportedValidators={request.reportsCount}
                 minimumValidators={request.minCount}
                 requestValidators={request.askCount}
               />
             </div>
           </div>
         | None => React.null
         }}
      </div>
    </>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
