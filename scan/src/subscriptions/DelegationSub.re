type t = {
  delegatorAddress: Address.t,
  validatorAddress: Address.t,
  shares: float,
};

type stake_t = {
  amount: Coin.t,
  delegatorAddress: Address.t,
  moniker: string,
  operatorAddress: Address.t,
  reward: Coin.t,
  sharePercentage: float,
  identity: string,
};

type stake_aggregate_t = {
  amount: Coin.t,
  reward: Coin.t,
};

module StakeConfig = [%graphql
  {|
  subscription Stake($limit: Int!, $offset: Int!, $delegator_address: String!)  {
    delegations_view(offset: $offset, limit: $limit, order_by: {amount: desc}, where: {delegator_address: {_eq: $delegator_address}}) @bsRecord  {
      amount @bsDecoder(fn: "GraphQLParser.coinExn")
      reward @bsDecoder(fn: "GraphQLParser.coinExn")
      sharePercentage: share_percentage @bsDecoder(fn: "GraphQLParser.floatWithDefault")
      delegatorAddress: delegator_address @bsDecoder(fn: "GraphQLParser.addressExn")
      operatorAddress: operator_address @bsDecoder(fn: "GraphQLParser.addressExn")
      moniker @bsDecoder(fn: "GraphQLParser.stringExn")
      identity @bsDecoder(fn: "GraphQLParser.stringExn")
    }
  }
  |}
];

module TotalStakeByDelegatorConfig = [%graphql
  {|
  subscription TotalStake($delegator_address: String!) {
    delegations_view_aggregate(where: {delegator_address: {_eq: $delegator_address}}){
      aggregate{
        sum{
          amount @bsDecoder(fn: "GraphQLParser.coinWithDefault")
          reward @bsDecoder(fn: "GraphQLParser.coinWithDefault")
        }
      }
    }
  }
  |}
];

module StakeByValidatorConfig = [%graphql
  {|
  subscription StakeByValidator($delegator_address: String!, $operator_address: String!) {
    delegations_view(where: {_and: {delegator_address: {_eq: $delegator_address}, operator_address: {_eq: $operator_address}}}) @bsRecord {
      amount @bsDecoder(fn: "GraphQLParser.coinExn")
      reward @bsDecoder(fn: "GraphQLParser.coinExn")
    }
  }
|}
];

module StakeCountByDelegatorConfig = [%graphql
  {|
  subscription CountByDelegator($delegator_address: String!) {
    delegations_view_aggregate(where: {delegator_address: {_eq: $delegator_address}}) {
      aggregate {
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

module DelegatorsByValidatorConfig = [%graphql
  {|
  subscription Stake($limit: Int!, $offset: Int!, $operator_address: String!)  {
    delegations_view(offset: $offset, limit: $limit, order_by: {amount: desc}, where: {operator_address: {_eq: $operator_address}}) @bsRecord  {
      amount @bsDecoder(fn: "GraphQLParser.coinExn")
      reward @bsDecoder(fn: "GraphQLParser.coinExn")
      sharePercentage: share_percentage @bsDecoder(fn: "GraphQLParser.floatWithDefault")
      delegatorAddress: delegator_address @bsDecoder(fn: "GraphQLParser.addressExn")
      operatorAddress: operator_address @bsDecoder(fn: "GraphQLParser.addressExn")
      moniker @bsDecoder(fn: "GraphQLParser.stringExn")
      identity @bsDecoder(fn: "GraphQLParser.stringExn")
    }
  }
  |}
];

module DelegatorCountConfig = [%graphql
  {|
    subscription DelegatorCount($operator_address: String!) {
      delegations_view_aggregate(where: {operator_address: {_eq: $operator_address}}) {
        aggregate {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

let getStakeList = (delegatorAddress, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      StakeConfig.definition,
      ~variables=
        StakeConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Sub.map(_, x => x##delegations_view);
};

let getTotalStakeByDelegator = delegatorAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      TotalStakeByDelegatorConfig.definition,
      ~variables=
        TotalStakeByDelegatorConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          (),
        ),
    );

  let delegatorInfoSub =
    result
    |> Sub.map(_, a =>
         (a##delegations_view_aggregate##aggregate |> Belt_Option.getExn)##sum
         |> Belt_Option.getExn
       );

  let%Sub delegatorInfo = delegatorInfoSub;
  {amount: delegatorInfo##amount, reward: delegatorInfo##reward} |> Sub.resolve;
};

let getStakeByValidator = (delegatorAddress, operatorAddress) => {
  let (result, _) = {
    ApolloHooks.useSubscription(
      StakeByValidatorConfig.definition,
      ~variables=
        StakeByValidatorConfig.makeVariables(
          ~operator_address=operatorAddress |> Address.toOperatorBech32,
          ~delegator_address=delegatorAddress |> Address.toBech32,
          (),
        ),
    );
  };
  result
  |> Sub.map(_, internal =>
       internal##delegations_view
       ->Belt_Array.get(0)
       ->Belt_Option.getWithDefault({
           amount: Coin.newUBANDFromAmount(0.),
           reward: Coin.newUBANDFromAmount(0.),
         })
     );
};

let getStakeCountByDelegator = delegatorAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      StakeCountByDelegatorConfig.definition,
      ~variables=
        StakeCountByDelegatorConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          (),
        ),
    );
  result
  |> Sub.map(_, x =>
       x##delegations_view_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
     );
};

let getDelegatorsByValidator = (validatorAddress, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      DelegatorsByValidatorConfig.definition,
      ~variables=
        DelegatorsByValidatorConfig.makeVariables(
          ~operator_address=validatorAddress |> Address.toOperatorBech32,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Sub.map(_, x => x##delegations_view);
};

let getDelegatorCountByValidator = validatorAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      DelegatorCountConfig.definition,
      ~variables=
        DelegatorCountConfig.makeVariables(
          ~operator_address=validatorAddress |> Address.toOperatorBech32,
          (),
        ),
    );
  result
  |> Sub.map(_, x =>
       x##delegations_view_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
     );
};
