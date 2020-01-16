let copyStringToClipboard: (string, string) => unit = [%bs.raw
  {|
function copyStringToClipboard (str,strMsg) {
   var element = document.createElement('textarea');
   element.value = str;
   element.setAttribute('readonly', '');
   element.style = {position: 'absolute', left: '-9999px'};
   document.body.appendChild(element);
   element.select();
   document.execCommand('copy');
   document.body.removeChild(element);
   alert(strMsg)
}
  |}
];

module Styles = {
  open Css;
  let mediumText = style([fontSize(`px(14)), lineHeight(`px(20))]);
  let tableLowerContainer = style([padding(`px(20)), background(Colors.lighterGray)]);
  let tableHeader = style([cursor(`pointer), width(`percent(100.0))]);
  let maxHeight20 = style([maxHeight(`px(20))]);
};

[@react.component]
let make = (~reqID) => {
  let proofOpt = ProofHook.get(~requestId=reqID, ());
  let (innerElement, evmProofHex, loading) =
    switch (proofOpt) {
    | Some({jsonProof, evmProofBytes}) => (
        <ReactHighlight>
          {jsonProof |> Js.Json.stringifyWithSpace(_, 2) |> React.string}
        </ReactHighlight>,
        evmProofBytes |> JsBuffer.toHex,
        false,
      )
    | None => ("Loading ..." |> React.string, "", true)
    };
  <div className=Styles.tableLowerContainer>
    <div
      className=Styles.tableHeader
      onClick={_ =>
        copyStringToClipboard(
          evmProofHex,
          {
            loading ? "Loading ... Please Wait" : "Copied!";
          },
        )
      }>
      <Row>
        <img src=Images.copy className=Styles.maxHeight20 />
        <HSpacing size=Spacing.md />
        <Text value="Copy proof for Ethereum" size=Text.Lg color=Colors.brightPurple />
      </Row>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.mediumText> innerElement </div>
    <VSpacing size=Spacing.lg />
  </div>;
};
