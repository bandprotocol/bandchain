module Styles = {
  open Css;

  let highlightsContainer = style([width(`percent(100.)), paddingBottom(Spacing.xl)]);

  let section = style([paddingTop(`px(45)), width(`percent(100.))]);

  let topicBar =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`row),
      justifyContent(`spaceBetween),
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
    style([
      position(`absolute),
      maxWidth(`px(984)),
      marginLeft(`auto),
      marginRight(`auto),
      width(`percent(100.)),
    ]);
};

[@react.component]
let make = () => {
  // Subscribe for latest 11 blocks here so both "LatestBlocks" and "ChainInfoHighLights"
  // share the same infomation.
  let lastest11BlocksSub = BlockSub.getList(~pageSize=11, ~page=1, ());
  let latestBlockSub = lastest11BlocksSub->Sub.map(blocks => blocks->Belt_Array.getExn(0));

  <div className=Styles.highlightsContainer>
    <ChainInfoHighlights latestBlockSub />
    <VSpacing size=Spacing.xl />
    // TODO: for next version
    // <div className=Styles.grayArea>
    //   <div className=Styles.bg>
    //     <img src=Images.bg />
    //     <img src=Images.bg className=Styles.imgRight />
    //   </div>
    //   <div className=Styles.grayAreaInner> <DataScriptsHighlights /> </div>
    // </div>
    // <div className=Styles.skip />
    <div className=Styles.section>
      <Row alignItems=`initial>
        <Col> <LatestBlocks blocksSub=lastest11BlocksSub /> </Col>
        <HSpacing size=Spacing.lg />
        <Col size=1.> <LatestTxTable /> </Col>
      </Row>
    </div>
  </div>;
};
