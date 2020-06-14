module Styles = {
  open Css;

  let vFlex = align => style([display(`flex), flexDirection(`row), alignItems(align)]);

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);

  let withWidth = w => style([width(`px(w))]);

  let fillLeft = style([marginLeft(`auto)]);
  let pagination = style([height(`px(50))]);
};

module TableHeader = {
  [@react.component]
  let make = () => {
    <THead>
      <Row>
        <Col> <HSpacing size=Spacing.lg /> </Col>
        <Col size=1.0>
          <Text block=true value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
        </Col>
        <Col size=2.75>
          <Text block=true value="TIMESTAMP" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
        </Col>
        <Col size=6.0>
          <Text block=true value="BLOCK HASH" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
        </Col>
        <Col size=1.5>
          <div className={Styles.vFlex(`flexEnd)}>
            <div className=Styles.fillLeft />
            <Text block=true value="TXN" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <Col> <HSpacing size=Spacing.lg /> </Col>
      </Row>
    </THead>;
  };
};

module Loading = {
  [@react.component]
  let make = (~withCount=true) => {
    <div>
      {withCount
         ? <>
             <Row> <LoadingCensorBar width=100 height=15 /> </Row>
             <VSpacing size=Spacing.lg />
           </>
         : React.null}
      {Belt_Array.make(
         5,
         <Row>
           <Col> <HSpacing size=Spacing.lg /> </Col>
           <Col size=1.5> <LoadingCensorBar width=75 height=15 /> </Col>
           <Col size=4.0> <LoadingCensorBar width=200 height=15 /> </Col>
           <Col size=3.0>
             <div className={Styles.withWidth(500)}>
               <LoadingCensorBar width=500 height=15 />
             </div>
           </Col>
           <Col size=1.5>
             <Row> <div className=Styles.fillLeft /> <LoadingCensorBar width=50 height=15 /> </Row>
           </Col>
           <Col> <HSpacing size=Spacing.lg /> </Col>
         </Row>,
       )
       ->Belt.Array.mapWithIndex((i, e) => {<TBody key={i |> string_of_int}> e </TBody>})
       ->React.array}
      <VSpacing size=Spacing.lg />
      <div className=Styles.pagination />
    </div>;
  };
};

module BlocksTable = {
  [@react.component]
  let make = (~blocks: array(BlockSub.t)) => {
    blocks
    ->Belt.Array.map(({height, timestamp, hash, txn}) => {
        <TBody key={hash |> Hash.toHex(~upper=true)}>
          <Row>
            <Col> <HSpacing size=Spacing.lg /> </Col>
            <Col size=1.5> <TypeID.Block id=height /> </Col>
            <Col size=4.0> <Timestamp time=timestamp code=true size=Text.Md /> </Col>
            <Col size=3.0>
              <div className={Styles.withWidth(500)}>
                <Text
                  value={hash |> Hash.toHex(~upper=true)}
                  block=true
                  code=true
                  ellipsis=true
                />
              </div>
            </Col>
            <Col size=1.5>
              <Row>
                <div className=Styles.fillLeft />
                <Text value={txn |> Format.iPretty} code=true />
              </Row>
            </Col>
            <Col> <HSpacing size=Spacing.lg /> </Col>
          </Row>
        </TBody>
      })
    ->React.array;
  };
};

module ProposedBlockCount = {
  [@react.component]
  let make = (~consensusAddress) => {
    let blocksCountSub = BlockSub.countByConsensusAddress(~address=consensusAddress, ());
    <Row>
      {switch (blocksCountSub) {
       | Data(blocksCount) =>
         <>
           <HSpacing size={`px(25)} />
           <Text value={blocksCount |> string_of_int} weight=Text.Bold />
           <HSpacing size={`px(5)} />
           <Text value="Blocks" />
         </>
       | _ => <LoadingCensorBar width=100 height=15 />
       }}
    </Row>;
  };
};

[@react.component]
let make = (~consensusAddress) =>
  {
    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 5;

    let blocksCountSub = BlockSub.countByConsensusAddress(~address=consensusAddress, ());
    let blocksSub =
      BlockSub.getListByConsensusAddress(~address=consensusAddress, ~pageSize, ~page, ());

    let%Sub blocksCount = blocksCountSub;
    let pageCount = Page.getPageCount(blocksCount, pageSize);

    <div className=Styles.tableWrapper>
      <ProposedBlockCount consensusAddress />
      <VSpacing size=Spacing.lg />
      <TableHeader />
      {switch (blocksSub) {
       | Data(blocks) =>
         blocks->Belt.Array.size > 0
           ? <>
               <BlocksTable blocks />
               <VSpacing size=Spacing.lg />
               <div className=Styles.pagination>
                 <Pagination
                   currentPage=page
                   pageCount
                   onPageChange={newPage => setPage(_ => newPage)}
                 />
               </div>
             </>
           : <div className=Styles.iconWrapper>
               <VSpacing size={`px(30)} />
               <img src=Images.noRequestIcon className=Styles.icon />
               <VSpacing size={`px(40)} />
               <Text block=true value="NO BLOCK" weight=Text.Regular color=Colors.blue4 />
               <VSpacing size={`px(15)} />
             </div>
       | _ => <Loading withCount=false />
       }}
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
