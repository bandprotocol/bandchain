open ValidatorSub.Mini;
open TxSub.Mini;

type resolve_status_t =
  | Pending
  | Success
  | Failure
  | Unknown;

let parseResolveStatus = x =>
  switch (x) {
  | "Pending" => Pending
  | "Success" => Success
  | "Failure" => Failure
  | _ => Unknown
  };

module Mini = {
  open TxSub.Mini;

  type oracle_script_internal_t = {
    id: ID.OracleScript.t,
    name: string,
    schema: string,
  };

  type aggregate_internal_t = {count: int};

  type aggregate_wrapper_intenal_t = {aggregate: option(aggregate_internal_t)};

  type request_internal = {
    id: ID.Request.t,
    requester: Address.t,
    clientID: string,
    requestTime: option(MomentRe.Moment.t),
    resolveTime: option(MomentRe.Moment.t),
    calldata: JsBuffer.t,
    oracleScript: oracle_script_internal_t,
    transaction: TxSub.Mini.t,
    reportsAggregate: aggregate_wrapper_intenal_t,
    minCount: int,
    resolveStatus: resolve_status_t,
    requestedValidatorsAggregate: aggregate_wrapper_intenal_t,
    result: option(JsBuffer.t),
  };

  type t = {
    id: ID.Request.t,
    requester: Address.t,
    clientID: string,
    requestTime: option(MomentRe.Moment.t),
    resolveTime: option(MomentRe.Moment.t),
    calldata: JsBuffer.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    txHash: Hash.t,
    txTimestamp: MomentRe.Moment.t,
    blockHeight: ID.Block.t,
    reportsCount: int,
    minCount: int,
    askCount: int,
    resolveStatus: resolve_status_t,
    result: option(JsBuffer.t),
  };

  module MultiMiniByDataSourceConfig = [%graphql
    {|
      subscription RequestsMiniByDataSource($id: bigint!, $limit: Int!, $offset: Int!) {
        raw_data_requests(
          where: {data_source_id: {_eq: $id}}
          limit: $limit
          offset: $offset
          order_by: {request_id: desc}
        ) {
          request @bsRecord {
            id @bsDecoder(fn: "ID.Request.fromJson")
            clientID: client_id
            requestTime: request_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
            resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
            requester @bsDecoder(fn: "Address.fromBech32")
            calldata @bsDecoder(fn: "GraphQLParser.buffer")
            oracleScript: oracle_script @bsRecord {
              id @bsDecoder(fn: "ID.OracleScript.fromJson")
              name
              schema
            }
            transaction @bsRecord {
              txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
              blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
              timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
            }
            reportsAggregate: reports_aggregate @bsRecord {
              aggregate @bsRecord {
                count @bsDecoder(fn: "Belt_Option.getExn")
              }
            }
            resolveStatus: resolve_status  @bsDecoder(fn: "parseResolveStatus")
            minCount: min_count @bsDecoder(fn: "GraphQLParser.int64")
            requestedValidatorsAggregate: requested_validators_aggregate @bsRecord {
              aggregate @bsRecord {
                count @bsDecoder(fn: "Belt_Option.getExn")
              }
            }
            result @bsDecoder(fn: "GraphQLParser.optionBuffer")
          }
        }
      }
    |}
  ];

  module MultiMiniByOracleScriptConfig = [%graphql
    {|
      subscription RequestsMiniByOracleScript($id: bigint!, $limit: Int!, $offset: Int!) {
        requests(
          where: {oracle_script_id: {_eq: $id}}
          limit: $limit
          offset: $offset
          order_by: {id: desc}
        ) {
          id @bsDecoder(fn: "ID.Request.fromJson")
          clientID: client_id
          requestTime: request_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
          resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
          requester @bsDecoder(fn: "Address.fromBech32")
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
          oracleScript: oracle_script @bsRecord {
            id @bsDecoder(fn: "ID.OracleScript.fromJson")
            name
            schema
          }
          transaction @bsRecord {
            txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
            timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
          }
          reportsAggregate: reports_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          resolveStatus: resolve_status  @bsDecoder(fn: "parseResolveStatus")
          minCount: min_count @bsDecoder(fn: "GraphQLParser.int64")
          requestedValidatorsAggregate: requested_validators_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          result @bsDecoder(fn: "GraphQLParser.optionBuffer")
        }
      }
    |}
  ];

  module MultiMiniByTxHashConfig = [%graphql
    {|
      subscription RequestsMiniByTxHashCon($tx_hash:bytea!) {
        requests(where: {tx_hash: {_eq: $tx_hash}}) {
          id @bsDecoder(fn: "ID.Request.fromJson")
          clientID: client_id
          requestTime: request_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
          resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
          requester @bsDecoder(fn: "Address.fromBech32")
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
          oracle_script @bsRecord {
            id @bsDecoder(fn: "ID.OracleScript.fromJson")
            name
            schema
          }
          transaction @bsRecord {
            txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
            timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
          }
          reportsAggregate: reports_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          resolveStatus: resolve_status  @bsDecoder(fn: "parseResolveStatus")
          minCount: min_count @bsDecoder(fn: "GraphQLParser.int64")
          requestedValidatorsAggregate: requested_validators_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          result @bsDecoder(fn: "GraphQLParser.optionBuffer")
        }
      }
    |}
  ];

