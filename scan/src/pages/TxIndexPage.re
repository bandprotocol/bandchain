module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(37))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.mediumGray),
    ]);

  let hashContainer = style([marginTop(`px(25)), marginBottom(`px(44))]);

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

  let correctLogo = style([width(`px(20)), marginLeft(`px(10))]);

  let seperatorLine =
    style([
      width(`percent(100.)),
      height(`pxFloat(1.4)),
      backgroundColor(Colors.lightGray),
      display(`flex),
    ]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);
};

[@react.component]
let make = (~txHash) => {
  let txOpt = TxHook.atHash(txHash);
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.txLogo className=Styles.logo />
          <Text
            value="TRANSACTION"
            weight=Text.Medium
            nowrap=true
            color=Colors.mediumGray
            spacing={Text.Em(0.06)}
            block=true
          />
          <div className=Styles.seperatedLine />
          <Text
            value="SUCCESS"
            weight=Text.Thin
            nowrap=true
            color=Colors.mediumGray
            spacing={Text.Em(0.06)}
            block=true
          />
          <img src=Images.success className=Styles.correctLogo />
        </div>
      </Col>
    </Row>
    <div className=Styles.hashContainer>
      <Text
        value={txHash |> Hash.toHex(~upper=true)}
        size=Text.Xxl
        weight=Text.Bold
        nowrap=true
        code=true
        color=Colors.mediumGray
      />
    </div>
    <Row>
      {switch (txOpt) {
       | Some(tx) =>
         <>
           <Col size=0.9> <InfoHL info={InfoHL.Height(tx.blockHeight)} header="BLOCK" /> </Col>
           <Col size=2.2>
             <InfoHL info={InfoHL.Timestamp(tx.timestamp)} header="TIMESTAMP" />
           </Col>
           <Col size=1.4>
             <InfoHL
               info={
                 InfoHL.Address(
                   "band17rprjgtj0krfw3wyl9creueej6ca9dc4dgxv6e" |> Address.fromBech32,
                 )
               }
               header="SENDER"
             />
           </Col>
         </>
       | None =>
         <>
           <Col size=0.9> <InfoHL info={InfoHL.Text("?")} header="BLOCK" /> </Col>
           <Col size=2.2> <InfoHL info={InfoHL.Text("?")} header="TIMESTAMP" /> </Col>
           <Col size=1.4> <InfoHL info={InfoHL.Text("?")} header="SENDER" /> </Col>
         </>
       }}
    </Row>
    <VSpacing size=Spacing.xl />
    <Row>
      {switch (txOpt) {
       | Some(tx) =>
         <>
           <Col size=1.35> <InfoHL info={InfoHL.Count(130082)} header="GAS USED" /> </Col>
           <Col size=1.> <InfoHL info={InfoHL.Count(200000)} header="GAS LIMIT" /> </Col>
           <Col size=1.>
             <InfoHL info={InfoHL.Float(0.000010)} header="GAS PRICE (BAND)" isLeft=false />
           </Col>
           <Col size=1.35>
             <InfoHL info={InfoHL.Float(0.13)} header="FEE (BAND)" isLeft=false />
           </Col>
         </>
       | None =>
         <>
           <Col size=0.9> <InfoHL info={InfoHL.Text("?")} header="BLOCK" /> </Col>
           <Col size=2.2> <InfoHL info={InfoHL.Text("?")} header="TIMESTAMP" /> </Col>
           <Col size=1.4> <InfoHL info={InfoHL.Text("?")} header="SENDER" /> </Col>
         </>
       }}
    </Row>
    <VSpacing size=Spacing.xxl />
    <div className=Styles.vFlex>
      <HSpacing size=Spacing.md />
      {switch (txOpt) {
       | Some(tx) =>
         <Text
           value={tx.messages |> Belt.List.length |> string_of_int}
           weight=Text.Semibold
           size=Text.Lg
         />
       | None => <Text value="?" weight=Text.Semibold size=Text.Lg />
       }}
      <HSpacing size=Spacing.md />
      <Text value="Messages" size=Text.Lg spacing={Text.Em(0.06)} />
    </div>
    <VSpacing size=Spacing.md />
    {switch (txOpt) {
     | Some(tx) =>
       <> <div className=Styles.seperatorLine /> <TxIndexPageTable messages={tx.messages} /> </>
     | None => <VSpacing size={`px(250)} />
     }}
  </div>;
};
