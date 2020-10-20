let renderBody = (reserveIndex, blockSub: ApolloHooks.Subscription.variant(BlockSub.t)) => {
  <TBody.Grid
    key={
      switch (blockSub) {
      | Data({height}) => height |> ID.Block.toString
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row.Grid alignItems=Row.Center>
      <Col.Grid col=Col.Two>
        {switch (blockSub) {
         | Data({height}) => <TypeID.Block id=height />
         | _ => <LoadingCensorBar width=65 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Four>
        {switch (blockSub) {
         | Data({hash}) =>
           <Text
             value={hash |> Hash.toHex(~upper=true)}
             weight=Text.Medium
             block=true
             code=true
             ellipsis=true
           />
         | _ => <LoadingCensorBar fullWidth=true height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Three>
        {switch (blockSub) {
         | Data({validator}) =>
           <ValidatorMonikerLink
             validatorAddress={validator.operatorAddress}
             moniker={validator.moniker}
             identity={validator.identity}
           />
         | _ => <LoadingCensorBar fullWidth=true height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.One>
        <div className={CssHelper.flexBox(~justify=`center, ())}>
          {switch (blockSub) {
           | Data({txn}) => <Text value={txn |> Format.iPretty} code=true weight=Text.Medium />
           | _ => <LoadingCensorBar width=40 height=15 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
          {switch (blockSub) {
           | Data({timestamp}) =>
             <Timestamp.Grid
               time=timestamp
               size=Text.Md
               weight=Text.Regular
               textAlign=Text.Right
             />
           | _ =>
             <>
               <LoadingCensorBar width=70 height=15 />
               <LoadingCensorBar width=80 height=15 mt=5 />
             </>
           }}
        </div>
      </Col.Grid>
    </Row.Grid>
  </TBody.Grid>;
};

let renderBodyMobile = (reserveIndex, blockSub: ApolloHooks.Subscription.variant(BlockSub.t)) => {
  switch (blockSub) {
  | Data({height, timestamp, validator, txn}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Block", Height(height)),
        ("Timestamp", Timestamp(timestamp)),
        (
          "Proposer",
          Validator(validator.operatorAddress, validator.moniker, validator.identity),
        ),
        ("Txn", Count(txn)),
      ]
      key={height |> ID.Block.toString}
      idx={height |> ID.Block.toString}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Block", Loading(70)),
        ("Timestamp", Loading(166)),
        ("Proposer", Loading(136)),
        ("Txn", Loading(20)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = () => {
  let blocksSub = BlockSub.getList(~pageSize=10, ~page=1, ());
  let isMobile = Media.isMobile();

  <Section>
    <div className=CssHelper.container id="blocksSection">
      <div className=CssHelper.mobileSpacing>
        <Row.Grid alignItems=Row.Center marginBottom=40 marginBottomSm=24>
          <Col.Grid col=Col.Twelve>
            <Heading value="All Blocks" size=Heading.H2 marginBottom=40 marginBottomSm=24 />
            {switch (blocksSub) {
             | Data(blocks) =>
               <Heading
                 value={
                   blocks
                   ->Belt.Array.get(0)
                   ->Belt.Option.mapWithDefault(0, ({height}) => height |> ID.Block.toInt)
                   ->Format.iPretty
                   ++ " In total"
                 }
                 size=Heading.H3
               />
             | _ => <LoadingCensorBar width=65 height=21 />
             }}
          </Col.Grid>
        </Row.Grid>
        {isMobile
           ? React.null
           : <THead.Grid>
               <Row.Grid alignItems=Row.Center>
                 <Col.Grid col=Col.Two>
                   <Text
                     block=true
                     value="Block"
                     size=Text.Md
                     weight=Text.Semibold
                     color=Colors.gray7
                   />
                 </Col.Grid>
                 <Col.Grid col=Col.Four>
                   <Text
                     block=true
                     value="Block Hash"
                     size=Text.Md
                     weight=Text.Semibold
                     color=Colors.gray7
                   />
                 </Col.Grid>
                 <Col.Grid col=Col.Three>
                   <Text
                     block=true
                     value="Proposer"
                     size=Text.Md
                     weight=Text.Semibold
                     color=Colors.gray7
                   />
                 </Col.Grid>
                 <Col.Grid col=Col.One>
                   <Text
                     block=true
                     value="Txn"
                     size=Text.Md
                     weight=Text.Semibold
                     color=Colors.gray7
                     align=Text.Center
                   />
                 </Col.Grid>
                 <Col.Grid col=Col.Two>
                   <Text
                     block=true
                     value="Timestamp"
                     size=Text.Md
                     weight=Text.Semibold
                     color=Colors.gray7
                     align=Text.Right
                   />
                 </Col.Grid>
               </Row.Grid>
             </THead.Grid>}
        {switch (blocksSub) {
         | Data(blocks) =>
           blocks
           ->Belt_Array.mapWithIndex((i, e) =>
               isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
             )
           ->React.array
         | _ =>
           Belt_Array.make(10, ApolloHooks.Subscription.NoData)
           ->Belt_Array.mapWithIndex((i, noData) =>
               isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
             )
           ->React.array
         }}
      </div>
    </div>
  </Section>;
};
