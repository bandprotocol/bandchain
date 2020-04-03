module Styles = {
  open Css;

  let resultWrapper = (w, h, paddingV, overflow_choice) =>
    style([
      width(w),
      height(h),
      display(`flex),
      flexDirection(`column),
      padding2(~v=paddingV, ~h=`zero),
      justifyContent(`center),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      overflow(overflow_choice),
    ]);

  let hFlex = h =>
    style([display(`flex), flexDirection(`row), alignItems(`center), height(h)]);

  let vFlex = (w, h) => style([display(`flex), flexDirection(`column), width(w), height(h)]);

  let pr = size => style([paddingRight(`px(size))]);
};

[@react.component]
let make = (~txResponse: BandWeb3.tx_response_t, ~schema: string) =>
  {
    let requestsByTxHashSub = RequestSub.Mini.getListByTxHash(txResponse.txHash);
    let%Sub requestsByTxHash = requestsByTxHashSub;
    let requestOpt = requestsByTxHash->Belt_Array.get(0);

    <>
      <VSpacing size=Spacing.lg />
      <div className={Styles.resultWrapper(`percent(100.), `auto, `px(30), `auto)}>
        <div className={Styles.hFlex(`auto)}>
          <HSpacing size=Spacing.lg />
          <div className={Styles.resultWrapper(`px(220), `px(12), `zero, `auto)}>
            <Text value="EXIT STATUS" size=Text.Sm color=Colors.gray6 weight=Text.Semibold />
          </div>
          <Text value={txResponse.success ? "0" : "1"} />
        </div>
        {switch (requestOpt) {
         | Some({id}) =>
           <>
             <VSpacing size=Spacing.lg />
             <div className={Styles.hFlex(`auto)}>
               <HSpacing size=Spacing.lg />
               <div className={Styles.resultWrapper(`px(220), `px(12), `zero, `auto)}>
                 <Text value="REQUEST ID" size=Text.Sm color=Colors.gray6 weight=Text.Semibold />
               </div>
               <TypeID.Request id />
             </div>
           </>
         | None => React.null
         }}
        <VSpacing size=Spacing.lg />
        <div className={Styles.hFlex(`auto)}>
          <HSpacing size=Spacing.lg />
          <div className={Styles.resultWrapper(`px(220), `px(12), `zero, `auto)}>
            <Text value="TX HASH" size=Text.Sm color=Colors.gray6 weight=Text.Semibold />
          </div>
          <TxLink txHash={txResponse.txHash} width=500 />
        </div>
        <VSpacing size=Spacing.lg />
        {switch (requestOpt) {
         | Some({id, result: Some(result)}) =>
           let outputKVs =
             Borsh.decode(schema, "Output", result)->Belt_Option.getWithDefault([||]);
           <>
             <div className={Styles.hFlex(`auto)}>
               <HSpacing size=Spacing.lg />
               <div
                 className={Styles.vFlex(`px(220), `px(20 * (outputKVs |> Belt_Array.size)))}>
                 <Text
                   value="OUTPUT"
                   size=Text.Sm
                   color=Colors.gray6
                   weight=Text.Semibold
                   height={Text.Px(20)}
                 />
               </div>
               <div className={Styles.vFlex(`auto, `auto)}>
                 {outputKVs->Belt_Array.map(((k, v)) =>
                    <div className={Styles.hFlex(`px(20))}>
                      <div className={Styles.vFlex(`px(220), `auto)}>
                        <Text value=k color=Colors.gray8 />
                      </div>
                      <div className={Styles.vFlex(`px(440), `auto)}>
                        <Text value=v code=true color=Colors.gray8 weight=Text.Bold />
                      </div>
                    </div>
                  )
                  |> React.array}
               </div>
             </div>
             <VSpacing size=Spacing.lg />
             <OracleScriptExecuteProof id />
           </>;
         | Some(request) =>
           <div className={Styles.hFlex(`auto)}>
             <HSpacing size=Spacing.lg />
             <div className={Styles.resultWrapper(`px(220), `px(12), `zero, `auto)}>
               <Text
                 value="WAITING FOR OUTPUT AND PROOF"
                 size=Text.Sm
                 color=Colors.gray6
                 weight=Text.Semibold
               />
             </div>
             <div
               className={Css.merge([
                 Styles.resultWrapper(`px(660), `px(40), `zero, `auto),
                 Styles.pr(40),
               ])}>
               <ProgressBar
                 reportedValidators={request.reportsCount}
                 minimumValidators={request.sufficientValidatorCount}
                 requestValidators={request.requestedValidatorsCount}
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
