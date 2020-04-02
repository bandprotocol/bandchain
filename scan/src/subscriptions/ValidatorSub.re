type status_t =
  | Bonded
  | Unbonded
  | Unbonding;

module Mini = {
  type t = {
    consensusAddress: string,
    operatorAddress: Address.t,
    moniker: string,
  };
};
module Validator = {
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
    moniker: string,
    identity: string,
    website: string,
  };

  type t = {
    avgResponseTime: int,
    isActive: bool,
    operatorAddress: Address.t,
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
    uptime: float,
    completedRequestCount: int,
    missedRequestCount: int,
    nodeStatus: node_status_t,
    delegators: list(delegator_t),
  };

  let toExternal = ({operatorAddress, moniker, identity, website}: internal_t) => {
    avgResponseTime: 2,
    isActive: true,
    operatorAddress,
    consensusPubKey:
      "bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj"
      |> PubKey.fromBech32,
    rewardDestinationAddress: "band17ljds2gj3kds234lkg",
    votingPower: 25.0,
    moniker,
    identity,
    website,
    details: "DETAILS",
    tokens: 100.00,
    commission: 100.00,
    bondedHeight: 1,
    uptime: 100.0,
    completedRequestCount: 23459,
    missedRequestCount: 20,
    nodeStatus: {
      uptime: 100.00,
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
          moniker
          identity
          website
        }
      }
  |}
  ];

  module MultiConfig = [%graphql
    {|
      subscription Validator($limit: Int!, $offset: Int!) {
        validators(limit: $limit, offset: $offset) @bsRecord {
          operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
          moniker
          identity
          website
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

  let getList = (~page=1, ~pageSize=10, ()) => {
    let offset = (page - 1) * pageSize;
    let (result, _) =
      ApolloHooks.useSubscription(
        MultiConfig.definition,
        ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ()),
      );
    result |> Sub.map(_, x => x##validators->Belt_Array.map(toExternal));
  };

  let count = () => {
    let (result, _) = ApolloHooks.useSubscription(ValidatorCountConfig.definition);
    result
    |> Sub.map(_, x => x##validators_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
  };
};

module GlobalInfo = {
  type t = {
    totalSupply: int,
    inflationRate: float,
    avgBlockTime: float,
  };

  let getGlobalInfo = _ => {
    {totalSupply: 10849023, inflationRate: 12.45, avgBlockTime: 2.59};
  };
};
