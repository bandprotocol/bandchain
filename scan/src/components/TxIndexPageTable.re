module Styles = {
  open Css;

  let addressContainer = width_ => style([width(`px(width_))]);

  let badgeContainer = style([display(`block), marginTop(`px(-4))]);

  let badge = color =>
    style([
      display(`inlineFlex),
      padding2(~v=`px(5), ~h=`px(8)),
      backgroundColor(color),
      borderRadius(`px(50)),
    ]);

  let hFlex = style([display(`flex), alignItems(`center)]);

  let topicContainer =
    style([
      display(`flex),
      justifyContent(`spaceBetween),
      width(`percent(100.)),
      lineHeight(`px(16)),
      alignItems(`center),
    ]);

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

  let failIcon = style([width(`px(16)), height(`px(16))]);

  let failedMessageDetails =
    style([
      display(`flex),
      width(`px(120)),
      alignItems(`center),
      justifyContent(`spaceBetween),
    ]);

  let separatorLine =
    style([
      borderStyle(`none),
      backgroundColor(Colors.gray9),
      height(`px(1)),
      margin2(~v=`px(10), ~h=`auto),
    ]);
};

let renderSend = (send: TxSub.Msg.Send.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
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
      <AmountRender coins={send.amount} pos=AmountRender.TxIndex />
    </div>
    <VSpacing size=Spacing.lg />
  </Col>;
};
// TODO: move it to file later.
let renderRequest = (request: TxSub.Msg.Request.t) => {
  let calldataKVsOpt = Obi.decode(request.schema, "input", request.calldata);
  <Col size=Styles.thirdCol alignSelf=Col.Start>
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
      <CopyButton
        data={request.calldata |> JsBuffer.toHex(~with0x=false)}
        title="Copy as bytes"
        width=125
      />
    </div>
    <VSpacing size=Spacing.md />
    {switch (calldataKVsOpt) {
     | Some(calldataKVs) =>
       <KVTable
         tableWidth=480
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
    <VSpacing size=Spacing.xl />
    <div className=Styles.topicContainer>
      <Text value="ASK COUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={request.askCount |> string_of_int} weight=Text.Bold />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="MIN COUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={request.minCount |> string_of_int} weight=Text.Bold />
    </div>
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.lg />
  </Col>;
};

let renderReport = (report: TxSub.Msg.Report.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
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
      headers=["EXTERNAL ID", "EXIT CODE", "VALUE"]
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
    <VSpacing size=Spacing.lg />
  </Col>;
};

let renderCreateDataSource = (dataSource: TxSub.Msg.CreateDataSource.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
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
    // <div className=Styles.topicContainer>
    //   <Text value="FEE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
    //   <AmountRender coins={dataSource.fee} pos=AmountRender.TxIndex />
    // </div>
    <VSpacing size=Spacing.md />
  </Col>;
};

let renderEditDataSource = (dataSource: TxSub.Msg.EditDataSource.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={dataSource.owner} />
      </div>
    </div>
    // TODO : Redesign edit messages
    // <VSpacing size=Spacing.lg />
    // <div className=Styles.topicContainer>
    //   <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
    //   <div className=Styles.hFlex>
    //     <TypeID.DataSource id={dataSource.id} />
    //     <HSpacing size=Spacing.sm />
    //     <Text value={dataSource.name} />
    //   </div>
    // </div>
    // <VSpacing size=Spacing.md />
    //  <div className=Styles.topicContainer>
    //    <Text value="FEE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
    //    <AmountRender coins={dataSource.fee} pos=AmountRender.TxIndex />
    //  </div>
    <VSpacing size=Spacing.md />
  </Col>;
};

let renderCreateOracleScript = (oracleScript: TxSub.Msg.CreateOracleScript.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
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
    <div className=Styles.topicContainer>
      <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={oracleScript.owner} />
      </div>
    </div>
    // TODO : Redesign edit messages
    // <VSpacing size=Spacing.lg />
    // <div className=Styles.topicContainer>
    //   <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
    //   <div className=Styles.hFlex>
    //     <TypeID.OracleScript id={oracleScript.id} />
    //     <HSpacing size=Spacing.sm />
    //     <Text value={oracleScript.name} />
    //   </div>
    // </div>
    <VSpacing size=Spacing.md />
  </Col>;
};

let renderAddReporter = (address: TxSub.Msg.AddReporter.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
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

let renderRemoveReporter = (address: TxSub.Msg.RemoveReporter.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
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
      <Text value="DELEGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AddressRender address={validator.delegatorAddress} />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AddressRender address={validator.validatorAddress} accountType=`validator />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="PUBLIC KEY" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <PubKeyRender pubKey={validator.publicKey} />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="MIN SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins=[validator.minSelfDelegation] />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins=[validator.selfDelegation] />
    </div>
  </Col>;
};

let renderEditValidator = (validator: TxSub.Msg.EditValidator.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="MONIKER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.moniker == Config.doNotModify ? "Unchanged" : validator.moniker}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="IDENTITY" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.identity == Config.doNotModify ? "Unchanged" : validator.identity}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="WEBSITE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={validator.website == Config.doNotModify ? "Unchanged" : validator.website}
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="DETAILS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className=Styles.detailContainer>
        <Text
          value={validator.details == Config.doNotModify ? "Unchanged" : validator.details}
          code=true
          height={Text.Px(16)}
          align=Text.Right
        />
      </div>
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="COMMISSION RATE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text
        value={
          switch (validator.commissionRate) {
          | Some(rate) => rate->Js.Float.toFixedWithPrecision(~digits=4) ++ "%"
          | None => "Unchanged"
          }
        }
        code=true
      />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AddressRender address={validator.sender} accountType=`validator />
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.topicContainer>
      <Text value="MIN SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      {switch (validator.minSelfDelegation) {
       | Some(minSelfDelegation') => <AmountRender coins=[minSelfDelegation'] />
       | None => <Text value="Unchanged" code=true />
       }}
    </div>
  </Col>;
};

let renderCreateClient = (info: TxSub.Msg.CreateClient.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
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
     | AcknowledgementMsg({acknowledgement}) =>
       <>
         <VSpacing size=Spacing.md />
         <div className=Styles.topicContainer>
           <Text value="ACKNOWLEDGEMENT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=acknowledgement code=true />
         </div>
       </>
     | TimeoutMsg({nextSequenceReceive}) =>
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
     | ConnectionOpenInitMsg({clientID}) =>
       <>
         <div className=Styles.topicContainer>
           <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=clientID code=true />
         </div>
       </>
     | ConnectionOpenTryMsg({clientID}) =>
       <>
         <div className=Styles.topicContainer>
           <Text value="CLIENT ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
           <Text value=clientID code=true />
         </div>
       </>
     | _ => React.null
     }}
  </Col>;
};

let renderDelegate = (delegation: TxSub.Msg.Delegate.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="DELEGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={delegation.delegatorAddress} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={delegation.validatorAddress} accountType=`validator />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins=[delegation.amount] pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderUndelegate = (delegation: TxSub.Msg.Undelegate.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="DELEGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={delegation.delegatorAddress} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={delegation.validatorAddress} accountType=`validator />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins=[delegation.amount] pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderRedelegate = (delegation: TxSub.Msg.Redelegate.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="DELEGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={delegation.delegatorAddress} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="SOURCE ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={delegation.validatorSourceAddress} accountType=`validator />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="DESTINATION ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={delegation.validatorDestinationAddress} accountType=`validator />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins=[delegation.amount] pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderWithdrawReward = (withdrawal: TxSub.Msg.WithdrawReward.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="DELEGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={withdrawal.delegatorAddress} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={withdrawal.validatorAddress} accountType=`validator />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins={withdrawal.amount} pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderUnjail = (unjail: TxSub.Msg.Unjail.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={unjail.address} accountType=`validator />
      </div>
    </div>
  </Col>;
};
let renderSubmitProposal = (proposal: TxSub.Msg.SubmitProposal.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="TITLE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={proposal.title} code=true />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="DESCRIPTION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={proposal.description} code=true />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="PROPOSER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={proposal.proposer} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins={proposal.initialDeposit} pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderDeposit = (deposit: TxSub.Msg.Deposit.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="DEPOSITOR" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={deposit.depositor} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="PROPOSAL ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <Text value={deposit.proposalID |> string_of_int} code=true />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins={deposit.amount} pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderSetWithdrawAddress = (set: TxSub.Msg.SetWithdrawAddress.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="DELEGATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={set.delegatorAddress} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="WITHDRAW ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={set.withdrawAddress} />
      </div>
    </div>
  </Col>;
};
let renderWithdrawCommission = (withdrawal: TxSub.Msg.WithdrawCommission.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={withdrawal.validatorAddress} accountType=`validator />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <AmountRender coins={withdrawal.amount} pos=AmountRender.TxIndex />
    </div>
  </Col>;
};

let renderVote = (vote: TxSub.Msg.Vote.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VOTER ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={vote.voterAddress} />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="PROPOSAL ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className=Styles.hFlex> <Text value={vote.proposalID |> string_of_int} /> </div>
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="OPTION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className=Styles.hFlex> <Text value={vote.option} /> </div>
    </div>
  </Col>;
};

let renderMultiSend = (tx: TxSub.Msg.MultiSend.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="Inputs" size=Text.Md weight=Text.Semibold spacing={Text.Em(0.06)} />
    </div>
    <VSpacing size=Spacing.lg />
    {tx.inputs
     ->Belt_List.mapWithIndex((idx, input) =>
         <div key={idx |> string_of_int}>
           <div className=Styles.topicContainer>
             <Text value="FROM" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
             <div className={Styles.addressContainer(300)}>
               <AddressRender address={input.address} />
             </div>
           </div>
           <VSpacing size=Spacing.lg />
           <div className=Styles.topicContainer>
             <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
             <AmountRender coins={input.coins} pos=AmountRender.TxIndex />
           </div>
           <VSpacing size=Spacing.lg />
           <VSpacing size=Spacing.md />
         </div>
       )
     ->Belt_List.toArray
     ->React.array}
    <VSpacing size=Spacing.lg />
    <div className=Styles.topicContainer>
      <Text value="Outputs" size=Text.Md weight=Text.Semibold spacing={Text.Em(0.06)} />
    </div>
    <VSpacing size=Spacing.lg />
    {tx.outputs
     ->Belt_List.mapWithIndex((idx, output) =>
         <div key={idx |> string_of_int}>
           <div className=Styles.topicContainer>
             <Text value="TO" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
             <div className={Styles.addressContainer(300)}>
               <AddressRender address={output.address} />
             </div>
           </div>
           <VSpacing size=Spacing.lg />
           <div className=Styles.topicContainer>
             <Text value="AMOUNT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
             <AmountRender coins={output.coins} pos=AmountRender.TxIndex />
           </div>
           <VSpacing size=Spacing.lg />
           <VSpacing size=Spacing.lg />
         </div>
       )
     ->Belt_List.toArray
     ->React.array}
  </Col>;
};

let renderActivate = (activate: TxSub.Msg.Activate.t) => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="VALIDATOR ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={activate.validatorAddress} accountType=`validator />
      </div>
    </div>
  </Col>;
};

let renderFailMessage = () => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="MESSAGE FAILED" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <div className=Styles.hFlex> <img src=Images.fail className=Styles.failIcon /> </div>
    </div>
  </Col>;
};

let renderUnknownMessage = () => {
  <Col size=Styles.thirdCol alignSelf=Col.Start>
    <div className=Styles.topicContainer>
      <Text value="UNKNOWN MESSAGE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
      <img src=Images.fail className=Styles.failIcon />
    </div>
  </Col>;
};

let renderBody = (msg: TxSub.Msg.t) =>
  switch (msg) {
  | SendMsg(send) => renderSend(send)
  | CreateDataSourceMsg(dataSource) => renderCreateDataSource(dataSource)
  | EditDataSourceMsg(dataSource) => renderEditDataSource(dataSource)
  | CreateOracleScriptMsg(oracleScript) => renderCreateOracleScript(oracleScript)
  | EditOracleScriptMsg(oracleScript) => renderEditOracleScript(oracleScript)
  | RequestMsg(request) => renderRequest(request)
  | ReportMsg(report) => renderReport(report)
  | AddReporterMsg(address) => renderAddReporter(address)
  | RemoveReporterMsg(address) => renderRemoveReporter(address)
  | CreateValidatorMsg(validator) => renderCreateValidator(validator)
  | EditValidatorMsg(validator) => renderEditValidator(validator)
  | CreateClientMsg(info) => renderCreateClient(info)
  | UpdateClientMsg(info) => renderUpdateClient(info)
  | SubmitClientMisbehaviourMsg(info) => renderSubmitClientMisbehaviour(info)
  | ConnectionOpenInitMsg(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenTryMsg(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenAckMsg(info) => renderConnectionVariant(msg, info.common)
  | ConnectionOpenConfirmMsg(info) => renderConnectionVariant(msg, info.common)
  | ChannelOpenInitMsg(info) => renderChannelVariant(info.common)
  | ChannelOpenTryMsg(info) => renderChannelVariant(info.common)
  | ChannelOpenAckMsg(info) => renderChannelVariant(info.common)
  | ChannelOpenConfirmMsg(info) => renderChannelVariant(info.common)
  | ChannelCloseInitMsg(info) => renderChannelVariant(info.common)
  | ChannelCloseConfirmMsg(info) => renderChannelVariant(info.common)
  | PacketMsg(info) => renderPacketVariant(msg, info.common)
  | AcknowledgementMsg(info) => renderPacketVariant(msg, info.common)
  | TimeoutMsg(info) => renderPacketVariant(msg, info.common)
  | DelegateMsg(delegation) => renderDelegate(delegation)
  | UndelegateMsg(delegation) => renderUndelegate(delegation)
  | RedelegateMsg(delegation) => renderRedelegate(delegation)
  | WithdrawRewardMsg(withdrawal) => renderWithdrawReward(withdrawal)
  | WithdrawCommissionMsg(withdrawal) => renderWithdrawCommission(withdrawal)
  | UnjailMsg(unjail) => renderUnjail(unjail)
  | SetWithdrawAddressMsg(set) => renderSetWithdrawAddress(set)
  | SubmitProposalMsg(proposal) => renderSubmitProposal(proposal)
  | DepositMsg(deposit) => renderDeposit(deposit)
  | VoteMsg(vote) => renderVote(vote)
  | MultiSendMsg(tx) => renderMultiSend(tx)
  | ActivateMsg(activate) => renderActivate(activate)
  | FailMsg(_) => renderFailMessage()
  | UnknownMsg => renderUnknownMessage()
  };

module THead = {
  [@react.component]
  let make = () => {
    <THead.Grid>
      <Row.Grid alignItems=Row.Center>
        <Col.Grid col=Col.Two>
          <Text block=true value="Message Type" weight=Text.Semibold color=Colors.gray7 />
        </Col.Grid>
        <Col.Grid col=Col.Three>
          <Text block=true value="Creator" weight=Text.Semibold color=Colors.gray7 />
        </Col.Grid>
        <Col.Grid col=Col.Seven>
          <Text block=true value="Detail" weight=Text.Semibold color=Colors.gray7 />
        </Col.Grid>
      </Row.Grid>
    </THead.Grid>;
  };
};

[@react.component]
let make = (~messages: list(TxSub.Msg.t)) => {
  <>
    <THead />
    {messages
     ->Belt.List.mapWithIndex((index, msg) => {
         let theme = msg |> TxSub.Msg.getBadgeTheme;
         //TODO: Change index to be uniqe something
         <TBody.Grid key={index |> string_of_int} paddingH={`px(24)}>
           <Row.Grid alignItems=Row.Start>
             <Col.Grid col=Col.Two>
               <div className=Styles.badgeContainer>
                 <div className={Styles.badge(theme.bgColor)}>
                   <Text value={theme.text} size=Text.Xs color={theme.textColor} />
                 </div>
                 {switch (msg) {
                  | CreateDataSourceMsg(dataSource) =>
                    <>
                      <hr className=Styles.separatorLine />
                      <TypeID.DataSource id={dataSource.id} />
                    </>
                  | EditDataSourceMsg(dataSource) =>
                    <>
                      <hr className=Styles.separatorLine />
                      <TypeID.DataSource id={dataSource.id} />
                    </>
                  | CreateOracleScriptMsg(oracleScript) =>
                    <>
                      <hr className=Styles.separatorLine />
                      <TypeID.OracleScript id={oracleScript.id} />
                    </>
                  | EditOracleScriptMsg(oracleScript) =>
                    <>
                      <hr className=Styles.separatorLine />
                      <TypeID.OracleScript id={oracleScript.id} />
                    </>
                  | RequestMsg(request) =>
                    <> <hr className=Styles.separatorLine /> <TypeID.Request id={request.id} /> </>
                  | _ => React.null
                  }}
               </div>
             </Col.Grid>
             <Col.Grid col=Col.Three>
               <div className={Styles.addressContainer(170)}>
                 <AddressRender address={msg |> TxSub.Msg.getCreator} />
               </div>
             </Col.Grid>
             <Col.Grid col=Col.Seven> {renderBody(msg)} </Col.Grid>
           </Row.Grid>
         </TBody.Grid>;
       })
     ->Array.of_list
     ->React.array}
  </>;
};

module Loading = {
  [@react.component]
  let make = () => {
    <>
      <THead />
      <TBody.Grid paddingH={`px(24)}>
        <Row.Grid>
          <Col.Grid col=Col.Two> <LoadingCensorBar width=95 height=25 /> </Col.Grid>
          <Col.Grid col=Col.Three>
            <VSpacing size=Spacing.sm />
            <LoadingCensorBar width=160 height=15 />
          </Col.Grid>
          <Col.Grid col=Col.Seven>
            <VSpacing size=Spacing.sm />
            <div className=Styles.topicContainer>
              <LoadingCensorBar width=60 height=15 />
              <LoadingCensorBar width=310 height=15 />
            </div>
            <VSpacing size=Spacing.lg />
            <div className=Styles.topicContainer>
              <LoadingCensorBar width=75 height=15 />
              <LoadingCensorBar width=240 height=15 />
            </div>
            <VSpacing size=Spacing.lg />
            <div className=Styles.topicContainer>
              <LoadingCensorBar width=95 height=15 />
              <LoadingCensorBar width=280 height=15 />
            </div>
          </Col.Grid>
        </Row.Grid>
      </TBody.Grid>
    </>;
  };
};
