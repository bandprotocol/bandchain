module Styles = {
  open Css;

  let nameContainer = style([width(`px(230))]);

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
    let numRevisionsSub = id |> DataSourceRevisionSub.count;
    let revisionsSub = id |> DataSourceRevisionSub.get;
    let%Sub numRevisions = numRevisionsSub;
    let%Sub revisions = revisionsSub;

    <div>
      {numRevisions > 0
         ? <>
             <THead>
               <Row alignItems=Row.Center>
                 <Col col=Col.Four>
                   <div className={CssHelper.flexBox()}>
                     <Text value={numRevisions |> string_of_int} weight=Text.Semibold />
                     <HSpacing size={`px(5)} />
                     <Text
                       value={numRevisions == 1 ? "Revision" : "Revisions"}
                       weight=Text.Semibold
                     />
                   </div>
                 </Col>
                 <Col col=Col.Three>
                   <Text block=true value="Timestamp" size=Text.Md weight=Text.Semibold />
                 </Col>
                 <Col col=Col.One>
                   <Text block=true value="Block" size=Text.Md weight=Text.Semibold />
                 </Col>
                 <Col col=Col.Four>
                   <Text block=true value="TX HASH" size=Text.Md weight=Text.Semibold />
                 </Col>
               </Row>
             </THead>
             {revisions
              ->Belt.Array.map(({name, transaction}) => {
                  <TBody
                    paddingH={`px(24)}
                    key={
                      switch (transaction) {
                      | Some(tx) => tx.hash |> Hash.toHex(~upper=true)
                      | None => "Genesis"
                      }
                    }>
                    <Row>
                      <Col> <HSpacing size=Spacing.lg /> </Col>
                      <Col col=Col.Four>
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
                      <Col col=Col.Three>
                        {switch (transaction) {
                         | Some(tx) =>
                           <TimeAgos time={tx.block.timestamp} size=Text.Md weight=Text.Medium />
                         | None => <Text value="Genesis" />
                         }}
                      </Col>
                      <Col col=Col.One>
                        {switch (transaction) {
                         | Some(tx) => <TypeID.Block id={tx.blockHeight} />
                         | None => <Text value="Genesis" />
                         }}
                      </Col>
                      <Col col=Col.Four>
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
