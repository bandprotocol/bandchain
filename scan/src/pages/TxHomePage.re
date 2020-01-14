module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50)), minHeight(`px(500))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
};

[@react.component]
let make = () => {
  <div className=Styles.pageContainer>
    <Row>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL TRANSACTIONS"
            weight=Text.Bold
            size=Text.Xl
            nowrap=true
            color=Colors.grayHeader
          />
          <div className=Styles.seperatedLine />
          <Text value="99,999 in total" />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <TxsTable />
    <VSpacing size=Spacing.lg />
    <LoadMore onClick={_ => ()} />
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.xl />
  </div>;
};
