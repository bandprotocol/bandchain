module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50)), minHeight(`px(500))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
};

[@react.component]
let make = () => {
  <div className=Styles.pageContainer>
    <Row>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL TRANSACTIONS"
            weight=Text.Bold
            size=Text.Xl
            nowrap=true
            color=Colors.grayHeader
          />
          <div className=Styles.seperatedLine />
          <Text value="99,999 in total" />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <THead>
      <Row>
        <Col> <div className=Transaction.Styles.txIcon /> </Col>
        <Col size=1.0>
          <div className=Transaction.Styles.hashCol>
            <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=1.1>
          <Text block=true value="TYPE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.65>
          <Text block=true value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=1.1>
          <Text block=true value="SENDER" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.5>
          <div className=Transaction.Styles.feeCol>
            <Text block=true value="FEE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
      </Row>
    </THead>
    {[
       (
         Transaction.DataRequest("ETH/USD Price Feed"),
         "0x123343020302302",
         503000,
         "0x1934838538483848348384",
         MomentRe.momentWithUnix(1578548371),
         0.13,
       ),
       (
         DataRequest("ETH/BTC Price Feed"),
         "0x123343020302302",
         503000,
         "0x1934838538483848348384",
         MomentRe.momentWithUnix(1578448371),
         0.13,
       ),
       (
         NewScript("Anime Episodes Ranking - WINTER 2020"),
         "0xd83ab82c9f838391283",
         503000,
         "0x1934838538483848348384",
         MomentRe.momentWithUnix(1578348371),
         0.1,
       ),
       (
         DataRequest("ETH/BTC Price Feed"),
         "0x123343020302302",
         503000,
         "0x1934838538483848348384",
         MomentRe.momentWithUnix(1578348371),
         0.13,
       ),
       (
         DataRequest("ETH/BTC Price Feed"),
         "0x123343020302302",
         503000,
         "0x1934838538483848348384",
         MomentRe.momentWithUnix(1578348371),
         0.13,
       ),
       (
         NewScript("Anime Episodes Ranking - WINTER 2020"),
         "0xd83ab82c9f838391283",
         503000,
         "0x1934838538483848348384",
         MomentRe.momentWithUnix(1568548371),
         0.1,
       ),
     ]
     ->Belt.List.map(((type_, hash, height, sender, timestamp, fee)) => {
         <TBody>
           <Row>
             <Col>
               <img src={type_->Transaction.txIcon} className=Transaction.Styles.txIcon />
             </Col>
             <Col size=1.0> {Transaction.renderTxHash(hash, timestamp)} </Col>
             <Col size=1.1> {type_ |> Transaction.renderDataType} </Col>
             <Col size=0.65> {height |> Transaction.renderHeight} </Col>
             <Col size=1.1> {sender |> Transaction.renderHash} </Col>
             <Col size=0.5> {fee |> Transaction.renderFee} </Col>
           </Row>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
    <VSpacing size=Spacing.lg />
    <LoadMore />
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.xl />
  </div>;
};
