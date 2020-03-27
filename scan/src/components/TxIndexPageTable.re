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

// TODO: move it to file later.
module CopyButton = {
  open Css;

  [@react.component]
  let make = (~data) => {
    <div
      className={style([
        backgroundColor(Colors.blue1),
        padding2(~h=`px(8), ~v=`px(4)),
        display(`flex),
        width(`px(103)),
        borderRadius(`px(6)),
        cursor(`pointer),
        boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(20, 32, 184, 0.2))),
      ])}
      onClick={_ => {Copy.copy(data |> JsBuffer.toHex(~with0x=false))}}>
      <img src=Images.copy className={Css.style([maxHeight(`px(12))])} />
      <HSpacing size=Spacing.sm />
      <Text value="Copy as bytes" size=Text.Sm block=true color=Colors.bandBlue nowrap=true />
    </div>;
  };
};

let renderSend = (msg, send: TxHook.Msg.Send.t) => {
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
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
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
          <Text value={send.amount |> TxHook.Coin.toCoinsString} weight=Text.Semibold code=true />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

// TODO: move it to file later.
let renderRequest = (msg, request: TxHook.Msg.Request.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=Styles.firstCol alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.orange1)}>
          <Text value="REQUEST DATA" size=Text.Sm spacing={Text.Em(0.07)} color=Colors.orange6 />
        </div>
        <VSpacing size=Spacing.md />
        <div className={Styles.badge(Colors.orange1)}>
          <TypeID.Request id={ID.Request.ID(request.id)} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="ORACLE SCRIPT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={ID.OracleScript.ID(request.oracleScriptID)} />
          <HSpacing size=Spacing.sm />
          <Text value="Mock oracle script name" />
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
        kv=[
          ("crypto_symbol", "BTC"),
          ("aggregation_method", "mean"),
          ("data_sources", "Binance v1, coingecko v1, coinmarketcap v1, band-validator"),
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

let renderReport = (msg, report: TxHook.Msg.Report.t) => {
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
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="REQUEST ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.Request id={ID.Request.ID(report.requestID)} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.sm />
      <div className=Styles.hFlex>
        <Text value="RAW DATA REPORTS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <HSpacing size=Spacing.md />
      </div>
      <VSpacing size=Spacing.md />
      <KVTable
        header=["EXTERNAL ID", "VALUE"]
        kv={
          report.dataSet
          |> Belt_List.map(_, rawReport =>
               (
                 rawReport.externalDataID |> string_of_int,
                 rawReport.data |> JsBuffer._toString(_, "UTF-8"),
               )
             )
        }
      />
      <VSpacing size=Spacing.lg />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderCreateDataSource = (msg, dataSource: TxHook.Msg.CreateDataSource.t) => {
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
          <TypeID.DataSource id={ID.DataSource.ID(dataSource.id)} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
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
          <TypeID.DataSource id={ID.DataSource.ID(dataSource.id)} />
          <HSpacing size=Spacing.sm />
          <Text value={dataSource.name} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="FEE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value={dataSource.fee |> TxHook.Coin.toCoinsString} weight=Text.Bold code=true />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderEditDataSource = (msg, dataSource: TxHook.Msg.EditDataSource.t) => {
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
          <TypeID.DataSource id={ID.DataSource.ID(dataSource.id)} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
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
          <TypeID.DataSource id={ID.DataSource.ID(dataSource.id)} />
          <HSpacing size=Spacing.sm />
          <Text value={dataSource.name} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="FEE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value={dataSource.fee |> TxHook.Coin.toCoinsString} weight=Text.Bold code=true />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderCreateOracleScript = (msg, oracleScript: TxHook.Msg.CreateOracleScript.t) => {
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
          <TypeID.OracleScript id={ID.OracleScript.ID(oracleScript.id)} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
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
          <TypeID.OracleScript id={ID.OracleScript.ID(oracleScript.id)} />
          <HSpacing size=Spacing.sm />
          <Text value={oracleScript.name} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderEditOracleScript = (msg, oracleScript: TxHook.Msg.EditOracleScript.t) => {
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
          <TypeID.OracleScript id={ID.OracleScript.ID(oracleScript.id)} />
        </div>
      </div>
    </Col>
    <Col size=Styles.secondCol alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
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
          <TypeID.OracleScript id={ID.OracleScript.ID(oracleScript.id)} />
          <HSpacing size=Spacing.sm />
          <Text value={oracleScript.name} />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderAddOracleAddress = (msg, address: TxHook.Msg.AddOracleAddress.t) => {
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
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="VALIDATOR" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={address.validator} code=true />
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="REPORTER ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <AddressRender address={address.reporterAddress} />
      </div>
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderRemoveOracleAddress = (msg, address: TxHook.Msg.RemoveOracleAddress.t) => {
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
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
      </div>
    </Col>
    <Col size=Styles.thirdCol alignSelf=Col.Start>
      <VSpacing size=Spacing.sm />
      <div className=Styles.topicContainer>
        <Text value="VALIDATOR" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <Text value={address.validator} code=true />
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="REPORTER ADDRESS" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <AddressRender address={address.reporterAddress} />
      </div>
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderCreateValidator = (msg, validator: TxHook.Msg.CreateValidator.t) => {
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
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
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
            value={validator.minSelfDelegation |> TxHook.Coin.toCoinsString}
            weight=Text.Semibold
            code=true
          />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text
            value={validator.selfDelegation |> TxHook.Coin.toCoinsString}
            weight=Text.Semibold
            code=true
          />
        </div>
      </div>
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};
/*
 let renderCreateValidator = (msg, validator: TxHook.Msg.EditValidator.t) => {
   <Row>
     <Col> <HSpacing size=Spacing.md /> </Col>
     <Col size=0.51 alignSelf=Col.Start>
       <div className=Styles.badgeContainer>
         <div className={Styles.badge(Colors.purple1)}>
           <Text
             value="Create VALIDATOR"
             size=Text.Sm
             spacing={Text.Em(0.07)}
             color=Colors.purple6
           />
         </div>
       </div>
     </Col>
     <Col size=0.6 alignSelf=Col.Start>
       <VSpacing size=Spacing.sm />
       <div className={Styles.addressContainer(170)}>
         <AddressRender address={msg |> TxHook.Msg.getCreator} />
       </div>
     </Col>
     <Col size=1.3 alignSelf=Col.Start>
       <VSpacing size=Spacing.sm />
       <div className=Styles.topicContainer>
         <Text
           value="MONIKER"
           size=Text.Sm
           weight=Text.Thin
           spacing={Text.Em(0.06)}
         />
         <Text value={validator.moniker} code=true />
       </div>
       <VSpacing size=Spacing.md />
       <div className=Styles.topicContainer>
         <Text
           value="IDENTITY"
           size=Text.Sm
           weight=Text.Thin
           spacing={Text.Em(0.06)}
         />
         <Text value={validator.identity} code=true />
       </div>
       <VSpacing size=Spacing.md />
       <div className=Styles.topicContainer>
         <Text
           value="WEBSITE"
           size=Text.Sm
           weight=Text.Thin
           spacing={Text.Em(0.06)}
         />
         <Text value={validator.website} code=true />
       </div>
       <VSpacing size=Spacing.md />
       <div className=Styles.topicContainer>
         <Text
           value="DETAILS"
           size=Text.Sm
           weight=Text.Thin
           spacing={Text.Em(0.06)}
         />
         <div className=Styles.detailContainer>
           <Text
             value={validator.details}
             code=true
             height={Text.Px(16)}
             align=Text.Right
           />
         </div>
       </div>
       <VSpacing size=Spacing.md />
       <div className=Styles.topicContainer>
         <Text
           value="COMMISSION RATE"
           size=Text.Sm
           weight=Text.Thin
           spacing={Text.Em(0.06)}
         />
         <Text
           value={
             validator.commissionRate->Js.Float.toFixedWithPrecision(~digits=4)
             ++ "%"
           }
           code=true
         />
       </div>
       <VSpacing size=Spacing.md />
       <div className=Styles.topicContainer>
         <Text
           value="MIN SELF DELEGATION"
           size=Text.Sm
           weight=Text.Thin
           spacing={Text.Em(0.06)}
         />
         <div className=Styles.hFlex>
           <Text
             value={validator.minSelfDelegation |> TxHook.Coin.toCoinsString}
             weight=Text.Semibold
             code=true
           />
         </div>
       </div>
     </Col>
     <Col> <HSpacing size=Spacing.md /> </Col>
   </Row>;
 };
 */

let renderEditValidator = (msg, validator: TxHook.Msg.EditValidator.t) => {
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
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
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
        <AddressRender address={validator.validatorAddress} validator=true />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="MIN SELF DELEGATION" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text
            value={validator.minSelfDelegation |> TxHook.Coin.toCoinsString}
            weight=Text.Semibold
            code=true
          />
        </div>
      </div>
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderBody = (msg: TxHook.Msg.t) => {
  switch (msg.action) {
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
  | Unknown => React.null
  };
};

[@react.component]
let make = (~messages: list(TxHook.Msg.t)) => {
  let messages =
    TxHook.Msg.[
      {
        action:
          AddOracleAddress({
            validator: "node-validator-2",
            reporterAddress: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
            sender: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
          }),
        events: [],
      },
      {
        action:
          RemoveOracleAddress({
            validator: "node-validator-2",
            reporterAddress: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
            sender: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
          }),
        events: [],
      },
      {
        action:
          CreateValidator({
            moniker: "Bitkrub-node-validator",
            identity: "Bitkrub The next generation digital asset exchange",
            website: "https://www.bitkrub.com/",
            details: "CEO Changpeng Zhao had previously founded Fusion Systems in 2005 in Shanghai; the company built high-frequency trading systems for brokers.",
            commissionRate: 3.5250,
            commissionMaxRate: 10.0000,
            commissionMaxChange: 0.1000,
            delegatorAddress: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
            validatorAddress:
              "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec" |> Address.fromBech32,
            publicKey:
              "bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj"
              |> PubKey.fromBech32,
            minSelfDelegation: [{denom: "uband", amount: 100.00}],
            selfDelegation: [{denom: "uband", amount: 150.00}],
            sender: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
          }),
        events: [],
      },
      {
        action:
          EditValidator({
            moniker: "Bitkrub-node-validator",
            identity: "Bitkrub The next generation digital asset exchange",
            website: "https://www.bitkrub.com/",
            details: "CEO Changpeng Zhao had previously founded Fusion Systems in 2005 in Shanghai; the company built high-frequency trading systems for brokers.",
            commissionRate: 3.5250,
            validatorAddress:
              "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec" |> Address.fromBech32,
            minSelfDelegation: [{denom: "uband", amount: 100.00}],
            sender: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
          }),
        events: [],
      },
      ...messages,
    ];
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
