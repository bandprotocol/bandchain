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

  let proofContainer =
    style([
      selector(
        "> button + button",
        [marginLeft(`px(24)), Media.mobile([marginLeft(`px(16))])],
      ),
    ]);
};

[@react.component]
let make = (~id: ID.Request.t) => {
  let (proofOpt, reload) = ProofHook.get(id);
  let isMobile = Media.isMobile();

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
    {switch (proofOpt) {
     | Some(proof) =>
       <div className={Css.merge([CssHelper.flexBox(), Styles.proofContainer])}>
         <CopyButton
           data={proof.evmProofBytes |> JsBuffer.toHex(~with0x=false)}
           title={isMobile ? "EVM" : "Copy EVM proof"}
           py=10
           px=14
         />
         <CopyButton
           data={
             proof.jsonProof->NonEVMProof.createProofFromJson |> JsBuffer.toHex(~with0x=false)
           }
           title={isMobile ? "non-EVM" : "Copy non-EVM proof"}
           py=10
           px=14
         />
       </div>
     | _ =>
       <div className={Styles.withWH(`percent(100.), `auto)}>
         <Loading width={`px(104)} />
       </div>
     }}
  </div>;
};
