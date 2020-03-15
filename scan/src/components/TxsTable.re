module Styles = {
  open Css;
  let hScale = 20;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let hashContainer = style([maxWidth(`px(140))]);
  let statusContainer =
    style([maxWidth(`px(95)), display(`flex), flexDirection(`row), alignItems(`center)]);
  let logo = style([width(`px(20)), marginLeft(`auto), marginRight(`px(15))]);

  let heightByMsgsNum = (numMsgs, mt) =>
    style([
      minHeight(numMsgs <= 1 ? `auto : `px(numMsgs * hScale)),
      marginTop(`px(numMsgs <= 1 ? 0 : mt)),
    ]);
};

[@react.component]
let make = (~txs: list(TxHook.Tx.t)) => {
  <>
    <THead>
      <Row>
        <HSpacing size={`px(20)} />
        <Col size=1.67>
          <div className=Styles.fullWidth>
            <Text value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <Col size=0.88>
          <div className=Styles.fullWidth>
            <Text value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <Col size=1.>
          <div className=Styles.fullWidth>
            <Text value="STATUS" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <Col size=1.25>
          <div className=Styles.fullWidth>
            <AutoSpacing dir="left" />
            <Text value="GAS FEE (BAND)" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
            <HSpacing size={`px(20)} />
          </div>
        </Col>
        <Col size=5.>
          <div className=Styles.fullWidth>
            <Text value="ACTIONS" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <HSpacing size={`px(20)} />
      </Row>
    </THead>
    {txs
     ->Belt.List.map(({blockHeight, hash, fee, messages, success}) => {
         let numMsgs = messages->Belt_List.size;
         <TBody key={hash |> Hash.toHex}>
           <Row minHeight={`px(30)}>
             <HSpacing size={`px(20)} />
             <Col size=1.67>
               <div className={Styles.heightByMsgsNum(numMsgs, 0)}>
                 <TxLink txHash=hash width=140 />
               </div>
             </Col>
             <Col size=0.88>
               <div className={Styles.heightByMsgsNum(numMsgs, -4)}>
                 <TypeID.Block id={ID.Block.ID(blockHeight)} />
               </div>
             </Col>
             <Col size=1.>
               <div className={Styles.heightByMsgsNum(numMsgs, -8)}>
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
               </div>
             </Col>
             <Col size=1.25>
               <div className={Styles.heightByMsgsNum(numMsgs, 0)}>
                 <div className=Styles.fullWidth>
                   <AutoSpacing dir="left" />
                   <Text
                     block=true
                     code=true
                     spacing={Text.Em(0.02)}
                     value={fee->TxHook.Coin.getBandAmountFromCoins->Format.fPretty}
                     weight=Text.Medium
                     ellipsis=true
                   />
                   <HSpacing size={`px(20)} />
                 </div>
               </div>
             </Col>
             <Col size=5.>
               {messages
                ->Belt_List.toArray
                ->Belt_Array.map(msg =>
                    <>
                      <VSpacing size=Spacing.sm />
                      <Msg msg success width=330 />
                      <VSpacing size=Spacing.sm />
                    </>
                  )
                ->React.array}
             </Col>
             <HSpacing size={`px(20)} />
           </Row>
         </TBody>;
       })
     ->Array.of_list
     ->React.array}
  </>;
};
