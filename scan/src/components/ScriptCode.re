module Styles = {
  open Css;

  let tableLowerContainer =
    style([
      padding(`px(20)),
      backgroundImage(
        `linearGradient((
          deg(0.0),
          [(`percent(0.0), Colors.white), (`percent(100.0), Colors.lighterGray)],
        )),
      ),
    ]);

  let codeTabHeader =
    style([lineHeight(`px(20)), borderBottom(`px(1), `solid, Colors.lightGray)]);

  let mediumText = style([fontSize(`px(14)), lineHeight(`px(20))]);

  let maxHeight20 = style([maxHeight(`px(20))]);
};

let renderCode = ((name, content)) => {
  <div key=name>
    <VSpacing size=Spacing.xl />
    <Row>
      <img src=Images.textDocument className=Styles.maxHeight20 />
      <HSpacing size=Spacing.md />
      <Text value=name size=Text.Lg color=Colors.grayHeader />
    </Row>
    <VSpacing size=Spacing.md />
    <div className=Styles.mediumText>
      <ReactHighlight> {content |> React.string} </ReactHighlight>
    </div>
  </div>;
};

[@react.component]
let make = (~codeHash) => {
  let codes = CodeHook.getCode(codeHash);
  let watchedValue = codes->Belt.Option.isSome ? codeHash |> Hash.toHex : "";

  React.useMemo1(
    () =>
      <div className=Styles.tableLowerContainer>
        <div className=Styles.codeTabHeader>
          <Row>
            <Col size=1.0>
              <Row>
                <Col size=1.0>
                  <Text value="Platform" color=Colors.darkGrayText size=Text.Lg />
                </Col>
                <Col size=1.0> <Text value="OWASM v0.1" size=Text.Lg /> </Col>
              </Row>
            </Col>
            <Col size=1.0>
              <Row>
                <Col size=1.0>
                  <Text value="Parameters" color=Colors.darkGrayText size=Text.Lg />
                </Col>
                <Col size=1.0> <Text value="2" size=Text.Lg /> </Col>
              </Row>
            </Col>
          </Row>
          <VSpacing size=Spacing.lg />
          <Row>
            <Col size=1.0>
              <Row>
                <Col size=1.0>
                  <Text value="Language" color=Colors.darkGrayText size=Text.Lg />
                </Col>
                <Col size=1.0> <Text value="Rust 1.39.0" size=Text.Lg /> </Col>
              </Row>
            </Col>
            <Col size=1.0> <div /> </Col>
          </Row>
          <VSpacing size=Spacing.lg />
        </div>
        {switch (codes) {
         | Some(codes') =>
           codes'
           ->Belt.List.map(({name, content}) => (name, content))
           ->Belt.List.toArray
           ->Belt.Array.map(renderCode)
           ->React.array
         | None => <Text value="Code Not Found" size=Text.Lg />
         }}
      </div>,
    [|watchedValue|],
  );
};
