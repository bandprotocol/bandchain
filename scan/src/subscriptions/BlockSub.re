type aggregate_t = {count: int};

type transactions_aggregate_t = {aggregate: option(aggregate_t)};

type internal_t = {
  height: int,
  hash: Hash.t,
  proposer: Address.t,
  timestamp: MomentRe.Moment.t,
  transactions_aggregate: transactions_aggregate_t,
};

type t = {
  height: int,
  hash: Hash.t,
  proposer: Address.t,
  timestamp: MomentRe.Moment.t,
  txn: int,
};

let toExternal = ({height, hash, proposer, timestamp, transactions_aggregate}) => {
  height,
  hash,
  proposer,
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
    blocks(limit: $limit, offset: $offset) @bsRecord {
      height @bsDecoder(fn: "GraphQLParser.int64")
      hash: block_hash @bsDecoder(fn: "GraphQLParser.hash")
      proposer @bsDecoder(fn: "Address.fromHex")
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
  subscription DataSource($height: Int!) {
    blocks_by_pk(height: $height) @bsRecord {
      height @bsDecoder(fn: "GraphQLParser.int64")
      hash: block_hash @bsDecoder(fn: "GraphQLParser.hash")
      proposer @bsDecoder(fn: "Address.fromHex")
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
  subscription DataSourcesCount {
    blocks_aggregate{
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
      ~variables=SingleConfig.makeVariables(~height, ()),
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

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(BlockCountConfig.definition);
  result
  |> Sub.map(_, x => x##blocks_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};
