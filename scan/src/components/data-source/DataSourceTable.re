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
  let requests = [RequestHook.Request.{}];

  <>
    <THead>
      <Row>
        <Col> <div className=TElement.Styles.msgIcon /> </Col>
        <Col size=1.3>
          <div className=TElement.Styles.hashContainer>
            <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=1.3>
          <Text block=true value="TYPE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.5>
          <div className=TElement.Styles.feeContainer>
            <Text block=true value="FEE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
      </Row>
    </THead>
    {requests
     ->Belt.List.map(({hash, timestamp, messages}) => {
         <div onClick={_ => Route.redirect(TxIndexPage(hash))}>
           <TBody key={hash |> Hash.toHex}>
             <Row>
               <Col>
                 <TElement elementType={messages->Belt.List.getExn(0)->TElement.Icon} />
               </Col>
               <Col size=1.3> <TElement elementType={TElement.TxHash(hash, timestamp)} /> </Col>
               <Col size=1.3> <TElement elementType={messages->TElement.TxTypeWithDetail} /> </Col>
               <Col size=0.5> <TElement elementType={0.->TElement.Fee} /> </Col>
             </Row>
           </TBody>
         </div>
       })
     ->Array.of_list
     ->React.array}
    <div className=Styles.seeMoreContainer onClick={_ => Route.redirect(TxHomePage)}>
      <Text value="SEE MORE" size=Text.Sm weight=Text.Bold block=true color=Colors.grayText />
    </div>
  </>;
};
