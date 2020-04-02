module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(40))]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);
  let headerContainer = style([lineHeight(`px(25))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

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

  let buttonWrapper = color =>
    style([
      backgroundColor(color),
      padding2(~h=`px(8), ~v=`px(4)),
      display(`flex),
      width(`px(103)),
      height(`px(25)),
      borderRadius(`px(6)),
      cursor(`pointer),
      alignItems(`center),
      justifyContent(`center),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(20, 32, 184, 0.2))),
    ]);

  let hFlex = h =>
    style([display(`flex), flexDirection(`row), alignItems(`center), height(h)]);

  let vFlex = (w, h) => style([display(`flex), flexDirection(`column), width(w), height(h)]);
};

let copyButton = (~data) => {
  <div
    className={Styles.buttonWrapper(Colors.blue1)}
    onClick={_ => {Copy.copy(data |> JsBuffer.toHex(~with0x=false))}}>
    <img src=Images.copy className=Styles.logo />
    <HSpacing size=Spacing.sm />
    <Text value="Copy Proof" size=Text.Sm block=true color=Colors.bandBlue nowrap=true />
  </div>;
};

let extLinkButton = () => {
  <a href="https://twitter.com/bandprotocol" target="_blank" rel="noopener">
    <div className={Styles.buttonWrapper(Colors.gray4)}>
      <img src=Images.externalLink className=Styles.logo />
      <HSpacing size=Spacing.sm />
      <Text value="What is Proof ?" size=Text.Sm block=true color=Colors.gray7 nowrap=true />
    </div>
  </a>;
};

[@react.component]
let make = (~txResponse: BandWeb3.tx_response_t, ~schema: string) =>
  {
    let paramsOutput =
      schema
      ->Borsh.extractFields("Output")
      ->Belt_Option.getWithDefault([||])
      ->Belt_Array.map(((paramName, paramType)) => Param.{paramName, paramType});

    Js.Console.log2("paramsOutput", paramsOutput);

    let requestsByTxHashSub = RequestSub.Mini.getListByTxHash(txResponse.txHash);
    let%Sub requestsByTxHash = requestsByTxHashSub;
    let requestOpt = requestsByTxHash->Belt_Array.get(0);

    let kvs = [["Price", "866825"], ["Random", "135730902915"]];
    let proof =
      "0x0000000000000000000434000000009024900000000000b0a0df0000000fd070a00b0becd989f8989af9c80000fd070a00b0becd989f8989af"
      |> JsBuffer.fromHex;

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
         | Some({result: Some(result)}) =>
           Js.Console.log2("result", result);
           <>
             <div className={Styles.hFlex(`auto)}>
               <HSpacing size=Spacing.lg />
               <div className={Styles.vFlex(`px(220), `px(20 * (kvs |> Belt_List.length)))}>
                 <Text
                   value="OUTPUT"
                   size=Text.Sm
                   color=Colors.gray6
                   weight=Text.Semibold
                   height={Text.Px(20)}
                 />
               </div>
               <div className={Styles.vFlex(`auto, `auto)}>
                 {kvs->Belt_List.map(entry =>
                    <div className={Styles.hFlex(`px(20))}>
                      <div className={Styles.vFlex(`px(220), `auto)}>
                        <Text value={entry->Belt_List.getExn(0)} color=Colors.gray8 />
                      </div>
                      <div className={Styles.vFlex(`px(440), `auto)}>
                        <Text
                          value={entry->Belt_List.getExn(1)}
                          code=true
                          color=Colors.gray8
                          weight=Text.Bold
                        />
                      </div>
                    </div>
                  )
                  |> Belt_List.toArray
                  |> React.array}
               </div>
             </div>
             <VSpacing size=Spacing.lg />
             <div className={Styles.hFlex(`auto)}>
               <HSpacing size=Spacing.lg />
               <div className={Styles.vFlex(`px(220), `auto)}>
                 <Text
                   value="PROOF OF VALIDITY"
                   size=Text.Sm
                   color=Colors.gray6
                   weight=Text.Semibold
                   height={Text.Px(15)}
                 />
               </div>
               <div className={Styles.vFlex(`px(660), `auto)}>
                 <Text
                   value={proof |> JsBuffer.toHex}
                   height={Text.Px(15)}
                   code=true
                   ellipsis=true
                 />
               </div>
             </div>
             <VSpacing size=Spacing.md />
             <div className={Styles.hFlex(`auto)}>
               <HSpacing size=Spacing.lg />
               <div className={Styles.vFlex(`px(220), `auto)} />
               {copyButton(~data=proof)}
               <HSpacing size=Spacing.lg />
               {extLinkButton()}
             </div>
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
             <div className={Styles.resultWrapper(`px(660), `px(40), `zero, `auto)}>
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
