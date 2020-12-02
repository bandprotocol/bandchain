module Styles = {
  open Css;

  let tableWrapper = style([Media.mobile([padding2(~v=`px(16), ~h=`zero)])]);
  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

//TODO: Will Remove After Doing on Validator index

module LoadingWithHeader = {
  [@react.component]
  let make = () => {
    <div className=Styles.tableWrapper>
      <THead>
        <Row alignItems=Row.Center>
          <Col col=Col.Two>
            <Text block=true value="Block" weight=Text.Semibold color=Colors.gray7 />
          </Col>
          <Col col=Col.Seven>
            <Text block=true value="Block Hash" weight=Text.Semibold color=Colors.gray7 />
          </Col>
          <Col col=Col.One>
            <Text block=true value="Txn" weight=Text.Semibold color=Colors.gray7 />
          </Col>
          <Col col=Col.Two>
            <Text
              block=true
              value="Timestamp"
              weight=Text.Semibold
              color=Colors.gray7
              align=Text.Right
            />
          </Col>
        </Row>
      </THead>
    </div>;
  };
};

let renderBody = (reserveIndex, blockSub: ApolloHooks.Subscription.variant(BlockSub.t)) => {
  <TBody
    key={
      switch (blockSub) {
      | Data({height}) => height |> ID.Block.toString
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row alignItems=Row.Center>
      <Col col=Col.Two>
        {switch (blockSub) {
         | Data({height}) => <TypeID.Block id=height />
         | _ => <LoadingCensorBar width=135 height=15 />
         }}
      </Col>
      <Col col=Col.Seven>
        {switch (blockSub) {
         | Data({hash}) =>
           <Text value={hash |> Hash.toHex(~upper=true)} block=true code=true ellipsis=true />

         | _ => <LoadingCensorBar width=522 height=15 />
         }}
      </Col>
      <Col col=Col.One>
        <div className={CssHelper.flexBox(~justify=`center, ())}>
          {switch (blockSub) {
           | Data({txn}) => <Text value={txn |> Format.iPretty} align=Text.Center />
           | _ => <LoadingCensorBar width=20 height=15 />
           }}
        </div>
      </Col>
      <Col col=Col.Two>
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
      </Col>
    </Row>
  </TBody>;
};

let renderBodyMobile = (reserveIndex, blockSub: ApolloHooks.Subscription.variant(BlockSub.t)) => {
  switch (blockSub) {
  | Data({height, timestamp, txn, hash}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Block", Height(height)),
        ("Block Hash", TxHash(hash, Media.isSmallMobile() ? 170 : 200)),
        ("Txn", Count(txn)),
        ("Timestamp", Timestamp(timestamp)),
      ]
      key={height |> ID.Block.toString}
      idx={height |> ID.Block.toString}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Block", Loading(Media.isSmallMobile() ? 170 : 200)),
        ("Block Hash", Loading(166)),
        ("Txn", Loading(20)),
        ("Timestamp", Loading(166)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = (~consensusAddress) => {
  let pageSize = 10;

  let blocksSub =
    BlockSub.getListByConsensusAddress(~address=consensusAddress, ~pageSize, ~page=1, ());

  let isMobile = Media.isMobile();
  <div className=Styles.tableWrapper>
    {isMobile
       ? React.null
       : <THead>
           <Row alignItems=Row.Center>
             <Col col=Col.Two>
               <Text block=true value="Block" weight=Text.Semibold color=Colors.gray7 />
             </Col>
             <Col col=Col.Seven>
               <Text block=true value="Block Hash" weight=Text.Semibold color=Colors.gray7 />
             </Col>
             <Col col=Col.One>
               <Text
                 block=true
                 value="Txn"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Center
               />
             </Col>
             <Col col=Col.Two>
               <Text
                 block=true
                 value="Timestamp"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Right
               />
             </Col>
           </Row>
         </THead>}
    {switch (blocksSub) {
     | Data(blocks) =>
       <>
         {blocks
          ->Belt_Array.mapWithIndex((i, e) =>
              isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
            )
          ->React.array}
       </>
     | _ =>
       Belt_Array.make(pageSize, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
         )
       ->React.array
     }}
  </div>;
};
