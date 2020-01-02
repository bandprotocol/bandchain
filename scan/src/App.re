[@bs.val] external require: string => string = "";
let logoSrc = require("./bandlogo.png");

module Styles = {
  open Css;

  let container =
    style([
      width(`percent(100.)),
      height(`percent(100.)),
      position(`relative),
    ]);

  let innerContainer =
    style([
      maxWidth(`px(1100)),
      marginLeft(`auto),
      marginRight(`auto),
      paddingLeft(Spacing.md),
      paddingRight(Spacing.md),
    ]);

  /* Nav links */
  let nav = style([paddingLeft(Spacing.md)]);
  let navContainer =
    style([
      paddingTop(Spacing.md),
      paddingBottom(Spacing.md),
      maxWidth(`px(1100)),
      marginLeft(`auto),
      marginRight(`auto),
      paddingLeft(Spacing.md),
      paddingRight(Spacing.md),
    ]);

  /* Main bar */
  let logo = style([width(`px(40)), height(`px(40))]);
  let mainBar = style([display(`flex), paddingTop(Spacing.lg)]);
  let version =
    style([
      display(`flex),
      borderRadius(`px(10)),
      backgroundColor(Colors.pinkLight),
      padding4(
        ~top=`px(0),
        ~bottom=`px(0),
        ~left=Spacing.sm,
        ~right=Spacing.sm,
      ),
      height(`px(20)),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(Spacing.xs),
      marginTop(`px(1)),
    ]);

  let uFlex = style([display(`flex), flexDirection(`row)]);

  let highlightsContainer =
    style([
      width(`percent(100.)),
      paddingTop(`px(40)),
      paddingBottom(Spacing.xl),
    ]);

  let section = style([paddingTop(`px(48)), width(`percent(100.))]);

  let hr = style([border(`px(1), `dashed, Css.hex("eeeeee"))]);

  let bg =
    style([
      width(`percent(100.)),
      height(`px(300)),
      left(`px(0)),
      bottom(`px(0)),
      position(`relative),
      background(hex("F6F3FA")),
      before([
        position(`absolute),
        contentRule(""),
        background(hex("F6F3FA")),
        width(`percent(100.)),
        height(`px(300)),
        transform(`skewY(`deg(6.))),
        zIndex(-1),
        top(`px(-150)),
      ]),
    ]);
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
      <div>
        <Text value=label size=Text.Sm weight=Text.Bold color=Colors.purple />
      </div>
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
  <div className=Styles.container>
    <div className=Styles.navContainer>
      <Row>
        <Col size=1.>
          <Text color=Colors.grayText value="Made with <3 by Band Protocol" />
        </Col>
        <Col>
          <Row justify=Row.Right>
            {[
               "Validators",
               "Blocks",
               "Transactions",
               "Request Scripts",
               "Data Providers",
               "OWASM Studio",
             ]
             ->Belt.List.map(v =>
                 <Col key=v>
                   <div className=Styles.nav>
                     <Text color=Colors.grayText value=v nowrap=true />
                   </div>
                 </Col>
               )
             ->Array.of_list
             ->React.array}
          </Row>
        </Col>
      </Row>
    </div>
    <div className=Styles.innerContainer>
      <div className=Styles.mainBar>
        <Row>
          <Col>
            <div className=Styles.uFlex>
              <Text value="D3N" size=Text.Xxl weight=Text.Bold nowrap=true />
              <div className=Styles.version>
                <Text
                  value="v1.0 TESTNET"
                  size=Text.Sm
                  color=Colors.pink
                  nowrap=true
                />
              </div>
            </div>
            <Text value="Data Request Explorer" nowrap=true />
          </Col>
        </Row>
        <SearchBar />
      </div>
      <div className=Styles.highlightsContainer>
        <Row>
          <Col>
            <Highlights
              label="BAND PRICE"
              value="$0.642"
              extra="@0.012 BTC"
              extraSuffixFunc={() =>
                <Text value="(+1.23%)" size=Text.Sm color=Colors.green />
              }
            />
          </Col>
          <Col size=3.>
            <Highlights
              label="MARKET CAP"
              value="$8,428,380.55"
              extra="12,356.012 BTC"
            />
          </Col>
          <HSpacing size=Spacing.xl />
          <Col> <img src=logoSrc className=Styles.logo /> </Col>
          <Col size=3.>
            <Highlights
              label="LATEST BLOCK"
              valuePrefixFunc={() =>
                <Text
                  value="# "
                  size=Text.Xxl
                  weight=Text.Bold
                  color=Colors.pink
                />
              }
              value="472,395"
              extra="7 seconds ago"
            />
          </Col>
          <Col>
            <Highlights
              label="ACTIVE VALIDATORS"
              value="4 Nodes"
              extra="431,324.98 BAND Bonded"
            />
          </Col>
        </Row>
        <VSpacing size=Spacing.xl />
        <VSpacing size=Spacing.lg />
        <DataScriptsHighlights />
        <div className=Styles.section>
          <Row alignItems=`initial>
            <Col size=1.>
              <VSpacing size=Spacing.md />
              <Text
                value="Latest Transactions"
                size=Text.Xl
                weight=Text.Bold
                block=true
              />
              <VSpacing size=Spacing.lg />
              {Transaction.renderHeader()}
              <Transaction
                type_={Transaction.DataRequest("ETH/USD Price Feed")}
                hash="0x128f12db1a99dce2937"
                fee="0.10 BAND"
                timestamp="2 days ago"
              />
              <Transaction
                type_={
                  Transaction.NewScript(
                    "Anime Episodes Ranking - WINTER 2020",
                  )
                }
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
              <Text
                value="Latest Blocks"
                size=Text.Xl
                weight=Text.Bold
                block=true
              />
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
      </div>
    </div>
    <div className=Styles.bg />
  </div>;
};