module Mini = {
  open TxSub.Mini;

  type oracle_script_internal_t = {
    id: ID.OracleScript.t,
    name: string,
  };

  type request_internal = {
    oracleScript: oracle_script_internal_t,
    transaction: TxSub.Mini.t,
  };

  type raw_data_requests_internal_t = {
    idInternal: ID.Request.t,
    request: request_internal,
  };

  type internal_t = {rawDataRequests: array(raw_data_requests_internal_t)};

  type t = {
    id: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    txHash: Hash.t,
    blockHeight: ID.Block.t,
    timestamp: MomentRe.Moment.t,
  };

  module MultiMiniByDataSourceConfig = [%graphql
    {|
      subscription RequestMiniByDataSource($id: bigint!, $limit: Int!, $offset: Int!) {
        data_sources_by_pk(id: $id) @bsRecord {
          rawDataRequests: raw_data_requests(limit: $limit, offset: $offset) @bsRecord {
            idInternal: request_id @bsDecoder(fn: "ID.Request.fromJson")
            request @bsRecord {
              oracleScript: oracle_script @bsRecord {
                id @bsDecoder(fn: "ID.OracleScript.fromJson")
                name
              }
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
      ({idInternal, request: {oracleScript, transaction: {txHash, blockHeight, timestamp}}}) =>
      {
        id: idInternal,
        oracleScriptID: oracleScript.id,
        oracleScriptName: oracleScript.name,
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
         | Some(data) => data |> toExternal
         | None => [||]
         }
       );
  };
};

module RequestCountByDataSourceConfig = [%graphql
  {|
    subscription RequestMiniByDataSourceCount($id: bigint!) {
      data_sources_by_pk(id: $id) {
        raw_data_requests_aggregate {
          aggregate {
            count @bsDecoder(fn: "Belt_Option.getExn")
          }
        }
      }
    }
  |}
];

let countByDataSource = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      RequestCountByDataSourceConfig.definition,
      ~variables=RequestCountByDataSourceConfig.makeVariables(~id=id |> ID.DataSource.toJson, ()),
    );
  result
  |> Sub.map(_, x => {
       {
         let%Opt dataSource = x##data_sources_by_pk;
         let%Opt aggregate = dataSource##raw_data_requests_aggregate##aggregate;
         Some(aggregate##count);
       }
       ->Belt_Option.getWithDefault(0)
     });
};
