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

  let seeAllContainer =
    style([alignItems(`center), justifyContent(`center), display(`flex), cursor(`pointer)]);

  let rightArrow = style([width(`px(13)), marginLeft(`px(5))]);

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
    style([
      position(`absolute),
      maxWidth(`px(984)),
      marginLeft(`auto),
      marginRight(`auto),
      width(`percent(100.)),
    ]);
};

/* SEE ALL btn */
let renderSeeAll = route =>
  <div className=Styles.seeAllContainer onClick={_ => Route.redirect(route)}>
    <Text block=true value="SEE ALL" size=Text.Sm weight=Text.Bold color=Colors.gray5 />
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
            {renderSeeAll(TxHomePage)}
          </div>
          <VSpacing size=Spacing.lg />
          <LatestTxTable />
        </Col>
        <HSpacing size=Spacing.xl />
        <HSpacing size=Spacing.lg />
        <Col>
          <VSpacing size=Spacing.md />
          <div className=Styles.topicBar>
            <Text value="Latest Blocks" size=Text.Xl weight=Text.Bold block=true />
            {renderSeeAll(BlockHomePage)}
          </div>
          <VSpacing size=Spacing.md />
          <LatestBlocks />
        </Col>
      </Row>
    </div>
  </div>;
};
