module Mini = {
  open TxSub.Mini;

  type request_internal = {
    oracleScriptIDInternal: ID.OracleScript.t,
    transaction: TxSub.Mini.t,
  };

  type raw_data_requests_internal = {
    requestIDInternal: ID.Request.t,
    request: request_internal,
  };

  type internal_t = {rawDataRequests: array(raw_data_requests_internal)};

  type t = {
    requestID: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
    txHash: Hash.t,
    blockHeight: ID.Block.t,
    timestamp: MomentRe.Moment.t,
  };

  module MultiMiniByDataSourceConfig = [%graphql
    {|
      subscription RequestMiniByDataSource($id: bigint!, $limit: Int!, $offset: Int!) {
        data_sources_by_pk(id: $id) @bsRecord {
          rawDataRequests: raw_data_requests(limit: $limit, offset: $offset) @bsRecord {
            requestIDInternal: request_id @bsDecoder(fn: "ID.Request.fromJson")
            request @bsRecord {
              oracleScriptIDInternal: oracle_script_id @bsDecoder(fn: "ID.OracleScript.fromJson")
              transaction @bsRecord {
                txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
                blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
                timestamp @bsDecoder(fn: "GraphQLParser.time")
              }
            }
          }
        }
      }
    |}
  ];

  let toExternal = ({rawDataRequests}) =>
    rawDataRequests->Belt_Array.map(
      (
        {
          requestIDInternal,
          request: {oracleScriptIDInternal, transaction: {txHash, blockHeight, timestamp}},
        },
      ) =>
      {
        requestID: requestIDInternal,
        oracleScriptID: oracleScriptIDInternal,
        txHash,
        blockHeight,
        timestamp,
      }
    );

  let getListByDataSource = (id, ~page, ~pageSize, ()) => {
    let offset = (page - 1) * pageSize;
    let (result, _) =
      ApolloHooks.useSubscription(
        MultiMiniByDataSourceConfig.definition,
        ~variables=
          MultiMiniByDataSourceConfig.makeVariables(
            ~id=id |> ID.DataSource.toJson,
            ~limit=pageSize,
            ~offset,
            (),
          ),
      );
    result
    |> Sub.map(_, x =>
         switch (x##data_sources_by_pk) {
         | Some(data) => Sub.resolve(data |> toExternal)
         | None => NoData
         }
       );
  };
};

// let count = () => {
//   let (result, _) = ApolloHooks.useSubscription(DataSourcesCountConfig.definition);
//   result
//   |> Sub.map(_, x => x##data_sources_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
// };
