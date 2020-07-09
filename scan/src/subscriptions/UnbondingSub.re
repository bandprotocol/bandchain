type t = {balance: Coin.t};

type unbonding_status_t = {
  completionTime: MomentRe.Moment.t,
  amount: Coin.t,
};

module SingleConfig = [%graphql
  {|
    subscription Unbonding($delegator_address: String!) {
        unbonding_delegations_aggregate(where: {delegator_address: {_eq: $delegator_address}}) {
          aggregate {
            sum {
              amount @bsDecoder(fn: "GraphQLParser.coinWithDefault")
            }
          }
        }
    }
|}
];

module MultiConfig = [%graphql
  {|
  subscription Unbonding($delegator_address: String!, $operator_address: String!) {
  unbonding_delegations(where: {_and: {delegator_address: {_eq: $delegator_address}, operator_address: {_eq: $operator_address}}}, order_by: {completion_time: asc}) @bsRecord {
    completionTime: completion_time @bsDecoder(fn: "GraphQLParser.timestamp")
    amount @bsDecoder(fn: "GraphQLParser.coin")
  }
  }
|}
];

module UnbondingByValidatorConfig = [%graphql
  {|
    subscription Unbonding($delegator_address: String!, $operator_address: String!) {
        unbonding_delegations_aggregate(where: {_and: {delegator_address: {_eq: $delegator_address}, operator_address: {_eq: $operator_address}}}) {
          aggregate {
            sum {
              amount @bsDecoder(fn: "GraphQLParser.coinWithDefault")
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
  unbondingInfo##amount |> Sub.resolve;
};

let getUnbondingBalanceByValidator = (delegatorAddress, operatorAddress) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      UnbondingByValidatorConfig.definition,
      ~variables=
        UnbondingByValidatorConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          ~operator_address=operatorAddress |> Address.toOperatorBech32,
          (),
        ),
    );

  let unbondingInfoSub =
    result
    |> Sub.map(_, a =>
         (a##unbonding_delegations_aggregate##aggregate |> Belt_Option.getExn)##sum
         |> Belt_Option.getExn
       );

  let%Sub unbondingInfo = unbondingInfoSub;
  unbondingInfo##amount |> Sub.resolve;
};

let getUnbondingList = (delegatorAddress, operatorAddress) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=
        MultiConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          ~operator_address=operatorAddress |> Address.toOperatorBech32,
          (),
        ),
    );
  result |> Sub.map(_, x => x##unbonding_delegations);
};
