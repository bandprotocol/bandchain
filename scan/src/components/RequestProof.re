module Styles = {
  open Css;

  let emptyContainer =
    style([
      height(`px(130)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      flexDirection(`column),
      backgroundColor(Colors.blueGray1),
    ]);

  let withWidth = w => style([width(`px(w))]);

  let topicContainer = h =>
    style([display(`flex), alignItems(`center), width(`percent(100.)), height(`px(h))]);

  let scriptContainer =
    style([
      fontSize(`px(12)),
      lineHeight(`px(20)),
      fontFamilies([
        `custom("IBM Plex Mono"),
        `custom("cousine"),
        `custom("sfmono-regular"),
        `custom("Consolas"),
        `custom("Menlo"),
        `custom("liberation mono"),
        `custom("ubuntu mono"),
        `custom("Courier"),
        `monospace,
      ]),
    ]);

  let padding = style([padding(`px(20))]);

  let loading = style([width(`px(65)), height(`px(20)), marginBottom(`px(16))]);
};

[@react.component]
let make = (~request: RequestSub.t) => {
  let (proofOpt, reload) = ProofHook.get(request.id);
  let (showProof, setShowProof) = React.useState(_ => false);
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

  switch (proofOpt) {
  | Some(proof) =>
    <>
      <div className={CssHelper.flexBox()}>
        <ShowProofButton showProof setShowProof />
        <HSpacing size={`px(24)} />
        <CopyButton.Modern
          data={proof.evmProofBytes |> JsBuffer.toHex(~with0x=false)}
          title={isMobile ? "EVM" : "Copy EVM proof"}
          width=155
          py=12
          px=20
        />
        <HSpacing size={`px(24)} />
        <CopyButton.Modern
          data={
            NonEVMProof.Request(request)->NonEVMProof.createProof
            |> JsBuffer.toHex(~with0x=false)
          }
          title={isMobile ? "non-EVM" : "Copy non-EVM proof"}
          width=180
          py=12
          px=20
        />
      </div>
      {showProof
         ? <>
             <VSpacing size=Spacing.lg />
             <div className=Styles.scriptContainer>
               <ReactHighlight className=Styles.padding>
                 {proof.jsonProof |> Js.Json.stringifyWithSpace(_, 2) |> React.string}
               </ReactHighlight>
             </div>
           </>
         : React.null}
    </>
  | None =>
    <div className=Styles.emptyContainer>
      <img src=Images.loadingCircles className=Styles.loading />
      <Heading
        size=Heading.H4
        value="Waiting for proof"
        align=Heading.Center
        weight=Heading.Regular
        color=Colors.bandBlue
      />
    </div>
  };
};
