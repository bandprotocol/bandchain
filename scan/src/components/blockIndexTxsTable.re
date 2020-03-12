module Styles = {
  open Css;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let container = style([width(`px(68))]);
  let hashContainer = style([maxWidth(`px(140))]);
  let paddingTopContainer = style([paddingTop(`px(5))]);
  let statusContainer =
    style([maxWidth(`px(95)), display(`flex), flexDirection(`row), alignItems(`center)]);
  let logo = style([width(`px(20)), marginLeft(`auto), marginRight(`px(15))]);
};
[@react.component]
let make = (~txs: list(TxHook.Tx.t)) => {
  // let x = txs.
  let sendMsg =
    TxHook.Msg.{
      action:
        Send({
          fromAddress: "band129umpweqxfywq0f2zdpgjcfnkhzcu8jyewxvyx" |> Address.fromBech32,
          toAddress: "band129umpweqxfywq0f2zdpgjcfnkhzcu8jyewxvyx" |> Address.fromBech32,
          amount: [{denom: "BAND", amount: 1.2}],
        }),
      events: [],
    };
  let createDataSource =
    TxHook.Msg.{
      action:
        CreateDataSource({
          id: 23,
          owner: "band129umpweqxfywq0f2zdpgjcfnkhzcu8jyewxvyx" |> Address.fromBech32,
          name: "CoinGecko V.2",
          fee: [{denom: "BAND", amount: 1.2}],
          executable: "band129umpweqxfywq0f2zdpgjcfnkhzcu8jyewxvyx" |> JsBuffer.fromBase64,
          sender: "band129umpweqxfywq0f2zdpgjcfnkhzcu8jyewxvyx" |> Address.fromBech32,
        }),
      events: [],
    };
  <>
    <THead>
      <Row>
        <HSpacing size={`px(20)} />
        <Col size=1.67>
          <div className=Styles.fullWidth>
            <Text value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.mediumLightGray />
          </div>
        </Col>
        <Col size=1.05>
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
        <Col> <div className=Styles.container /> </Col>
        <Col size=5.>
          <div className=Styles.fullWidth>
            <Text value="ACTIONS" size=Text.Sm weight=Text.Bold color=Colors.mediumLightGray />
          </div>
        </Col>
        <HSpacing size={`px(20)} />
      </Row>
    </THead>
    {txs
     ->Belt.List.map(({blockHeight, hash, timestamp, fee, gasUsed, messages, sender}) => {
         <TBody key={hash |> Hash.toHex}>
           <Row>
             <HSpacing size={`px(20)} />
             <Col size=1.67 alignSelf=Col.FlexStart>
               <div className={Css.merge([Styles.hashContainer, Styles.paddingTopContainer])}>
                 <Text
                   block=true
                   code=true
                   spacing={Text.Em(0.02)}
                   value={hash |> Hash.toHex(~upper=true)}
                   weight=Text.Medium
                   ellipsis=true
                 />
               </div>
             </Col>
             <Col size=1.05 alignSelf=Col.FlexStart>
               <div className={Css.merge([Styles.fullWidth, Styles.paddingTopContainer])}>
                 <AutoSpacing dir="left" />
                 <Text
                   block=true
                   code=true
                   spacing={Text.Em(0.02)}
                   value={(fee.amount /. 10000.0)->Format.fPretty}
                   weight=Text.Medium
                   ellipsis=true
                 />
                 <HSpacing size={`px(20)} />
               </div>
             </Col>
             <Col> <div className=Styles.container /> </Col>
             <Col size=5. alignSelf=Col.FlexStart>
               {messages
                ->Belt.List.map(msg => {<> <Msg msg width=330 /> <VSpacing size=Spacing.md /> </>})
                ->Belt.List.toArray
                ->React.array}
             </Col>
             <HSpacing size={`px(20)} />
           </Row>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
  </>;
};
