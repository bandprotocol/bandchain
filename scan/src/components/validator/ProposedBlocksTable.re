module Styles = {
  open Css;

  let vFlex = align => style([display(`flex), flexDirection(`row), alignItems(align)]);

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let fixWidth = style([width(`px(230))]);

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
let make = () => {
  let blocks = BlockSub.getList(~pageSize=5, ~page=1, ())->Sub.default([||])->Belt_List.fromArray;

  let numProposedBlocks = 241;

  <div className=Styles.tableWrapper>
    <Row>
      <HSpacing size={`px(25)} />
      <Text value={numProposedBlocks |> string_of_int} weight=Text.Bold />
      <HSpacing size={`px(5)} />
      <Text value="Blocks" />
    </Row>
    <VSpacing size=Spacing.lg />
    {blocks->Belt_List.length > 0
       ? <>
           <THead>
             <Row>
               <Col> <HSpacing size=Spacing.lg /> </Col>
               <Col size=1.0>
                 <Text block=true value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
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
                   <Text block=true value="TXN" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
                 </div>
               </Col>
               <Col> <HSpacing size=Spacing.lg /> </Col>
             </Row>
           </THead>
           {blocks
            ->Belt.List.map(({height, timestamp, hash, txn}) => {
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
            ->Array.of_list
            ->React.array}
         </>
       : <div className=Styles.iconWrapper>
           <VSpacing size={`px(30)} />
           <img src=Images.noRequestIcon className=Styles.icon />
           <VSpacing size={`px(40)} />
           <Text block=true value="NO BLOCK" weight=Text.Regular color=Colors.blue4 />
           <VSpacing size={`px(15)} />
         </div>}
  </div>;
};
