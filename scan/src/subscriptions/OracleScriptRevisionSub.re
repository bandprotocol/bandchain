open TxSub.Mini;

type t = {
  name: string,
  transaction: option(TxSub.Mini.t),
};

// module RevisionsConfig = [%graphql
//   {|
//   subscription OracleScriptRevisions($id: bigint!) {
//     oracle_script_revisions(
//       where: {oracle_script_id: {_eq: $id}}
//       order_by: {revision_number: desc}
//     ) @bsRecord{
//       name
//       transaction @bsRecord {
//         txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
//         blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJsonString")
//         timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
//       }
//     }
//   }
// |}
// ];

// module RevisionCountConfig = [%graphql
//   {|
//   subscription OracleScriptRevisionCount($id: bigint!) {
//     oracle_script_revisions_aggregate(where: {oracle_script_id: {_eq: $id}}) {
//       aggregate {
//         count @bsDecoder(fn: "Belt_Option.getExn")
//       }
//     }
//   }
// |}
// ];

let get = id => {
  // let (result, _) =
  //   ApolloHooks.useSubscription(
  //     RevisionsConfig.definition,
  //     ~variables=RevisionsConfig.makeVariables(~id=id |> ID.OracleScript.toJson, ()),
  //   );
  // result |> Sub.map(_, x => x##oracle_script_revisions);
  Sub.resolve([|
    {name: "sdsd", transaction: None},
  |]);
};

let count = id => {
  // let (result, _) =
  //   ApolloHooks.useSubscription(
  //     RevisionCountConfig.definition,
  //     ~variables=RevisionCountConfig.makeVariables(~id=id |> ID.OracleScript.toJson, ()),
  //   );
  // result
  // |> Sub.map(_, x =>
  //      x##oracle_script_revisions_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
  //    );
  Sub.resolve(
    2,
  );
};
