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
  let blocksOpt = BlockHook.latest(~limit=5, ());
  let blocks = blocksOpt->Belt.Option.getWithDefault([]);

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
               <Col size=0.75>
                 <div className=TElement.Styles.hashContainer>
                   <Text
                     block=true
                     value="BLOCK"
                     size=Text.Sm
                     weight=Text.Bold
                     color=Colors.gray6
                   />
                 </div>
               </Col>
               <Col size=1.>
                 <Text block=true value="AGE" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
               </Col>
               <Col size=4.>
                 <Text
                   block=true
                   value="BLOCK HASH"
                   size=Text.Sm
                   weight=Text.Bold
                   color=Colors.gray6
                 />
               </Col>
               <Col size=1.>
                 <div className={Styles.vFlex(`flexEnd)}>
                   <div className=Styles.fillLeft />
                   <Text block=true value="TXN" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
                 </div>
               </Col>
               <Col> <HSpacing size=Spacing.lg /> </Col>
             </Row>
           </THead>
           {blocks
            ->Belt.List.map(({height, timestamp, hash, numTxs}) => {
                <TBody key={hash |> Hash.toHex(~upper=true)}>
                  <Row>
                    <Col> <HSpacing size=Spacing.lg /> </Col>
                    <Col size=0.75> <TypeID.Block id={ID.Block.ID(height)} /> </Col>
                    <Col size=1.>
                      <TimeAgos time=timestamp size=Text.Md weight=Text.Medium />
                    </Col>
                    <Col size=4.>
                      <div className={Styles.withWidth(500)}>
                        <Text
                          value={hash |> Hash.toHex(~upper=true)}
                          weight=Text.Medium
                          block=true
                          code=true
                          ellipsis=true
                        />
                      </div>
                    </Col>
                    <Col size=1.>
                      <Row>
                        <div className=Styles.fillLeft />
                        <Text value={numTxs |> Format.iPretty} code=true weight=Text.Medium />
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
