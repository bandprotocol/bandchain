module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);
  let hFlex = style([display(`flex), alignItems(`center)]);

  let center = style([justifyContent(center)]);

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
let make = (~reqID) =>
  {
    let requestSub = RequestSub.get(reqID);
    let blockCountSub = BlockSub.count();

    let%Sub request = requestSub;
    let%Sub blockCount = blockCountSub;

    let numReport = request.reports |> Belt_Array.size;
    let remainingBlock =
      blockCount >= request.expirationHeight ? 0 : request.expirationHeight - blockCount;
    let calldataKVs =
      Borsh.decode(request.oracleScript.schema, "Input", request.calldata)
      ->Belt_Option.getWithDefault([||]);

    <>
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
            <Timestamp
              time={request.transaction.timestamp}
              size=Text.Md
              weight=Text.Thin
              spacing={Text.Em(0.06)}
              height={Text.Px(18)}
              upper=true
            />
          </div>
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <div className=Styles.vFlex> <TypeID.Request id={request.id} position=TypeID.Title /> </div>
      <VSpacing size=Spacing.xl />
      <Row>
        <Col size=2.8>
          <InfoHL
            info={
              InfoHL.OracleScript(request.oracleScript.oracleScriptID, request.oracleScript.name)
            }
            header="ORACLE SCRIPT"
          />
        </Col>
        <Col size=3.2>
          <InfoHL header="SENDER" info={InfoHL.Address(request.requester, 280)} />
        </Col>
        <Col size=4.0>
          <InfoHL header="TX HASH" info={InfoHL.TxHash(request.transaction.txHash, 385)} />
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <Row>
        <Col>
          <InfoHL
            info={
              InfoHL.ValidatorsMini(
                request.requestedValidators->Belt_Array.map(({validator}) => validator),
              )
            }
            header="REQUEST TO VALIDATORS"
          />
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
              requestValidators={request.requestedValidators->Belt_Array.size}
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
               | RequestSub.Success =>
                 <>
                   <HSpacing size=Spacing.sm />
                   <Text value=", success" weight=Text.Regular color=Colors.gray8 />
                   <HSpacing size=Spacing.sm />
                   <img src=Images.success className={Styles.logo(20, 0)} />
                 </>
               | RequestSub.Failure =>
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
               | RequestSub.Pending =>
                 <>
                   <HSpacing size=Spacing.sm />
                   <Text
                     value={j|($remainingBlock blocks remaining)|j}
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
                value="CALLDATA"
                size=Text.Sm
                weight=Text.Semibold
                spacing={Text.Em(0.06)}
                color=Colors.gray6
              />
              <HSpacing size=Spacing.md />
              <CopyButton data={request.calldata} />
            </div>
          </Col>
        </div>
        <KVTable
          tableWidth=880
          theme=KVTable.RequestMiniTable
          rows={
            calldataKVs
            ->Belt_Array.map(({fieldName, fieldValue}) =>
                [KVTable.Value(fieldName), KVTable.Value(fieldValue)]
              )
            ->Belt_List.fromArray
          }
        />
        <VSpacing size=Spacing.lg />
        <div className={Styles.topicContainer(40)}>
          <Col size=1.>
            <div className=Styles.hFlex>
              <Text
                value="EXTERNAL DATA"
                size=Text.Sm
                weight=Text.Semibold
                spacing={Text.Em(0.06)}
                color=Colors.gray6
              />
              <HSpacing size=Spacing.md />
            </div>
          </Col>
        </div>
        <KVTable
          tableWidth=880
          theme=KVTable.RequestMiniTable
          headers=["EXTERNAL ID", "DATA SOURCE", "PARAM"]
          rows={
            request.rawDataRequests
            ->Belt_Array.map(({externalID, dataSource, calldata}) =>
                [
                  KVTable.Value(externalID |> string_of_int),
                  KVTable.DataSource(dataSource.dataSourceID, dataSource.name),
                  KVTable.Value(calldata |> JsBuffer.toUTF8),
                ]
              )
            ->Belt_List.fromArray
          }
        />
        {switch (request.result) {
         | Some(result) =>
           let resultKVs =
             Borsh.decode(request.oracleScript.schema, "Output", result)
             ->Belt_Option.getWithDefault([||]);
           <>
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
                   <CopyButton data=result />
                 </div>
               </Col>
             </div>
             <KVTable
               tableWidth=880
               theme=KVTable.RequestMiniTable
               rows={
                 resultKVs
                 ->Belt_Array.map(({fieldName, fieldValue}) =>
                     [KVTable.Value(fieldName), KVTable.Value(fieldValue)]
                   )
                 ->Belt_List.fromArray
               }
             />
           </>;
         | None => React.null
         }}
        // {numReport >= request.sufficientValidatorCount
        //    ? {
        //      <RequestProof requestID={request.id} />;
        //    }
        //    : React.null}
        <VSpacing size=Spacing.xl />
        <div className=Styles.seperatedLongLine />
        <VSpacing size=Spacing.md />
        <div className={Styles.topicContainer(50)}>
          <Col size=1.>
            <div className=Styles.hFlex>
              <Text value="Data Report from" weight=Text.Regular color=Colors.gray7 />
              <HSpacing size=Spacing.sm />
              <Text value={numReport |> Format.iPretty} weight=Text.Semibold color=Colors.gray7 />
              <HSpacing size=Spacing.sm />
              <Text value="of" weight=Text.Regular color=Colors.gray7 />
              <HSpacing size=Spacing.sm />
              <Text
                value={request.requestedValidators->Belt_Array.size |> Format.iPretty}
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
               theme=KVTable.RequestMiniTable
               sizes=[0.92, 0.73, 2., 0.63, 2.4]
               isRights=[false, false, false, true, true]
               headers=["REPORT BY", "BLOCK", "TX HASH", "EXTERNAL ID", "VALUE"]
               rows={
                 request.reports
                 ->Belt_Array.map(report =>
                     [
                       KVTable.Validator(report.validatorByValidator),
                       KVTable.Block(report.transaction.blockHeight),
                       KVTable.TxHash(report.transaction.txHash),
                       KVTable.Values(
                         report.reportDetails
                         ->Belt_Array.map(({externalID}) => externalID |> Format.iPretty)
                         ->Belt_List.fromArray,
                       ),
                       KVTable.Values(
                         report.reportDetails
                         ->Belt_Array.map(({data}) => data |> JsBuffer.toUTF8)
                         ->Belt_List.fromArray,
                       ),
                     ]
                   )
                 ->Belt_List.fromArray
               }
             />
           : <div className={Styles.topicContainer(200)}>
               <Col size=1.>
                 <div className={Css.merge([Styles.center, Styles.hFlex])}>
                   <img src=Images.noReport className={Styles.logo(80, 0)} />
                 </div>
                 <VSpacing size=Spacing.xl />
                 <div className={Css.merge([Styles.center, Styles.hFlex])}>
                   <Text value="NO REPORT" weight=Text.Regular color=Colors.blue4 />
                 </div>
               </Col>
             </div>}
      </div>
    </>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
