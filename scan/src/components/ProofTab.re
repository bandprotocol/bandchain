module Styles = {
  open Css;
  let mediumText = style([fontSize(`px(14)), lineHeight(`px(20))]);
  let tableLowerContainer = style([padding(`px(20)), background(Colors.lighterGray)]);
  let tableHeader = style([cursor(`pointer), width(`percent(100.0))]);
  let maxHeight20 = style([maxHeight(`px(20))]);
};

module CopyComponent = {
  [@react.component]
  let make = (~evmProofHex) => {
    let (copying, setCopying) = React.useState(_ => false);

    React.useEffect1(
      () => {
        let timeoutId = Js.Global.setTimeout(() => setCopying(_ => false), 1000);
        Some(() => Js.Global.clearTimeout(timeoutId));
      },
      [|copying|],
    );

    <div
      className=Styles.tableHeader
      onClick={_ => {
        Copy.copy(
          evmProofHex |> JsBuffer.fromHex |> JsBuffer.toHex(~with0x=true),
          // Window.alert(evmProofHex->String.length > 0 ? "Copied!" : "Loading ... Please Wait");
        )
      }}>
      <Row>
        <img src=Images.copy className=Styles.maxHeight20 />
        <HSpacing size=Spacing.md />
        <Text value="Copy proof for Ethereum" size=Text.Lg color=Colors.brightPurple />
      </Row>
    </div>;
  };
};

[@react.component]
let make = (~reqID) => {
  let proofOpt = ProofHook.get(~requestId=reqID, ());

  let (innerElement, evmProofHex) =
    switch (proofOpt) {
    | Some({jsonProof, evmProofBytes}) => (
        <ReactHighlight>
          {jsonProof |> Js.Json.stringifyWithSpace(_, 2) |> React.string}
        </ReactHighlight>,
        evmProofBytes |> JsBuffer.toHex,
      )
    | None => ("Loading ..." |> React.string, "")
    };
  <div className=Styles.tableLowerContainer>
    <CopyComponent evmProofHex />
    <VSpacing size=Spacing.lg />
    <div className=Styles.mediumText> innerElement </div>
    <VSpacing size=Spacing.lg />
  </div>;
};
