module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(50))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let fixHeight = style([height(`px(40))]);

  let logo = style([width(`px(30)), marginRight(`px(10))]);

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

  let dataContainer =
    style([display(`flex), border(`px(1), `solid, `hex("EEEEEE")), flexDirection(`column)]);

  let topBoxContainer =
    style([
      display(`flex),
      background(Colors.white),
      padding(`px(24)),
      border(`px(1), `solid, `hex("EEEEEE")),
      borderBottom(`px(0), `solid, `hex("EEEEEE")),
      flexDirection(`column),
    ]);

  let flexStart = style([alignItems(`flexStart)]);
  let subHeaderContainer = style([display(`flex), flex(`num(1.))]);
  let detailContainer = style([display(`flex), flex(`num(3.5))]);

  let tableHeader =
    style([
      backgroundColor(Colors.white),
      padding3(~top=`px(30), ~h=`px(20), ~bottom=`px(20)),
    ]);
  let tableLowerContainer = style([padding(`px(20)), background(Colors.lighterGray)]);
};

[@react.component]
let make = (~reqID, ~hashtag: Route.request_tab_t) => {
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className={Css.merge([Styles.vFlex, Styles.fixHeight])}>
          <img src=Images.dataRequest className=Styles.logo />
          <Text
            value="DATA REQUEST"
            weight=Text.Semibold
            size=Text.Lg
            nowrap=true
            color=Colors.grayHeader
            block=true
          />
          <HSpacing size=Spacing.sm />
          <div className=Styles.seperatedLine />
          <Text value={j|#$reqID|j} />
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
    <VSpacing size=Spacing.xl />
    <div className=Styles.dataContainer>
      <div className=Styles.topBoxContainer>
        <div className=Styles.vFlex>
          <div className=Styles.subHeaderContainer>
            <Text value="Request ID" size=Text.Xl color=Colors.darkGrayText />
          </div>
          <div className=Styles.detailContainer> <Text value=reqID size=Text.Lg /> </div>
        </div>
        <VSpacing size=Spacing.xl />
        <div className=Styles.vFlex>
          <div className=Styles.subHeaderContainer>
            <Text value="Status" size=Text.Xl color=Colors.darkGrayText />
          </div>
          <div className=Styles.detailContainer> <RequestStatus reqID /> </div>
        </div>
        <VSpacing size=Spacing.xl />
        <div className=Styles.vFlex>
          <div className=Styles.subHeaderContainer>
            <Text value="Targeted Block" size=Text.Xl color=Colors.darkGrayText />
          </div>
          <div className=Styles.detailContainer>
            <Text value="1,329" size=Text.Lg weight=Text.Semibold />
            <HSpacing size=Spacing.sm />
            <Text value=" (2 blocks remaining)" size=Text.Lg />
          </div>
        </div>
        <VSpacing size=Spacing.xl />
        <div className={Css.merge([Styles.vFlex, Styles.flexStart])}>
          <div className=Styles.subHeaderContainer>
            <Text value="Parameters" size=Text.Xl color=Colors.darkGrayText />
          </div>
          <div className=Styles.detailContainer> <Parameters /> </div>
        </div>
      </div>
      <div className=Styles.tableHeader>
        <Row>
          <TabButton
            active={hashtag == RequestReportStatus}
            text="Data Report Status"
            route={Route.RequestIndexPage(reqID, RequestReportStatus)}
          />
          <HSpacing size=Spacing.lg />
          <TabButton
            active={hashtag == RequestProof}
            text="Proof of Validaity"
            route={Route.RequestIndexPage(reqID, RequestProof)}
          />
        </Row>
      </div>
      {switch (hashtag) {
       | RequestReportStatus =>
         <div className=Styles.tableLowerContainer>
           <Text
             value="Data Report from 3 Validators (Completed 3/4)"
             color=Colors.grayHeader
             size=Text.Lg
           />
           <VSpacing size=Spacing.lg />
           <ReportTable />
           <VSpacing size=Spacing.lg />
         </div>
       | RequestProof => <div> {"TODO1" |> React.string} </div>
       }}
    </div>
    <VSpacing size=Spacing.xxl />
  </div>;
};
