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
    let%Sub blocks = blocksSub;

    let pageCount = Page.getPageCount(blocksCount, pageSize);

    <div className=Styles.tableWrapper>
      <Row>
        <HSpacing size={`px(25)} />
        <Text value={blocksCount |> string_of_int} weight=Text.Bold />
        <HSpacing size={`px(5)} />
        <Text value="Blocks" />
      </Row>
      <VSpacing size=Spacing.lg />
      {blocks->Belt.Array.size > 0
         ? <>
             <THead>
               <Row>
                 <Col> <HSpacing size=Spacing.lg /> </Col>
                 <Col size=1.0>
                   <Text
                     block=true
                     value="BLOCK"
                     size=Text.Sm
                     weight=Text.Bold
                     color=Colors.gray6
                   />
                 </Col>
                 <Col size=2.75>
                   <Text
                     block=true
                     value="TIMESTAMP"
                     size=Text.Sm
                     weight=Text.Bold
                     color=Colors.gray6
                   />
                 </Col>
                 <Col size=6.0>
                   <Text
                     block=true
                     value="BLOCK HASH"
                     size=Text.Sm
                     weight=Text.Bold
                     color=Colors.gray6
                   />
                 </Col>
                 <Col size=1.5>
                   <div className={Styles.vFlex(`flexEnd)}>
                     <div className=Styles.fillLeft />
                     <Text
                       block=true
                       value="TXN"
                       size=Text.Sm
                       weight=Text.Bold
                       color=Colors.gray6
                     />
                   </div>
                 </Col>
                 <Col> <HSpacing size=Spacing.lg /> </Col>
               </Row>
             </THead>
             {blocks
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
              ->React.array}
             <VSpacing size=Spacing.lg />
             <Pagination
               currentPage=page
               pageCount
               onPageChange={newPage => setPage(_ => newPage)}
             />
           </>
         : <div className=Styles.iconWrapper>
             <VSpacing size={`px(30)} />
             <img src=Images.noRequestIcon className=Styles.icon />
             <VSpacing size={`px(40)} />
             <Text block=true value="NO BLOCK" weight=Text.Regular color=Colors.blue4 />
             <VSpacing size={`px(15)} />
           </div>}
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
