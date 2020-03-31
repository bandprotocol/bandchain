module Mini = {
  open TxSub.Mini;

  type oracle_script_internal_t = {
    id: ID.OracleScript.t,
    name: string,
  };

  type request_internal = {
    id: ID.Request.t,
    requester: Address.t,
    oracleScript: oracle_script_internal_t,
    transaction: TxSub.Mini.t,
  };

  type t = {
    id: ID.Request.t,
    requester: Address.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    txHash: Hash.t,
    blockHeight: ID.Block.t,
    timestamp: MomentRe.Moment.t,
  };

  module MultiMiniByDataSourceConfig = [%graphql
    {|
      subscription RequestMiniByDataSource($id: bigint!, $limit: Int!, $offset: Int!) {
        data_sources_by_pk(id: $id) {
          raw_data_requests(limit: $limit, offset: $offset) {
            request @bsRecord {
              id @bsDecoder(fn: "ID.Request.fromJson")
              requester @bsDecoder(fn: "Address.fromBech32")
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

  module MultiMiniByOracleScriptConfig = [%graphql
    {|
      subscription RequestMiniByOracleScript($id: bigint!, $limit: Int!, $offset: Int!) {
        oracle_scripts_by_pk(id: $id) {
          requests(limit: $limit, offset: $offset) {
            id @bsDecoder(fn: "ID.Request.fromJson")
            requester @bsDecoder(fn: "Address.fromBech32")
            oracle_script @bsRecord {
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
    |}
  ];

  let toExternal =
      ({id, requester, oracleScript, transaction: {txHash, blockHeight, timestamp}}) => {
    id,
    requester,
    oracleScriptID: oracleScript.id,
    oracleScriptName: oracleScript.name,
    txHash,
    blockHeight,
    timestamp,
  };

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
         | Some(dataSource) =>
           dataSource##raw_data_requests->Belt_Array.map(y => y##request |> toExternal)
         | None => [||]
         }
       );
  };

  let getListByOracleScript = (id, ~page, ~pageSize, ()) => {
    let offset = (page - 1) * pageSize;
    let (result, _) =
      ApolloHooks.useSubscription(
        MultiMiniByOracleScriptConfig.definition,
        ~variables=
          MultiMiniByOracleScriptConfig.makeVariables(
            ~id=id |> ID.OracleScript.toJson,
            ~limit=pageSize,
            ~offset,
            (),
          ),
      );
    result
    |> Sub.map(_, x =>
         switch (x##oracle_scripts_by_pk) {
         | Some(oracleScript) =>
           oracleScript##requests
           ->Belt_Array.map(y =>
               {
                 id: y##id,
                 requester: y##requester,
                 oracleScriptID: y##oracle_script.id,
                 oracleScriptName: y##oracle_script.name,
                 txHash: y##transaction.txHash,
                 blockHeight: y##transaction.blockHeight,
                 timestamp: y##transaction.timestamp,
               }
             )
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

module RequestCountByOracleScriptConfig = [%graphql
  {|
    subscription RequestMiniByDataSourceCount($id: bigint!) {
      oracle_scripts_by_pk(id: $id) {
        requests_aggregate {
          aggregate {
            count @bsDecoder(fn: "Belt_Option.getExn")
          }
        }
      }
    }
  |}
];

let countByOracleScript = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      RequestCountByOracleScriptConfig.definition,
      ~variables=
        RequestCountByOracleScriptConfig.makeVariables(~id=id |> ID.OracleScript.toJson, ()),
    );
  result
  |> Sub.map(_, x => {
       {
         let%Opt oracleScript = x##oracle_scripts_by_pk;
         let%Opt aggregate = oracleScript##requests_aggregate##aggregate;
         Some(aggregate##count);
       }
       ->Belt_Option.getWithDefault(0)
     });
};

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
