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
let make = (~messages: list(TxHook.Msg.t)) => {
  <>
    <THead>
      <Row>
        <Col> <div className=TElement.Styles.msgIcon /> </Col>
        <Col size=0.5>
          <div className=TElement.Styles.hashContainer>
            <Text
              block=true
              value="MESSAGE TYPE"
              size=Text.Sm
              weight=Text.Bold
              color=Colors.grayText
            />
          </div>
        </Col>
        <Col size=1.0>
          <Text block=true value="DETAIL" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=1.3>
          <div className=TElement.Styles.feeContainer>
            <Text block=true value="CREATOR" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        // <Col size=0.5>
        //   <div className=TElement.Styles.feeCol>
        //     <Text block=true value="STATUS" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        //   </div>
        // </Col>
        <Col size=0.5>
          <div className=TElement.Styles.feeContainer>
            <Text block=true value="FEE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
      </Row>
    </THead>
    {messages
     ->Belt.List.map(msg => {
         <TBody>
           <Row>
             <Col size=0.3> <TElement elementType={msg->TElement.Icon} /> </Col>
             <Col size=0.5> <TElement elementType={msg->TElement.TxType} /> </Col>
             <Col size=1.0>
               <TElement elementType={msg->TxHook.Msg.getDescription->TElement.Detail} />
             </Col>
             <Col size=1.3>
               <TElement elementType={msg->TxHook.Msg.getCreator->TElement.Address} />
             </Col>
             // <Col size=0.5> <TElement elementType={"PENDING DATA"->TElement.Status} /> </Col>
             <Col size=0.5> <TElement elementType={0.->TElement.Fee} /> </Col>
           </Row>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
  </>;
};
