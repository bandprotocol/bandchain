type t = {
  name: string,
  transaction: option(TxSub.Mini.t),
};

// module RevisionsConfig = [%graphql
//   {|
//   subscription DataSourceRevisions($id: bigint!) {
//     data_source_revisions(
//       where: {data_source_id: {_eq: $id}}
//       order_by: {revision_number: desc}
//     ) @bsRecord{
//       name
//       transaction @bsRecord {
//         txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
//         blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
//         timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
//       }
//     }
//   }
// |}
// ];

// module RevisionCountConfig = [%graphql
//   {|
//   subscription DataSourceRevisionCount($id: bigint!) {
//     data_source_revisions_aggregate(where: {data_source_id: {_eq: $id}}) {
//       aggregate {
//         count @bsDecoder(fn: "Belt_Option.getExn")
//       }
//     }
//   }
// |}
// ];

let get = _ => {
  // let (result, _) =
  //   ApolloHooks.useSubscription(
  //     RevisionsConfig.definition,
  //     ~variables=RevisionsConfig.makeVariables(~id=id |> ID.DataSource.toJson, ()),
  //   );
  // result |> Sub.map(_, x => x##data_source_revisions);
  Sub.resolve([|
    {name: "sdsd", transaction: None},
  |]);
};

let count = _ => {
  // let (result, _) =
  //   ApolloHooks.useSubscription(
  //     RevisionCountConfig.definition,
  //     ~variables=RevisionCountConfig.makeVariables(~id=id |> ID.DataSource.toJson, ()),
  //   );
  // result
  // |> Sub.map(_, x =>
  //      x##data_source_revisions_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
  //    );
  Sub.resolve(
    1,
  );
};
