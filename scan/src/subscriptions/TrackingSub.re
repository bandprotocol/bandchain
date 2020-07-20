type t = {chainID: string};

module Config = [%graphql
  {|
  subscription Tracking {
    tracking @bsRecord {
      chainID: chain_id
    }
  }
|}
];

let use = () => {
  let (result, _) = ApolloHooks.useSubscription(Config.definition);
  result |> Sub.map(_, internal => internal##tracking |> Belt.Array.getExn(_, 0));
};
