module Styles = {
  open Css;

  let tableLowerContainer = style([position(`relative)]);
  let tableWrapper =
    style([padding(`px(24)), Media.mobile([padding2(~v=`px(20), ~h=`zero)])]);
  let codeImage = style([width(`px(20)), marginRight(`px(10))]);

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
      top(`px(10)),
      right(`px(10)),
      zIndex(2),
      Media.mobile([right(`zero), top(`px(-38))]),
    ]);
  let titleSpacing = style([marginBottom(`px(8))]);
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
        <Row marginBottom=24>
          <Col.Grid col=Col.Six colSm=Col.Six>
            <div className={Css.merge([CssHelper.flexBox(), Styles.titleSpacing])}>
              <Heading size=Heading.H5 value="Platform" />
              <HSpacing size=Spacing.xs />
              <CTooltip
                tooltipPlacementSm=CTooltip.BottomLeft
                tooltipText="The platform to which to generate the code for">
                <Icon name="fal fa-info-circle" size=10 />
              </CTooltip>
            </div>
            <Text
              value="OWASM v0.1"
              weight=Text.Regular
              size=Text.Lg
              block=true
              color=Colors.gray7
            />
          </Col.Grid>
          <Col.Grid col=Col.Six colSm=Col.Six>
            <div className={Css.merge([CssHelper.flexBox(), Styles.titleSpacing])}>
              <Heading size=Heading.H5 value="Language" />
              <HSpacing size=Spacing.xs />
              <CTooltip tooltipText="The programming language">
                <Icon name="fal fa-info-circle" size=10 />
              </CTooltip>
            </div>
            <Text
              value="Rust 1.40.0"
              weight=Text.Regular
              size=Text.Lg
              block=true
              color=Colors.gray7
            />
          </Col.Grid>
        </Row>
        <Row marginBottom=24 marginBottomSm=12>
          <Col.Grid>
            <div className={CssHelper.flexBox()}>
              <Icon name="fal fa-file" size=16 />
              <HSpacing size=Spacing.sm />
              <Text
                value="src/logic.rs"
                weight=Text.Semibold
                size=Text.Lg
                block=true
                color=Colors.gray7
              />
            </div>
          </Col.Grid>
        </Row>
        <div className=Styles.tableLowerContainer>
          <div className=Styles.copyContainer> <CopyButton data=code title="Copy Code" /> </div>
          {code |> renderCode}
        </div>
      </div>,
    );
  }
  |> Belt.Option.getWithDefault(_, React.null);
