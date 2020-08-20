module Styles = {
  open Css;
  let labelWrapper = style([flexShrink(0.), flexGrow(0.), flexBasis(`px(220))]);
  let resultBox = style([padding(`px(20))]);
  let withWH = (w, h) =>
    style([
      width(w),
      height(h),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
    ]);
};

[@react.component]
let make = (~id: ID.Request.t, ~requestOpt: option(RequestSub.Mini.t)) => {
  let (proofOpt, reload) = ProofHook.get(id);

  React.useEffect1(
    () => {
      let intervalID =
        Js.Global.setInterval(
          () =>
            if (proofOpt == None) {
              reload((), ());
            },
          2000,
        );
      Some(() => Js.Global.clearInterval(intervalID));
    },
    [|proofOpt|],
  );

  <div className={Css.merge([CssHelper.flexBox(), Styles.resultBox])}>
    <div className=Styles.labelWrapper>
      <Text
        value="Proof of validaty"
        color=Colors.gray6
        weight=Text.Regular
        height={Text.Px(15)}
      />
    </div>
    {switch (proofOpt, requestOpt) {
     | (Some(proof), Some({result: Some(_)})) =>
       <div className={CssHelper.flexBox()}>
         <CopyButton data={proof.evmProofBytes} title="Copy EVM proof" width=115 />
         <HSpacing size=Spacing.md />
         <CopyButton
           data={NonEVMProof.RequestMini(requestOpt->Belt_Option.getExn)->NonEVMProof.createProof}
           title="Copy non-EVM proof"
           width=130
         />
         <HSpacing size=Spacing.lg />
         <ExtLinkButton link="https://docs.bandchain.org/" description="What is proof ?" />
       </div>
     | _ =>
       <div className={Styles.withWH(`percent(100.), `auto)}>
         <img src=Images.loadingCircles className={Styles.withWH(`px(104), `px(30))} />
       </div>
     }}
  </div>;
};
