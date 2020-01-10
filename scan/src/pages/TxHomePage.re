module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50)), minHeight(`px(500))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
};

[@react.component]
let make = () => {
  let txs: list(TxHook.Tx.t) = [
    {
      blockHeight: 120339,
      hash: "0x10302032123901239101230" |> Hash.fromHex,
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

  <div className=Styles.pageContainer>
    <Row>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL TRANSACTIONS"
            weight=Text.Bold
            size=Text.Xl
            nowrap=true
            color=Colors.grayHeader
          />
          <div className=Styles.seperatedLine />
          <Text value="99,999 in total" />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <THead>
      <Row>
        <Col> <div className=TElement.Styles.msgIcon /> </Col>
        <Col size=1.1>
          <div className=TElement.Styles.hashContainer>
            <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=1.1>
          <Text block=true value="TYPE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.65>
          <Text block=true value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=1.1>
          <Text block=true value="SENDER" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.5>
          <div className=TElement.Styles.feeContainer>
            <Text block=true value="FEE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
      </Row>
    </THead>
    {txs
     ->Belt.List.mapWithIndex((idx, {blockHeight, hash, timestamp, gasUsed, messages}) => {
         <TBody key={idx |> string_of_int}>
           <Row>
             <Col> <TElement elementType={messages->Belt.List.getExn(0)->TElement.Icon} /> </Col>
             <Col size=1.1> <TElement elementType={TElement.TxHash(hash, timestamp)} /> </Col>
             <Col size=1.1> <TElement elementType={messages->TElement.TxType} /> </Col>
             <Col size=0.65> <TElement elementType={TElement.Height(blockHeight)} /> </Col>
             <Col size=1.1> <TElement elementType={hash->TElement.Hash} /> </Col>
             <Col size=0.5> <TElement elementType={TElement.Fee(gasUsed, true)} /> </Col>
           </Row>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
    <VSpacing size=Spacing.lg />
    <LoadMore />
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.xl />
  </div>;
};
