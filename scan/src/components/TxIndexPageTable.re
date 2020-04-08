module Styles = {
  open Css;

  let addressContainer = width_ => style([width(`px(width_))]);

  let badgeContainer = style([display(`block)]);

  let badge = color =>
    style([
      display(`inlineFlex),
      padding2(~v=`px(8), ~h=`px(10)),
      backgroundColor(color),
      borderRadius(`px(15)),
    ]);

  let hFlex = style([display(`flex), alignItems(`center)]);

  let topicContainer =
    style([display(`flex), justifyContent(`spaceBetween), width(`percent(100.))]);

  let detailContainer = style([display(`flex), maxWidth(`px(360)), justifyContent(`flexEnd)]);

  let firstCol = 0.45;
  let secondCol = 0.50;
  let thirdCol = 1.20;
};

let renderSend = (msg, send: TxSub.Msg.Send.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.blue1)}>
          <Text value="SEND TOKEN" size=Text.Sm spacing={Text.Em(0.07)} color=Colors.blue7 />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="FROM" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={send.fromAddress} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="TO" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={send.toAddress} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value={send.amount |> Coin.toCoinsString} weight=Text.Semibold code=true />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

// TODO: move it to file later.
let renderRequest = (msg, request: TxSub.Msg.Request.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.orange1)}>
          <Text value="REQUEST DATA" size=Text.Sm spacing={Text.Em(0.07)} color=Colors.orange6 />
        </div>
        <VSpacing size=Spacing.md />
        <div className={Styles.badge(Colors.orange1)}> <TypeID.Request id={request.id} /> </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="ORACLE SCRIPT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={request.oracleScriptID} />
          <HSpacing size=Spacing.sm />
          <Text value={request.oracleScriptName} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.md />
      <div className=Styles.hFlex>
        <Text value="CALLDATA" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <HSpacing size=Spacing.md />
        <CopyButton data={request.calldata} />
      </div>
      <VSpacing size=Spacing.md />
      // TODO: Mock calldata
      <KVTable
        tableWidth=480
        rows=[
          [KVTable.Value("crypto_symbol"), KVTable.Value("BTC")],
          [KVTable.Value("aggregation_method"), KVTable.Value("mean")],
          [
            KVTable.Value("data_sources"),
            KVTable.Value("Binance v1, coingecko v1, coinmarketcap v1, band-validator"),
          ],
        ]
      />
      <VSpacing size=Spacing.xl />
      <div className=Styles.topicContainer>
        <Text
          value="REQUEST VALIDATOR COUNT"
          size=Text.Sm
          weight=Text.Thin
          spacing={Text.Em(0.06)}
        />
        <Text value={request.requestedValidatorCount |> string_of_int} weight=Text.Bold />
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text
          value="SUFFICIENT VALIDATOR COUNT"
          size=Text.Sm
          weight=Text.Thin
          spacing={Text.Em(0.06)}
        />
        <Text value={request.sufficientValidatorCount |> string_of_int} weight=Text.Bold />
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="REPORT PERIOD" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value={request.expiration |> string_of_int} weight=Text.Bold code=true />
          <HSpacing size=Spacing.sm />
          <Text value="Blocks" code=true />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderReport = (msg, report: TxSub.Msg.Report.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.orange1)}>
          <Text value="REPORT DATA" size=Text.Sm spacing={Text.Em(0.07)} color=Colors.orange6 />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="REQUEST ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex> <TypeID.Request id={report.requestID} /> </div>
      </div>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.sm />
      <div className=Styles.hFlex>
        <Text value="RAW DATA REPORTS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <HSpacing size=Spacing.md />
      </div>
      <VSpacing size=Spacing.md />
      <KVTable
        tableWidth=480
        headers=["EXTERNAL ID", "VALUE"]
        rows={
          report.dataSet
          |> Belt_List.map(_, rawReport =>
               [
                 KVTable.Value(rawReport.externalDataID |> string_of_int),
                 KVTable.Value(rawReport.data |> JsBuffer._toString(_, "UTF-8")),
               ]
             )
        }
      />
      <VSpacing size=Spacing.lg />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderCreateDataSource = (msg, dataSource: TxSub.Msg.CreateDataSource.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <VSpacing size=Spacing.sm />
        <div className={Styles.badge(Colors.yellow1)}>
          <Text
            value="NEW DATA SOURCE"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.yellow5
          />
        </div>
        <VSpacing size=Spacing.sm />
        <div className={Styles.badge(Colors.yellow1)}>
          <TypeID.DataSource id={dataSource.id} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <Col> <VSpacing size=Spacing.md /> </Col>
      <div className=Styles.topicContainer>
        <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={dataSource.owner} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.DataSource id={dataSource.id} />
          <HSpacing size=Spacing.sm />
          <Text value={dataSource.name} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="FEE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value={dataSource.fee |> Coin.toCoinsString} weight=Text.Bold code=true />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderEditDataSource = (msg, dataSource: TxSub.Msg.EditDataSource.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <VSpacing size=Spacing.sm />
        <div className={Styles.badge(Colors.yellow1)}>
          <Text
            value="EDIT DATA SOURCE"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.yellow5
          />
        </div>
        <VSpacing size=Spacing.sm />
        <div className={Styles.badge(Colors.yellow1)}>
          <TypeID.DataSource id={dataSource.id} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <Col> <VSpacing size=Spacing.md /> </Col>
      <div className=Styles.topicContainer>
        <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={dataSource.owner} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.DataSource id={dataSource.id} />
          <HSpacing size=Spacing.sm />
          <Text value={dataSource.name} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="FEE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value={dataSource.fee |> Coin.toCoinsString} weight=Text.Bold code=true />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderCreateOracleScript = (msg, oracleScript: TxSub.Msg.CreateOracleScript.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <VSpacing size=Spacing.sm />
        <div className={Styles.badge(Colors.pink1)}>
          <Text
            value="NEW ORACLE SCRIPT"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.pink6
          />
        </div>
        <VSpacing size=Spacing.sm />
        <div className={Styles.badge(Colors.pink1)}>
          <TypeID.OracleScript id={oracleScript.id} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <Col> <VSpacing size=Spacing.md /> </Col>
      <div className=Styles.topicContainer>
        <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={oracleScript.owner} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={oracleScript.id} />
          <HSpacing size=Spacing.sm />
          <Text value={oracleScript.name} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderEditOracleScript = (msg, oracleScript: TxSub.Msg.EditOracleScript.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <VSpacing size=Spacing.sm />
        <div className={Styles.badge(Colors.pink1)}>
          <Text
            value="EDIT ORACLE SCRIPT"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.pink6
          />
        </div>
        <VSpacing size=Spacing.sm />
        <div className={Styles.badge(Colors.pink1)}>
          <TypeID.OracleScript id={oracleScript.id} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <Col> <VSpacing size=Spacing.md /> </Col>
      <div className=Styles.topicContainer>
        <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={oracleScript.owner} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={oracleScript.id} />
          <HSpacing size=Spacing.sm />
          <Text value={oracleScript.name} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderAddOracleAddress = (msg, address: TxSub.Msg.AddOracleAddress.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.purple1)}>
          <Text
            value="ADD ORACLE ADDRESS"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.purple6
          />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="VALIDATOR" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={address.validatorMoniker} code=true />
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="REPORTER ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <AddressRender address={address.reporter} />
      </div>
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderRemoveOracleAddress = (msg, address: TxSub.Msg.RemoveOracleAddress.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.purple1)}>
          <Text
            value="REMOVE ORACLE ADDRESS"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.purple6
          />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="VALIDATOR" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={address.validatorMoniker} code=true />
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="REPORTER ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <AddressRender address={address.reporter} />
      </div>
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderCreateValidator = (msg, validator: TxSub.Msg.CreateValidator.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.purple1)}>
          <Text
            value="CREATE VALIDATOR"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.purple6
          />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="MONIKER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={validator.moniker} code=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="IDENTITY" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={validator.identity} code=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="WEBSITE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={validator.website} code=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="DETAILS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.detailContainer>
          <Text value={validator.details} code=true height={Text.Px(16)} align=Text.Right />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="COMMISSION RATE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text
          value={validator.commissionRate->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"}
          code=true
        />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="COMMISSION MAX RATE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text
          value={validator.commissionMaxRate->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"}
          code=true
        />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text
          value="COMMISSION MAX CHANGE"
          size=Text.Sm
          weight=Text.Thin
          spacing={Text.Em(0.06)}
        />
        <Text
          value={validator.commissionMaxChange->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"}
          code=true
        />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="DELAGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <AddressRender address={validator.delegatorAddress} />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <AddressRender address={validator.validatorAddress} validator=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="PUBLIC KEY" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <PubKeyRender pubKey={validator.publicKey} />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="MIN SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text
            value={validator.minSelfDelegation |> Js.Float.toString}
            weight=Text.Semibold
            code=true
          />
          <HSpacing size=Spacing.sm />
          <Text value="BAND" weight=Text.Thin code=true />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value="150.00" weight=Text.Semibold code=true />
          <HSpacing size=Spacing.sm />
          <Text value="BAND" weight=Text.Thin code=true />
        </div>
      </div>
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderEditValidator = (msg, validator: TxSub.Msg.EditValidator.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.purple1)}>
          <Text
            value="EDIT VALIDATOR"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.purple6
          />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxSub.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="MONIKER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={validator.moniker} code=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="IDENTITY" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={validator.identity} code=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="WEBSITE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={validator.website} code=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="DETAILS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.detailContainer>
          <Text value={validator.details} code=true height={Text.Px(16)} align=Text.Right />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="COMMISSION RATE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text
          value={validator.commissionRate->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"}
          code=true
        />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <AddressRender address={validator.sender} validator=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="MIN SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value="100.00" weight=Text.Semibold code=true />
          <HSpacing size=Spacing.sm />
          <Text value="BAND" weight=Text.Thin code=true />
        </div>
      </div>
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderBody = (msg: TxSub.Msg.t) => {
  switch (msg) {
  | Send(send) => renderSend(msg, send)
  | CreateDataSource(dataSource) => renderCreateDataSource(msg, dataSource)
  | EditDataSource(dataSource) => renderEditDataSource(msg, dataSource)
  | CreateOracleScript(oracleScript) => renderCreateOracleScript(msg, oracleScript)
  | EditOracleScript(oracleScript) => renderEditOracleScript(msg, oracleScript)
  | Request(request) => renderRequest(msg, request)
  | Report(report) => renderReport(msg, report)
  | AddOracleAddress(address) => renderAddOracleAddress(msg, address)
  | RemoveOracleAddress(address) => renderRemoveOracleAddress(msg, address)
  | CreateValidator(validator) => renderCreateValidator(msg, validator)
  | EditValidator(validator) => renderEditValidator(msg, validator)
  | FailMessage(_) => "Failed msg" |> React.string
  | Unknown => React.null
  };
};

[@react.component]
let make = (~messages: list(TxSub.Msg.t)) => {
  <>
    <THead>
      <Row>
        <Col> <HSpacing size=Spacing.md /> </Col>
        <Col size=Styles.firstCol>
          <Text
            block=true
            value="MESSAGE TYPE"
            size=Text.Sm
            weight=Text.Semibold
            spacing={Text.Em(0.1)}
            color=Colors.gray5
          />
        </Col>
        <Col size=Styles.secondCol>
          <div>
            <Text
              block=true
              value="CREATOR"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray5
              spacing={Text.Em(0.1)}
            />
          </div>
        </Col>
        <Col size=Styles.thirdCol>
          <div>
            <Text
              block=true
              value="DETAIL"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray5
              spacing={Text.Em(0.1)}
            />
          </div>
        </Col>
        <Col> <HSpacing size=Spacing.md /> </Col>
      </Row>
    </THead>
    {messages
     ->Belt.List.mapWithIndex((index, msg) => {
         <TBody key={index |> string_of_int}> {renderBody(msg)} </TBody>
       })
     ->Array.of_list
     ->React.array}
  </>;
};
