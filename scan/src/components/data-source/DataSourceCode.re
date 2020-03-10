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
let make = () => {
  let code = {f|#!/bin/bash

#Cryptocurrency price endpoint: https://www.coingecko.com/api/documentations/v3
URL="https://api.coingecko.com/api/v3/simple/price?ids=$1&cs_currencies=usd"
KEY=".$1.usd"

# Performs data fetching and parses the result
curl -s -X GET $URL -H "accept: application/json" | jq -r ".[\"$1"\].usd"|f};

  React.useMemo1(
    () => <div className=Styles.tableLowerContainer> {code |> renderCode} </div>,
    [||],
  );
};
