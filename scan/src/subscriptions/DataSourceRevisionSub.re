type t = {
  name: string,
  timestamp: MomentRe.Moment.t,
  height: int,
  txHash: Hash.t,
};

module Config = [%graphql
  {|
  subscription DataSourceRevisions($id: bigint!) {
    data_source_revisions(where: {data_source_id: {_eq: $id}}) @bsRecord{
      name
      timestamp @bsDecoder(fn: "GraphQLParser.time")
      height: block_height @bsDecoder(fn: "GraphQLParser.int64")
      txHash: tx_hash@bsDecoder(fn: "GraphQLParser.hash")
    }
  }
|}
];

let get = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      Config.definition,
      ~variables=Config.makeVariables(~id=id |> ID.DataSource.toJson, ()),
    );
  result |> Sub.map(_, x => x##data_source_revisions);
};
