module Styles = {
  open Css;

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let txContainer = style([width(`px(300)), cursor(`pointer)]);

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
let make = (~revisions: list(OracleScriptHook.OracleScript.revision_t)) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageCount = 10;

  let numRevision = revisions |> Belt_List.size;

  <div className=Styles.tableWrapper>
    <Row>
      <HSpacing size={`px(25)} />
      <Text value={numRevision |> string_of_int} weight=Text.Bold />
      <HSpacing size={`px(5)} />
      <Text value="Revisions" />
    </Row>
    <VSpacing size=Spacing.md />
    <VSpacing size=Spacing.sm />
    {numRevision > 0
       ? <>
           <THead>
             <Row>
               <Col> <HSpacing size=Spacing.md /> </Col>
               <Col size=3.>
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
               <Col size=2.>
                 <Text
                   block=true
                   value="AGE"
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
            ->Belt.List.map(({name, timestamp, height, txHash}) => {
                <TBody key={txHash |> Hash.toHex(~upper=true)}>
                  <Row>
                    <Col> <HSpacing size=Spacing.md /> </Col>
                    <Col size=3.>
                      <Text block=true value=name weight=Text.Medium color=Colors.gray7 />
                    </Col>
                    <Col size=2.>
                      <TimeAgos time=timestamp size=Text.Md weight=Text.Medium />
                    </Col>
                    <Col size=1.6> <TypeID.Block id={ID.Block.ID(height)} /> </Col>
                    <Col size=3.9>
                      <div
                        className=Styles.txContainer
                        onClick={_ => Route.redirect(Route.TxIndexPage(txHash))}>
                        <Text
                          block=true
                          value={txHash |> Hash.toHex(~upper=true)}
                          weight=Text.Medium
                          code=true
                          color=Colors.gray7
                          ellipsis=true
                          nowrap=true
                        />
                      </div>
                    </Col>
                    <Col> <HSpacing size=Spacing.lg /> </Col>
                  </Row>
                </TBody>
              })
            ->Array.of_list
            ->React.array}
           <VSpacing size=Spacing.lg />
           <VSpacing size=Spacing.lg />
           <Pagination
             currentPage=page
             pageCount
             onPageChange={newPage => setPage(_ => newPage)}
           />
         </>
       : <div className=Styles.iconWrapper>
           <VSpacing size={`px(30)} />
           <img src=Images.noRevisionIcon className=Styles.icon />
           <VSpacing size={`px(40)} />
           <Text block=true value="NO REVISION" weight=Text.Regular color=Colors.blue4 />
           <VSpacing size={`px(15)} />
         </div>}
  </div>;
};
