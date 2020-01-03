module Styles = {
  open Css;

  let logo = style([width(`px(40)), height(`px(40))]);

  let highlightsContainer =
    style([width(`percent(100.)), paddingTop(`px(40)), paddingBottom(Spacing.xl)]);

  let section = style([paddingTop(`px(48)), width(`percent(100.))]);
};

module Highlights = {
  open Belt.Option;

  module Styles = {
    open Css;

    let highlights = style([textAlign(`center)]);
  };

  [@react.component]
  let make = (~label, ~value, ~valuePrefixFunc=?, ~extra, ~extraSuffixFunc=?) => {
    <div className=Styles.highlights>
      <div> <Text value=label size=Text.Sm weight=Text.Bold color=Colors.purple /> </div>
      <div className={Css.style([Css.marginTop(Spacing.sm)])}>
        {valuePrefixFunc->map(v => v())->getWithDefault(React.string(""))}
        <Text value size=Text.Xxl weight=Text.Bold />
      </div>
      <div>
        <Text value=extra size=Text.Sm />
        {extraSuffixFunc->map(v => v())->getWithDefault(React.string(""))}
      </div>
    </div>;
  };
};

[@react.component]
let make = () => {
  <div className=Styles.highlightsContainer>
    <Row>
      <Col>
        <Highlights
          label="BAND PRICE"
          value="$0.642"
          extra="@0.012 BTC"
          extraSuffixFunc={() => <Text value="(+1.23%)" size=Text.Sm color=Colors.green />}
        />
      </Col>
      <Col size=3.>
        <Highlights label="MARKET CAP" value="$8,428,380.55" extra="12,356.012 BTC" />
      </Col>
      <HSpacing size=Spacing.xl />
      <Col> <img src=Images.bandLogo className=Styles.logo /> </Col>
      <Col size=3.>
        <Highlights
          label="LATEST BLOCK"
          valuePrefixFunc={() =>
            <Text value="# " size=Text.Xxl weight=Text.Bold color=Colors.pink />
          }
          value="472,395"
          extra="7 seconds ago"
        />
      </Col>
      <Col>
        <Highlights label="ACTIVE VALIDATORS" value="4 Nodes" extra="431,324.98 BAND Bonded" />
      </Col>
    </Row>
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
