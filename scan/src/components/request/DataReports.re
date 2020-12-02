module Styles = {
  open Css;

  let reportCard = isLast =>
    style([
      padding(`px(24)),
      isLast ? borderBottomStyle(`none) : borderBottom(`px(1), `solid, Colors.gray9),
      Media.mobile([padding2(~v=`px(24), ~h=`px(16))]),
    ]);

  let reportsTable =
    style([
      padding2(~v=`px(16), ~h=`px(24)),
      paddingBottom(`px(1)),
      marginTop(`px(24)),
      backgroundColor(Colors.profileBG),
      transition(~duration=200, "all"),
      height(`auto),
      Media.mobile([padding(`zero), backgroundColor(Colors.white)]),
    ]);

  let mobileCard =
    style([backgroundColor(Colors.profileBG), boxShadow(`none), marginTop(`px(8))]);

  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

[@react.component]
let make = (~reports: array(RequestSub.report_t)) => {
  let isMobile = Media.isMobile();

  reports->Belt.Array.size > 0
    ? reports
      ->Belt.Array.mapWithIndex(
          (
            idx,
            {
              reportValidator: {operatorAddress, moniker, identity},
              transactionOpt,
              reportDetails,
            },
          ) => {
          <div
            key={operatorAddress |> Address.toOperatorBech32}
            className={Styles.reportCard(idx == reports->Belt.Array.size - 1)}>
            <Row marginBottom=24>
              <Col.Grid col=Col.Six mbSm=24>
                <Heading value="Report by" size=Heading.H5 />
                <VSpacing size={`px(8)} />
                <ValidatorMonikerLink
                  validatorAddress=operatorAddress
                  moniker
                  identity
                  width={`percent(100.)}
                  avatarWidth=20
                />
              </Col.Grid>
              <Col.Grid col=Col.Six>
                <Heading value="TX Hash" size=Heading.H5 />
                <VSpacing size={`px(8)} />
                {switch (transactionOpt) {
                 | Some({hash}) => <TxLink txHash=hash width=280 />
                 | None => <Text value="Genesis Transaction" />
                 }}
              </Col.Grid>
            </Row>
            <div className=Styles.reportsTable>
              {isMobile
                 ? React.null
                 : <Row alignItems=Row.Center marginBottom=16>
                     <Col.Grid col=Col.Three>
                       <Text value="External ID" weight=Text.Medium />
                     </Col.Grid>
                     <Col.Grid col=Col.Three>
                       <Text value="Exit Code" weight=Text.Medium />
                     </Col.Grid>
                     <Col.Grid col=Col.Six> <Text value="Value" weight=Text.Medium /> </Col.Grid>
                   </Row>}
              {reportDetails
               ->Belt.Array.map(({externalID, exitCode, data}) => {
                   isMobile
                     ? <MobileCard
                         values=InfoMobileCard.[
                           ("External ID", Text(externalID)),
                           ("Exit Code", Text(exitCode)),
                           ("Value", Text(data |> JsBuffer.toUTF8)),
                         ]
                         key={externalID ++ exitCode}
                         idx={externalID ++ exitCode}
                         styles=Styles.mobileCard
                       />
                     : <Row alignItems=Row.Start marginBottom=16 key=externalID>
                         <Col.Grid col=Col.Three>
                           <Text value=externalID weight=Text.Medium />
                         </Col.Grid>
                         <Col.Grid col=Col.Three>
                           <Text value=exitCode weight=Text.Medium />
                         </Col.Grid>
                         <Col.Grid col=Col.Six>
                           <Text
                             value={data |> JsBuffer.toUTF8}
                             weight=Text.Medium
                             breakAll=true
                           />
                         </Col.Grid>
                       </Row>
                 })
               ->React.array}
            </div>
          </div>
        })
      ->React.array
    : <EmptyContainer height={`px(250)} backgroundColor=Colors.blueGray1>
        <img src=Images.noSource className=Styles.noDataImage />
        <Heading
          size=Heading.H4
          value="No Report"
          align=Heading.Center
          weight=Heading.Regular
          color=Colors.bandBlue
        />
      </EmptyContainer>;
};
