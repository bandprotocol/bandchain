module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(50))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);
};

[@react.component]
let make = (~dataSourceID, ~hashtag: Route.data_source_tab_t) => {
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.dataSourceLogo className=Styles.logo />
          <Text
            value="DATA SOURCE"
            weight=Text.Semibold
            size=Text.Md
            nowrap=true
            color=Colors.grayHeader
            block=true
          />
          <div className=Styles.seperatedLine />
          <Text
            value="Last updated 4 hours ago"
            size=Text.Md
            weight=Text.Regular
            color=Colors.grayHeader
            nowrap=true
          />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.md />
    <div className=Styles.vFlex>
      <Text
        value="#D253"
        size=Text.Xxl
        weight=Text.Bold
        color=Colors.brightOrange
        nowrap=true
        code=true
      />
      <HSpacing size=Spacing.md />
      <Text value="CoinGecko V.2" size=Text.Xxl weight=Text.Bold nowrap=true />
    </div>
    <VSpacing size=Spacing.xl />
    <Row>
      <Col size=1.> <InfoHL header="OWNER" info={InfoHL.MockAddress(Colors.grayHeader)} /> </Col>
      <Col size=0.8> <InfoHL info={InfoHL.Fee(1000.)} header="REQUEST FEE" /> </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.xxl />
  </div>;
};
