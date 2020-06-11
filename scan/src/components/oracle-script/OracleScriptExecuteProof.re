module Styles = {
  open Css;

  let hFlex = h =>
    style([display(`flex), flexDirection(`row), alignItems(`center), height(h)]);

  let vFlex = (w, h) => style([display(`flex), flexDirection(`column), width(w), height(h)]);

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

  <>
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
      {switch (proofOpt, requestOpt) {
       | (Some(proof), Some({result: Some(_)})) =>
         <div className={Styles.hFlex(`auto)}>
           <CopyButton data={proof.evmProofBytes} title="Copy EVM proof" width=115 />
           <HSpacing size=Spacing.md />
           <CopyButton
             data={
               NonEVMProof.RequestMini(requestOpt->Belt_Option.getExn)->NonEVMProof.createProof
             }
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
    </div>
  </>;
};
