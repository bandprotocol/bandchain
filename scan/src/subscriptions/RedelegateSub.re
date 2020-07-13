type t = {balance: Coin.t};

type validator_t = {
  operatorAddress: Address.t,
  moniker: string,
  identity: string,
};

type redelegate_list_t = {
  amount: Coin.t,
  completionTime: MomentRe.Moment.t,
  validator: validator_t,
};

module RedelegationByDelegatorConfig = [%graphql
  {|
    subscription UnbondingByDelegator($limit: Int!, $offset: Int!, $delegator_address: String!, $current_time: timestamp) {
      accounts_by_pk(address: $delegator_address) {
        redelegations(offset: $offset, limit: $limit, order_by: {completion_time: asc}, where: {completion_time: {_gte: $current_time}}) @bsRecord{
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

module RedelegateCountByDelegatorConfig = [%graphql
  {|
    subscription UnbondingCountByDelegator($delegator_address: String!, $current_time: timestamp) {
      accounts_by_pk(address: $delegator_address) {
        redelegations_aggregate(where: {completion_time: {_gte: $current_time}}) {
          aggregate{
            count @bsDecoder(fn: "Belt_Option.getExn")
          }
        }
      }
    }
  |}
];

let getRedelegationByDelegator = (delegatorAddress, currentTime, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      RedelegationByDelegatorConfig.definition,
      ~variables=
        RedelegationByDelegatorConfig.makeVariables(
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
       | Some(x') => x'##redelegations
       | None => [||]
       }
     });
};

let getRedelegateCountByDelegator = (delegatorAddress, currentTime) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      RedelegateCountByDelegatorConfig.definition,
      ~variables=
        RedelegateCountByDelegatorConfig.makeVariables(
          ~delegator_address=delegatorAddress |> Address.toBech32,
          ~current_time=currentTime |> Js.Json.string,
          (),
        ),
    );
  result
  |> Sub.map(_, x => {
       switch (x##accounts_by_pk) {
       | Some(x') =>
         x'##redelegations_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
       | None => 10
       }
     });
};
