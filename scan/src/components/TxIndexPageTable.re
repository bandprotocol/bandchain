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
};

// TODO: move it to file later.
module CopyButton = {
  open Css;

  [@react.component]
  let make = () => {
    <div
      className={style([
        backgroundColor(Colors.fadeBlue),
        padding2(~h=`px(8), ~v=`px(4)),
        display(`flex),
        width(`px(103)),
        borderRadius(`px(6)),
        cursor(`pointer),
        boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(20, 32, 184, 0.2))),
      ])}>
      <img src=Images.copy className={Css.style([maxHeight(`px(12))])} />
      <HSpacing size=Spacing.sm />
      <Text value="Copy as bytes" size=Text.Sm block=true color=Colors.brightBlue nowrap=true />
    </div>;
  };
};

let renderSend = (msg, send: TxHook.Msg.Send.t) => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=0.4 alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.blue1)}>
          <Text value="SEND TOKEN" size=Text.Sm spacing={Text.Em(0.07)} color=Colors.blue7 />
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
          <Text value={send.amount |> TxHook.Coin.toCoinsString} weight=Text.Bold code=true />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

// TODO: move it to file later.
let renderRequest = (msg: TxHook.Msg.t, request: TxHook.Msg.Request.t) => {
  let requestID =
    msg.events->TxHook.Event.getValueOfKey("request.id")->Belt_Option.getWithDefault("0")
    |> int_of_string;
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=0.4 alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.fadeOrange)}>
          <Text
            value="REQUEST DATA"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.darkOrange
          />
        </div>
        <VSpacing size=Spacing.md />
        <div className={Styles.badge(Colors.fadeOrange)}>
          <TypeID.Request id={ID.Request.ID(requestID)} />
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
        <Text value="ORACLE SCRIPT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={ID.OracleScript.ID(request.oracleScriptID)} />
          <HSpacing size=Spacing.sm />
          <Text value="Mean Platinum Price" />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.md />
      <div className=Styles.hFlex>
        <Text value="CALL DATA" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <HSpacing size=Spacing.md />
        <CopyButton />
      </div>
      <VSpacing size=Spacing.md />
      // TODO: Mock call data
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
    <Col size=0.4 alignSelf=Col.Start>
      <div className=Styles.badgeContainer>
        <div className={Styles.badge(Colors.fadeOrange)}>
          <Text
            value="REPORT DATA"
            size=Text.Sm
            spacing={Text.Em(0.07)}
            color=Colors.darkOrange
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
        <Text value="REQUEST ID" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.Request id={ID.Request.ID(report.requestID)} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.sm />
      <div className=Styles.hFlex>
        <Text value="RAW DATA REPORT" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <HSpacing size=Spacing.md />
        <CopyButton />
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

let renderCreateDataSource = msg => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=0.4 alignSelf=Col.Start>
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
          <TypeID.DataSource id={ID.DataSource.ID(123)} />
        </div>
      </div>
    </Col>
    <Col size=0.6 alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
      </div>
    </Col>
    <Col size=1.3 alignSelf=Col.Start>
      <Col> <VSpacing size=Spacing.md /> </Col>
      <div className=Styles.topicContainer>
        <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(100)}>
          <AddressRender address={msg |> TxHook.Msg.getCreator} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.DataSource id={ID.DataSource.ID(123)} />
          <HSpacing size=Spacing.sm />
          <Text value="Binance Crypto Price" />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="FEE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value="1.5" weight=Text.Bold code=true />
          <HSpacing size=Spacing.sm />
          <Text value="BAND" code=true />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderEditDataSource = msg => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=0.4 alignSelf=Col.Start>
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
          <TypeID.DataSource id={ID.DataSource.ID(123)} />
        </div>
      </div>
    </Col>
    <Col size=0.6 alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(170)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
      </div>
    </Col>
    <Col size=1.3 alignSelf=Col.Start>
      <Col> <VSpacing size=Spacing.md /> </Col>
      <div className=Styles.topicContainer>
        <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={msg |> TxHook.Msg.getCreator} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.DataSource id={ID.DataSource.ID(123)} />
          <HSpacing size=Spacing.sm />
          <Text value="Binance Crypto Price" />
        </div>
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.topicContainer>
        <Text value="FEE" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <Text value="1.5" weight=Text.Bold code=true />
          <HSpacing size=Spacing.sm />
          <Text value="BAND" code=true />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderCreateOracleScript = msg => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=0.4 alignSelf=Col.Start>
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
          <TypeID.OracleScript id={ID.OracleScript.ID(123)} />
        </div>
      </div>
    </Col>
    <Col size=0.6 alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
      </div>
    </Col>
    <Col size=1.3 alignSelf=Col.Start>
      <Col> <VSpacing size=Spacing.md /> </Col>
      <div className=Styles.topicContainer>
        <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={msg |> TxHook.Msg.getCreator} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={ID.OracleScript.ID(123)} />
          <HSpacing size=Spacing.sm />
          <Text value="Binance Crypto Price" />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderEditOracleScript = msg => {
  <Row>
    <Col> <HSpacing size=Spacing.md /> </Col>
    <Col size=0.4 alignSelf=Col.Start>
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
          <TypeID.OracleScript id={ID.OracleScript.ID(123)} />
        </div>
      </div>
    </Col>
    <Col size=0.6 alignSelf=Col.Start>
      <VSpacing size=Spacing.md />
      <div className={Styles.addressContainer(300)}>
        <AddressRender address={msg |> TxHook.Msg.getCreator} />
      </div>
    </Col>
    <Col size=1.3 alignSelf=Col.Start>
      <Col> <VSpacing size=Spacing.md /> </Col>
      <div className=Styles.topicContainer>
        <Text value="OWNER" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className={Styles.addressContainer(300)}>
          <AddressRender address={msg |> TxHook.Msg.getCreator} />
        </div>
      </div>
      <VSpacing size=Spacing.lg />
      <div className=Styles.topicContainer>
        <Text value="NAME" size=Text.Sm weight=Text.Thin spacing={Text.Em(0.06)} />
        <div className=Styles.hFlex>
          <TypeID.OracleScript id={ID.OracleScript.ID(123)} />
          <HSpacing size=Spacing.sm />
          <Text value="Binance Crypto Price" />
        </div>
      </div>
      <VSpacing size=Spacing.md />
    </Col>
    <Col> <HSpacing size=Spacing.md /> </Col>
  </Row>;
};

let renderBody = (msg: TxHook.Msg.t) => {
  switch (msg.action) {
  | Send(send) => renderSend(msg, send)
  | CreateDataSource(_) => renderCreateDataSource(msg)
  | EditDataSource(_) => renderEditDataSource(msg)
  | CreateOracleScript(_) => renderCreateOracleScript(msg)
  | EditOracleScript(_) => renderEditOracleScript(msg)
  | Request(request) => renderRequest(msg, request)
  | Report(report) => renderReport(msg, report)
  | Unknown => React.null
  };
};

[@react.component]
let make = (~messages: list(TxHook.Msg.t)) => {
  <>
    <THead>
      <Row>
        <Col> <HSpacing size=Spacing.md /> </Col>
        <Col size=0.4>
          <Text
            block=true
            value="MESSAGE TYPE"
            size=Text.Sm
            weight=Text.Semibold
            spacing={Text.Em(0.1)}
            color=Colors.grayText
          />
        </Col>
        <Col size=0.6>
          <div>
            <Text
              block=true
              value="CREATOR"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.grayText
              spacing={Text.Em(0.1)}
            />
          </div>
        </Col>
        <Col size=1.3>
          <div>
            <Text
              block=true
              value="DETAIL"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.grayText
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
