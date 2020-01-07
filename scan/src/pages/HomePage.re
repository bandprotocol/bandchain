module Styles = {
  open Css;

  let highlightsContainer =
    style([width(`percent(100.)), paddingTop(`px(40)), paddingBottom(Spacing.xl)]);

  let section = style([paddingTop(`px(48)), width(`percent(100.))]);

  let topicBar =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`row),
      justifyContent(`spaceBetween),
    ]);

  let seeAllContainer = style([alignItems(`center), justifyContent(`center), display(`flex)]);

  let rightArrow = style([width(`px(13)), marginLeft(`px(5))]);

  let seeMoreContainer =
    style([
      width(`percent(100.)),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      backgroundColor(white),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      height(`px(30)),
      cursor(`pointer),
    ]);

  let skip = style([marginTop(`px(380))]);

  let bg =
    style([
      width(`percent(100.)),
      zIndex(0),
      position(`absolute),
      top(`px(-30)),
      display(`flex),
      flexDirection(`row),
    ]);

  let imgRight = style([marginLeft(`auto), transform(`scaleX(-1.0))]);

  let grayArea =
    style([
      display(`flex),
      alignItems(`center),
      justifyContent(`center),
      overflow(`hidden),
      backgroundColor(`hex("fafafa")),
      borderBottom(`px(1), `solid, `hex("eeeeee")),
      borderTop(`px(1), `solid, `hex("eeeeee")),
      position(`absolute),
      left(`zero),
      minHeight(`px(380)),
      width(`percent(100.)),
    ]);

  let grayAreaInner =
    style([position(`absolute), maxWidth(`px(1100)), marginLeft(`auto), marginRight(`auto)]);
};

/* SEE ALL btn */
let renderSeeAll = () =>
  <div className=Styles.seeAllContainer>
    <Text block=true value="SEE ALL" size=Text.Sm weight=Text.Bold color=Colors.grayText />
    <img src=Images.rightArrow className=Styles.rightArrow />
  </div>;

[@react.component]
let make = () => {
  <div className=Styles.highlightsContainer>
    <ChainInfoHighlights />
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.lg />
    <div className=Styles.grayArea>
      <div className=Styles.bg>
        <img src=Images.bg />
        <img src=Images.bg className=Styles.imgRight />
      </div>
      <div className=Styles.grayAreaInner> <DataScriptsHighlights /> </div>
    </div>
    <div className=Styles.skip />
    <div className=Styles.section>
      <Row alignItems=`initial>
        <Col size=1.>
          <VSpacing size=Spacing.md />
          <div className=Styles.topicBar>
            <Text value="Latest Transactions" size=Text.Xl weight=Text.Bold block=true />
            {renderSeeAll()}
          </div>
          <VSpacing size=Spacing.lg />
          {Transaction.renderHeader()}
          <Transaction
            type_={Transaction.DataRequest("ETH/USD Price Feed")}
            hash="0x128f12db1a99dce2937"
            fee="0.10 BAND"
            timestamp="2 days ago"
          />
          <Transaction
            type_={Transaction.NewScript("Anime Episodes Ranking - WINTER 2020")}
            hash="0xd83ab82c9f838391283"
            fee="0.10 BAND"
            timestamp="3 days ago"
          />
          <Transaction
            type_={Transaction.DataRequest("ETH/BTC Price Feed")}
            hash="0xc83128273823dce2937"
            fee="0.10 BAND"
            timestamp="2 days ago"
          />
          <Transaction
            type_={Transaction.DataRequest("ETH/USD Price Feed")}
            hash="0xd293f12db1a99dceabb"
            fee="0.10 BAND"
            timestamp="2 days ago"
          />
          <Transaction
            type_={Transaction.DataRequest("BTC/USD Price Feed")}
            hash="0x128f12db1a99dce2937"
            fee="0.10 BAND"
            timestamp="2 days ago"
          />
          <Transaction
            type_={Transaction.NewScript("BTC/USD Price Feed")}
            hash="0xabcdef1234deadbeef2"
            fee="0.10 BAND"
            timestamp="2 days ago"
          />
          <div className=Styles.seeMoreContainer>
            <Text
              value="SEE MORE"
              size=Text.Sm
              weight=Text.Bold
              block=true
              color=Colors.grayText
            />
          </div>
        </Col>
        <HSpacing size=Spacing.xl />
        <HSpacing size=Spacing.lg />
        <Col>
          <VSpacing size=Spacing.md />
          <div className=Styles.topicBar>
            <Text value="Latest Blocks" size=Text.Xl weight=Text.Bold block=true />
            {renderSeeAll()}
          </div>
          <VSpacing size=Spacing.md />
          <LatestBlocks
            blocks=[
              LatestBlocks.{id: 472395, proposer: "Stake.us"},
              LatestBlocks.{id: 472394, proposer: "Stake.us"},
              LatestBlocks.{id: 472393, proposer: "Stake.us"},
              LatestBlocks.{id: 472392, proposer: "Stake.us"},
              LatestBlocks.{id: 472391, proposer: "Stake.us"},
              LatestBlocks.{id: 472390, proposer: "Stake.us"},
              LatestBlocks.{id: 472389, proposer: "Stake.us"},
              LatestBlocks.{id: 472388, proposer: "Stake.us"},
              LatestBlocks.{id: 472387, proposer: "Stake.us"},
            ]
          />
        </Col>
      </Row>
    </div>
  </div>;
};
