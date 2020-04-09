type t = {
  delegatorAddress: Address.t,
  validatorAddress: Address.t,
  shares: float,
};

type stake_t = {
  amount: option(float),
  delegatorAddress: Address.t,
  validatorAddress: Address.t,
};

module SingleConfig = [%graphql
  {|
      subscription Validator($delegator_address: String!, $validator_address: String!) {
        delegations_by_pk(delegator_address: $delegator_address, validator_address: $validator_address) @bsRecord {
            delegatorAddress: delegator_address @bsDecoder(fn: "Address.fromBech32")
            validatorAddress: validator_address @bsDecoder(fn: "Address.fromBech32")
            shares @bsDecoder(fn: "float_of_string")
        }
      }
  |}
];

module MultiConfig = [%graphql
  {|
  subscription Delegation($delegator_address: String!)  {
    delegations(where: {delegator_address: {_eq: $delegator_address}}) @bsRecord  {
      delegatorAddress: delegator_address @bsDecoder(fn: "Address.fromBech32")
      validatorAddress: validator_address @bsDecoder(fn: "Address.fromBech32")
      shares @bsDecoder(fn: "float_of_string")
    }
  }
  |}
];

module StakeConfig = [%graphql
  {|
  subscription Stake($delegator_address: String!)  {
    delegations_view(where: {delegator_address: {_eq: $delegator_address}}) @bsRecord  {
      amount @bsDecoder(fn: "GraphQLParser.numberOpt")
      delegatorAddress: delegator_address @bsDecoder(fn: "GraphQLParser.addressOpt")
      validatorAddress: validator_address @bsDecoder(fn: "GraphQLParser.addressOpt")
    }
  }
  |}
];

module TotalStakeConfig = [%graphql
  {|
  subscription TotalStake($delegator_address: String!) {
    delegations_view_aggregate(where: {delegator_address: {_eq: $delegator_address}}){
      aggregate{
        sum{
          amount @bsDecoder(fn: "GraphQLParser.numberOpt")
        }
      }
    }
  }
  |}
];

let get = (delegatorAddress, validatorAddress) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=
        SingleConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          ~validator_address=validatorAddress |> Address.toOperatorBech32,
          (),
        ),
    );
  result |> Sub.map(_, x => x##delegations_by_pk);
};

let getList = delegatorAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=
        MultiConfig.makeVariables(~delegator_address=delegatorAddress |> Address.toBech32, ()),
    );
  result |> Sub.map(_, internal => internal##delegations);
};

let getStake = delegatorAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      StakeConfig.definition,
      ~variables=
        StakeConfig.makeVariables(~delegator_address=delegatorAddress |> Address.toBech32, ()),
    );
  result |> Sub.map(_, x => x##delegations_view);
};

let getTotalStake = delegatorAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      TotalStakeConfig.definition,
      ~variables=
        TotalStakeConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          (),
        ),
    );
  result
  |> Sub.map(_, a =>
       (
         (a##delegations_view_aggregate##aggregate |> Belt_Option.getExn)##sum |> Belt_Option.getExn
       )##amount
       |> Belt_Option.getWithDefault(_, 0.0)
     );
};
