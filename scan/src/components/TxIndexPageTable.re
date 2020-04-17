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

  let hashContainer =
    style([
      display(`flex),
      maxWidth(`px(350)),
      justifyContent(`flexEnd),
      wordBreak(`breakAll),
    ]);

  let firstCol = 0.45;
  let secondCol = 0.50;
  let thirdCol = 1.20;
};

let renderSend = (send: TxSub.Msg.Send.t) => {
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
  </Col>;
};

// TODO: move it to file later.
let renderRequest = (request: TxSub.Msg.Request.t) => {
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
    <VSpacing size=Spacing.lg />
  </Col>;
};

let renderReport = (report: TxSub.Msg.Report.t) => {
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
  </Col>;
};

let renderCreateDataSource = (dataSource: TxSub.Msg.CreateDataSource.t) => {
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
  </Col>;
};

let renderEditDataSource = (dataSource: TxSub.Msg.EditDataSource.t) => {
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
  </Col>;
};

let renderCreateOracleScript = (oracleScript: TxSub.Msg.CreateOracleScript.t) => {
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
  </Col>;
};

let renderEditOracleScript = (oracleScript: TxSub.Msg.EditOracleScript.t) => {
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
  </Col>;
};

let renderAddOracleAddress = (address: TxSub.Msg.AddOracleAddress.t) => {
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
  </Col>;
};

let renderRemoveOracleAddress = (address: TxSub.Msg.RemoveOracleAddress.t) => {
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
  </Col>;
};

let renderCreateValidator = (validator: TxSub.Msg.CreateValidator.t) => {
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
      <Text value="COMMISSION MAX CHANGE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
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
  </Col>;
};

let renderEditValidator = (validator: TxSub.Msg.EditValidator.t) => {
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
  </Col>;
};

let renderCreateClient = (info: TxSub.Msg.CreateClient.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <VSpacing size=Spacing.sm />
    <div className=Styles.topicContainer>
      <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.clientID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="TRUSTING PERIOD" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.trustingPeriod |> MomentRe.Duration.toISOString} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="UNBOUNDING PERIOD" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.unbondingPeriod |> MomentRe.Duration.toISOString} code=true />
    </div>
  </Col>;
};

let renderUpdateClient = (info: TxSub.Msg.UpdateClient.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <VSpacing size=Spacing.sm />
    <div className=Styles.topicContainer>
      <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.clientID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text
        value="VALIDATOR HASH"
        size=Text.Sm
        weight=Text.Thin
        height={Text.Px(16)}
        spacing={Text.Em(0.06)}
      />
      <div className=Styles.hashContainer>
        <Text
          value={info.validatorHash |> Hash.toHex}
          code=true
          height={Text.Px(16)}
          align=Text.Right
        />
      </div>
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text
        value="PREVIOUS VALIDATOR HASH"
        size=Text.Sm
        weight=Text.Thin
        height={Text.Px(16)}
        spacing={Text.Em(0.06)}
      />
      <div className=Styles.hashContainer>
        <Text
          value={info.prevValidatorHash |> Hash.toHex}
          code=true
          height={Text.Px(16)}
          align=Text.Right
        />
      </div>
    </div>
  </Col>;
};

let renderSubmitClientMisbehaviour = (info: TxSub.Msg.SubmitClientMisbehaviour.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <VSpacing size=Spacing.sm />
    <div className=Styles.topicContainer>
      <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.clientID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={info.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text
        value="VALIDATOR HASH"
        size=Text.Sm
        weight=Text.Thin
        height={Text.Px(16)}
        spacing={Text.Em(0.06)}
      />
      <div className=Styles.hashContainer>
        <Text
          value={info.validatorHash |> Hash.toHex}
          code=true
          height={Text.Px(16)}
          align=Text.Right
        />
      </div>
    </div>
  </Col>;
};

let renderPacketVariant = (msg: TxSub.Msg.t, common: TxSub.Msg.Packet.common_t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <VSpacing size=Spacing.sm />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="SEQUENCE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.sequence |> string_of_int} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="SOURCE PORT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.sourcePort} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="SOURCE CHANNEL" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.sourceChannel} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="DESTINATION PORT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.destinationPort} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="DESTINATION CHANNEL" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.destinationChannel} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="TIMEOUT HEIGHT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.timeoutHeight |> string_of_int} code=true />
    </div>
    {switch (msg) {
     | Acknowledgement({acknowledgement}) =>
       <>
         <VSpacing size=Spacing.md />
         <div className=Styles.topicContainer>
           <Text value="ACKNOWLEDGEMENT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=acknowledgement code=true />
         </div>
       </>
     | Timeout({nextSequenceReceive}) =>
       <>
         <VSpacing size=Spacing.md />
         <div className=Styles.topicContainer>
           <Text value="ACKNOWLEDGEMENT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value={nextSequenceReceive |> string_of_int} code=true />
         </div>
       </>
     | _ => React.null
     }}
  </Col>;
};

let renderChannelVariant = (common: TxSub.Msg.ChannelCommon.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <VSpacing size=Spacing.sm />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="PORT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.portID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CHANNEL ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.channelID} code=true />
    </div>
  </Col>;
};