  let toExternal =
      (
        {
          id,
          requester,
          clientID,
          requestTime,
          resolveTime,
          calldata,
          oracleScript,
          transaction: {txHash, blockHeight, timestamp: txTimestamp},
          reportsAggregate,
          minCount,
          resolveStatus,
          requestedValidatorsAggregate,
          result,
        },
      ) => {
    id,
    requester,
    clientID,
    requestTime,
    resolveTime,
    calldata,
    oracleScriptID: oracleScript.id,
    oracleScriptName: oracleScript.name,
    txHash,
    txTimestamp,
    blockHeight,
    reportsCount:
      reportsAggregate.aggregate->Belt_Option.map(({count}) => count)->Belt_Option.getExn,
    minCount,
    askCount:
      requestedValidatorsAggregate.aggregate
      ->Belt_Option.map(({count}) => count)
      ->Belt_Option.getExn,
    resolveStatus,
    result,
  };

  let getListByTxHash = (txHash: Hash.t) => {
    let (result, _) =
      ApolloHooks.useSubscription(
        MultiMiniByTxHashConfig.definition,
        ~variables=
          MultiMiniByTxHashConfig.makeVariables(
            ~tx_hash=txHash |> Hash.toHex |> (x => "\x" ++ x) |> Js.Json.string,
            (),
          ),
      );
    result
    |> Sub.map(_, x =>
         x##requests
         ->Belt_Array.map(y =>
             {
               id: y##id,
               requester: y##requester,
               clientID: y##clientID,
               requestTime: y##requestTime,
               resolveTime: y##resolveTime,
               calldata: y##calldata,
               oracleScriptID: y##oracle_script.id,
               oracleScriptName: y##oracle_script.name,
               txHash: y##transaction.txHash,
               txTimestamp: y##transaction.timestamp,
               blockHeight: y##transaction.blockHeight,
               reportsCount:
                 y##reportsAggregate.aggregate
                 ->Belt_Option.map(({count}) => count)
                 ->Belt_Option.getExn,
               minCount: y##minCount,
               askCount:
                 y##requestedValidatorsAggregate.aggregate
                 ->Belt_Option.map(({count}) => count)
                 ->Belt_Option.getExn,
               resolveStatus: y##resolveStatus,
               result: y##result,
             }
           )
       );
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
    result |> Sub.map(_, x => x##raw_data_requests->Belt_Array.map(y => y##request |> toExternal));
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
         x##requests
         ->Belt_Array.map(y =>
             {
               id: y##id,
               requester: y##requester,
               clientID: y##clientID,
               requestTime: y##requestTime,
               resolveTime: y##resolveTime,
               calldata: y##calldata,
               oracleScriptID: y##oracleScript.id,
               oracleScriptName: y##oracleScript.name,
               txHash: y##transaction.txHash,
               txTimestamp: y##transaction.timestamp,
               blockHeight: y##transaction.blockHeight,
               reportsCount:
                 y##reportsAggregate.aggregate
                 ->Belt_Option.map(({count}) => count)
                 ->Belt_Option.getExn,
               minCount: y##minCount,
               askCount:
                 y##requestedValidatorsAggregate.aggregate
                 ->Belt_Option.map(({count}) => count)
                 ->Belt_Option.getExn,
               resolveStatus: y##resolveStatus,
               result: y##result,
             }
           )
       );
  };
};

