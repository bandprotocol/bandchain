type t = {balance: Coin.t};

type unbonding_status_t = {
  completionTime: MomentRe.Moment.t,
  amount: Coin.t,
};

module SingleConfig = [%graphql
  {|
    subscription Unbonding($delegator_address: String!) {
      accounts_by_pk(address: $delegator_address){
        unbonding_delegations_aggregate {
          aggregate {
            sum {
              amount @bsDecoder(fn: "GraphQLParser.coinWithDefault")
            }
          }
        }
      }
    }
|}
];

module MultiConfig = [%graphql
  {|
  subscription Unbonding($delegator_address: String!, $operator_address: String!) {
  accounts_by_pk(address: $delegator_address) {
    unbonding_delegations(where: {validator: {operator_address: {_eq: $operator_address}}}) @bsRecord {
      completionTime: completion_time @bsDecoder(fn: "GraphQLParser.timestamp")
      amount @bsDecoder(fn: "GraphQLParser.coin")
    }
  }
  }
|}
];

module UnbondingByValidatorConfig = [%graphql
  {|
    subscription Unbonding($delegator_address: String!, $operator_address: String!) {
      accounts_by_pk(address: $delegator_address) {
        unbonding_delegations_aggregate(where: {validator: {operator_address: {_eq: $operator_address}}}) {
          aggregate {
            sum {
              amount @bsDecoder(fn: "GraphQLParser.coinWithDefault")
            }
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
         (
           (a##accounts_by_pk |> Belt.Option.getExn)##unbonding_delegations_aggregate##aggregate
           |> Belt_Option.getExn
         )##sum
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
    |> Sub.map(_, a => {
         switch (a##accounts_by_pk) {
         | Some(account) =>
           (
             (account##unbonding_delegations_aggregate##aggregate |> Belt_Option.getExn)##sum
             |> Belt_Option.getExn
           )##amount
         | None => Coin.newUBANDFromAmount(0.)
         }
       });

  let%Sub unbondingInfo = unbondingInfoSub;
  unbondingInfo |> Sub.resolve;
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
  result
  |> Sub.map(_, x => {
       switch (x##accounts_by_pk) {
       | Some(x') => x'##unbonding_delegations
       | None => [||]
       }
     });
};
