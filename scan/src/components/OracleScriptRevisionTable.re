module Styles = {
  open Css;

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let nameContainer = style([width(`px(160))]);

  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);
};

[@react.component]
let make = (~id) =>
  {
    let numRevisionsSub = id |> OracleScriptRevisionSub.count;
    let revisionsSub = id |> OracleScriptRevisionSub.get;
    let%Sub numRevisions = numRevisionsSub;
    let%Sub revisions = revisionsSub;

    <div className=Styles.tableWrapper>
      <Row>
        <HSpacing size={`px(25)} />
        <Text value={numRevisions |> string_of_int} weight=Text.Bold />
        <HSpacing size={`px(5)} />
        <Text value="Revisions" />
      </Row>
      <VSpacing size=Spacing.md />
      <VSpacing size=Spacing.sm />
      {numRevisions > 0
         ? <>
             <THead>
               <Row>
                 <Col> <HSpacing size=Spacing.md /> </Col>
                 <Col size=2.5>
                   <div className=TElement.Styles.hashContainer>
                     <Text
                       block=true
                       value="NAME"
                       size=Text.Sm
                       weight=Text.Semibold
                       color=Colors.gray6
                     />
                   </div>
                 </Col>
                 <Col size=3.3>
                   <Text
                     block=true
                     value="TIMESTAMP"
                     size=Text.Sm
                     weight=Text.Semibold
                     color=Colors.gray6
                   />
                 </Col>
                 <Col size=1.6>
                   <Text
                     block=true
                     value="BLOCK"
                     size=Text.Sm
                     weight=Text.Semibold
                     color=Colors.gray6
                   />
                 </Col>
                 <Col size=3.9>
                   <Text
                     block=true
                     value="TX HASH"
                     size=Text.Sm
                     weight=Text.Semibold
                     color=Colors.gray6
                   />
                 </Col>
                 <Col> <HSpacing size=Spacing.lg /> </Col>
               </Row>
             </THead>
             {revisions
              ->Belt.Array.map(({name, transaction}) => {
                  <TBody
                    key={
                      switch (transaction) {
                      | Some(tx) => tx.hash |> Hash.toHex(~upper=true)
                      | None => "Genesis"
                      }
                    }>
                    <Row>
                      <Col> <HSpacing size=Spacing.md /> </Col>
                      <Col size=2.5>
                        <div className=Styles.nameContainer>
                          <Text
                            block=true
                            value=name
                            weight=Text.Medium
                            color=Colors.gray7
                            nowrap=true
                            ellipsis=true
                          />
                        </div>
                      </Col>
                      <Col size=3.3>
                        {switch (transaction) {
                         | Some(tx) =>
                           <Timestamp
                             time={tx.block.timestamp}
                             weight=Text.Regular
                             code=true
                             size=Text.Md
                             height={Text.Px(16)}
                           />
                         | None => <Text value="Genesis" weight=Text.Regular code=true />
                         }}
                      </Col>
                      <Col size=1.6>
                        {switch (transaction) {
                         | Some(tx) => <TypeID.Block id={tx.blockHeight} />
                         | None => <Text value="Genesis" />
                         }}
                      </Col>
                      <Col size=3.9>
                        {switch (transaction) {
                         | Some(tx) => <TxLink txHash={tx.hash} width=300 weight=Text.Medium />
                         | None =>
                           <Text
                             block=true
                             value="Genesis transaction"
                             weight=Text.Medium
                             code=true
                             color=Colors.gray7
                             ellipsis=true
                             nowrap=true
                           />
                         }}
                      </Col>
                      <Col> <HSpacing size=Spacing.lg /> </Col>
                    </Row>
                  </TBody>
                })
              ->React.array}
           </>
         : <div className=Styles.iconWrapper>
             <VSpacing size={`px(30)} />
             <img src=Images.noRevisionIcon className=Styles.icon />
             <VSpacing size={`px(40)} />
             <Text block=true value="NO REVISION" weight=Text.Regular color=Colors.blue4 />
             <VSpacing size={`px(15)} />
           </div>}
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