module RequestCountByDataSourceConfig = [%graphql
  {|
    subscription RequestsMiniCountByDataSource($id: bigint!) {
      raw_data_requests_aggregate(where: {data_source_id: {_eq: $id}}) {
        aggregate {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

module RequestCountByOracleScriptConfig = [%graphql
  {|
    subscription RequestsCountMiniByOracleScript($id: bigint!) {
      requests_aggregate(where: {oracle_script_id: {_eq: $id}}) {
        aggregate {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

type report_detail_t = {
  externalID: int,
  data: JsBuffer.t,
};

type report_t = {
  reporter: Address.t,
  transaction: TxSub.Mini.t,
  reportDetails: array(report_detail_t),
  validatorByValidator: ValidatorSub.Mini.t,
};

type data_source_internal_t = {
  dataSourceID: ID.DataSource.t,
  name: string,
};

type oracle_script_internal_t = {
  oracleScriptID: ID.OracleScript.t,
  name: string,
  schema: string,
};

type raw_data_request_t = {
  externalID: int,
  dataSource: data_source_internal_t,
  calldata: JsBuffer.t,
};

type requested_validator_internal_t = {validator: ValidatorSub.Mini.t};

type t = {
  id: ID.Request.t,
  clientID: string,
  requestTime: option(MomentRe.Moment.t),
  resolveTime: option(MomentRe.Moment.t),
  oracleScript: oracle_script_internal_t,
  calldata: JsBuffer.t,
  requestedValidators: array(requested_validator_internal_t),
  minCount: int,
  expirationHeight: int,
  resolveStatus: resolve_status_t,
  requester: Address.t,
  transaction: TxSub.Mini.t,
  rawDataRequests: array(raw_data_request_t),
  reports: array(report_t),
  result: option(JsBuffer.t),
};

module SingleRequestConfig = [%graphql
  {|
    subscription Request($id: bigint!) {
      requests_by_pk(id: $id) @bsRecord {
        id @bsDecoder(fn: "ID.Request.fromJson")
        clientID: client_id
        requestTime: request_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
        resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
        oracleScript: oracle_script @bsRecord {
          oracleScriptID:id @bsDecoder(fn: "ID.OracleScript.fromJson")
          name
          schema
        }
        calldata @bsDecoder(fn: "GraphQLParser.buffer")
        requestedValidators: requested_validators @bsRecord {
          validator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
          }
        }
        minCount: min_count @bsDecoder(fn: "GraphQLParser.int64")
        expirationHeight: expiration_height @bsDecoder(fn: "GraphQLParser.int64")
        resolveStatus: resolve_status  @bsDecoder(fn: "parseResolveStatus")
        requester @bsDecoder(fn: "Address.fromBech32")
        transaction @bsRecord {
          txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
          blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
          timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
        }
        rawDataRequests: raw_data_requests @bsRecord {
          externalID: external_id @bsDecoder(fn: "GraphQLParser.int64")
          dataSource: data_source @bsRecord {
            dataSourceID: id @bsDecoder(fn: "ID.DataSource.fromJson")
            name
          }
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
        }
        reports @bsRecord {
          reporter @bsDecoder(fn: "Address.fromBech32")
          transaction @bsRecord {
            txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
            timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
          }
          reportDetails: report_details @bsRecord {
            externalID: external_id @bsDecoder(fn: "GraphQLParser.int64")
            data @bsDecoder(fn: "GraphQLParser.buffer")
          }
          validatorByValidator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
          }
        }
        result @bsDecoder(fn: "GraphQLParser.optionBuffer")
      }
    }
  |}
];

module MultiRequestConfig = [%graphql
  {|
    subscription Requests($limit: Int!, $offset: Int!) {
      requests(limit: $limit, offset: $offset, order_by: {id: desc}) @bsRecord {
        id @bsDecoder(fn: "ID.Request.fromJson")
        clientID: client_id
        requestTime: request_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
        resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.optionTimeS")
        oracleScript: oracle_script @bsRecord {
          oracleScriptID:id @bsDecoder(fn: "ID.OracleScript.fromJson")
          name
          schema
        }
        calldata @bsDecoder(fn: "GraphQLParser.buffer")
        requestedValidators: requested_validators @bsRecord {
          validator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
          }
        }
        minCount: min_count @bsDecoder(fn: "GraphQLParser.int64")
        expirationHeight: expiration_height @bsDecoder(fn: "GraphQLParser.int64")
        resolveStatus: resolve_status  @bsDecoder(fn: "parseResolveStatus")
        requester @bsDecoder(fn: "Address.fromBech32")
        transaction @bsRecord {
          txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
          blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
          timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
        }
        rawDataRequests: raw_data_requests @bsRecord {
          externalID: external_id @bsDecoder(fn: "GraphQLParser.int64")
          dataSource: data_source @bsRecord {
            dataSourceID: id @bsDecoder(fn: "ID.DataSource.fromJson")
            name
          }
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
        }
        reports @bsRecord {
          reporter @bsDecoder(fn: "Address.fromBech32")
          transaction @bsRecord {
            txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
            timestamp @bsDecoder(fn: "GraphQLParser.timeMS")
          }
          reportDetails: report_details @bsRecord {
            externalID: external_id @bsDecoder(fn: "GraphQLParser.int64")
            data @bsDecoder(fn: "GraphQLParser.buffer")
          }
          validatorByValidator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
          }
        }
        result @bsDecoder(fn: "GraphQLParser.optionBuffer")
      }
    }
  |}
];

module RequestCountConfig = [%graphql
  {|
  subscription RequestCount {
    requests_aggregate {
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
      SingleRequestConfig.definition,
      ~variables=SingleRequestConfig.makeVariables(~id=id |> ID.Request.toJson, ()),
    );
  switch (result) {
  | ApolloHooks.Subscription.Data(data) =>
    switch (data##requests_by_pk) {
    | Some(x) => ApolloHooks.Subscription.Data(x)
    | None => NoData
    }
  | Loading => Loading
  | Error(e) => Error(e)
  | NoData => NoData
  };
};

let getList = (~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiRequestConfig.definition,
      ~variables=MultiRequestConfig.makeVariables(~limit=pageSize, ~offset, ()),
    );
  result |> Sub.map(_, x => x##requests);
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(RequestCountConfig.definition);
  result
  |> Sub.map(_, x => x##requests_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};

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
         let%Opt aggregate = x##requests_aggregate##aggregate;
         Some(aggregate##count);
       }
       ->Belt_Option.getExn
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
         let%Opt aggregate = x##raw_data_requests_aggregate##aggregate;
         Some(aggregate##count);
       }
       ->Belt_Option.getExn
     });
};
