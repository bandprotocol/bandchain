module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(20))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);
  let hFlex = style([display(`flex), alignItems(`center)]);

  let minWidth = mw => style([minWidth(`px(mw))]);

  let topicContainer = h =>
    style([display(`flex), alignItems(`center), width(`percent(100.)), height(`px(h))]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);
  let headerContainer = style([lineHeight(`px(25))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let fillRight = style([marginRight(`auto)]);

  let lowerPannel =
    style([
      width(`percent(100.)),
      height(`px(540)),
      padding(`px(30)),
      display(`flex),
      flexDirection(`column),
      backgroundColor(Colors.white),
      boxShadows([
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), Css.rgba(0, 0, 0, 0.1)),
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(12), Css.rgba(0, 0, 0, 0.03)),
      ]),
      borderRadius(`px(10)),
    ]);
};

[@react.component]
let make = (~reqID) => {
  let requestOpt = RequestHook.get(reqID);
  let requestValidators =
    switch (React.useContext(GlobalContext.context), requestOpt) {
    | (Some(info), Some(request)) =>
      info.validators
      ->Belt_List.keep(validator =>
          request.requestedValidators
          ->Belt_List.has(validator.operatorAddress, (a, b) => a->Address.isEqual(b))
        )
    | _ => []
    };

  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.requestLogo className=Styles.logo />
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
          {switch (requestOpt) {
           | Some(request) =>
             <TimeAgos
               time={request.requestedAtTime}
               size=Text.Md
               weight=Text.Thin
               spacing={Text.Em(0.06)}
               height={Text.Px(18)}
               upper=true
             />

           | None =>
             <Text
               value="???"
               size=Text.Md
               weight=Text.Thin
               spacing={Text.Em(0.06)}
               height={Text.Px(18)}
             />
           }}
        </div>
      </Col>
    </Row>
    {switch (requestOpt) {
     | Some(request) =>
       <>
         <VSpacing size=Spacing.xl />
         <div className=Styles.vFlex>
           <TypeID.Request id={request.id} position=TypeID.Title />
         </div>
         <VSpacing size=Spacing.xl />
         <Row>
           <Col size=2.8>
             <InfoHL
               info={InfoHL.OracleScript(request.oracleScriptID, request.oracleScriptName)}
               header="ORACLE SCRIPT"
             />
           </Col>
           <Col size=3.2>
             <InfoHL header="SENDER" info={InfoHL.Address(request.requester, 280)} />
           </Col>
           <Col size=4.0>
             <InfoHL header="TX HASH" info={InfoHL.TxHash(request.txHash, 385)} />
           </Col>
         </Row>
         <VSpacing size=Spacing.xl />
         <Row>
           <Col>
             <InfoHL info={InfoHL.Validators(requestValidators)} header="REQUEST TO VALIDATORS" />
           </Col>
         </Row>
       </>
     | None => React.null
     }}
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
          <ProgressBar reportedValidators=3 minimumValidators=4 totalValidators=5 />
        </Col>
        <Col size=1.5>
          <div className=Styles.hFlex>
            <div className=Styles.fillRight />
            <Text value={0 |> string_of_int} weight=Text.Bold code=true color=Colors.gray8 /> // Mock
            <HSpacing size=Spacing.sm />
            <Text value="Reported" weight=Text.Regular code=true color=Colors.gray8 />
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
            <TypeID.Block id={ID.Block.ID(35135)} /> // Mock
            <HSpacing size=Spacing.sm />
            <Text
              value="(5 blocks remaining)" // Mock
              weight=Text.Regular
              code=true
              color=Colors.gray8
            />
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
      </div>
      <KVTable
        headers=["EXTERNAL ID", "DATA SOURCE", "PARAM"]
        rows=[
          [
            KVTable.Value("1"),
            KVTable.DataSource(ID.DataSource.ID(12), "Mock Data Source"),
            KVTable.Value("BTC"),
          ],
          [
            KVTable.Value("2"),
            KVTable.DataSource(ID.DataSource.ID(12), "Mock Data Source"),
            KVTable.Value("BTC"),
          ],
          [
            KVTable.Value("3"),
            KVTable.DataSource(ID.DataSource.ID(12), "Mock Data Source"),
            KVTable.Value("BTC"),
          ],
        ]
      />
    </div>
  </div>;
};
