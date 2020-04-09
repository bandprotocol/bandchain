module Styles = {
  open Css;

  let title = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let vFlex = style([display(`flex), flexDirection(`column), width(`percent(100.))]);
  let hFlex = style([display(`flex), flexDirection(`row), width(`percent(100.))]);
};

[@react.component]
let make = () => {
  let (page, setPage) = React.useState(_ => 1);
  // let pageSize = 10;
  let pageCount = 5;
  let packets = IBCSub.getMockList();
  <>
    <Row>
      <Col> <img src=Images.ibcLogo className=Styles.logo /> </Col>
      <Col>
        <div className=Styles.title>
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
      <div className=Styles.hFlex>
        <Col> <HSpacing size=Spacing.md /> </Col>
        <Col size=9.7>
          <Text
            block=true
            value="PACKET"
            size=Text.Sm
            weight=Text.Semibold
            color=Colors.gray5
            spacing={Text.Em(0.1)}
          />
        </Col>
        <Col size=18.3>
          <Text
            block=true
            value="PEER INFO"
            size=Text.Sm
            weight=Text.Semibold
            color=Colors.gray5
            spacing={Text.Em(0.1)}
          />
        </Col>
        <Col size=8.1>
          <Text
            block=true
            value="BLOCK"
            size=Text.Sm
            weight=Text.Semibold
            color=Colors.gray5
            spacing={Text.Em(0.1)}
          />
        </Col>
        <Col size=44.2>
          <Text
            block=true
            value="DETAIL"
            size=Text.Sm
            weight=Text.Semibold
            color=Colors.gray5
            spacing={Text.Em(0.1)}
          />
        </Col>
        <Col> <HSpacing size=Spacing.md /> </Col>
      </div>
    </THead>
    {packets
     ->Belt_Array.mapWithIndex((i, {direction, chainID, chennel, port, blockHeight, packet}) => {
         <TBody key={i |> string_of_int}>
           <Row>
             <Col> <HSpacing size=Spacing.md /> </Col>
             <Col size=9.7> <Text value="test" /> </Col>
             <Col size=18.3> <Text value="test" /> </Col>
             <Col size=8.1> <TypeID.Block id=blockHeight /> </Col>
             <Col size=44.2> <Packet packet /> </Col>
             <Col> <HSpacing size=Spacing.md /> </Col>
           </Row>
         </TBody>
       })
     ->React.array}
    <VSpacing size=Spacing.lg />
    <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
    <VSpacing size=Spacing.lg />
  </>;
};
