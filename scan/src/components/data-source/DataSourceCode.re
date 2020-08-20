module Styles = {
  open Css;

  let tableLowerContainer =
    style([
      padding(`px(20)),
      position(`relative),
      Media.mobile([padding2(~v=`px(20), ~h=`zero)]),
    ]);

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
  let copyContainer =
    style([
      position(`absolute),
      top(`px(30)),
      right(`px(30)),
      zIndex(2),
      Media.mobile([
        position(`static),
        display(`flex),
        justifyContent(`flexEnd),
        marginBottom(`px(8)),
      ]),
    ]);
};

let renderCode = content => {
  <div className=Styles.scriptContainer>
    <ReactHighlight className=Styles.padding> {content |> React.string} </ReactHighlight>
  </div>;
};

[@react.component]
let make = (~executable) => {
  let code = executable |> JsBuffer.toUTF8;
  React.useMemo1(
    () =>
      <div className=Styles.tableLowerContainer>
        <div className=Styles.copyContainer> <CopyButton.Modern data=code title="Copy Code" /> </div>
        {code |> renderCode}
      </div>,
    [||],
  );
};
