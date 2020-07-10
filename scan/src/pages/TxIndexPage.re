module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let header =
    style([display(`flex), flexDirection(`row), alignItems(`center), height(`px(50))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let hashContainer =
    style([
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      marginTop(`px(25)),
      marginBottom(`px(44)),
    ]);

  let correctLogo = style([width(`px(20)), marginLeft(`px(10))]);

  let logo = style([minWidth(`px(50)), marginRight(`px(10))]);

  let notfoundContainer =
    style([
      width(`percent(100.)),
      minHeight(`px(450)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
      paddingLeft(`px(50)),
      paddingRight(`px(50)),
      justifyContent(`center),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(0, 0, 0, 0.1))),
    ]);
  let notfoundLogo = style([width(`px(180)), marginRight(`px(10))]);
};

module TxNotFound = {
  [@react.component]
  let make = () => {
    <>
      <VSpacing size=Spacing.lg />
      <div className=Styles.notfoundContainer>
        <Col> <img src=Images.notFoundBg className=Styles.notfoundLogo /> </Col>
        <VSpacing size=Spacing.md />
        <Text
          value="Sorry, we are unable to retrieve information on this transaction hash."
          size=Text.Lg
          color=Colors.blueGray6
        />
        <VSpacing size=Spacing.lg />
        <Text
          value="Note: Transactions usually take 5-10 seconds to appear."
          size=Text.Lg
          color=Colors.blueGray6
        />
      </div>
    </>;
  };
};

[@react.component]
let make = (~txHash) => {
  let txSub = TxSub.get(txHash);

  switch (txSub) {
  | Loading
  | Data(_) =>
    <>
      <Row justify=Row.Between>
        <div className=Styles.header>
          <img src=Images.txLogo className=Styles.logo />
          <Text
            value="TRANSACTION"
            weight=Text.Medium
            nowrap=true
            color=Colors.gray7
            spacing={Text.Em(0.06)}
            block=true
          />
          <div className=Styles.seperatedLine />
          {switch (txSub) {
           | Data({success}) =>
             <>
               <Text
                 value={success ? "SUCCESS" : "FAILED"}
                 weight=Text.Thin
                 nowrap=true
                 color=Colors.gray7
                 spacing={Text.Em(0.06)}
                 block=true
               />
               <img src={success ? Images.success : Images.fail} className=Styles.correctLogo />
             </>
           | _ =>
             <>
               <LoadingCensorBar width=60 height=15 />
               <HSpacing size=Spacing.sm />
               <LoadingCensorBar width=20 height=20 radius=20 />
             </>
           }}
        </div>
      </Row>
      <div className=Styles.hashContainer>
        {switch (txSub) {
         | Data(_) =>
           <>
             <Text
               value={txHash |> Hash.toHex(~upper=true)}
               size=Text.Xxl
               weight=Text.Bold
               nowrap=true
               code=true
               color=Colors.gray7
             />
             <HSpacing size=Spacing.sm />
             <CopyRender width=15 message={txHash |> Hash.toHex(~upper=true)} />
           </>
         | _ => <LoadingCensorBar width=700 height=20 />
         }}
      </div>
      <Row>
        <Col size=0.9>
          {switch (txSub) {
           | Data({blockHeight}) => <InfoHL info={InfoHL.Height(blockHeight)} header="BLOCK" />
           | _ => <InfoHL info={InfoHL.Loading(75)} header="BLOCK" />
           }}
        </Col>
        <Col size=2.2>
          {switch (txSub) {
           | Data({timestamp}) =>
             <InfoHL info={InfoHL.Timestamp(timestamp)} header="TIMESTAMP" />
           | _ => <InfoHL info={InfoHL.Loading(400)} header="TIMESTAMP" />
           }}
        </Col>
        <Col size=1.4>
          {switch (txSub) {
           | Data({sender}) => <InfoHL info={InfoHL.Address(sender, 290)} header="SENDER" />
           | _ => <InfoHL info={InfoHL.Loading(295)} header="SENDER" />
           }}
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <Row>
        <Col size=1.35>
          {switch (txSub) {
           | Data({gasUsed}) => <InfoHL info={InfoHL.Count(gasUsed)} header="GAS USED" />
           | _ => <InfoHL info={InfoHL.Loading(75)} header="GAS USED" />
           }}
        </Col>
        <Col size=1.>
          {switch (txSub) {
           | Data({gasLimit}) => <InfoHL info={InfoHL.Count(gasLimit)} header="GAS LIMIT" />
           | _ => <InfoHL info={InfoHL.Loading(75)} header="GAS LIMIT" />
           }}
        </Col>
        <Col size=1.>
          {switch (txSub) {
           | Data({gasFee, gasLimit}) =>
             <InfoHL
               info={
                 InfoHL.Float(
                   (gasFee |> Coin.getBandAmountFromCoins) /. (gasLimit |> float_of_int) *. 1e6,
                 )
               }
               header="GAS PRICE (UBAND)"
               isLeft=false
             />
           | _ => <InfoHL info={InfoHL.Loading(75)} header="GAS PRICE (BAND)" isLeft=false />
           }}
        </Col>
        <Col size=1.35>
          {switch (txSub) {
           | Data({gasFee}) =>
             <InfoHL
               info={InfoHL.Float(gasFee |> Coin.getBandAmountFromCoins)}
               header="FEE (BAND)"
               isLeft=false
             />
           | _ => <InfoHL info={InfoHL.Loading(75)} header="FEE (BAND)" isLeft=false />
           }}
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <Row>
        <Col>
          {switch (txSub) {
           | Data({memo}) => <InfoHL info={InfoHL.Description(memo)} header="MEMO" />
           | _ => <InfoHL info={InfoHL.Loading(75)} header="MEMO" />
           }}
        </Col>
      </Row>
      {switch (txSub) {
       | Data({success, rawLog, messages}) =>
         <>
           {success ? React.null : <> <VSpacing size=Spacing.xl /> <TxError.Full msg=rawLog /> </>}
           <VSpacing size=Spacing.xxl />
           <div className=Styles.vFlex>
             <HSpacing size=Spacing.md />
             <Text
               value={messages |> Belt.List.length |> string_of_int}
               weight=Text.Semibold
               size=Text.Lg
             />
             <HSpacing size=Spacing.md />
             <Text value="Messages" size=Text.Lg spacing={Text.Em(0.06)} />
           </div>
           <VSpacing size=Spacing.md />
           <TxIndexPageTable messages />
         </>
       | _ =>
         <>
           <VSpacing size=Spacing.xxl />
           <div className=Styles.vFlex>
             <HSpacing size=Spacing.md />
             <LoadingCensorBar width=100 height=20 />
           </div>
           <VSpacing size=Spacing.md />
           <TxIndexPageTable.Loading />
         </>
       }}
    </>
  | _ => <TxNotFound />
  };
};
