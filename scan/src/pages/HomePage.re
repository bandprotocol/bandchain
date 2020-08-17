[@react.component]
let make = () => {
  // Subscribe for latest 11 blocks here so both "LatestBlocks" and "ChainInfoHighLights"
  // share the same infomation.
  let pageSize = Media.isMobile() ? 7 : 11;
  let lastest11BlocksSub = BlockSub.getList(~pageSize, ~page=1, ());
  let latestBlockSub = lastest11BlocksSub->Sub.map(blocks => blocks->Belt_Array.getExn(0));

  <>
    <Section bg=Colors.highlightBg ptSm=0>
      <div className=CssHelper.container> <ChainInfoHighlights latestBlockSub /> </div>
    </Section>
    <Section pt=40 pb=40 ptSm=24 pbSm=24 bg=Colors.bg>
      <div className=CssHelper.container>
        <Row.Grid>
          <Col.Grid col=Col.Six> <TotalRequestsGraph /> </Col.Grid>
          <Col.Grid col=Col.Six> <LatestRequests /> </Col.Grid>
        </Row.Grid>
      </div>
    </Section>
    <Section>
      <div className=CssHelper.container>
        <Row alignItems=`initial wrap=true>
          <Col size={Media.isMobile() ? 1. : 0.}>
            <LatestBlocks blocksSub=lastest11BlocksSub />
          </Col>
          <Col size=1.> <LatestTxTable /> </Col>
        </Row>
      </div>
    </Section>
  </>;
};
