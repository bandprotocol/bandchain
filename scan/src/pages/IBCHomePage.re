module Styles = {
  open Css;

  let title = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);
  let dirArrow = dir =>
    style([
      width(`px(20)),
      transforms([
        translateX(`px(5)),
        translateY(`px(-5)),
        rotate(
          `deg(
            switch (dir) {
            | IBCSub.Incoming => 0.
            | IBCSub.Outgoing => 180.
            },
          ),
        ),
      ]),
    ]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let badge = color =>
    style([
      display(`inlineFlex),
      backgroundColor(color),
      alignItems(`center),
      maxHeight(`px(16)),
      padding(`px(5)),
      borderRadius(`px(16)),
      transform(translateY(`px(-5))),
    ]);

  let vFlex = style([display(`flex), flexDirection(`column), width(`percent(100.))]);
  let hFlex = style([display(`flex), flexDirection(`row), width(`percent(100.))]);

  let minWidth = x => style([minWidth(`px(x))]);
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
         <TBody key={i |> string_of_int} paddingV={`px(20)}>
           <Row alignItems=`flexStart>
             <Col> <HSpacing size=Spacing.md /> </Col>
             <Col size=9.7>
               <div className={Styles.badge(Colors.orange1)}>
                 <Text
                   value={
                     switch (packet) {
                     | Request(_) => "ORACLE REQUEST"
                     | Response(_) => "ORACLE RESPONSE"
                     | Unknown => "Unknown"
                     }
                   }
                   size=Text.Xs
                   color=Colors.orange6
                   spacing={Text.Em(0.07)}
                 />
               </div>
               <VSpacing size=Spacing.md />
               <div className={Styles.badge(Colors.blue1)}>
                 <Text
                   value={
                     switch (direction) {
                     | Incoming => "INCOMING"
                     | Outgoing => "OUTGOING"
                     }
                   }
                   size=Text.Xs
                   color=Colors.blue7
                   spacing={Text.Em(0.07)}
                 />
               </div>
               <VSpacing size=Spacing.md />
               <img src=Images.ibcDirArrow className={Styles.dirArrow(direction)} />
             </Col>
             <Col size=18.3>
               <div className=Styles.hFlex>
                 <Text
                   value={j|ChainID:‌‌ ‌‌ |j}
                   size=Text.Sm
                   code=true
                   height={Text.Px(16)}
                 />
                 <Text value=chainID size=Text.Sm code=true height={Text.Px(16)} />
               </div>
               <VSpacing size=Spacing.md />
               <div className=Styles.hFlex>
                 <Text
                   value={j|Channel:‌‌ ‌‌ |j}
                   size=Text.Sm
                   code=true
                   height={Text.Px(16)}
                 />
                 <Text value=chennel size=Text.Sm code=true height={Text.Px(16)} />
               </div>
               <VSpacing size=Spacing.md />
               <div className=Styles.hFlex>
                 <Text
                   value={j|Port:‌‌ ‌‌ ‌‌ ‌‌ ‌‌ |j}
                   size=Text.Sm
                   code=true
                   height={Text.Px(16)}
                 />
                 <Text value=port size=Text.Sm code=true height={Text.Px(16)} />
               </div>
             </Col>
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
