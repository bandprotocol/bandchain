module Mini = {
  type t = {
    consensusAddress: string,
    operatorAddress: Address.t,
    moniker: string,
  };
};

type node_status_t = {
  uptime: float,
  avgResponseTime: int,
};

type delegator_t = {
  delegator: string,
  sharePercentage: float,
  amount: int,
};

type internal_t = {
  operatorAddress: Address.t,
  consensusAddress: Address.t,
  moniker: string,
  identity: string,
  website: string,
  tokens: float,
  commissionRate: float,
  consensusPubKey: PubKey.t,
  bondedHeight: int,
  jailed: bool,
  electedCount: float,
  votedCount: float,
  details: string,
};

type t = {
  avgResponseTime: int,
  isActive: bool,
  operatorAddress: Address.t,
  consensusAddress: Address.t,
  consensusPubKey: PubKey.t,
  rewardDestinationAddress: string,
  votingPower: float,
  moniker: string,
  identity: string,
  website: string,
  details: string,
  tokens: float,
  commission: float,
  bondedHeight: int,
  completedRequestCount: int,
  missedRequestCount: int,
  nodeStatus: node_status_t,
  delegators: list(delegator_t),
};

let toExternal =
    (
      {
        operatorAddress,
        consensusAddress,
        moniker,
        identity,
        website,
        tokens,
        commissionRate,
        consensusPubKey,
        bondedHeight,
        jailed,
        electedCount,
        votedCount,
        details,
      }: internal_t,
    ) => {
  avgResponseTime: 2,
  isActive: !jailed,
  operatorAddress,
  consensusAddress,
  consensusPubKey,
  rewardDestinationAddress: "band17ljds2gj3kds234lkg",
  votingPower: tokens,
  moniker,
  identity,
  website,
  details,
  tokens,
  commission: commissionRate *. 100.,
  bondedHeight,
  completedRequestCount: 23459,
  missedRequestCount: 20,
  nodeStatus: {
    uptime: votedCount /. electedCount *. 100.,
    avgResponseTime: 2,
  },
  delegators: [
    {
      delegator: "bandvaloper1cg26m90y3wk50p9dn8pema8zmaa22plx3ensxr",
      sharePercentage: 12.0,
      amount: 12,
    },
  ],
};

module SingleConfig = [%graphql
  {|
      subscription Validator($operator_address: String!) {
        validators_by_pk(operator_address: $operator_address) @bsRecord {
          operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
          consensusAddress: consensus_address @bsDecoder(fn: "Address.fromHex")
          moniker
          identity
          website
          tokens @bsDecoder(fn: "GraphQLParser.floatWithMillionDivision")
          commissionRate: commission_rate @bsDecoder(fn: "float_of_string")
          consensusPubKey: consensus_pubkey @bsDecoder(fn: "PubKey.fromBech32")
          bondedHeight: bonded_height @bsDecoder(fn: "GraphQLParser.int64")
          jailed
          votedCount: voted_count @bsDecoder(fn: "float_of_int")
          electedCount: elected_count @bsDecoder(fn: "float_of_int")
          details
        }
      }
  |}
];

module MultiConfig = [%graphql
  {|
      subscription Validator($limit: Int!, $offset: Int!, $jailed: Boolean!) {
        validators(limit: $limit, offset: $offset, where: {jailed: {_eq: $jailed}}, order_by: {tokens: desc}) @bsRecord {
          operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
          consensusAddress: consensus_address @bsDecoder(fn: "Address.fromHex")
          moniker
          identity
          website
          tokens @bsDecoder(fn: "GraphQLParser.floatWithMillionDivision")
          commissionRate: commission_rate @bsDecoder(fn: "float_of_string")
          consensusPubKey: consensus_pubkey @bsDecoder(fn: "PubKey.fromBech32")
          bondedHeight: bonded_height @bsDecoder(fn: "GraphQLParser.int64")
          jailed
          votedCount: voted_count @bsDecoder(fn: "float_of_int")
          electedCount: elected_count @bsDecoder(fn: "float_of_int")
          details
        }
      }
  |}
];

module TotalBondedAmountConfig = [%graphql
  {|
  subscription TotalBondedAmount{
    validators_aggregate{
      aggregate{
        sum{
          tokens @bsDecoder(fn: "GraphQLParser.numberWithDefault")
        }
      }
    }
  }
  |}
];

module ValidatorCountConfig = [%graphql
  {|
    subscription Validator {
      validators_aggregate{
        aggregate{
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

module ValidatorCountByJailedConfig = [%graphql
  {|
    subscription Validator($jailed: Boolean!) {
      validators_aggregate(where: {jailed: {_eq: $jailed}}) {
        aggregate{
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

let get = operator_address => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=
        SingleConfig.makeVariables(
          ~operator_address=operator_address |> Address.toOperatorBech32,
          (),
        ),
    );
  let%Sub x = result;
  switch (x##validators_by_pk) {
  | Some(data) => Sub.resolve(data |> toExternal)
  | None => NoData
  };
};

let getList = (~page, ~pageSize, ~isActive, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ~jailed=!isActive, ()),
    );
  result |> Sub.map(_, x => x##validators->Belt_Array.map(toExternal));
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(ValidatorCountConfig.definition);
  result
  |> Sub.map(_, x => x##validators_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};

let countByActive = isActive => {
  let (result, _) =
    ApolloHooks.useSubscription(
      ValidatorCountByJailedConfig.definition,
      ~variables=ValidatorCountByJailedConfig.makeVariables(~jailed=!isActive, ()),
    );
  result
  |> Sub.map(_, x => x##validators_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};

let getTotalBondedAmount = () => {
  let (result, _) = ApolloHooks.useSubscription(TotalBondedAmountConfig.definition);
  result
  |> Sub.map(_, a =>
       ((a##validators_aggregate##aggregate |> Belt_Option.getExn)##sum |> Belt_Option.getExn)##tokens
     );
};

module GlobalInfo = {
  type t = {
    totalSupply: int,
    inflationRate: float,
    avgBlockTime: float,
  };

  let getGlobalInfo = _ => {
    {totalSupply: 300849023, inflationRate: 12.45, avgBlockTime: 2.59};
  };
};
