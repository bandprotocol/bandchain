type t = {balance: Coin.t};

type unbonding_status_t = {
  completionTime: MomentRe.Moment.t,
  amount: Coin.t,
};

type validator_t = {
  operatorAddress: Address.t,
  moniker: string,
  identity: string,
};

type unbonding_list_t = {
  amount: Coin.t,
  completionTime: MomentRe.Moment.t,
  validator: validator_t,
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
  subscription Unbonding($delegator_address: String!, $operator_address: String!, $completion_time: timestamp) {
  accounts_by_pk(address: $delegator_address) {
    unbonding_delegations(order_by: {completion_time: asc}, where: {_and: {completion_time: {_gte: $completion_time}, validator: {operator_address: {_eq: $operator_address}}}}) @bsRecord {
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

module UnbondingByDelegatorConfig = [%graphql
  {|
    subscription UnbondingByDelegator($limit: Int!, $offset: Int!, $delegator_address: String!, $current_time: timestamp) {
      accounts_by_pk(address: $delegator_address) {
        unbonding_delegations(offset: $offset, limit: $limit, order_by: {completion_time: asc}, where: {completion_time: {_gte: $current_time}}) @bsRecord{
          amount @bsDecoder(fn: "GraphQLParser.coin")
          completionTime: completion_time @bsDecoder(fn: "GraphQLParser.timestamp")
          validator @bsRecord{
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
            identity
          }
        }
      }
    }
  |}
];

module UnbondingCountByDelegatorConfig = [%graphql
  {|
    subscription UnbondingCountByDelegator($delegator_address: String!, $current_time: timestamp) {
      accounts_by_pk(address: $delegator_address) {
        unbonding_delegations_aggregate(where: {completion_time: {_gte: $current_time}}) {
          aggregate{
            count @bsDecoder(fn: "Belt_Option.getExn")
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

let getUnbondingByDelegator = (delegatorAddress, currentTime, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      UnbondingByDelegatorConfig.definition,
      ~variables=
        UnbondingByDelegatorConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          ~limit=pageSize,
          ~current_time=currentTime |> Js.Json.string,
          ~offset,
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

let getUnbondingCountByDelegator = (delegatorAddress, currentTime) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      UnbondingCountByDelegatorConfig.definition,
      ~variables=
        UnbondingCountByDelegatorConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          ~current_time=currentTime |> Js.Json.string,
          (),
        ),
    );
  result
  |> Sub.map(_, x => {
       switch (x##accounts_by_pk) {
       | Some(x') =>
         x'##unbonding_delegations_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
       | None => 0
       }
     });
};

let getUnbondingList = (delegatorAddress, operatorAddress, completionTime) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=
        MultiConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          ~operator_address=operatorAddress |> Address.toOperatorBech32,
          ~completion_time=completionTime |> Js.Json.string,
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