let renderConnectionVariant = (msg: TxSub.Msg.t, common: TxSub.Msg.ConnectionCommon.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <VSpacing size=Spacing.sm />
    <div className=Styles.topicContainer>
      <Text value="CHAIN ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.chainID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="CONNECTION ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={common.connectionID} code=true />
    </div>
    <VSpacing size=Spacing.md />
    {switch (msg) {
     | ConnectionOpenInit({clientID, consensusHeight}) =>
       <>
         <div className=Styles.topicContainer>
           <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=clientID code=true />
         </div>
         <VSpacing size=Spacing.md />
         <div className=Styles.topicContainer>
           <Text value="CONSENSUS HEIGHT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <TypeID.Block id=consensusHeight />
         </div>
       </>
     | ConnectionOpenTry({clientID, consensusHeight}) =>
       <>
         <div className=Styles.topicContainer>
           <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=clientID code=true />
         </div>
         <VSpacing size=Spacing.md />
         <div className=Styles.topicContainer>
           <Text value="CONSENSUS HEIGHT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <TypeID.Block id=consensusHeight />
         </div>
       </>
     | ConnectionOpenAck({consensusHeight}) =>
       <div className=Styles.topicContainer>
         <Text value="CONSENSUS HEIGHT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
         <TypeID.Block id=consensusHeight />
       </div>
     | _ => React.null
     }}
  </Col>;
};

let renderBody = (msg: TxSub.Msg.t) => {
  switch (msg) {
  | Send(send) => renderSend(send)
  | CreateDataSource(dataSource) => renderCreateDataSource(dataSource)
  | EditDataSource(dataSource) => renderEditDataSource(dataSource)
  | CreateOracleScript(oracleScript) => renderCreateOracleScript(oracleScript)
  | EditOracleScript(oracleScript) => renderEditOracleScript(oracleScript)
  | Request(request) => renderRequest(request)
  | Report(report) => renderReport(report)
  | AddOracleAddress(address) => renderAddOracleAddress(address)
  | RemoveOracleAddress(address) => renderRemoveOracleAddress(address)
  | CreateValidator(validator) => renderCreateValidator(validator)
  | EditValidator(validator) => renderEditValidator(validator)
  | CreateClient(info) => renderCreateClient(info)
  | UpdateClient(info) => renderUpdateClient(info)
  | SubmitClientMisbehaviour(info) => renderSubmitClientMisbehaviour(info)
  | ConnectionOpenInit(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenTry(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenAck(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenConfirm(info) => renderConnectionVariant(msg, info.common)
  | ChannelOpenInit(info) => renderChannelVariant(info.common)
  | ChannelOpenTry(info) => renderChannelVariant(info.common)
  | ChannelOpenAck(info) => renderChannelVariant(info.common)
  | ChannelOpenConfirm(info) => renderChannelVariant(info.common)
  | ChannelCloseInit(info) => renderChannelVariant(info.common)
  | ChannelCloseConfirm(info) => renderChannelVariant(info.common)
  | Packet(info) => renderPacketVariant(msg, info.common)
  | Acknowledgement(info) => renderPacketVariant(msg, info.common)
  | Timeout(info) => renderPacketVariant(msg, info.common)
  | FailMessage(_) => "Failed msg" |> React.string
  | _ => React.null
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
         let theme = msg |> TxSub.Msg.getBadgeTheme;
         <TBody key={index |> string_of_int}>
           <Row>
             <Col> <HSpacing size=Spacing.md /> </Col>
             <Col size=Styles.firstCol alignSelf=Col.Start>
               <div className=Styles.badgeContainer>
                 <div className={Styles.badge(theme.bgColor)}>
                   <Text
                     value={theme.text}
                     size=Text.Sm
                     spacing={Text.Em(0.07)}
                     color={theme.textColor}
                   />
                 </div>
                 <VSpacing size=Spacing.sm />
                 {switch (msg) {
                  | CreateDataSource(dataSource) =>
                    <div className={Styles.badge(theme.bgColor)}>
                      <TypeID.DataSource id={dataSource.id} />
                    </div>
                  | EditDataSource(dataSource) =>
                    <div className={Styles.badge(theme.bgColor)}>
                      <TypeID.DataSource id={dataSource.id} />
                    </div>
                  | CreateOracleScript(oracleScript) =>
                    <div className={Styles.badge(theme.bgColor)}>
                      <TypeID.OracleScript id={oracleScript.id} />
                    </div>
                  | EditOracleScript(oracleScript) =>
                    <div className={Styles.badge(theme.bgColor)}>
                      <TypeID.OracleScript id={oracleScript.id} />
                    </div>
                  | Request(request) =>
                    <div className={Styles.badge(theme.bgColor)}>
                      <TypeID.Request id={request.id} />
                    </div>
                  | _ => React.null
                  }}
               </div>
             </Col>
             <Col size=Styles.secondCol alignSelf=Col.Start>
               <VSpacing size=Spacing.sm />
               <div className={Styles.addressContainer(170)}>
                 <AddressRender address={msg |> TxSub.Msg.getCreator} />
               </div>
             </Col>
             {renderBody(msg)}
             <Col> <HSpacing size=Spacing.md /> </Col>
           </Row>
         </TBody>;
       })
     ->Array.of_list
     ->React.array}
  </>;
};
