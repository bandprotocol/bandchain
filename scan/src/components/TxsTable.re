module Styles = {
  open Css;
  let hScale = 20;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let hashContainer = style([maxWidth(`px(140))]);
  let statusContainer =
    style([maxWidth(`px(95)), display(`flex), flexDirection(`row), alignItems(`center)]);
  let logo = style([width(`px(20)), marginLeft(`auto), marginRight(`px(15))]);
};

[@react.component]
let make = (~txs: array(TxSub.t)) => {
  <>
    <THead>
      <Row>
        <HSpacing size={`px(20)} />
        <Col size=1.6>
          <div className=Styles.fullWidth>
            <Text
              value="TX HASH"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col size=0.88>
          <div className=Styles.fullWidth>
            <Text
              value="BLOCK"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col size=1.>
          <div className=Styles.fullWidth>
            <Text
              value="STATUS"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col size=1.25>
          <div className=Styles.fullWidth>
            <AutoSpacing dir="left" />
            <Text
              value="GAS FEE (BAND)"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
            <HSpacing size={`px(20)} />
          </div>
        </Col>
        <Col size=5.>
          <div className=Styles.fullWidth>
            <Text
              value="ACTIONS"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <HSpacing size={`px(20)} />
      </Row>
    </THead>
    {txs
     ->Belt_Array.map(({blockHeight, txHash, gasFee, messages, success}) => {
         <TBody key={txHash |> Hash.toHex}>
           <Row minHeight={`px(30)} alignItems=`flexStart>
             <HSpacing size={`px(20)} />
             <Col size=1.67> <VSpacing size=Spacing.sm /> <TxLink txHash width=140 /> </Col>
             <Col size=0.88> <VSpacing size=Spacing.sm /> <TypeID.Block id=blockHeight /> </Col>
             <Col size=1.>
               <VSpacing size={`px(4)} />
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
               <VSpacing size={`px(7)} />
               <div className=Styles.fullWidth>
                 <AutoSpacing dir="left" />
                 <Text
                   block=true
                   code=true
                   spacing={Text.Em(0.02)}
                   value={gasFee->Coin.getBandAmountFromCoins->Format.fPretty}
                   weight=Text.Medium
                   ellipsis=true
                 />
                 <HSpacing size={`px(20)} />
               </div>
             </Col>
             <Col size=5.>
               {messages
                ->Belt_List.toArray
                ->Belt_Array.mapWithIndex((i, msg) =>
                    <React.Fragment key={(txHash |> Hash.toHex) ++ (i |> string_of_int)}>
                      <VSpacing size=Spacing.sm />
                      <Msg msg width=450 />
                      <VSpacing size=Spacing.sm />
                    </React.Fragment>
                  )
                ->React.array}
             </Col>
             <HSpacing size={`px(20)} />
           </Row>
         </TBody>
       })
     ->React.array}
  </>;
};
