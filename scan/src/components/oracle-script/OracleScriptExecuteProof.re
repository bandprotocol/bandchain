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
let make = (~id: ID.Request.t) => {
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
      {switch (proofOpt) {
       | Some(proof) =>
         <div className={Styles.vFlex(`px(660), `auto)}>
           <Text
             value={proof.evmProofBytes |> JsBuffer.toHex}
             height={Text.Px(15)}
             code=true
             ellipsis=true
           />
         </div>
       | None =>
         <div className={Styles.withWH(`percent(100.), `auto)}>
           <img src=Images.loadingCircles className={Styles.withWH(`px(104), `px(30))} />
         </div>
       }}
    </div>
    <VSpacing size=Spacing.md />
    {switch (proofOpt) {
     | Some(proof) =>
       <div className={Styles.hFlex(`auto)}>
         <HSpacing size=Spacing.lg />
         <div className={Styles.vFlex(`px(220), `auto)} />
         <CopyButton data={proof.evmProofBytes} />
         <HSpacing size=Spacing.lg />
         <ExtLinkButton link="https://docs.bandchain.org/" description="What is proof ?" />
       </div>
     | None => React.null
     }}
  </>;
};
