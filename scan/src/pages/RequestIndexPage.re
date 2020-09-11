module Styles = {
  open Css;

  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      padding(`px(24)),
      Media.mobile([padding(`px(16))]),
    ]);

  let setPaddingBottom = (~pb, ~pbSm, ()) =>
    style([paddingBottom(`px(pb)), Media.mobile([paddingBottom(`px(pbSm))])]);

  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
  let infoIcon = style([width(`px(12)), height(`px(12)), display(`block)]);

  let emptyContainer =
    style([
      height(`px(200)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      flexDirection(`column),
      backgroundColor(Colors.blueGray1),
    ]);

  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);

  let kvTableContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      padding(`px(1)),
    ]);

  let kvTableHeader =
    style([
      padding4(~top=`px(16), ~left=`px(24), ~right=`px(24), ~bottom=`px(12)),
      Media.mobile([padding2(~v=`px(14), ~h=`px(12))]),
    ]);

  let kvTableMobile =
    style([
      margin4(~top=`zero, ~left=`px(12), ~right=`px(12), ~bottom=`px(12)),
      backgroundColor(Colors.profileBG),
      boxShadow(`none),
    ]);

  let seperatedLine =
    style([
      width(`percent(100.)),
      height(`px(1)),
      backgroundColor(Colors.gray9),
      margin2(~v=`px(24), ~h=`zero),
    ]);

  let addressContainer = style([Media.mobile([width(`px(260))])]);

  let dataSourceContainer = style([width(`percent(100.))]);

  let validatorReportStatus = style([marginBottom(`px(13))]);

  let loading = style([width(`px(65)), height(`px(20)), marginBottom(`px(16))]);
};

module ValidatorReportStatus = {
  [@react.component]
  let make = (~moniker, ~isReport, ~resolveStatus) => {
    <div
      className={Css.merge([
        CssHelper.flexBox(~align=`center, ~wrap=`nowrap, ()),
        Styles.validatorReportStatus,
      ])}>
      {switch (isReport, resolveStatus) {
       | (true, _) => <Icon name="fas fa-check-circle" color=Colors.green4 />
       | (false, RequestSub.Expired) => <Icon name="fas fa-times-circle" color=Colors.red4 />
       | (false, _) => <Icon name="fas fa-clock" color=Colors.yellowAccent2 />
       }}
      <HSpacing size=Spacing.sm />
      <Text value=moniker color=Colors.gray7 ellipsis=true />
    </div>;
  };
};

module KVTableContainer = {
  module TableHeader = {
    [@react.component]
    let make = () => {
      <THead.Grid>
        <Row.Grid alignItems=Row.Center>
          <Col.Grid col=Col.Three>
            <Heading value="Key" size=Heading.H5 weight=Heading.Medium />
          </Col.Grid>
          <Col.Grid col=Col.Nine>
            <Heading value="Value" size=Heading.H5 weight=Heading.Medium />
          </Col.Grid>
        </Row.Grid>
      </THead.Grid>;
    };
  };

  module Loading = {
    [@react.component]
    let make = () => {
      let isMobile = Media.isMobile();

      isMobile
        ? <MobileCard
            values=InfoMobileCard.[("Key", Loading(60)), ("Value", Loading(60))]
            idx="1"
            styles=Styles.kvTableMobile
          />
        : <>
            <TableHeader />
            <TBody.Grid paddingH={`px(24)}>
              <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
                <Col.Grid col=Col.Three> <LoadingCensorBar width=60 height=15 /> </Col.Grid>
                <Col.Grid col=Col.Nine> <LoadingCensorBar width=100 height=15 /> </Col.Grid>
              </Row.Grid>
            </TBody.Grid>
            <TBody.Grid paddingH={`px(24)}>
              <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
                <Col.Grid col=Col.Three> <LoadingCensorBar width=60 height=15 /> </Col.Grid>
                <Col.Grid col=Col.Nine> <LoadingCensorBar width=100 height=15 /> </Col.Grid>
              </Row.Grid>
            </TBody.Grid>
          </>;
    };
  };

  [@react.component]
  let make = (~decodesOpt) => {
    switch (decodesOpt) {
    | Some(decodes) =>
      let isMobile = Media.isMobile();
      isMobile
        ? decodes
          ->Belt.Array.map(({Obi.fieldName, fieldValue}) =>
              <MobileCard
                values=InfoMobileCard.[("Key", Text(fieldName)), ("Value", Text(fieldValue))]
                key={fieldName ++ fieldValue}
                idx={fieldName ++ fieldValue}
                styles=Styles.kvTableMobile
              />
            )
          ->React.array
        : <>
            <TableHeader />
            {decodes
             ->Belt.Array.map(({Obi.fieldName, fieldValue}) => {
                 <TBody.Grid key={fieldName ++ fieldValue} paddingH={`px(24)}>
                   <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
                     <Col.Grid col=Col.Three>
                       <Text value=fieldName color=Colors.gray7 weight=Text.Thin/>
                     </Col.Grid>
                     <Col.Grid col=Col.Nine>
                       <Text value=fieldValue color=Colors.gray7 weight=Text.Thin breakAll=true/>
                     </Col.Grid>
                   </Row.Grid>
                 </TBody.Grid>
               })
             ->React.array}
          </>;
    | None =>
      <div className=Styles.emptyContainer>
        <img src=Images.noSource className=Styles.noDataImage />
        <Heading
          size=Heading.H4
          value="Schema not found"
          align=Heading.Center
          weight=Heading.Regular
          color=Colors.bandBlue
        />
      </div>
    };
  };
};

