module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(50))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(27)), marginRight(`px(10))]);

  let sourceContainer = style([marginTop(`px(15))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let codeVerifiedBadge =
    style([
      backgroundColor(`hex("D7FFEC")),
      borderRadius(`px(6)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      padding4(~top=`px(10), ~bottom=`px(10), ~left=`px(13), ~right=`px(13)),
    ]);

  let checkLogo = style([marginRight(`px(10))]);

  let tableContainer = style([border(`px(1), `solid, Colors.lightGray)]);

  let tableHeader = style([backgroundColor(Colors.white)]);

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
};

[@react.component]
let make = (~codeHash, ~hashtag) => {
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.newScript className=Styles.logo />
          <Text
            value="DATA REQUEST SCRIPT"
            weight=Text.Semibold
            size=Text.Lg
            nowrap=true
            color=Colors.grayHeader
            block=true
          />
          <HSpacing size=Spacing.sm />
          <div className=Styles.seperatedLine />
          <Text value="CREATED 96 DAYS AGO" />
        </div>
      </Col>
      <Col>
        <div className=Styles.codeVerifiedBadge>
          <img src=Images.checkIcon className=Styles.checkLogo />
          <Text value="Code Verified" size=Text.Lg weight=Text.Semibold color=Colors.darkGreen />
        </div>
      </Col>
    </Row>
    <div className=Styles.sourceContainer>
      <Text value="ETH/USD Median Price" size=Text.Xxl weight=Text.Bold nowrap=true />
    </div>
    <VSpacing size=Spacing.xl />
    <InfoHL
      info={InfoHL.DataSources(["CoinMarketCap", "CryptoCompare", "Binance"])}
      header="DATA SOURCES"
    />
    <VSpacing size=Spacing.xl />
    <Row>
      <Col>
        <InfoHL
          info={InfoHL.Hash("0x012030123901923912391293", Colors.brightPurple)}
          header="SCRIPT HASH"
        />
      </Col>
      <HSpacing size=Spacing.xl />
      <HSpacing size=Spacing.xl />
      <Col>
        <InfoHL
          info={InfoHL.Hash("0x92392392392939239293293923", Colors.brightPurple)}
          header="CREATOR"
        />
      </Col>
    </Row>
    <VSpacing size={Css.px(25)} />
    <div className=Styles.tableContainer>
      <div className=Styles.tableHeader> {"test" |> React.string} </div>
      <div className=Styles.tableLowerContainer> <TxsTable /> </div>
    </div>
  </div>;
};
