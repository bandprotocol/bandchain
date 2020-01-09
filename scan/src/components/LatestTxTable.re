module Styles = {
  open Css;

  let seeMoreContainer =
    style([
      width(`percent(100.)),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      backgroundColor(white),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      height(`px(30)),
      cursor(`pointer),
    ]);
};

[@react.component]
let make = () => {
  <>
    <THead>
      <Row>
        <Col> <div className=Transaction.Styles.txIcon /> </Col>
        <Col size=1.3>
          <div className=Transaction.Styles.hashCol>
            <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=1.3>
          <Text
            block=true
            value="SOURCE & TYPE"
            size=Text.Sm
            weight=Text.Bold
            color=Colors.grayText
          />
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
         0.13,
         MomentRe.momentWithUnix(1578348371),
       ),
       (
         DataRequest("ETH/BTC Price Feed"),
         "0x123343020302302",
         0.13,
         MomentRe.momentWithUnix(1578348171),
       ),
       (
         NewScript("Anime Episodes Ranking - WINTER 2020"),
         "0xd83ab82c9f838391283",
         0.1,
         MomentRe.momentWithUnix(1578337371),
       ),
       (
         DataRequest("ETH/BTC Price Feed"),
         "0x123343020302302",
         0.13,
         MomentRe.momentWithUnix(1578348171),
       ),
       (
         DataRequest("ETH/BTC Price Feed"),
         "0x123343020302302",
         0.13,
         MomentRe.momentWithUnix(1578348171),
       ),
       (
         NewScript("Anime Episodes Ranking - WINTER 2020"),
         "0xd83ab82c9f838391283",
         0.1,
         MomentRe.momentWithUnix(1578337371),
       ),
     ]
     ->Belt.List.map(((type_, hash, fee, timestamp)) => {
         <TBody>
           <Row>
             <Col>
               <img src={type_->Transaction.txIcon} className=Transaction.Styles.txIcon />
             </Col>
             <Col size=1.3> {Transaction.renderTxHash(hash, timestamp)} </Col>
             <Col size=1.3> {type_ |> Transaction.renderDataType} </Col>
             <Col size=0.5> {fee |> Transaction.renderFee} </Col>
           </Row>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
    <div className=Styles.seeMoreContainer>
      <Text value="SEE MORE" size=Text.Sm weight=Text.Bold block=true color=Colors.grayText />
    </div>
  </>;
};
