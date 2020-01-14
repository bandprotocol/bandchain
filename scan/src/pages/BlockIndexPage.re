module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(50))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let txTypeBadge =
    style([
      paddingLeft(`px(12)),
      paddingRight(`px(12)),
      paddingTop(`px(5)),
      paddingBottom(`px(5)),
      backgroundColor(Colors.lightBlue),
      borderRadius(`px(15)),
    ]);

  let msgAmount =
    style([borderRadius(`percent(50.)), padding(`px(3)), backgroundColor(Colors.lightGray)]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let addressContainer = style([marginTop(`px(15))]);

  let checkLogo = style([marginRight(`px(10))]);

  let seperatorLine =
    style([
      width(`percent(100.)),
      height(`pxFloat(1.4)),
      backgroundColor(Colors.lightGray),
      display(`flex),
    ]);
};

[@react.component]
let make = (~height) => {
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="BLOCK"
            weight=Text.Semibold
            size=Text.Lg
            nowrap=true
            color=Colors.grayHeader
            block=true
          />
          <div className=Styles.seperatedLine />
          <Text value="51 MINUTES AGO" />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.lg />
    <div className=Styles.vFlex>
      <Text value="#" size=Text.Xxl weight=Text.Semibold color=Colors.brightPurple />
      <HSpacing size=Spacing.xs />
      <Text value=height size=Text.Xxl weight=Text.Semibold />
    </div>
    <VSpacing size=Spacing.lg />
    <Row>
      <Col size=1.> <InfoHL info={InfoHL.Count(1)} header="TRANSACTIONS" /> </Col>
      <Col size=2.5>
        <InfoHL
          info={InfoHL.Hash("0xe38475F47166d30A6e4E2E2C37e4B75E88Aa8b5B", Colors.grayHeader)}
          header="PROPOSED BY"
        />
      </Col>
      <Col size=2.>
        <InfoHL
          info={InfoHL.Timestamp(MomentRe.momentWithUnix(1578052800))}
          header="TIMESTAMP"
        />
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <div className=Styles.seperatorLine />
    <TxsTable />
    <VSpacing size=Spacing.lg />
    <LoadMore onClick={_ => ()} />
    <VSpacing size=Spacing.xl />
  </div>;
};
