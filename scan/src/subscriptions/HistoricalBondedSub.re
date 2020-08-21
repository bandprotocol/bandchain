type t = {
  t: int,
  y: int,
};

module HistoricalConfig = [%graphql
  {|
  subscription HistoricalBondedToken($operator_address: String!) {
    historical_bonded_token_on_validators(where: {validator: {operator_address: {_eq: $operator_address}}}) {
      bonded_tokens
      timestamp
    }
  }
|}
];

let get = operatorAddress => {
  let (resultSub, _) =
    ApolloHooks.useSubscription(
      HistoricalConfig.definition,
      ~variables=
        HistoricalConfig.makeVariables(
          ~operator_address=operatorAddress |> Address.toOperatorBech32,
          (),
        ),
    );

  let%Sub result = resultSub;
  let x = result##historical_bonded_token_on_validators;
  Sub.resolve(
    x->Belt.Array.map(each => {
      Js.Json.{
        t: each##timestamp |> GraphQLParser.timestamp |> MomentRe.Moment.toUnix,
        y:
          (each##bonded_tokens |> decodeString |> Belt.Option.getExn |> int_of_string) / 1_000_000,
      }
    }),
  );
};