[@react.component]
let make = (~reqID) => {
  let requestSub = RequestSub.get(reqID);
  let isMobile = Media.isMobile();

  <Section>
    <div className=CssHelper.container>
      <Row.Grid marginBottom=40 marginBottomSm=16>
        <Col.Grid>
          <Heading value="Data Request" size=Heading.H4 marginBottom=40 marginBottomSm=24 />
          {switch (requestSub) {
           | Data({id}) => <TypeID.Request id position=TypeID.Title />
           | _ => <LoadingCensorBar width=150 height=23 />
           }}
        </Col.Grid>
      </Row.Grid>
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div
            className={Css.merge([
              Styles.infoContainer,
              Styles.setPaddingBottom(~pb=11, ~pbSm=5, ()),
            ])}>
            <Heading
              value="Request Info"
              size=Heading.H4
              style=Styles.infoHeader
              marginBottom=24
            />
            <Row.Grid marginBottom=24>
              <Col.Grid col=Col.Six mbSm=24>
                <Heading value="Oracle Scripts" size=Heading.H5 />
                <VSpacing size={`px(8)} />
                {switch (requestSub) {
                 | Data({oracleScript: {oracleScriptID, name}}) =>
                   <div className={CssHelper.flexBox()}>
                     <TypeID.OracleScript id=oracleScriptID position=TypeID.Subtitle />
                     <HSpacing size=Spacing.sm />
                     <Text value=name size=Text.Lg />
                   </div>
                 | _ => <LoadingCensorBar width=200 height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six>
                <Heading value="Sender" size=Heading.H5 />
                <VSpacing size={`px(8)} />
                {switch (requestSub) {
                 | Data({requester}) =>
                   <div className=Styles.addressContainer>
                     <AddressRender address=requester position=AddressRender.Subtitle />
                   </div>
                 | _ => <LoadingCensorBar width=200 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
            <Row.Grid>
              <Col.Grid col=Col.Six mbSm=24>
                <Heading value="TX Hash" size=Heading.H5 />
                <VSpacing size=Spacing.sm />
                {switch (requestSub) {
                 | Data({transaction: {hash}}) =>
                   <TxLink txHash=hash width={isMobile ? 260 : 360} />
                 | _ => <LoadingCensorBar width=200 height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six>
                <Heading value="Fee" size=Heading.H5 />
                <VSpacing size=Spacing.sm />
                {switch (requestSub) {
                 | Data({transaction: {gasFee}}) =>
                   <Text
                     block=true
                     value={
                       (gasFee |> Coin.getBandAmountFromCoins |> Format.fPretty(~digits=2))
                       ++ " BAND"
                     }
                     size=Text.Lg
                     color=Colors.gray7
                   />
                 | _ => <LoadingCensorBar width=200 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
            <div className=Styles.seperatedLine />
            <Row.Grid marginBottom=24>
              <Col.Grid col=Col.Six mbSm=24>
                <Heading value="Report Status" size=Heading.H5 marginBottom=8 />
                {switch (requestSub) {
                 | Data({minCount, requestedValidators, reports}) =>
                   <ProgressBar
                     reportedValidators={reports->Belt.Array.size}
                     minimumValidators=minCount
                     requestValidators={requestedValidators->Belt.Array.size}
                   />
                 | _ => <LoadingCensorBar width=200 height=15 />
                 }}
              </Col.Grid>
              <Col.Grid col=Col.Six>
                <Heading value="Resolve Status" size=Heading.H5 marginBottom=8 />
                {switch (requestSub) {
                 | Data({resolveStatus}) =>
                   <RequestStatus resolveStatus display=RequestStatus.Full />
                 | _ => <LoadingCensorBar width=200 height=15 />
                 }}
              </Col.Grid>
            </Row.Grid>
            <Row.Grid>
              <Col.Grid>
                <Heading value="Request to" size=Heading.H5 marginBottom=8 />
                <Row.Grid wrap=true>
                  {switch (requestSub) {
                   | Data({requestedValidators, resolveStatus, reports}) =>
                     requestedValidators
                     ->Belt.Array.map(({validator: {moniker, consensusAddress}}) => {
                         let isReport =
                           reports->Belt.Array.some(({reportValidator}) =>
                             consensusAddress == reportValidator.consensusAddress
                           );
                         <Col.Grid col=Col.Three colSm=Col.Six key=moniker>
                           <ValidatorReportStatus moniker isReport resolveStatus />
                         </Col.Grid>;
                       })
                     ->React.array
                   | _ =>
                     <Col.Grid>
                       <LoadingCensorBar width=200 height=15 />
                       <VSpacing size={`px(isMobile ? 5 : 11)} />
                     </Col.Grid>
                   }}
                </Row.Grid>
              </Col.Grid>
            </Row.Grid>
          </div>
        </Col.Grid>
      </Row.Grid>
      // Calldata
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div className=Styles.kvTableContainer>
            <div
              className={Css.merge([
                CssHelper.flexBox(~justify=`spaceBetween, ()),
                Styles.kvTableHeader,
              ])}>
              <div className={CssHelper.flexBox()}>
                <Heading value="Calldata" size=Heading.H4 />
                <HSpacing size=Spacing.xs />
                <CTooltip tooltipText="The input parameters associated with the request">
                  <Icon name="fal fa-info-circle" size=10 />
                </CTooltip>
              </div>
              {switch (requestSub) {
               | Data({calldata}) =>
                 <CopyButton
                   data={calldata |> JsBuffer.toHex(~with0x=false)}
                   title="Copy as bytes"
                   width=125
                 />
               | _ => <LoadingCensorBar width=125 height=28 />
               }}
            </div>
            {switch (requestSub) {
             | Data({oracleScript: {schema}, calldata}) =>
               let decodesOpt = Obi.decode(schema, "input", calldata);
               <KVTableContainer decodesOpt />;
             | _ => <KVTableContainer.Loading />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      // Result
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div className=Styles.kvTableContainer>
            <div
              className={Css.merge([
                CssHelper.flexBox(~justify=`spaceBetween, ()),
                Styles.kvTableHeader,
              ])}>
              <div className={CssHelper.flexBox()}>
                <Heading value="Result" size=Heading.H4 />
                <HSpacing size=Spacing.xs />
                <CTooltip tooltipText="The final result of the request">
                  <Icon name="fal fa-info-circle" size=10 />
                </CTooltip>
              </div>
              {switch (requestSub) {
               | Data({result: resultOpt, resolveStatus}) =>
                 switch (resultOpt, resolveStatus) {
                 | (Some(result), Success) =>
                   <CopyButton
                     data={result |> JsBuffer.toHex(~with0x=false)}
                     title="Copy as bytes"
                     width=125
                   />
                 | (_, _) => React.null
                 }
               | _ => <LoadingCensorBar width=125 height=28 />
               }}
            </div>
            {switch (requestSub) {
             | Data({oracleScript: {schema}, result: resultOpt, resolveStatus}) =>
               switch (resolveStatus, resultOpt) {
               | (RequestSub.Success, Some(result)) =>
                 let decodesOpt = Obi.decode(schema, "output", result);
                 <KVTableContainer decodesOpt />;
               | (Pending, _) =>
                 <div className=Styles.emptyContainer>
                   <img src=Images.loadingCircles className=Styles.loading />
                   <Heading
                     size=Heading.H4
                     value="Waiting for result"
                     align=Heading.Center
                     weight=Heading.Regular
                     color=Colors.bandBlue
                   />
                 </div>
               | (_, _) =>
                 <div className=Styles.emptyContainer>
                   <img src=Images.noSource className=Styles.noDataImage />
                   <Heading
                     size=Heading.H4
                     value="This request hasn't resolved"
                     align=Heading.Center
                     weight=Heading.Regular
                     color=Colors.bandBlue
                   />
                 </div>
               }
             | _ => <KVTableContainer.Loading />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      // Proof
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div className=Styles.kvTableContainer>
            <div className={Css.merge([Styles.kvTableHeader, CssHelper.flexBox()])}>
              <Heading value="Proof of validity" size=Heading.H4 />
            </div>
            // TODO: add later
            // <ExtLinkButton link="https://docs.bandchain.org/" description="What is proof ?" />
            {switch (requestSub) {
             | Data(request) =>
               switch (request.resolveStatus) {
               | Success => <RequestProof request />
               | Pending =>
                 <div className=Styles.emptyContainer>
                   <img src=Images.loadingCircles className=Styles.loading />
                   <Heading
                     size=Heading.H4
                     value="Waiting for result"
                     align=Heading.Center
                     weight=Heading.Regular
                     color=Colors.bandBlue
                   />
                 </div>
               | _ =>
                 <div className=Styles.emptyContainer>
                   <img src=Images.noSource className=Styles.noDataImage />
                   <Heading
                     size=Heading.H4
                     value="This request hasn't resolved"
                     align=Heading.Center
                     weight=Heading.Regular
                     color=Colors.bandBlue
                   />
                 </div>
               }
             | _ => <LoadingCensorBar fullWidth=true height=100 />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      // External Data Table
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div className=Styles.kvTableContainer>
            <div className=Styles.kvTableHeader>
              <div className={CssHelper.flexBox()}>
                <Heading value="External Data" size=Heading.H4 />
                <HSpacing size=Spacing.xs />
                <CTooltip
                  tooltipText="Data reported by the validators by querying the data sources">
                  <Icon name="fal fa-info-circle" size=10 />
                </CTooltip>
              </div>
            </div>
            {isMobile
               ? React.null
               : <THead.Grid>
                   <Row.Grid alignItems=Row.Center>
                     <Col.Grid col=Col.Three>
                       <Heading value="External ID" size=Heading.H5 weight=Heading.Medium />
                     </Col.Grid>
                     <Col.Grid col=Col.Four>
                       <Heading value="Data Source" size=Heading.H5 weight=Heading.Medium />
                     </Col.Grid>
                     <Col.Grid col=Col.Five>
                       <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
                         <Heading value="Param" size=Heading.H5 weight=Heading.Medium />
                       </div>
                     </Col.Grid>
                   </Row.Grid>
                 </THead.Grid>}
            {switch (requestSub) {
             | Data({rawDataRequests}) =>
               rawDataRequests
               ->Belt.Array.map(({externalID, dataSource: {dataSourceID, name}, calldata}) => {
                   isMobile
                     ? <MobileCard
                         values=InfoMobileCard.[
                           ("External ID", Text(externalID)),
                           ("Data Source", DataSource(dataSourceID, name)),
                           ("Param", Text(calldata |> JsBuffer.toUTF8)),
                         ]
                         key={externalID ++ name}
                         idx={externalID ++ name}
                         styles=Styles.kvTableMobile
                       />
                     : <TBody.Grid key=externalID paddingH={`px(24)}>
                         <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
                           <Col.Grid col=Col.Three>
                             <Text value=externalID color=Colors.gray7 weight=Text.Thin />
                           </Col.Grid>
                           <Col.Grid col=Col.Four>
                             <div className={CssHelper.flexBox()}>
                               <TypeID.DataSource id=dataSourceID position=TypeID.Text />
                               <HSpacing size=Spacing.sm />
                               <Text value=name color=Colors.gray7 weight=Text.Thin />
                             </div>
                           </Col.Grid>
                           <Col.Grid col=Col.Five>
                             <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
                               <Text
                                 value={calldata->JsBuffer.toUTF8}
                                 color=Colors.gray7
                                 weight=Text.Thin
                                 align=Text.Right
                               />
                             </div>
                           </Col.Grid>
                         </Row.Grid>
                       </TBody.Grid>
                 })
               ->React.array
             | _ =>
               isMobile
                 ? <MobileCard
                     values=InfoMobileCard.[
                       ("External ID", Loading(60)),
                       ("Data Source", Loading(60)),
                       ("Param", Loading(60)),
                     ]
                     idx="1"
                     styles=Styles.kvTableMobile
                   />
                 : <TBody.Grid paddingH={`px(24)}>
                     <Row.Grid alignItems=Row.Center minHeight={`px(30)}>
                       <Col.Grid col=Col.Three> <LoadingCensorBar width=60 height=15 /> </Col.Grid>
                       <Col.Grid col=Col.Four> <LoadingCensorBar width=100 height=15 /> </Col.Grid>
                       <Col.Grid col=Col.Five>
                         <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
                           <LoadingCensorBar width=50 height=15 />
                         </div>
                       </Col.Grid>
                     </Row.Grid>
                   </TBody.Grid>
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
      // Data report
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <div className=Styles.kvTableContainer>
            <div className=Styles.kvTableHeader>
              <Heading value="Data Report" size=Heading.H4 />
            </div>
            {switch (requestSub) {
             | Data({reports}) => <DataReports reports />
             | _ => <LoadingCensorBar fullWidth=true height=200 />
             }}
          </div>
        </Col.Grid>
      </Row.Grid>
    </div>
  </Section>;
};
