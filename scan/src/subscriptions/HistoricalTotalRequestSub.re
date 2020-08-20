type t = {
  t: int,
  y: int,
};

module HistoricalConfig = [%graphql
  {|
  subscription HistoricalBondedToken {
    request_count_per_day(order_by: {timestamp: asc}) {
      count
      timestamp
    }
  }
|}
];

let get = () => {
  let (resultSub, _) = ApolloHooks.useSubscription(HistoricalConfig.definition);

  let%Sub result = resultSub;
  let x =
    result##request_count_per_day
    ->Belt.Array.map(each => {
        {
          t:
            each##timestamp
            |> GraphQLParser.timestampOpt
            |> Belt.Option.getExn
            |> MomentRe.Moment.toUnix,
          y: each##count |> Belt.Option.getExn,
        }
      });

  let cumulativeSum =
    x->Belt.Array.reduceWithIndex([||], (a, each, i) =>
      a->Belt.Array.concat([|
        a->Belt.Array.length > 0 ? {t: each.t, y: each.y + a[i - 1].y} : each,
      |])
    );

  Sub.resolve(cumulativeSum);
};
