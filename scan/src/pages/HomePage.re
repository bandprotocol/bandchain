[@react.component]
let make = () => {
  // Subscribe for latest 11 blocks here so both "LatestBlocks" and "ChainInfoHighLights"
  // share the same infomation.
  let pageSize = Media.isMobile() ? 7 : 11;
  let lastest11BlocksSub = BlockSub.getList(~pageSize, ~page=1, ());
  let latestBlockSub = lastest11BlocksSub->Sub.map(blocks => blocks->Belt_Array.getExn(0));
  let ({ThemeContext.theme}, _) = React.useContext(ThemeContext.context);

  <>
    <Section pt=40 pb=40 ptSm=24 pbSm=24 bg={theme.mainBg}>
      <div className=CssHelper.container>
        <ChainInfoHighlights latestBlockSub />
        <Row>
          <Col col=Col.Six mbSm=24> <TotalRequestsGraph /> </Col>
          <Col col=Col.Six> <LatestRequests /> </Col>
        </Row>
      </div>
    </Section>
    <Section pt=40 pb=80 pbSm=40 bg=Colors.white>
      <div className=CssHelper.container>
        <Row>
          <Col col=Col.Four> <LatestBlocks blocksSub=lastest11BlocksSub /> </Col>
          <Col col=Col.Eight> <LatestTxTable /> </Col>
        </Row>
      </div>
    </Section>
  </>;
};
