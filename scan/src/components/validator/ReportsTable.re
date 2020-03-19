module Styles = {
  open Css;

  let vFlex = align => style([display(`flex), flexDirection(`row), alignItems(align)]);

  let tableWrapper = style([padding2(~v=`px(9), ~h=`px(15))]);

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

  let hFlex = style([display(`flex), alignItems(`center)]);
};

[@react.component]
let make = () => {
  let (page, setPage) = React.useState(_ => 1);

  // TODO: Mock to use
  let pageCount = 1;

  let reports = [
    (
      324,
      "6F45B0D19B46F144CDD7ACA9674E2AD8E8F8C15EF56CA073749B2ACD7DF7739D",
      234,
      "Mean Crypto Price",
      [23, 12, 35],
      [1, 2, 3],
      [123213123, 123123132, 123123123],
    ),
    (
      324,
      "6F45B0D19B46F144CDD7ACA9674E2AD8E8F8C15EF56CA073749B2ACD7DF7739D",
      234,
      "Mean Crypto Price",
      [23, 12, 35],
      [1, 2, 3],
      [123213123, 123123132, 123123123],
    ),
    (
      324,
      "6F45B0D19B46F144CDD7ACA9674E2AD8E8F8C15EF56CA073749B2ACD7DF7739D",
      234,
      "Mean Crypto Price",
      [23, 12, 35],
      [1, 2, 3],
      [123213123, 123123132, 123123123],
    ),
  ];
  <div className=Styles.tableWrapper>
    <Row>
      <HSpacing size={`px(25)} />
      <Text value={reports |> Belt_List.length |> string_of_int} weight=Text.Bold />
      <HSpacing size={`px(5)} />
      <Text value="Reports" />
    </Row>
    <VSpacing size=Spacing.lg />
    {reports->Belt_List.length > 0
       ? <>
           <THead>
             <Row>
               <Col> <HSpacing size=Spacing.md /> </Col>
               <Col size=1.>
                 <Text
                   block=true
                   value="REQUEST"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </Col>
               <Col size=2.>
                 <Text
                   block=true
                   value="TX HASH"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </Col>
               <Col size=2.>
                 <Text
                   block=true
                   value="ORACLE SCRIPT"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </Col>
               <Col size=1.>
                 <Text
                   block=true
                   value="DATA SOURCE"
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.05)}
                 />
               </Col>
               <Col size=1.5>
                 <div className={Styles.vFlex(`flexEnd)}>
                   <div className=Styles.fillLeft />
                   <Text
                     block=true
                     value="EXTERNAL ID"
                     size=Text.Sm
                     weight=Text.Semibold
                     color=Colors.gray6
                     spacing={Text.Em(0.05)}
                   />
                 </div>
               </Col>
               <Col size=2.2>
                 <div className={Styles.vFlex(`flexEnd)}>
                   <div className=Styles.fillLeft />
                   <Text
                     block=true
                     value="VALUE"
                     size=Text.Sm
                     weight=Text.Semibold
                     color=Colors.gray6
                     spacing={Text.Em(0.05)}
                   />
                 </div>
               </Col>
               <Col> <HSpacing size=Spacing.lg /> </Col>
             </Row>
           </THead>
           {reports
            ->Belt.List.map(
                (
                  (
                    requestID,
                    hash,
                    oracleScriptID,
                    oracleScriptDescription,
                    dataSourceIDs,
                    externalDataIDs,
                    values,
                  ),
                ) => {
                <TBody key=hash>
                  <Row>
                    <Col> <HSpacing size=Spacing.md /> </Col>
                    <Col size=1. alignSelf=Col.Start>
                      <TypeID.Request id={ID.Request.ID(requestID)} />
                    </Col>
                    <Col size=2. alignSelf=Col.Start>
                      <div className={Styles.withWidth(140)}>
                        <Text value=hash block=true code=true ellipsis=true />
                      </div>
                    </Col>
                    <Col size=2. alignSelf=Col.Start>
                      <Row>
                        <TypeID.OracleScript id={ID.OracleScript.ID(oracleScriptID)} />
                        <HSpacing size=Spacing.sm />
                        <HSpacing size=Spacing.xs />
                        <Text value=oracleScriptDescription block=true code=true />
                      </Row>
                    </Col>
                    <Col size=1.>
                      {dataSourceIDs
                       ->Belt_List.map(dataSourceID => {
                           <>
                             <Row>
                               <TypeID.DataSource id={ID.DataSource.ID(dataSourceID)} />
                             </Row>
                             <VSpacing size=Spacing.sm />
                             <VSpacing size=Spacing.xs />
                           </>
                         })
                       ->Belt.List.toArray
                       ->React.array}
                    </Col>
                    <Col size=1.5>
                      {externalDataIDs
                       ->Belt_List.map(externalDataID => {
                           <>
                             <div className={Styles.vFlex(`flexEnd)}>
                               <Row>
                                 <div className=Styles.fillLeft />
                                 <Text
                                   value={externalDataID |> string_of_int}
                                   block=true
                                   code=true
                                 />
                               </Row>
                             </div>
                             <VSpacing size=Spacing.md />
                           </>
                         })
                       ->Belt.List.toArray
                       ->React.array}
                    </Col>
                    <Col size=2.2>
                      {values
                       ->Belt_List.map(value => {
                           <>
                             <div className={Styles.vFlex(`flexEnd)}>
                               <Row>
                                 <div className=Styles.fillLeft />
                                 <Text value={value |> string_of_int} block=true code=true />
                               </Row>
                             </div>
                             <VSpacing size=Spacing.md />
                           </>
                         })
                       ->Belt.List.toArray
                       ->React.array}
                    </Col>
                    <Col> <HSpacing size=Spacing.lg /> </Col>
                  </Row>
                </TBody>
              })
            ->Array.of_list
            ->React.array}
           //  <Pagination currentPage="1" pageCount onPageChange={newPage => setPage(_ => newPage)} />
         </>
       : <div className=Styles.iconWrapper>
           <VSpacing size={`px(30)} />
           <img src=Images.noRequestIcon className=Styles.icon />
           <VSpacing size={`px(40)} />
           <Text block=true value="NO REPORTS" weight=Text.Regular color=Colors.blue4 />
           <VSpacing size={`px(15)} />
         </div>}
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.sm />
    <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
  </div>;
};
