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
      hash: "0x103020321239012392139023101230" |> Hash.fromHex,
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Report({
          requestId: 2,
          data: "0x8812381238823182318",
          validator: "0x913932993219329319213" |> Address.fromHex,
        }),
      ],
    },
    {
      blockHeight: 120338,
      hash: "0x12391291239123921392139" |> Hash.fromHex,
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
    },
    {
      blockHeight: 120337,
      hash: "0x12391291239123921392139" |> Hash.fromHex,
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Store({
          code: "0x1923912912391293219",
          owner: "0x949494949499494949494" |> Address.fromHex,
        }),
      ],
    },
    {
      blockHeight: 120337,
      hash: "0x12391291239123921392139" |> Hash.fromHex,
      timestamp: MomentRe.momentWithUnix(1293912392),
      gasWanted: 0,
      gasUsed: 0,
      messages: [
        Request({
          codeHash: "0x91238123812838123" |> Hash.fromHex,
          params: "0x8238233288238238",
          reportPeriod: 23,
          sender: "0x99329329239239923923" |> Address.fromHex,
        }),
      ],
    },
  ];

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
    {txs
     ->Belt.List.mapWithIndex((idx, {hash, timestamp, gasUsed, messages}) => {
         <TBody key={idx |> string_of_int}>
           <Row>
             <Col> <TElement elementType={messages->Belt.List.getExn(0)->TElement.Icon} /> </Col>
             <Col size=1.3> <TElement elementType={TElement.TxHash(hash, timestamp)} /> </Col>
             <Col size=1.3> <TElement elementType={messages->TElement.TxType} /> </Col>
             <Col size=0.5> <TElement elementType={TElement.Fee(gasUsed, true)} /> </Col>
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
