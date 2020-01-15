module Styles = {
  open Css;

  let txhash = style([marginLeft(`px(20))]);
};

[@react.component]
let make = () => {
  let txs: list(TxHook.Tx.t) = [
    {
      sender: "0x498968C2B945Ac37b78414f66167b0786E522636" |> Address.fromHex,
      blockHeight: 120339,
      hash: "0x103020321239012391012300" |> Hash.fromHex,
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Report({
          requestId: 2,
          data: "0x88123812388231823180" |> JsBuffer.fromHex,
          validator: "0x9139329932193293192130" |> Address.fromHex,
        }),
      ],
      events: [],
    },
    {
      sender: "0x498968C2B945Ac37b78414f66167b0786E522636" |> Address.fromHex,
      blockHeight: 120338,
      hash: "0x123912912391239213921390" |> Hash.fromHex,
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Send({
          fromAddress: "0x923239239923923923" |> Address.fromHex,
          toAddress: "0x233262363262363263" |> Address.fromHex,
          amount: [{denom: "BAND", amount: 12.4}, {denom: "UATOM", amount: 10000.3}],
        }),
      ],
      events: [],
    },
    {
      sender: "0x498968C2B945Ac37b78414f66167b0786E522636" |> Address.fromHex,
      blockHeight: 120337,
      hash: "0x123912912391239213921390" |> Hash.fromHex,
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Store({
          code: "0x19239129123912932190" |> JsBuffer.fromHex,
          owner: "0x949494949499494949494" |> Address.fromHex,
        }),
      ],
      events: [],
    },
  ];
  <>
    <THead>
      <Row>
        <Col> <div className=Styles.txhash /> </Col>
        <Col size=1.0>
          <div className=TElement.Styles.hashContainer>
            <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=0.35>
          <Text block=true value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.4>
          <Text block=true value="AGE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=1.0>
          <Text block=true value="FROM" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.6> <div /> </Col>
        <Col size=0.9>
          <div className=TElement.Styles.feeContainer>
            <Text block=true value="VALUE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
      </Row>
    </THead>
    {txs
     ->Belt.List.mapWithIndex((idx, {blockHeight, hash, timestamp, gasUsed, messages}) => {
         <TBody key={idx |> string_of_int} height=100>
           <Row alignItems=Css.flexStart>
             <Col> <div className=Styles.txhash /> </Col>
             <Col size=1.0> <TElement elementType={hash->TElement.HashWithLink} /> </Col>
             <Col size=0.35> <TElement elementType={TElement.Height(blockHeight)} /> </Col>
             <Col size=0.4> <TElement elementType={timestamp->TElement.Timestamp} /> </Col>
             <Col size=1.0>
               <TElement elementType={hash->TElement.Hash} />
               <VSpacing size=Spacing.sm />
               <TElement elementType={"(CoinGecko DataProvider)"->TElement.Source} />
             </Col>
             <Col size=0.6>
               {["CoinMarketCap", "CryptoCompare", "Binance"]
                ->Belt.List.map(source =>
                    <>
                      <TElement elementType={source->TElement.Source} />
                      <VSpacing size=Spacing.sm />
                    </>
                  )
                ->Array.of_list
                ->React.array}
             </Col>
             <Col size=0.9>
               {["0x0000008332", "0x0000008332", "0x0000008332"]
                ->Belt.List.map(source =>
                    <>
                      <TElement elementType={source->TElement.Value} />
                      <VSpacing size=Spacing.sm />
                    </>
                  )
                ->Array.of_list
                ->React.array}
             </Col>
           </Row>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
  </>;
};
