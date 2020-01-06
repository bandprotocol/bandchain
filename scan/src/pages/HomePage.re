module Styles = {
  open Css;

  let highlightsContainer =
    style([width(`percent(100.)), paddingTop(`px(40)), paddingBottom(Spacing.xl)]);

  let section = style([paddingTop(`px(48)), width(`percent(100.))]);
};

[@react.component]
let make = () => {
  <div className=Styles.highlightsContainer>
    <ChainInfoHighlights />
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.lg />
    <DataScriptsHighlights />
    <div className=Styles.section>
      <Row alignItems=`initial>
        <Col size=1.>
          <VSpacing size=Spacing.md />
          <Text value="Latest Transactions" size=Text.Xl weight=Text.Bold block=true />
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
        </Col>
        <HSpacing size=Spacing.xl />
        <HSpacing size=Spacing.lg />
        <Col>
          <VSpacing size=Spacing.md />
          <Text value="Latest Blocks" size=Text.Xl weight=Text.Bold block=true />
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
