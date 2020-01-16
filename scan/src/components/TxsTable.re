module Styles = {
  open Css;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
};

[@react.component]
let make = (~txs: list(TxHook.Tx.t)) => {
  <>
    <THead>
      <Row>
        <Col> <div className=TElement.Styles.msgIcon /> </Col>
        <Col size=0.8>
          <div className=TElement.Styles.hashContainer>
            <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=0.5>
          <Text block=true value="TYPE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.2>
          <Text block=true value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.8>
          <Text block=true value="SENDER" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.3>
          <div className=TElement.Styles.feeContainer>
            <Text block=true value="FEE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
      </Row>
    </THead>
    {txs
     ->Belt.List.map(({blockHeight, hash, timestamp, gasUsed, messages, sender}) => {
         <TBody key={hash |> Hash.toHex}>
           <div
             className=Styles.fullWidth onClick={_ => Route.TxIndexPage(hash) |> Route.redirect}>
             <Row>
               <Col>
                 <TElement elementType={messages->Belt.List.getExn(0)->TElement.Icon} />
               </Col>
               <Col size=0.8> <TElement elementType={TElement.TxHash(hash, timestamp)} /> </Col>
               <Col size=0.5>
                 <TElement
                   elementType={messages->Belt.List.getExn(0)->TElement.TxTypeWithDetail}
                 />
               </Col>
               <Col size=0.2> <TElement elementType={TElement.Height(blockHeight)} /> </Col>
               <Col size=0.8> <TElement elementType={sender->TElement.Address} /> </Col>
               <Col size=0.3> <TElement elementType={0.->TElement.Fee} /> </Col>
             </Row>
           </div>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
  </>;
};
