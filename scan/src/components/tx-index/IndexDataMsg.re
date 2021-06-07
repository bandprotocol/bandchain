module CreateDataSourceMsg = {
  [@react.component]
  let make = (~dataSource: TxSub.Msg.CreateDataSource.success_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={dataSource.owner} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.DataSource position=TypeID.Subtitle id={dataSource.id} />
          <HSpacing size=Spacing.sm />
          <Text value={dataSource.name} size=Text.Lg />
        </div>
      </Col>
      <Col col=Col.Six>
        <Heading value="Treasury" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={dataSource.treasury} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Fee" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={dataSource.fee} />
      </Col>
    </Row>;
  };
};

module CreateDataSourceFailMsg = {
  [@react.component]
  let make = (~dataSource: TxSub.Msg.CreateDataSource.fail_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={dataSource.owner} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <Text value={dataSource.name} size=Text.Lg />
      </Col>
      <Col col=Col.Six>
        <Heading value="Treasury" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={dataSource.treasury} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Fee" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={dataSource.fee} />
      </Col>
    </Row>;
  };
};

module EditDataSourceMsg = {
  [@react.component]
  let make = (~dataSource: TxSub.Msg.EditDataSource.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={dataSource.owner} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.DataSource position=TypeID.Subtitle id={dataSource.id} />
          {dataSource.name == Config.doNotModify
             ? React.null
             : <> <HSpacing size=Spacing.sm /> <Text value={dataSource.name} size=Text.Lg /> </>}
        </div>
      </Col>
      <Col col=Col.Six>
        <Heading value="Fee" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={dataSource.fee} />
      </Col>
    </Row>;
  };
};

module CreateOracleScriptMsg = {
  [@react.component]
  let make = (~oracleScript: TxSub.Msg.CreateOracleScript.success_t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={oracleScript.owner} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.OracleScript position=TypeID.Subtitle id={oracleScript.id} />
          <HSpacing size=Spacing.sm />
          <Text value={oracleScript.name} />
        </div>
      </Col>
    </Row>;
  };
};

module CreateOracleScriptFailMsg = {
  [@react.component]
  let make = (~oracleScript: TxSub.Msg.CreateOracleScript.fail_t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={oracleScript.owner} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <Text value={oracleScript.name} size=Text.Lg />
      </Col>
    </Row>;
  };
};

module EditOracleScriptMsg = {
  [@react.component]
  let make = (~oracleScript: TxSub.Msg.EditOracleScript.t) => {
    <Row>
      <Col col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={oracleScript.owner} />
      </Col>
      <Col col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.OracleScript position=TypeID.Subtitle id={oracleScript.id} />
          {oracleScript.name == Config.doNotModify
             ? React.null
             : <> <HSpacing size=Spacing.sm /> <Text value={oracleScript.name} size=Text.Lg /> </>}
        </div>
      </Col>
    </Row>;
  };
};

module RequestMsg = {
  [@react.component]
  let make = (~request: TxSub.Msg.Request.success_t) => {
    let calldataKVsOpt = Obi.decode(request.schema, "input", request.calldata);
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={request.sender} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Request ID" size=Heading.H5 marginBottom=8 />
        <TypeID.Request position=TypeID.Subtitle id={request.id} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Oracle Script" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.OracleScript position=TypeID.Subtitle id={request.oracleScriptID} />
          <HSpacing size=Spacing.sm />
          <Text value={request.oracleScriptName} size=Text.Lg />
        </div>
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Fee Limit" size=Heading.H5 marginBottom=8 />
        <AmountRender coins={request.feeLimit} />
      </Col>
      <Col mb=24>
        <div
          className={Css.merge([CssHelper.flexBox(~justify=`spaceBetween, ()), CssHelper.mb()])}>
          <Heading value="Calldata" size=Heading.H5 />
          <CopyButton
            data={request.calldata |> JsBuffer.toHex(~with0x=false)}
            title="Copy as bytes"
            width=125
          />
        </div>
        {switch (calldataKVsOpt) {
         | Some(calldataKVs) =>
           <KVTable
             rows={
               calldataKVs
               ->Belt_Array.map(({fieldName, fieldValue}) =>
                   [KVTable.Value(fieldName), KVTable.Value(fieldValue)]
                 )
               ->Belt_List.fromArray
             }
           />
         | None =>
           <Text
             value="Could not decode calldata."
             spacing={Text.Em(0.02)}
             nowrap=true
             ellipsis=true
             code=true
             block=true
             size=Text.Lg
           />
         }}
      </Col>
      <Col col=Col.Six mbSm=24>
        <Heading value="Request Validator Count" size=Heading.H5 marginBottom=8 />
        <Text value={request.askCount |> string_of_int} size=Text.Lg />
      </Col>
      <Col col=Col.Six>
        <Heading value="Sufficient Validator Count" size=Heading.H5 marginBottom=8 />
        <Text value={request.minCount |> string_of_int} size=Text.Lg />
      </Col>
    </Row>;
  };
};

module RequestFailMsg = {
  [@react.component]
  let make = (~request: TxSub.Msg.Request.fail_t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={request.sender} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Oracle Script" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.OracleScript position=TypeID.Subtitle id={request.oracleScriptID} />
        </div>
      </Col>
      <Col mb=24>
        <Heading value="Calldata" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <Text value={request.calldata |> JsBuffer.toHex} color=Colors.gray7 size=Text.Lg />
          <HSpacing size=Spacing.sm />
          <CopyRender width=14 message={request.calldata |> JsBuffer.toHex} />
        </div>
      </Col>
      <Col col=Col.Six mbSm=24>
        <Heading value="Request Validator Count" size=Heading.H5 marginBottom=8 />
        <Text value={request.askCount |> string_of_int} size=Text.Lg />
      </Col>
      <Col col=Col.Six>
        <Heading value="Sufficient Validator Count" size=Heading.H5 marginBottom=8 />
        <Text value={request.minCount |> string_of_int} size=Text.Lg />
      </Col>
    </Row>;
  };
};

module ReportMsg = {
  [@react.component]
  let make = (~report: TxSub.Msg.Report.t) => {
    <Row>
      <Col col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender position=AddressRender.Subtitle address={report.reporter} />
      </Col>
      <Col col=Col.Six mb=24>
        <Heading value="Request ID" size=Heading.H5 marginBottom=8 />
        <TypeID.Request position=TypeID.Subtitle id={report.requestID} />
      </Col>
      <Col>
        <Heading value="Raw Data Report" size=Heading.H5 marginBottom=8 />
        <KVTable
          headers=["External Id", "Exit Code", "Value"]
          rows={
            report.rawReports
            |> Belt_List.map(_, rawReport =>
                 [
                   KVTable.Value(rawReport.externalDataID |> string_of_int),
                   KVTable.Value(rawReport.exitCode |> string_of_int),
                   KVTable.Value(rawReport.data |> JsBuffer.toUTF8),
                 ]
               )
          }
        />
      </Col>
    </Row>;
  };
};
