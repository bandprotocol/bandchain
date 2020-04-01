module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(20)), minHeight(`px(500))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);

  let fullWidth = style([width(`percent(100.0)), display(`flex)]);

  let progressBarContainer = style([maxWidth(`px(300))]);

  let resolveStatusContainer = style([justifyContent(`center)]);

  let loadingContainer =
    style([
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      height(`px(200)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
    ]);
};

[@react.component]
let make = () =>
  {
    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 10;

    let requestCount = 100;
    let requests = [
      {
        RequestHook.Request.id: ID.Request.ID(4),
        oracleScriptID: ID.OracleScript.ID(3),
        oracleScriptName: "Oracle script 3",
        calldata: "AAAAAAAAV0M=" |> JsBuffer.fromBase64,
        requestedValidators: [
          "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
        ],
        sufficientValidatorCount: 1,
        expirationHeight: 3000,
        resolveStatus: Unknown,
        requester: "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
        txHash: "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
        requestedAtHeight: 40000,
        requestedAtTime: MomentRe.momentNow(),
        rawDataRequests: [],
        reports: [],
        result: Some("AAAAAAAAV0M=" |> JsBuffer.fromBase64),
      },
      {
        id: ID.Request.ID(3),
        oracleScriptID: ID.OracleScript.ID(2),
        oracleScriptName: "Oracle script 2",
        calldata: "AAAAAAAAV0M=" |> JsBuffer.fromBase64,
        requestedValidators: [
          "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
          "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
          "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
          "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
        ],
        sufficientValidatorCount: 3,
        expirationHeight: 3000,
        resolveStatus: Success,
        requester: "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
        txHash: "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
        requestedAtHeight: 40000,
        requestedAtTime: MomentRe.momentNow(),
        rawDataRequests: [],
        reports: [
          {
            reporter: "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
            txHash:
              "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
            reportedAtHeight: 40000,
            reportedAtTime: MomentRe.momentNow(),
            values: [],
          },
          {
            reporter: "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
            txHash:
              "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
            reportedAtHeight: 40000,
            reportedAtTime: MomentRe.momentNow(),
            values: [],
          },
          {
            reporter: "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
            txHash:
              "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
            reportedAtHeight: 40000,
            reportedAtTime: MomentRe.momentNow(),
            values: [],
          },
          {
            reporter: "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
            txHash:
              "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
            reportedAtHeight: 40000,
            reportedAtTime: MomentRe.momentNow(),
            values: [],
          },
        ],
        result: Some("AAAAAAAAV0M=" |> JsBuffer.fromBase64),
      },
      {
        id: ID.Request.ID(2),
        oracleScriptID: ID.OracleScript.ID(1),
        oracleScriptName: "name",
        calldata: "AAAAAAAAV0M=" |> JsBuffer.fromBase64,
        requestedValidators: [
          "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
          "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
        ],
        sufficientValidatorCount: 2,
        expirationHeight: 3000,
        resolveStatus: Failure,
        requester: "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
        txHash: "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
        requestedAtHeight: 40000,
        requestedAtTime: MomentRe.momentNow(),
        rawDataRequests: [],
        reports: [],
        result: Some("AAAAAAAAV0M=" |> JsBuffer.fromBase64),
      },
      {
        id: ID.Request.ID(1),
        oracleScriptID: ID.OracleScript.ID(2),
        oracleScriptName: "Oracle script 2",
        calldata: "AAAAAAAAV0M=" |> JsBuffer.fromBase64,
        requestedValidators: [
          "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
          "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
          "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
        ],
        sufficientValidatorCount: 2,
        expirationHeight: 3000,
        resolveStatus: Open,
        requester: "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
        txHash: "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
        requestedAtHeight: 40000,
        requestedAtTime: MomentRe.momentNow(),
        rawDataRequests: [],
        reports: [
          {
            reporter: "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
            txHash:
              "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
            reportedAtHeight: 40000,
            reportedAtTime: MomentRe.momentNow(),
            values: [],
          },
        ],
        result: Some("AAAAAAAAV0M=" |> JsBuffer.fromBase64),
      },
    ];

    let pageCount = Page.getPageCount(requestCount, pageSize);

    <div className=Styles.pageContainer>
      <Row>
        <Col>
          <div className=Styles.vFlex>
            <img src=Images.requestLogo className=Styles.logo />
            <Text
              value="ALL REQUESTS"
              weight=Text.Medium
              size=Text.Md
              spacing={Text.Em(0.06)}
              height={Text.Px(15)}
              nowrap=true
              color=Colors.gray7
              block=true
            />
            <div className=Styles.seperatedLine />
            <Text
              value={requestCount->string_of_int ++ " In total"}
              size=Text.Md
              weight=Text.Thin
              spacing={Text.Em(0.06)}
              color=Colors.gray7
              nowrap=true
            />
          </div>
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <>
        <THead>
          <Row>
            <Col> <HSpacing size=Spacing.lg /> </Col>
            <Col size=0.5>
              <Text
                block=true
                value="REQUEST ID"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.gray5
                spacing={Text.Em(0.1)}
              />
            </Col>
            <Col size=0.78>
              <Text
                block=true
                value="AGE"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.gray5
                spacing={Text.Em(0.1)}
              />
            </Col>
            <Col size=1.15>
              <Text
                block=true
                value="ORACLE SCRIPTS"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.gray5
                spacing={Text.Em(0.1)}
              />
            </Col>
            <Col size=1.9>
              <Text
                block=true
                value="REPORT STATUS"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.gray5
                spacing={Text.Em(0.1)}
              />
            </Col>
            <Col size=0.72 justifyContent=Col.End>
              <Text
                block=true
                value="RESOLVE STATUS"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.gray5
                align=Text.Right
                spacing={Text.Em(0.1)}
              />
            </Col>
            <Col> <HSpacing size=Spacing.xl /> </Col>
          </Row>
        </THead>
        {requests
         ->Belt_List.map(
             (
               {
                 id,
                 requestedAtTime,
                 oracleScriptID,
                 oracleScriptName,
                 requestedValidators,
                 sufficientValidatorCount,
                 reports,
                 resolveStatus,
               },
             ) => {
             <TBody key={id |> ID.Request.toString}>
               <div className=Styles.fullWidth>
                 <Row minHeight={`px(35)}>
                   <Col> <HSpacing size=Spacing.lg /> </Col>
                   <Col size=0.5> <TElement elementType={TElement.Request(id)} /> </Col>
                   <Col size=0.78>
                     <TElement elementType={requestedAtTime->TElement.Timestamp} />
                   </Col>
                   <Col size=1.15>
                     <TElement
                       elementType={TElement.OracleScript(oracleScriptID, oracleScriptName)}
                     />
                   </Col>
                   <Col size=1.9>
                     <div className=Styles.progressBarContainer>
                       <ProgressBar
                         reportedValidators={reports |> Belt_List.length}
                         minimumValidators=sufficientValidatorCount
                         totalValidators={requestedValidators |> Belt_List.length}
                       />
                     </div>
                   </Col>
                   <Col size=0.72 justifyContent=Col.End>
                     <TElement elementType={resolveStatus->TElement.RequestStatus} />
                   </Col>
                   <Col> <HSpacing size=Spacing.xl /> </Col>
                 </Row>
               </div>
             </TBody>
           })
         ->Belt_List.toArray
         ->React.array}
      </>
      <VSpacing size=Spacing.lg />
      <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
      <VSpacing size=Spacing.lg />
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
