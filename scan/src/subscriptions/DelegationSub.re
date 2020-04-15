type t = {
  delegatorAddress: Address.t,
  validatorAddress: Address.t,
  shares: float,
};

type stake_t = {
  amount: float,
  sharePercentage: float,
  delegatorAddress: Address.t,
  validatorAddress: Address.t,
};

module StakeConfig = [%graphql
  {|
  subscription Stake($limit: Int!, $offset: Int!, $delegator_address: String!)  {
    delegations_view(offset: $offset, limit: $limit, order_by: {amount: desc}, where: {delegator_address: {_eq: $delegator_address}}) @bsRecord  {
      amount @bsDecoder(fn: "GraphQLParser.numberExn")
      sharePercentage: share_percentage @bsDecoder(fn: "GraphQLParser.floatExn")
      delegatorAddress: delegator_address @bsDecoder(fn: "GraphQLParser.addressExn")
      validatorAddress: validator_address @bsDecoder(fn: "GraphQLParser.addressExn")
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
          amount @bsDecoder(fn: "GraphQLParser.numberWithDefault")
        }
      }
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
  subscription Stake($limit: Int!, $offset: Int!, $validator_address: String!)  {
    delegations_view(offset: $offset, limit: $limit, order_by: {amount: desc}, where: {validator_address: {_eq: $validator_address}}) @bsRecord  {
      amount @bsDecoder(fn: "GraphQLParser.numberExn")
      sharePercentage: share_percentage @bsDecoder(fn: "GraphQLParser.floatExn")
      delegatorAddress: delegator_address @bsDecoder(fn: "GraphQLParser.addressExn")
      validatorAddress: validator_address @bsDecoder(fn: "GraphQLParser.addressExn")
    }
  }
  |}
];

module DelegatorCountConfig = [%graphql
  {|
    subscription DelegatorCount($validator_address: String!) {
      delegations_view_aggregate(where: {validator_address: {_eq: $validator_address}}) {
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
  result
  |> Sub.map(_, a =>
       (
         (a##delegations_view_aggregate##aggregate |> Belt_Option.getExn)##sum |> Belt_Option.getExn
       )##amount
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
} /* }*/;

let getDelegatorsByValidator = (validatorAddress, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      DelegatorsByValidatorConfig.definition,
      ~variables=
        DelegatorsByValidatorConfig.makeVariables(
          ~validator_address=validatorAddress |> Address.toOperatorBech32,
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
          ~validator_address=validatorAddress |> Address.toOperatorBech32,
          (),
        ),
    );
  result
  |> Sub.map(_, x =>
       x##delegations_view_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
     );
};
