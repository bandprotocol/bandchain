module Styles = {
  open Css;

  let hFlex = style([display(`flex), alignItems(`center)]);

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
};

[@react.component]
let make = (~requestID) => {
  let proofOpt = ProofHook.get(requestID);
  let (showProof, setShowProof) = React.useState(_ => false);

  switch (proofOpt) {
  | Some(proof) =>
    <>
      <VSpacing size=Spacing.lg />
      <div className={Styles.topicContainer(40)}>
        <Col size=1.>
          <Text
            value="PROOF OF VALIDITY"
            size=Text.Sm
            weight=Text.Semibold
            spacing={Text.Em(0.06)}
            color=Colors.gray6
          />
        </Col>
        <Col size=1.>
          <div className={Styles.withWidth(700)}>
            <Text
              value={proof.evmProofBytes |> JsBuffer.toHex}
              weight=Text.Medium
              color=Colors.gray7
              block=true
              code=true
              ellipsis=true
            />
          </div>
        </Col>
      </div>
      <div className={Styles.topicContainer(20)}>
        <Col size=1.> React.null </Col>
        <Col size=1.>
          <div className={Styles.withWidth(700)}>
            <div className=Styles.hFlex>
              <ShowProofButton showProof setShowProof />
              <HSpacing size=Spacing.md />
              <CopyButton data={proof.evmProofBytes} />
              <HSpacing size=Spacing.md />
              <ExtLinkButton link="https://docs.bandchain.org/" description="What is proof ?" />
            </div>
          </div>
        </Col>
      </div>
      {showProof
         ? <>
             <VSpacing size=Spacing.lg />
             <div className=Styles.scriptContainer>
               <ReactHighlight>
                 <div className=Styles.padding>
                   {proof.jsonProof |> Js.Json.stringifyWithSpace(_, 2) |> React.string}
                 </div>
               </ReactHighlight>
             </div>
           </>
         : React.null}
    </>
  | None => React.null
  };
};
