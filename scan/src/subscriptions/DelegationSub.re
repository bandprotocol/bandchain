type t = {
  delegatorAddress: Address.t,
  validatorAddress: Address.t,
  shares: float,
};

type stake_t = {
  amount: float,
  delegatorAddress: Address.t,
  validatorAddress: Address.t,
};

module StakeConfig = [%graphql
  {|
  subscription Stake($limit: Int!, $offset: Int!, $delegator_address: String!)  {
    delegations_view(offset: $offset, limit: $limit, order_by: {amount: desc}, where: {delegator_address: {_eq: $delegator_address}}) @bsRecord  {
      amount @bsDecoder(fn: "GraphQLParser.numberExn")
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
          amount @bsDecoder(fn: "GraphQLParser.numberExn")
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

//TODO: Change and use for Delegator Sub

// module SingleConfig = [%graphql
//   {|
//       subscription Validator($delegator_address: String!, $validator_address: String!) {
//         delegations_by_pk(delegator_address: $delegator_address, validator_address: $validator_address) @bsRecord {
//             delegatorAddress: delegator_address @bsDecoder(fn: "Address.fromBech32")
//             validatorAddress: validator_address @bsDecoder(fn: "Address.fromBech32")
//             shares @bsDecoder(fn: "float_of_string")
//         }
//       }
//   |}
// ];

// module MultiConfig = [%graphql
//   {|
//   subscription Delegation($delegator_address: String!)  {
//     delegations(where: {delegator_address: {_eq: $delegator_address}}) @bsRecord  {
//       delegatorAddress: delegator_address @bsDecoder(fn: "Address.fromBech32")
//       validatorAddress: validator_address @bsDecoder(fn: "Address.fromBech32")
//       shares @bsDecoder(fn: "float_of_string")
//     }
//   }
//   |}
// ];
// let get = (delegatorAddress, validatorAddress) => {
//   let (result, _) =
//     ApolloHooks.useSubscription(
//       SingleConfig.definition,
//       ~variables=
//         SingleConfig.makeVariables(
//           ~delegator_address=delegatorAddress |> Address.toBech32,
//           ~validator_address=validatorAddress |> Address.toOperatorBech32,
//           (),
//         ),
//     );
//   result |> Sub.map(_, x => x##delegations_by_pk);
// };

// let getList = delegatorAddress => {
//   let (result, _) =
//     ApolloHooks.useSubscription(
//       MultiConfig.definition,
//       ~variables=
//         MultiConfig.makeVariables(~delegator_address=delegatorAddress |> Address.toBech32, ()),
//     );
//   result |> Sub.map(_, internal => internal##delegations);
