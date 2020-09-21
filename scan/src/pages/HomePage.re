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
          <Col.Grid col=Col.Six mbSm=24> <TotalRequestsGraph /> </Col.Grid>
          <Col.Grid col=Col.Six> <LatestRequests /> </Col.Grid>
        </Row.Grid>
      </div>
    </Section>
    <Section pt=40 pb=80 pbSm=40 bg=Colors.white>
      <div className=CssHelper.container>
        <Row.Grid>
          <Col.Grid col=Col.Four> <LatestBlocks blocksSub=lastest11BlocksSub /> </Col.Grid>
          <Col.Grid col=Col.Eight> <LatestTxTable /> </Col.Grid>
        </Row.Grid>
      </div>
    </Section>
  </>;
};
