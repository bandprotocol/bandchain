module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(50))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let txTypeBadge =
    style([
      paddingLeft(`px(12)),
      paddingRight(`px(12)),
      paddingTop(`px(5)),
      paddingBottom(`px(5)),
      backgroundColor(Colors.lightBlue),
      borderRadius(`px(15)),
    ]);

  let msgAmount =
    style([borderRadius(`percent(50.)), padding(`px(3)), backgroundColor(Colors.lightGray)]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let addressContainer = style([marginTop(`px(15))]);

  let successBadge =
    style([
      backgroundColor(`hex("D7FFEC")),
      borderRadius(`px(6)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      width(`px(120)),
      height(`px(40)),
    ]);

  let checkLogo = style([marginRight(`px(10))]);

  let seperatorLine =
    style([
      width(`percent(100.)),
      height(`pxFloat(1.4)),
      backgroundColor(Colors.lightGray),
      display(`flex),
    ]);
};

[@react.component]
let make = (~txHash) => {
  let txOpt = TxHook.atHash(txHash);
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="TRANSACTION"
            weight=Text.Semibold
            size=Text.Lg
            nowrap=true
            color=Colors.grayHeader
            block=true
          />
          <HSpacing size=Spacing.sm />
          <div className=Styles.txTypeBadge>
            <Text value="DATA REQUEST" block=true size=Text.Sm />
          </div>
          <HSpacing size=Spacing.sm />
          <div className=Styles.msgAmount> <Text value="+1" block=true size=Text.Sm /> </div>
          <div className=Styles.seperatedLine />
          <Text value="51 MINUTES AGO" />
        </div>
      </Col>
      <Col>
        <div className=Styles.successBadge>
          <img src=Images.checkIcon className=Styles.checkLogo />
          <Text value="Success" size=Text.Lg weight=Text.Semibold color=Colors.darkGreen />
        </div>
      </Col>
    </Row>
    <div className=Styles.addressContainer>
      <Text
        value={txHash |> Hash.toHex(~with0x=true)}
        size=Text.Xxl
        weight=Text.Bold
        nowrap=true
      />
    </div>
    <VSpacing size=Spacing.xl />
    <Row>
      {switch (txOpt) {
       | Some(tx) =>
         <>
           <Col size=1.> <InfoHL info={InfoHL.Height(tx.blockHeight)} header="HEIGHT" /> </Col>
           <Col size=1.>
             <InfoHL info={InfoHL.Count(tx.messages |> Belt_List.size)} header="MESSAGES" />
           </Col>
           <Col size=2.>
             <InfoHL info={InfoHL.Timestamp(tx.timestamp)} header="TIMESTAMP" />
           </Col>
           <Col size=2.5> <InfoHL info={InfoHL.Text("FREE")} header="FEE" /> </Col>
         </>
       | None =>
         <>
           <Col size=1.> <InfoHL info={InfoHL.Text("?")} header="HEIGHT" /> </Col>
           <Col size=1.> <InfoHL info={InfoHL.Text("?")} header="MESSAGES" /> </Col>
           <Col size=2.> <InfoHL info={InfoHL.Text("?")} header="TIMESTAMP" /> </Col>
           <Col size=2.5> <InfoHL info={InfoHL.Text("?")} header="FEE" /> </Col>
         </>
       }}
    </Row>
    <VSpacing size=Spacing.xl />
    {switch (txOpt) {
     | Some(tx) =>
       <> <div className=Styles.seperatorLine /> <TxIndexPageTable messages={tx.messages} /> </>
     | None => <VSpacing size={`px(250)} />
     }}
  </div>;
};
