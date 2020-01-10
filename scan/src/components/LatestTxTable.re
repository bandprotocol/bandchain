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
  /* mock tx */
  let txs: list(TxHook.Tx.t) = [
    {
      blockHeight: 120339,
      hash: "0x12391291239123921392139",
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Report({requestId: 2, data: "0x8812381238823182318", validator: "CoinMarketCap"}),
      ],
    },
    {
      blockHeight: 120338,
      hash: "0x12391291239123921392139",
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Send({
          fromAddress: "0x923239239923923923",
          toAddress: "0x233262363262363263",
          amount: [{denom: "BAND", amount: 12.4}, {denom: "UATOM", amount: 10000.3}],
        }),
      ],
    },
    {
      blockHeight: 120337,
      hash: "0x12391291239123921392139",
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [Store({code: "0x1923912912391293219", owner: "0x949494949499494949494"})],
    },
    {
      blockHeight: 120337,
      hash: "0x12391291239123921392139",
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Request({
          codeHash: "0x91238123812838123",
          params: "0x8238233288238238",
          reportPeriod: 23,
          sender: "0x99329329239239923923",
        }),
      ],
    },
  ];

  <>
    <THead>
      <Row>
        <Col> <div className=TElement.Styles.msgIcon /> </Col>
        <Col size=1.3>
          <div className=TElement.Styles.hashCol>
            <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=1.3>
          <Text block=true value="TYPE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.5>
          <div className=TElement.Styles.feeCol>
            <Text block=true value="FEE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
      </Row>
    </THead>
    {txs
     ->Belt.List.map(({hash, timestamp, gasUsed, messages}) => {
         <TBody>
           <Row>
             <Col> <TElement elementType={messages->Belt.List.getExn(0)->TElement.Icon} /> </Col>
             <Col size=1.3> <TElement elementType={TElement.TxHash(hash, timestamp)} /> </Col>
             <Col size=1.3> <TElement elementType={messages->TElement.TxType} /> </Col>
             <Col size=0.5> <TElement elementType={gasUsed->TElement.Fee} /> </Col>
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
