[@react.component]
let make = () => {
  let pageSize = 20;

  let txsSub = TxSub.getList(~pageSize, ~page=1, ());
  let latestTxsSub = TxSub.getList(~pageSize=1, ~page=1, ());

  let isMobile = Media.isMobile();

  <Section>
    <div className=CssHelper.container id="transactionsSection">
      <Row alignItems=Row.Center marginBottom=40 marginBottomSm=24>
        <Col col=Col.Twelve>
          <Heading value="All Transactions" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
          {switch (latestTxsSub) {
           | Data(txs) =>
             <Heading
               value={
                 //  HACK: decrease tx count for guanyu testnet3 only
                 (
                   txs->Belt.Array.get(0)->Belt.Option.mapWithDefault(0, ({id}) => id)
                   - 2207294
                   |> Format.iPretty
                 )
                 ++ " In total"
               }
               size=Heading.H3
             />
           | _ => <LoadingCensorBar width=65 height=21 />
           }}
        </Col>
      </Row>
      {isMobile
         ? React.null
         : <THead>
             <Row alignItems=Row.Center>
               <Col col=Col.Two>
                 <Text block=true value="TX Hash" weight=Text.Semibold color=Colors.gray7 />
               </Col>
               <Col col=Col.One>
                 <Text block=true value="Block" weight=Text.Semibold color=Colors.gray7 />
               </Col>
               <Col col=Col.One>
                 <Text
                   block=true
                   value="Status"
                   size=Text.Md
                   weight=Text.Semibold
                   color=Colors.gray7
                   align=Text.Center
                 />
               </Col>
               <Col col=Col.Two>
                 <Text
                   block=true
                   value="Gas Fee (BAND)"
                   weight=Text.Semibold
                   color=Colors.gray7
                   align=Text.Center
                 />
               </Col>
               <Col col=Col.Six>
                 <Text block=true value="Actions" weight=Text.Semibold color=Colors.gray7 />
               </Col>
             </Row>
           </THead>}
      <TxsTable txsSub />
    </div>
  </Section>;
};
