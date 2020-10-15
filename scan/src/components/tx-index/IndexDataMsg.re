module CreateDataSourceMsg = {
  [@react.component]
  let make = (~dataSource: TxSub.Msg.CreateDataSource.success_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={dataSource.owner} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.DataSource id={dataSource.id} />
          <HSpacing size=Spacing.sm />
          <Text value={dataSource.name} />
        </div>
      </Col.Grid>
    </Row.Grid>;
  };
};

module CreateDataSourceFailMsg = {
  [@react.component]
  let make = (~dataSource: TxSub.Msg.CreateDataSource.fail_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={dataSource.owner} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <Text value={dataSource.name} />
      </Col.Grid>
    </Row.Grid>;
  };
};

module EditDataSourceMsg = {
  [@react.component]
  let make = (~dataSource: TxSub.Msg.EditDataSource.t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={dataSource.owner} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.DataSource id={dataSource.id} />
          <HSpacing size=Spacing.sm />
          <Text value={dataSource.name} />
        </div>
      </Col.Grid>
    </Row.Grid>;
  };
};

module CreateOracleScriptMsg = {
  [@react.component]
  let make = (~oracleScript: TxSub.Msg.CreateOracleScript.success_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={oracleScript.owner} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.OracleScript id={oracleScript.id} />
          <HSpacing size=Spacing.sm />
          <Text value={oracleScript.name} />
        </div>
      </Col.Grid>
    </Row.Grid>;
  };
};

module CreateOracleScriptFailMsg = {
  [@react.component]
  let make = (~oracleScript: TxSub.Msg.CreateOracleScript.fail_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={oracleScript.owner} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <Text value={oracleScript.name} />
      </Col.Grid>
    </Row.Grid>;
  };
};

module EditOracleScriptMsg = {
  [@react.component]
  let make = (~oracleScript: TxSub.Msg.EditOracleScript.t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={oracleScript.owner} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Name" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.OracleScript id={oracleScript.id} />
          <HSpacing size=Spacing.sm />
          <Text value={oracleScript.name} />
        </div>
      </Col.Grid>
    </Row.Grid>;
  };
};

module RequestMsg = {
  [@react.component]
  let make = (~request: TxSub.Msg.Request.success_t) => {
    let calldataKVsOpt = Obi.decode(request.schema, "input", request.calldata);
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={request.sender} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Oracle Script" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.OracleScript id={request.oracleScriptID} />
          <HSpacing size=Spacing.sm />
          <Text value={request.oracleScriptName} />
        </div>
      </Col.Grid>
      <Col.Grid mb=24>
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
           />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Request Validator Count" size=Heading.H5 marginBottom=8 />
        <Text value={request.askCount |> string_of_int} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Sufficient Validator Count" size=Heading.H5 marginBottom=8 />
        <Text value={request.minCount |> string_of_int} />
      </Col.Grid>
    </Row.Grid>;
  };
};

module RequestFailMsg = {
  [@react.component]
  let make = (~request: TxSub.Msg.Request.fail_t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={request.sender} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Oracle Script" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <TypeID.OracleScript id={request.oracleScriptID} />
        </div>
      </Col.Grid>
      <Col.Grid mb=24>
        <Heading value="Calldata" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}>
          <Text value={request.calldata |> JsBuffer.toHex} color=Colors.gray7 />
          <HSpacing size=Spacing.sm />
          <CopyRender width=14 message={request.calldata |> JsBuffer.toHex} />
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Six mbSm=24>
        <Heading value="Request Validator Count" size=Heading.H5 marginBottom=8 />
        <Text value={request.askCount |> string_of_int} />
      </Col.Grid>
      <Col.Grid col=Col.Six>
        <Heading value="Sufficient Validator Count" size=Heading.H5 marginBottom=8 />
        <Text value={request.minCount |> string_of_int} />
      </Col.Grid>
    </Row.Grid>;
  };
};

module ReportMsg = {
  [@react.component]
  let make = (~report: TxSub.Msg.Report.t) => {
    <Row.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Owner" size=Heading.H5 marginBottom=8 />
        <AddressRender address={report.reporter} />
      </Col.Grid>
      <Col.Grid col=Col.Six mb=24>
        <Heading value="Request ID" size=Heading.H5 marginBottom=8 />
        <div className={CssHelper.flexBox()}> <TypeID.Request id={report.requestID} /> </div>
      </Col.Grid>
      <Col.Grid>
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
      </Col.Grid>
    </Row.Grid>;
  };
};
