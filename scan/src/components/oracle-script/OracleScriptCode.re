module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(8))]);
  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);
  let codeImage = style([width(`px(20)), marginRight(`px(10))]);
  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

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
    <ReactHighlight className=Styles.padding> {content |> React.string} </ReactHighlight>
  </div>;
};

[@react.component]
let make = (~url: string) =>
  {
    let%Opt code = CodeHook.getCode(url);
    Some(
      <div className=Styles.tableWrapper>
        <>
          <VSpacing size={`px(10)} />
          <Row>
            <HSpacing size={`px(15)} />
            <Col>
              <div> <Text value="Platform" /> </div>
              <VSpacing size={`px(5)} />
              <div> <Text value="OWASM v0.1" code=true weight=Text.Semibold /> </div>
            </Col>
            <HSpacing size={`px(370)} />
            <Col>
              <div> <Text value="Language" /> </div>
              <VSpacing size={`px(5)} />
              <div> <Text value="Rust 1.40.0" code=true weight=Text.Semibold /> </div>
            </Col>
          </Row>
          <VSpacing size={`px(35)} />
          <div className=Styles.tableLowerContainer>
            <div className=Styles.vFlex>
              <img src=Images.code className=Styles.codeImage />
              <Text value="src/logic.rs" size=Text.Lg color=Colors.gray7 />
            </div>
            <VSpacing size=Spacing.lg />
            code->renderCode
          </div>
        </>
      </div>,
    );
  }
  |> Belt.Option.getWithDefault(_, React.null);
