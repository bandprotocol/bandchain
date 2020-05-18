type t = {balance: Coin.t};

module SingleConfig = [%graphql
  {|
    subscription Unbonding($delegator_address: String!) {
        unbonding_delegations_aggregate(where: {delegator_address: {_eq: $delegator_address}}) {
          aggregate {
            sum {
              balance @bsDecoder(fn: "GraphQLParser.coinWithDefault")
            }
          }
        }
    }
|}
];

let getUnbondingBalance = delegatorAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=
        SingleConfig.makeVariables(~delegator_address=delegatorAddress |> Address.toBech32, ()),
    );

  let unbondingInfoSub =
    result
    |> Sub.map(_, a =>
         (a##unbonding_delegations_aggregate##aggregate |> Belt_Option.getExn)##sum
         |> Belt_Option.getExn
       );

  let%Sub unbondingInfo = unbondingInfoSub;
  unbondingInfo##balance |> Sub.resolve;
};
