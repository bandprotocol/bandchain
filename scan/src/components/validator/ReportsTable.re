module Styles = {
  open Css;

  let vFlex = align => style([display(`flex), flexDirection(`row), alignItems(align)]);

  let tableWrapper = style([padding2(~v=`px(9), ~h=`px(15))]);

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
let make = (~address) =>
  // TODO: Mockssssssssssss
  {
    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 5;

    let reportsSub =
      ReportSub.ValidatorReport.getListByValidator(
        ~page,
        ~pageSize,
        ~validator={
          address |> Address.toOperatorBech32;
        },
      );
    let totalReportsSub = ReportSub.ValidatorReport.count();

    let%Sub totalReports = totalReportsSub;
    let%Sub reports = reportsSub;

    let pageCount = Page.getPageCount(totalReports, pageSize);

    <div className=Styles.tableWrapper>
      <Row>
        <HSpacing size={`px(25)} />
        <Text value={reports |> Belt_Array.length |> string_of_int} weight=Text.Bold />
        <HSpacing size={`px(5)} />
        <Text value="Reports" />
      </Row>
      <VSpacing size=Spacing.lg />
      {reports->Belt_Array.length > 0
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
                 <Col size=2.3>
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
              ->Belt.Array.map(({txHash, request, reportDetails}) => {
                  <TBody key={txHash |> Hash.toBase64}>
                    <Row>
                      <Col> <HSpacing size=Spacing.md /> </Col>
                      <Col size=1. alignSelf=Col.Start> <TypeID.Request id={request.id} /> </Col>
                      <Col size=2. alignSelf=Col.Start>
                        // TODO: Check TXHASH STYLINGS

                          <div className={Styles.withWidth(140)}>
                            <TxLink txHash width=110 />
                          </div>
                        </Col>
                      <Col size=2.3 alignSelf=Col.Start>
                        <Row>
                          <TypeID.OracleScript id={request.oracleScript.id} />
                          <HSpacing size=Spacing.sm />
                          <HSpacing size=Spacing.xs />
                          <div className={Styles.withWidth(140)}>
                            <Text
                              value={request.oracleScript.name}
                              block=true
                              code=true
                              ellipsis=true
                            />
                          </div>
                        </Row>
                      </Col>
                      <Col size=1.>
                        {reportDetails
                         ->Belt_Array.map(({dataSourceID}) => dataSourceID)
                         ->Belt_Array.map(dataSourceID => {
                             <>
                               <Row> <TypeID.DataSource id=dataSourceID /> </Row>
                               <VSpacing size=Spacing.sm />
                               <VSpacing size=Spacing.xs />
                             </>
                           })
                         ->React.array}
                      </Col>
                      <Col size=1.5>
                        {reportDetails
                         ->Belt_Array.map(({externalID}) => externalID)
                         ->Belt_Array.map(externalDataID => {
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
                         ->React.array}
                      </Col>
                      <Col size=2.2>
                        {reportDetails
                         ->Belt_Array.map(({data}) => data)
                         ->Belt_Array.map(value => {
                             <>
                               <div className={Styles.vFlex(`flexEnd)}>
                                 <Row>
                                   <div className=Styles.fillLeft />
                                   <Text value={value |> JsBuffer.toUTF8} block=true code=true />
                                 </Row>
                               </div>
                               <VSpacing size=Spacing.md />
                             </>
                           })
                         ->React.array}
                      </Col>
                      <Col> <HSpacing size=Spacing.lg /> </Col>
                    </Row>
                  </TBody>
                })
              ->React.array}
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
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
