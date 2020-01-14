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
  <div className=Styles.pageContainer>
    <Row>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL DATA REQUEST SCRIPTS"
            weight=Text.Bold
            size=Text.Xl
            nowrap=true
            color=Colors.grayHeader
          />
          <div className=Styles.seperatedLine />
          <Text value="30 in total" />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <THead>
      <Row>
        <Col> <HSpacing size=Spacing.xl /> </Col>
        <Col size=1.1>
          <div className=TElement.Styles.hashContainer>
            <Text block=true value="NAME" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=1.1>
          <Text
            block=true
            value="SCRIPT HASH"
            size=Text.Sm
            weight=Text.Bold
            color=Colors.grayText
          />
        </Col>
        <Col size=0.65>
          <Text
            block=true
            value="CREATED AT"
            size=Text.Sm
            weight=Text.Bold
            color=Colors.grayText
          />
        </Col>
        <Col size=1.1>
          <Text block=true value="CREATOR" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.5>
          <div className=TElement.Styles.feeContainer>
            <Text
              block=true
              value="QUERY FEE"
              size=Text.Sm
              weight=Text.Bold
              color=Colors.grayText
            />
          </div>
        </Col>
      </Row>
    </THead>
    {[
       (
         "ETH/USD Price Feed",
         "0x1923182381238123812383" |> Hash.fromHex,
         MomentRe.momentWithUnix(1578661269),
         "0x238283328823823823" |> Address.fromHex,
         0,
       ),
       (
         "ETH/BTC Price Feed",
         "0x1923182381238123812383" |> Hash.fromHex,
         MomentRe.momentWithUnix(1578661269),
         "0x238283328823823823" |> Address.fromHex,
         0,
       ),
       (
         "ETH/USD Price Feed",
         "0x1923182381238123812383" |> Hash.fromHex,
         MomentRe.momentWithUnix(1578661269),
         "0x238283328823823823" |> Address.fromHex,
         0,
       ),
     ]
     ->Belt.List.mapWithIndex((idx, (name, scriptHash, timestamp, creator, fee)) => {
         <TBody key={idx |> string_of_int}>
           <Row>
             <Col> <HSpacing size=Spacing.xl /> </Col>
             <Col size=1.1> <TElement elementType={name->TElement.Name} /> </Col>
             <Col size=1.1> <TElement elementType={scriptHash->TElement.Hash} /> </Col>
             <Col size=0.65> <TElement elementType={timestamp->TElement.Timestamp} /> </Col>
             <Col size=1.1> <TElement elementType={creator->TElement.Address} /> </Col>
             <Col size=0.5> <TElement elementType={fee->TElement.Fee} /> </Col>
           </Row>
         </TBody>
       })
     ->Array.of_list
     ->React.array}
    <VSpacing size=Spacing.lg />
    <LoadMore onClick={_ => ()} />
    <VSpacing size=Spacing.xl />
    <VSpacing size=Spacing.xl />
  </div>;
};
