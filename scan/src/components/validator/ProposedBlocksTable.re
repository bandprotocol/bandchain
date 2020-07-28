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
  let pagination = style([height(`px(30))]);
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
  let make = () => {
    <>
      {Belt_Array.make(
         10,
         <Row>
           <Col> <HSpacing size=Spacing.lg /> </Col>
           <Col size=1.5> <LoadingCensorBar width=50 height=20 /> </Col>
           <Col size=4.0> <LoadingCensorBar width=150 height=20 /> </Col>
           <Col size=3.0>
             <div className={Styles.withWidth(500)}>
               <LoadingCensorBar width=460 height=20 />
             </div>
           </Col>
           <Col size=1.5>
             <Row> <div className=Styles.fillLeft /> <LoadingCensorBar width=50 height=20 /> </Row>
           </Col>
           <Col> <HSpacing size=Spacing.lg /> </Col>
         </Row>,
       )
       ->Belt.Array.mapWithIndex((i, e) => {<TBody key={i |> string_of_int}> e </TBody>})
       ->React.array}
    </>;
  };
};

module LoadingWithHeader = {
  [@react.component]
  let make = () => {
    <div className=Styles.tableWrapper> <TableHeader /> <Loading /> </div>;
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
let make = (~consensusAddress) => {
  let blocksSub =
    BlockSub.getListByConsensusAddress(~address=consensusAddress, ~pageSize=10, ~page=1, ());

  <div className=Styles.tableWrapper>
    <TableHeader />
    {switch (blocksSub) {
     | Data(blocks) =>
       blocks->Belt.Array.size > 0
         ? <> <BlocksTable blocks /> </>
         : <div className=Styles.iconWrapper>
             <VSpacing size={`px(30)} />
             <img src=Images.noRequestIcon className=Styles.icon />
             <VSpacing size={`px(40)} />
             <Text block=true value="NO BLOCK" weight=Text.Regular color=Colors.blue4 />
             <VSpacing size={`px(15)} />
           </div>

     | _ => <Loading />
     }}
  </div>;
};
