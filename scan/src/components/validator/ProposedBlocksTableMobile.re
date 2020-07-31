module Styles = {
  open Css;

  let blockWrapper = style([padding2(~v=`px(20), ~h=`zero)]);
};

let renderBodyMobile = (reserveIndex, blockSub: ApolloHooks.Subscription.variant(BlockSub.t)) => {
  switch (blockSub) {
  | Data({height, timestamp, txn}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("BLOCK", Height(height)),
        ("TIMESTAMP", Timestamp(timestamp)),
        ("TXN", Count(txn)),
      ]
      key={height |> ID.Block.toString}
      idx={height |> ID.Block.toString}
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("BLOCK", Loading(70)),
        ("TIMESTAMP", Loading(166)),
        ("TXN", Loading(20)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

module Loading = {
  [@react.component]
  let make = () => {
    Belt_Array.make(10, ApolloHooks.Subscription.NoData)
    ->Belt_Array.mapWithIndex((i, noData) => renderBodyMobile(i, noData))
    ->React.array;
  };
};

[@react.component]
let make = (~consensusAddress) => {
  let blocksSub =
    BlockSub.getListByConsensusAddress(~address=consensusAddress, ~pageSize=10, ~page=1, ());
  <div className=Styles.blockWrapper>
    {switch (blocksSub) {
     | Data(blocks) =>
       blocks
       ->Belt_Array.mapWithIndex((i, e) => renderBodyMobile(i, Sub.resolve(e)))
       ->React.array
     | _ => <Loading />
     }}
  </div>;
};
