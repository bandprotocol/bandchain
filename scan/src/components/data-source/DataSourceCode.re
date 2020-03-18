module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(20))]);

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

let renderCode = content => {
  <div className=Styles.scriptContainer>
    <ReactHighlight>
      <div className=Styles.padding> {content |> React.string} </div>
    </ReactHighlight>
  </div>;
};

[@react.component]
let make = (~executable) => {
  let code = executable |> JsBuffer._toString(_, "UTF-8");

  React.useMemo1(
    () => <div className=Styles.tableLowerContainer> {code |> renderCode} </div>,
    [||],
  );
};
