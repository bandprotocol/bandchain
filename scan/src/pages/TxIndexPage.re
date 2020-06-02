module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

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

  let seperatorLine =
    style([
      width(`percent(100.)),
      height(`pxFloat(1.4)),
      backgroundColor(Colors.gray4),
      display(`flex),
    ]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

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
          value="Sorry, We are unable to locate this TxnHash"
          size=Text.Lg
          color=Colors.blueGray6
        />
        <VSpacing size=Spacing.lg />
        <Text
          value="Note that: If you have just submitted a transaction please wait for at least 5-10 seconds."
          size=Text.Lg
          color=Colors.blueGray6
        />
      </div>
    </>;
  };
};

[@react.component]
let make = (~txHash) =>
  {
    let txSub = TxSub.get(txHash);
    let%Sub tx = txSub;
    <>
      <Row justify=Row.Between>
        <Col>
          <div className=Styles.vFlex>
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
            <Text
              value={tx.success ? "SUCCESS" : "FAILED"}
              weight=Text.Thin
              nowrap=true
              color=Colors.gray7
              spacing={Text.Em(0.06)}
              block=true
            />
            <img src={tx.success ? Images.success : Images.fail} className=Styles.correctLogo />
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
          color=Colors.gray7
        />
        <HSpacing size=Spacing.sm />
        <CopyRender width=15 message={txHash |> Hash.toHex(~upper=true)} />
      </div>
      <Row>
        <Col size=0.9> <InfoHL info={InfoHL.Height(tx.blockHeight)} header="BLOCK" /> </Col>
        <Col size=2.2> <InfoHL info={InfoHL.Timestamp(tx.timestamp)} header="TIMESTAMP" /> </Col>
        <Col size=1.4> <InfoHL info={InfoHL.Address(tx.sender, 290)} header="SENDER" /> </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <Row>
        <Col size=1.35> <InfoHL info={InfoHL.Count(tx.gasUsed)} header="GAS USED" /> </Col>
        <Col size=1.> <InfoHL info={InfoHL.Count(tx.gasLimit)} header="GAS LIMIT" /> </Col>
        <Col size=1.>
          <InfoHL
            info={
              InfoHL.Float(
                (tx.gasFee |> Coin.getBandAmountFromCoins) /. (tx.gasLimit |> float_of_int) *. 1e6,
              )
            }
            header="GAS PRICE (UBAND)"
            isLeft=false
          />
        </Col>
        <Col size=1.35>
          <InfoHL
            info={InfoHL.Float(tx.gasFee |> Coin.getBandAmountFromCoins)}
            header="FEE (BAND)"
            isLeft=false
          />
        </Col>
      </Row>
      {tx.success
         ? React.null : <> <VSpacing size=Spacing.xl /> <TxError.Full msg={tx.rawLog} /> </>}
      <VSpacing size=Spacing.xxl />
      <div className=Styles.vFlex>
        <HSpacing size=Spacing.md />
        <Text
          value={tx.messages |> Belt.List.length |> string_of_int}
          weight=Text.Semibold
          size=Text.Lg
        />
        <HSpacing size=Spacing.md />
        <Text value="Messages" size=Text.Lg spacing={Text.Em(0.06)} />
      </div>
      <VSpacing size=Spacing.md />
      <div className=Styles.seperatorLine />
      <TxIndexPageTable messages={tx.messages} />
    </>
    |> Sub.resolve;
  }
  |> Sub.default(_, <TxNotFound />);
