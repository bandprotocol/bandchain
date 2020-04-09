module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);
};

[@react.component]
let make = () => {
  <>
    <Row>
      <Col> <img src=Images.ibcLogo className=Styles.logo /> </Col>
      <Col>
        <div className=Styles.vFlex>
          <Text
            value="ALL IBC PACKETS"
            weight=Text.Semibold
            color=Colors.gray7
            nowrap=true
            spacing={Text.Em(0.06)}
          />
          <div className=Styles.seperatedLine />
          // TODO: replace this mock
          <Text value={(999 |> Format.iPretty) ++ " in total"} />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <THead>
      <Row>
        <Col> <HSpacing size=Spacing.xl /> </Col>
        <Col size=0.5>
          <div className=TElement.Styles.hashContainer>
            <Text
              block=true
              value="PACKET"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray5
              spacing={Text.Em(0.1)}
            />
          </div>
        </Col>
        <Col size=0.5>
          <Text
            block=true
            value="PEER INFO"
            size=Text.Sm
            weight=Text.Semibold
            color=Colors.gray5
            spacing={Text.Em(0.1)}
          />
        </Col>
        <Col size=1.>
          <Text
            block=true
            value="BLOCK"
            size=Text.Sm
            weight=Text.Semibold
            color=Colors.gray5
            spacing={Text.Em(0.1)}
          />
        </Col>
        <Col size=0.4>
          <Text
            block=true
            value="DETAIL"
            size=Text.Sm
            weight=Text.Semibold
            color=Colors.gray5
            spacing={Text.Em(0.1)}
          />
        </Col>
        <Col> <HSpacing size=Spacing.xl /> </Col>
      </Row>
    </THead>
    // {dataSources
    //  ->Belt_Array.map(({id, name, timestamp, owner, fee}) => {
    //      <TBody key=name>
    //        <div className=Styles.fullWidth>
    //          <Row>
    //            <Col> <HSpacing size=Spacing.xl /> </Col>
    //            <Col size=0.5> <TElement elementType={TElement.DataSource(id, name)} /> </Col>
    //            <Col size=0.5> <TElement elementType={timestamp->TElement.Timestamp} /> </Col>
    //            <Col size=1.> <TElement elementType={owner->TElement.Address} /> </Col>
    //            <Col size=0.4>
    //              <TElement elementType={fee->Coin.getBandAmountFromCoins->TElement.Fee} />
    //            </Col>
    //            <Col> <HSpacing size=Spacing.xl /> </Col>
    //          </Row>
    //        </div>
    //      </TBody>
    //    })
    //  ->React.array}
    <VSpacing size=Spacing.lg />
    // <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
    <VSpacing size=Spacing.lg />
  </>;
};
