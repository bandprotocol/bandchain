open ValidatorSub.Mini;

type aggregate_t = {count: int};

type transactions_aggregate_t = {aggregate: option(aggregate_t)};

type internal_t = {
  height: ID.Block.t,
  hash: Hash.t,
  validator: ValidatorSub.Mini.t,
  timestamp: MomentRe.Moment.t,
  transactions_aggregate: transactions_aggregate_t,
};

type t = {
  height: ID.Block.t,
  hash: Hash.t,
  validator: ValidatorSub.Mini.t,
  timestamp: MomentRe.Moment.t,
  txn: int,
};

let toExternal = ({height, hash, validator, timestamp, transactions_aggregate}) => {
  height,
  hash,
  validator,
  timestamp,
  txn:
    switch (transactions_aggregate.aggregate) {
    | Some(aggregate) => aggregate.count
    | _ => 0
    },
};

module MultiConfig = [%graphql
  {|
  subscription Blocks($limit: Int!, $offset: Int!) {
    blocks(limit: $limit, offset: $offset, order_by: {height: desc}) @bsRecord {
      height @bsDecoder(fn: "ID.Block.fromJson")
      hash: block_hash @bsDecoder(fn: "GraphQLParser.hash")
      validator @bsRecord {
        consensusAddress: consensus_address
        operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
        moniker
      }
      timestamp @bsDecoder(fn: "GraphQLParser.time")
      transactions_aggregate @bsRecord {
        aggregate @bsRecord {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  }
|}
];

module MultiConsensusAddressConfig = [%graphql
  {|
  subscription BlocksByConsensusAddress($limit: Int!, $offset: Int!, $address: String!) {
    blocks(limit: $limit, offset: $offset, order_by: {height: desc}, where: {proposer: {_eq: $address}}) @bsRecord {
      height @bsDecoder(fn: "ID.Block.fromJson")
      hash: block_hash @bsDecoder(fn: "GraphQLParser.hash")
      validator @bsRecord {
        consensusAddress: consensus_address
        operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
        moniker
      }
      timestamp @bsDecoder(fn: "GraphQLParser.time")
      transactions_aggregate @bsRecord {
        aggregate @bsRecord {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  }
|}
];

module SingleConfig = [%graphql
  {|
  subscription Block($height: bigint!) {
    blocks_by_pk(height: $height) @bsRecord {
      height @bsDecoder(fn: "ID.Block.fromJson")
      hash: block_hash @bsDecoder(fn: "GraphQLParser.hash")
      validator @bsRecord {
        consensusAddress: consensus_address
        operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
        moniker
      }
      timestamp @bsDecoder(fn: "GraphQLParser.time")
      transactions_aggregate @bsRecord {
        aggregate @bsRecord {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  },
|}
];

module BlockCountConfig = [%graphql
  {|
  subscription BlocksCount {
    blocks_aggregate{
      aggregate{
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

module BlockCountConsensusAddressConfig = [%graphql
  {|
  subscription BlocksCountByConsensusAddress($address: String!) {
    blocks_aggregate(where: {proposer: {_eq: $address}}) {
      aggregate{
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

let get = height => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=SingleConfig.makeVariables(~height=height |> ID.Block.toJson, ()),
    );
  let%Sub x = result;
  switch (x##blocks_by_pk) {
  | Some(data) => Sub.resolve(data |> toExternal)
  | None => NoData
  };
};

let getList = (~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ()),
    );
  result |> Sub.map(_, internal => internal##blocks->Belt_Array.map(toExternal));
};

let getListByConsensusAddress = (~address, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConsensusAddressConfig.definition,
      ~variables=
        MultiConsensusAddressConfig.makeVariables(
          ~address=address |> Address.toHex(~upper=true),
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Sub.map(_, internal => internal##blocks->Belt_Array.map(toExternal));
};

let getLatest = () => {
  let%Sub blocks = getList(~pageSize=1, ~page=1, ());
  switch (blocks->Belt_Array.get(0)) {
  | Some(latestBlock) => latestBlock |> Sub.resolve
  | None => NoData
  };
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(BlockCountConfig.definition);
  result
  |> Sub.map(_, x => x##blocks_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};

let countByConsensusAddress = (~address, ()) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      BlockCountConsensusAddressConfig.definition,
      ~variables=
        BlockCountConsensusAddressConfig.makeVariables(
          ~address=address |> Address.toHex(~upper=true),
          (),
        ),
    );
  result
  |> Sub.map(_, x => x##blocks_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};
