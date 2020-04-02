module Styles = {
  open Css;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let container = style([width(`px(68))]);
  let hashContainer = style([maxWidth(`px(140)), cursor(`pointer)]);
  let paddingTopContainer = style([paddingTop(`px(5))]);
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
        <Col size=1.67>
          <div className=Styles.fullWidth>
            <Text value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <Col size=1.05>
          <div className=Styles.fullWidth>
            <AutoSpacing dir="left" />
            <Text value="GAS FEE (BAND)" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
            <HSpacing size={`px(20)} />
          </div>
        </Col>
        <Col> <div className=Styles.container /> </Col>
        <Col size=5.>
          <div className=Styles.fullWidth>
            <Text value="ACTIONS" size=Text.Sm weight=Text.Bold color=Colors.gray6 />
          </div>
        </Col>
        <HSpacing size={`px(20)} />
      </Row>
    </THead>
    {txs
     ->Belt_Array.map(({txHash, gasFee, messages, success}) => {
         <TBody key={txHash |> Hash.toHex}>
           <Row>
             <HSpacing size={`px(20)} />
             <Col size=1.67 alignSelf=Col.Start>
               <div
                 className={Css.merge([Styles.hashContainer, Styles.paddingTopContainer])}
                 onClick={_ => Route.redirect(Route.TxIndexPage(txHash))}>
                 <Text
                   block=true
                   code=true
                   spacing={Text.Em(0.02)}
                   value={txHash |> Hash.toHex(~upper=true)}
                   weight=Text.Medium
                   ellipsis=true
                 />
               </div>
             </Col>
             <Col size=1.05 alignSelf=Col.Start>
               <div className={Css.merge([Styles.fullWidth, Styles.paddingTopContainer])}>
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
             <Col> <div className=Styles.container /> </Col>
             <Col size=5. alignSelf=Col.Start>
               {messages
                ->Belt.List.map(msg => {
                    <> <Msg msg width=530 success /> <VSpacing size=Spacing.md /> </>
                  })
                ->Belt.List.toArray
                ->React.array}
             </Col>
             <HSpacing size={`px(20)} />
           </Row>
         </TBody>
       })
     ->React.array}
  </>;
};
