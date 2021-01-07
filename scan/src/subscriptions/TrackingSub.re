type t = {
  chainID: string,
  replayOffset: int,
};

module Config = [%graphql
  {|
  subscription Tracking {
    tracking @bsRecord {
      chainID: chain_id
      replayOffset: replay_offset
    }
  }
|}
];

let use = () => {
  let (result, _) = ApolloHooks.useSubscription(Config.definition);
  result |> Sub.map(_, internal => internal##tracking |> Belt.Array.getExn(_, 0));
};
