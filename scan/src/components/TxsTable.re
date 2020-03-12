module Styles = {
  open Css;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let hashContainer = style([maxWidth(`px(140))]);
  let statusContainer =
    style([maxWidth(`px(95)), display(`flex), flexDirection(`row), alignItems(`center)]);
  let logo = style([width(`px(20)), marginLeft(`auto), marginRight(`px(15))]);
};

[@react.component]
let make = (~txs: list(TxHook.Tx.t)) => {
  <>
    <THead>
      <Row>
        <HSpacing size={`px(20)} />
        <Col size=1.67>
          <div className=Styles.fullWidth>
            <Text value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.mediumLightGray />
          </div>
        </Col>
        <Col size=0.88>
          <div className=Styles.fullWidth>
            <Text value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.mediumLightGray />
          </div>
        </Col>
        <Col size=1.>
          <div className=Styles.fullWidth>
            <Text value="STATUS" size=Text.Sm weight=Text.Bold color=Colors.mediumLightGray />
          </div>
        </Col>
        <Col size=1.25>
          <div className=Styles.fullWidth>
            <AutoSpacing dir="left" />
            <Text
              value="GAS FEE (BAND)"
              size=Text.Sm
              weight=Text.Bold
              color=Colors.mediumLightGray
            />
            <HSpacing size={`px(20)} />
          </div>
        </Col>
        <Col size=5.>
          <div className=Styles.fullWidth>
            <Text value="ACTIONS" size=Text.Sm weight=Text.Bold color=Colors.mediumLightGray />
          </div>
        </Col>
        <HSpacing size={`px(20)} />
      </Row>
    </THead>
    {txs
     ->Belt.List.map(({blockHeight, hash, fee, messages, success}) => {
         <TBody key={hash |> Hash.toHex}>
           <Row>
             <HSpacing size={`px(20)} />
             <Col size=1.67> <TxLink txHash=hash width=140 /> </Col>
             <Col size=0.88> <TypeID.Block id={ID.Block.ID(blockHeight)} /> </Col>
             <Col size=1.>
               <div className=Styles.statusContainer>
                 <Text
                   block=true
                   code=true
                   spacing={Text.Em(0.02)}
                   value={success ? "success" : "fail"}
                   weight=Text.Medium
                   ellipsis=true
                 />
                 <img src={success ? Images.success : Images.fail} className=Styles.logo />
               </div>
             </Col>
             <Col size=1.25>
               <div className=Styles.fullWidth>
                 <AutoSpacing dir="left" />
                 <Text
                   block=true
                   code=true
                   spacing={Text.Em(0.02)}
                   value={fee->Format.fPretty}
                   weight=Text.Medium
                   ellipsis=true
                 />
                 <HSpacing size={`px(20)} />
               </div>
             </Col>
             <Col size=5.> <Msg msg={messages->Belt_List.getExn(0)} success width=330 /> </Col>
             <HSpacing size={`px(20)} />
           </Row>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
  </>;
};
