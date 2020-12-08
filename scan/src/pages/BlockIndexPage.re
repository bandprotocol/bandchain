module Styles = {
  open Css;

  let proposerContainer = style([width(`fitContent)]);

  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding(`px(24)),
      Media.mobile([padding(`px(16))]),
    ]);

  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
};

[@react.component]
let make = (~height) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;
  let isMobile = Media.isMobile();

  let blockSub = BlockSub.get(height);
  let txsSub = TxSub.getListByBlockHeight(height, ~pageSize, ~page, ());

  <Section>
    <div className=CssHelper.container>
      <Row marginBottom=40 marginBottomSm=16>
        <Col>
          <Heading value="Block" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
          {switch (blockSub) {
           | Data({height}) => <TypeID.Block id=height position=TypeID.Title />
           | _ => <LoadingCensorBar width=100 height=15 />
           }}
        </Col>
      </Row>
      <Row marginBottom=24>
        <Col>
          <div className=Styles.infoContainer>
            <Heading value="Information" size=Heading.H4 style=Styles.infoHeader marginBottom=24 />
            <Row marginBottom=24>
              <Col>
                <Heading value="Block Hash" size=Heading.H5 />
                <VSpacing size={`px(8)} />
                {switch (blockSub) {
                 | Data({hash}) =>
                   <Text
                     value={hash |> Hash.toHex(~upper=true)}
                     code=true
                     block=true
                     size=Text.Lg
                     breakAll=true
                   />
                 | _ => <LoadingCensorBar width={isMobile ? 200 : 350} height=15 />
                 }}
              </Col>
            </Row>
            <Row marginBottom=24>
              <Col col=Col.Six mbSm=16>
                <Heading value="Transaction" size=Heading.H5 />
                <VSpacing size={`px(8)} />
                {switch (blockSub) {
                 | Data({txn}) => <Text value={txn |> string_of_int} size=Text.Lg />
                 | _ => <LoadingCensorBar width=40 height=15 />
                 }}
              </Col>
              <Col col=Col.Six>
                <Heading value="Timestamp" size=Heading.H5 />
                <VSpacing size={`px(8)} />
                {switch (blockSub) {
                 | Data({timestamp}) =>
                   <div className={CssHelper.flexBox()}>
                     <Text
                       value={
                         timestamp
                         |> MomentRe.Moment.format(Config.timestampDisplayFormat)
                         |> String.uppercase_ascii
                       }
                       size=Text.Lg
                       color=Colors.gray6
                     />
                     <HSpacing size=Spacing.sm />
                     <TimeAgos
                       time=timestamp
                       prefix="("
                       suffix=")"
                       size=Text.Md
                       weight=Text.Thin
                       color=Colors.gray8
                     />
                   </div>
                 | _ => <LoadingCensorBar width=200 height=15 />
                 }}
              </Col>
            </Row>
            <Row>
              <Col>
                <Heading value="Proposer" size=Heading.H5 />
                <VSpacing size={`px(8)} />
                {switch (blockSub) {
                 | Data({validator: {operatorAddress, moniker, identity}}) =>
                   <div className=Styles.proposerContainer>
                     <ValidatorMonikerLink validatorAddress=operatorAddress moniker identity />
                   </div>
                 | _ => <LoadingCensorBar width=200 height=15 />
                 }}
              </Col>
            </Row>
          </div>
        </Col>
      </Row>
      <BlockIndexTxsTable txsSub />
      {switch (blockSub) {
       | Data({txn}) =>
         let pageCount = Page.getPageCount(txn, pageSize);

         <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />;
       | _ => React.null
       }}
    </div>
  </Section>;
};
