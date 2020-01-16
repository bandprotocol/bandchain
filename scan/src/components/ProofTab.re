module Styles = {
  open Css;

  let mediumText = style([fontSize(`px(14)), lineHeight(`px(20))]);
};

[@react.component]
let make = (~reqID) =>
  {
    let proofOpt = ProofHook.get(~requestId=reqID, ());
    let%Opt proof = proofOpt;

    Some(
      <div className=Styles.mediumText>
        <ReactHighlight>
          {proof |> Js.Json.stringifyWithSpace(_, 2) |> React.string}
        </ReactHighlight>
      </div>,
    );
  }
  ->Belt.Option.getWithDefault(
      <div className=Styles.mediumText> {"Loading ..." |> React.string} </div>,
    );
