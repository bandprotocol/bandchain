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
let make = () =>
  {
    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 10;

    let packetsSub = IBCSub.getList(~pageSize, ~page, ());
    let packetsCountSub = IBCSub.count();

    let%Sub packets = packetsSub;
    let%Sub packetsCount = packetsCountSub;

    let pageCount = Page.getPageCount(packetsCount, pageSize);

    <Section>
      <div className=CssHelper.container>
        <Row alignItems=Row.Center marginBottom=40 marginBottomSm=24>
          <Col col=Col.Twelve>
            <Heading value="All IBC Packets" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
            <Heading value={(packetsCount |> Format.iPretty) ++ " in total"} size=Heading.H3 />
          </Col>
        </Row>
        <THead>
          <Row alignItems=Row.Center>
            <Col col=Col.Two>
              <Text block=true value="Packet" weight=Text.Semibold color=Colors.gray7 />
            </Col>
            <Col col=Col.Three>
              <Text block=true value="Perr Info" weight=Text.Semibold color=Colors.gray7 />
            </Col>
            <Col col=Col.Two>
              <Text block=true value="Block" weight=Text.Semibold color=Colors.gray7 />
            </Col>
            <Col col=Col.Five>
              <Text block=true value="Detail" weight=Text.Semibold color=Colors.gray7 />
            </Col>
          </Row>
        </THead>
        {packets
         ->Belt_Array.mapWithIndex(
             (
               i,
               {
                 direction,
                 srcChannel,
                 srcPort,
                 chainID,
                 dstChannel,
                 dstPort,
                 blockHeight,
                 packet,
               },
             ) => {
             <TBody key={i |> string_of_int} paddingH={`px(24)}>
               <Row>
                 <Col col=Col.Two>
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
                 <Col col=Col.Three>
                   <div className=Styles.hFlex>
                     <Text
                       value={j|ChainID:‌‌ |j}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                     <Text value=chainID size=Text.Sm code=true height={Text.Px(16)} />
                   </div>
                   <VSpacing size=Spacing.md />
                   <div className=Styles.hFlex>
                     <Text
                       value={j|Src Channel:‌‌ |j}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                     <Text
                       value={srcChannel}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                   </div>
                   <VSpacing size=Spacing.md />
                   <div className=Styles.hFlex>
                     <Text
                       value={j|Src Port:‌‌ |j}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                     <Text
                       value={srcPort}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                   </div>
                   <VSpacing size=Spacing.md />
                   <div className=Styles.hFlex>
                     <Text
                       value={j|Dst Channel:‌‌ |j}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                     <Text
                       value={dstChannel}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                   </div>
                   <VSpacing size=Spacing.md />
                   <div className=Styles.hFlex>
                     <Text
                       value={j|Dst Port:‌‌ |j}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                     <Text
                       value={dstPort}
                       size=Text.Sm
                       code=true
                       height={Text.Px(16)}
                     />
                   </div>
                 </Col>
                 <Col col=Col.Two> <TypeID.Block id=blockHeight /> </Col>
                 <Col col=Col.Five>
                   {switch (packet) {
                    | Request({oracleScriptID}) => <Packet packet oracleScriptID />
                    | Response({oracleScriptID}) => <Packet packet oracleScriptID />
                    | Unknown => React.null
                    }}
                 </Col>
               </Row>
             </TBody>
           })
         ->React.array}
        <VSpacing size=Spacing.lg />
        <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
        <VSpacing size=Spacing.lg />
      </div>
    </Section>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
