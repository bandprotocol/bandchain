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
  let code = {f|[package]
name = "crypto_price"
version = "0.1.0"
authors = ["Band Protocol <dev@bandprotocol.com>"]
edition = "2018"

# See more keys and their definitions at https://doc.rust-lang.org/cargo/reference/manifest.html

[lib]
crate-type = ["cdylib"]

[dependencies]
owasm = { path = "../.." }
|f};

  React.useMemo1(
    () =>
      <div className=Styles.tableLowerContainer>
        <Text value="Cargo.toml" />
        <VSpacing size=Spacing.md />
        {code |> renderCode}
      </div>,
    [||],
  );
};
