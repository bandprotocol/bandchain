type t = {
  name: string,
  timestamp: MomentRe.Moment.t,
  height: int,
  txHash: Hash.t,
};

module RevisionsConfig = [%graphql
  {|
  subscription DataSourceRevisions($id: bigint!) {
    data_source_revisions(
      where: {data_source_id: {_eq: $id}}
      order_by: {revision_number: desc}
    ) @bsRecord{
      name
      timestamp @bsDecoder(fn: "GraphQLParser.time")
      height: block_height @bsDecoder(fn: "GraphQLParser.int64")
      txHash: tx_hash@bsDecoder(fn: "GraphQLParser.hash")
    }
  }
|}
];

module RevisionCountConfig = [%graphql
  {|
  subscription DataSourceRevisionCount($id: bigint!) {
    data_source_revisions_aggregate(where: {data_source_id: {_eq: $id}}) {
      aggregate {
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

let get = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      RevisionsConfig.definition,
      ~variables=RevisionsConfig.makeVariables(~id=id |> ID.DataSource.toJson, ()),
    );
  result |> Sub.map(_, x => x##data_source_revisions);
};

let count = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      RevisionCountConfig.definition,
      ~variables=RevisionCountConfig.makeVariables(~id=id |> ID.DataSource.toJson, ()),
    );
  result
  |> Sub.map(_, x =>
       x##data_source_revisions_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
     );
};
