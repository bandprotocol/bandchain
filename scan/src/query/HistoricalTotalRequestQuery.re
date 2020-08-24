type t = {
  t: int,
  y: int,
};

module HistoricalConfig = [%graphql
  {|
  query HistoricalRequest {
    request_count_per_days(order_by: {date: asc}) {
      count
      date
    }
  }
|}
];

let get = () => {
  let (resultQuery, _) = ApolloHooks.useQuery(HistoricalConfig.definition);

  let%Query result = resultQuery;
  let x =
    result##request_count_per_days
    ->Belt.Array.map(each => {
        {t: each##date |> GraphQLParser.timestamp |> MomentRe.Moment.toUnix, y: each##count}
      });

  let (cumulativeSum, _) =
    x->Belt.Array.reduce(
      ([||], 0),
      ((ls, sum), each) => {
        let newSum = sum + each.y;
        (ls->Belt.Array.concat([|{t: each.t, y: newSum}|]), newSum);
      },
    );

  Query.resolve(cumulativeSum);
};
