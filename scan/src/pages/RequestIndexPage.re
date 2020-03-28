module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(20))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);
  let hFlex = style([display(`flex), alignItems(`center)]);

  let centerHFlex = style([display(`flex), alignItems(`center), justifyContent(`center)]);

  let minWidth = mw => style([minWidth(`px(mw)), maxWidth(`px(mw))]);

  let topicContainer = h =>
    style([display(`flex), alignItems(`center), width(`percent(100.)), height(`px(h))]);

  let logo = (w, mr) => style([width(`px(w)), marginRight(`px(mr))]);
  let headerContainer = style([lineHeight(`px(25))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let seperatedLongLine =
    style([width(`percent(100.)), height(`px(2)), backgroundColor(Colors.gray4)]);

  let fillRight = style([marginRight(`auto)]);

  let lowerPannel =
    style([
      width(`percent(100.)),
      padding(`px(30)),
      display(`flex),
      flexDirection(`column),
      backgroundColor(Colors.white),
      boxShadows([
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), Css.rgba(0, 0, 0, 0.1)),
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(12), Css.rgba(0, 0, 0, 0.03)),
      ]),
      borderRadius(`px(4)),
    ]);
};

[@react.component]
let make = (~reqID) => {
  let requestOpt = RequestHook.get(reqID);
  let (requestValidators, totalValidators) =
    switch (React.useContext(GlobalContext.context), requestOpt) {
    | (Some(info), Some(request)) => (
        info.validators
        ->Belt_List.keep(validator =>
            request.requestedValidators
            ->Belt_List.has(validator.operatorAddress, (a, b) => a->Address.isEqual(b))
          ),
        info.validators |> Belt_List.length,
      )
    | _ => ([], 0)
    };

  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.requestLogo className={Styles.logo(50, 10)} />
          <Text
            value="DATA REQUEST"
            weight=Text.Medium
            size=Text.Md
            spacing={Text.Em(0.06)}
            height={Text.Px(15)}
            nowrap=true
            color=Colors.gray7
            block=true
          />
          <div className=Styles.seperatedLine />
          {switch (requestOpt) {
           | Some(request) =>
             <TimeAgos
               time={request.requestedAtTime}
               size=Text.Md
               weight=Text.Thin
               spacing={Text.Em(0.06)}
               height={Text.Px(18)}
               upper=true
             />

           | None =>
             <Text
               value="???"
               size=Text.Md
               weight=Text.Thin
               spacing={Text.Em(0.06)}
               height={Text.Px(18)}
             />
           }}
        </div>
      </Col>
    </Row>
    {switch (requestOpt) {
     | Some(request) =>
       let numReport = request.reports |> Belt_List.length;
       <>
         <VSpacing size=Spacing.xl />
         <div className=Styles.vFlex>
           <TypeID.Request id={request.id} position=TypeID.Title />
         </div>
         <VSpacing size=Spacing.xl />
         <Row>
           <Col size=2.8>
             <InfoHL
               info={InfoHL.OracleScript(request.oracleScriptID, request.oracleScriptName)}
               header="ORACLE SCRIPT"
             />
           </Col>
           <Col size=3.2>
             <InfoHL header="SENDER" info={InfoHL.Address(request.requester, 280)} />
           </Col>
           <Col size=4.0>
             <InfoHL header="TX HASH" info={InfoHL.TxHash(request.txHash, 385)} />
           </Col>
         </Row>
         <VSpacing size=Spacing.xl />
         <Row>
           <Col>
             <InfoHL info={InfoHL.Validators(requestValidators)} header="REQUEST TO VALIDATORS" />
           </Col>
         </Row>
         <VSpacing size=Spacing.xl />
         <div className=Styles.lowerPannel>
           <div className={Styles.topicContainer(50)}>
             <Col size=1.1>
               <Text
                 value="REPORT STATUS"
                 size=Text.Sm
                 weight=Text.Semibold
                 spacing={Text.Em(0.06)}
                 color=Colors.gray6
               />
             </Col>
             <Col size=5.>
               <ProgressBar
                 reportedValidators=numReport
                 minimumValidators={request.sufficientValidatorCount}
                 totalValidators
               />
             </Col>
             <Col size=1.5>
               <div className=Styles.hFlex>
                 <div className=Styles.fillRight />
                 <Text
                   value={numReport |> string_of_int}
                   weight=Text.Bold
                   code=true
                   color=Colors.gray8
                 />
                 <HSpacing size=Spacing.sm />
                 <Text value="Reported" weight=Text.Regular code=true color=Colors.gray8 />
                 {switch (request.resolveStatus) {
                  | RequestHook.Request.Success =>
                    <>
                      <HSpacing size=Spacing.sm />
                      <Text value=", success" weight=Text.Regular color=Colors.gray8 />
                      <HSpacing size=Spacing.sm />
                      <img src=Images.success className={Styles.logo(20, 0)} />
                    </>
                  | RequestHook.Request.Failure =>
                    <>
                      <HSpacing size=Spacing.sm />
                      <Text value=", failure" weight=Text.Regular color=Colors.gray8 />
                      <HSpacing size=Spacing.sm />
                      <img src=Images.fail className={Styles.logo(20, 0)} />
                    </>
                  | _ => React.null
                  }}
               </div>
             </Col>
           </div>
           <div className={Styles.topicContainer(50)}>
             <Col size=1.>
               <Text
                 value="EXPIRATION BLOCK"
                 size=Text.Sm
                 weight=Text.Semibold
                 spacing={Text.Em(0.06)}
                 color=Colors.gray6
               />
             </Col>
             <Col size=1.>
               <div className=Styles.hFlex>
                 <div className=Styles.fillRight />
                 <TypeID.Block id={ID.Block.ID(request.expirationHeight)} />
                 {switch (request.resolveStatus) {
                  | RequestHook.Request.Open =>
                    <>
                      <HSpacing size=Spacing.sm />
                      <Text
                        value="(5 blocks remaining)" // Mock
                        weight=Text.Regular
                        code=true
                        color=Colors.gray8
                      />
                    </>
                  | _ => React.null
                  }}
               </div>
             </Col>
           </div>
           <VSpacing size=Spacing.sm />
           <div className={Styles.topicContainer(40)}>
             <Col size=1.>
               <div className=Styles.hFlex>
                 <Text
                   value="EXPIRATION DATA"
                   size=Text.Sm
                   weight=Text.Semibold
                   spacing={Text.Em(0.06)}
                   color=Colors.gray6
                 />
                 <HSpacing size=Spacing.md />
                 // Mock
                 <CopyButton data={"aaaa" |> JsBuffer.fromHex} />
               </div>
             </Col>
           </div>
           <KVTable
             tableWidth=880
             theme=KVTable.THEME_2
             headers=["EXTERNAL ID", "DATA SOURCE", "PARAM"]
             rows=[
               [
                 KVTable.Value("1"),
                 KVTable.DataSource(ID.DataSource.ID(12), "Mock Data Source"),
                 KVTable.Value("BTC"),
               ],
               [
                 KVTable.Value("2"),
                 KVTable.DataSource(ID.DataSource.ID(12), "Mock Data Source"),
                 KVTable.Value("BTC"),
               ],
               [
                 KVTable.Value("3"),
                 KVTable.DataSource(ID.DataSource.ID(12), "Mock Data Source"),
                 KVTable.Value("BTC"),
               ],
             ]
           />
           <VSpacing size=Spacing.lg />
           <div className={Styles.topicContainer(40)}>
             <Col size=1.>
               <div className=Styles.hFlex>
                 <Text
                   value="CALLDATA"
                   size=Text.Sm
                   weight=Text.Semibold
                   spacing={Text.Em(0.06)}
                   color=Colors.gray6
                 />
                 <HSpacing size=Spacing.md />
                 // Mock
                 <CopyButton data={"aaaa" |> JsBuffer.fromHex} />
               </div>
             </Col>
           </div>
           <KVTable
             tableWidth=880
             theme=KVTable.THEME_2
             rows=[
               [KVTable.Value("crypto_symbol"), KVTable.Value("Bitcoin")],
               [KVTable.Value("method"), KVTable.Value("median")],
             ]
           />
           {numReport > 0
              ? <>
                  <VSpacing size=Spacing.lg />
                  <div className={Styles.topicContainer(40)}>
                    <Col size=1.>
                      <div className=Styles.hFlex>
                        <Text
                          value="RESULT"
                          size=Text.Sm
                          weight=Text.Semibold
                          spacing={Text.Em(0.06)}
                          color=Colors.gray6
                        />
                        <HSpacing size=Spacing.md />
                        // Mock
                        <CopyButton data={"aaaa" |> JsBuffer.fromHex} />
                      </div>
                    </Col>
                  </div>
                  <KVTable
                    tableWidth=880
                    theme=KVTable.THEME_2
                    rows=[
                      [KVTable.Value("price"), KVTable.Value("861200")],
                      [KVTable.Value("timestamp"), KVTable.Value("1583383759")],
                    ]
                  />
                </>
              : React.null}
           {numReport >= request.sufficientValidatorCount
              ? {
                <RequestProof requestID={request.id} />;
              }
              : React.null}
           <VSpacing size=Spacing.xl />
           <div className=Styles.seperatedLongLine />
           <VSpacing size=Spacing.md />
           <div className={Styles.topicContainer(50)}>
             <Col size=1.>
               <div className=Styles.hFlex>
                 <Text value="Data Report from" weight=Text.Regular color=Colors.gray7 />
                 <HSpacing size=Spacing.sm />
                 <Text
                   value={numReport |> Format.iPretty}
                   weight=Text.Semibold
                   color=Colors.gray7
                 />
                 <HSpacing size=Spacing.sm />
                 <Text value="of" weight=Text.Regular color=Colors.gray7 />
                 <HSpacing size=Spacing.sm />
                 <Text
                   value={request.sufficientValidatorCount |> Format.iPretty}
                   weight=Text.Semibold
                   color=Colors.gray7
                 />
                 <HSpacing size=Spacing.sm />
                 <Text
                   value={request.sufficientValidatorCount > 1 ? "Validators" : "Validator"}
                   weight=Text.Regular
                   color=Colors.gray7
                 />
               </div>
             </Col>
             <Col size=1.>
               <div className=Styles.hFlex>
                 <div className=Styles.fillRight />
                 <Text
                   value={
                     (
                       numReport < request.sufficientValidatorCount
                         ? request.sufficientValidatorCount - numReport : 0
                     )
                     |> Format.iPretty
                   }
                   weight=Text.Semibold
                   color=Colors.gray7
                 />
                 <HSpacing size=Spacing.sm />
                 <Text
                   value={
                     request.sufficientValidatorCount > 1
                       ? "Validators Required" : "Validator Required"
                   }
                   weight=Text.Regular
                   color=Colors.gray7
                 />
               </div>
             </Col>
           </div>
           {numReport > 0
              ? <KVTable
                  tableWidth=880
                  theme=KVTable.THEME_2
                  sizes=[0.92, 0.73, 2., 0.63, 2.4]
                  isRights=[false, false, false, true, true]
                  headers=["REPORT BY", "BLOCK", "TX HASH", "EXTERNAL ID", "VALUE"]
                  rows={
                    request.reports
                    ->Belt_List.map(report =>
                        [
                          KVTable.Value("Paul Node"),
                          KVTable.Block(ID.Block.ID(report.reportedAtHeight)),
                          KVTable.TxHash(report.txHash),
                          KVTable.Values(
                            report.values
                            ->Belt_List.map(({externalDataID}) =>
                                externalDataID |> Format.iPretty
                              ),
                          ),
                          KVTable.Values(
                            report.values->Belt_List.map(({data}) => data |> JsBuffer.toHex),
                          ),
                        ]
                      )
                  }
                />
              : <div className={Styles.topicContainer(200)}>
                  <Col size=1.>
                    <div className=Styles.centerHFlex>
                      <img src=Images.noReport className={Styles.logo(80, 0)} />
                    </div>
                    <VSpacing size=Spacing.xl />
                    <div className=Styles.centerHFlex>
                      <Text value="NO REPORT" weight=Text.Regular color=Colors.blue4 />
                    </div>
                  </Col>
                </div>}
         </div>
       </>;
     | None => React.null
     }}
  </div>;
};
