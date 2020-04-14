open ValidatorSub.Mini;
open TxSub.Mini;

module Mini = {
  open TxSub.Mini;

  type oracle_script_internal_t = {
    id: ID.OracleScript.t,
    name: string,
  };

  type aggregate_internal_t = {count: int};

  type aggregate_wrapper_intenal_t = {aggregate: option(aggregate_internal_t)};

  type request_internal = {
    id: ID.Request.t,
    requester: Address.t,
    oracleScript: oracle_script_internal_t,
    transaction: TxSub.Mini.t,
    reportsAggregate: aggregate_wrapper_intenal_t,
    sufficientValidatorCount: int,
    requestedValidatorsAgregate: aggregate_wrapper_intenal_t,
    result: option(JsBuffer.t),
  };

  type t = {
    id: ID.Request.t,
    requester: Address.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    txHash: Hash.t,
    blockHeight: ID.Block.t,
    timestamp: MomentRe.Moment.t,
    reportsCount: int,
    sufficientValidatorCount: int,
    requestedValidatorsCount: int,
    result: option(JsBuffer.t),
  };

  let optionBuffer = Belt_Option.map(_, GraphQLParser.buffer);

  module MultiMiniByDataSourceConfig = [%graphql
    {|
      subscription RequestsMiniByDataSource($id: bigint!, $limit: Int!, $offset: Int!) {
        raw_data_requests(
          where: {data_source_id: $id}
          limit: $limit
          offset: $offset
          order_by: {request_id: desc}
        ) {
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
            reportsAggregate: reports_aggregate @bsRecord {
              aggregate @bsRecord {
                count @bsDecoder(fn: "Belt_Option.getExn")
              }
            }
            sufficientValidatorCount: sufficient_validator_count @bsDecoder(fn: "GraphQLParser.int64")
            requestedValidatorsAgregate: requested_validators_aggregate @bsRecord {
              aggregate @bsRecord {
                count @bsDecoder(fn: "Belt_Option.getExn")
              }
            }
            result @bsDecoder(fn: "optionBuffer")
          }
        }
      }
    |}
  ];

  module MultiMiniByOracleScriptConfig = [%graphql
    {|
      subscription RequestsMiniByOracleScript($id: bigint!, $limit: Int!, $offset: Int!) {
        requests(
          where: {oracle_script_id: $id}
          limit: $limit
          offset: $offset
          order_by: {id: desc}
        ) {
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
          reportsAggregate: reports_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          sufficientValidatorCount: sufficient_validator_count @bsDecoder(fn: "GraphQLParser.int64")
          requestedValidatorsAgregate: requested_validators_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          result @bsDecoder(fn: "optionBuffer")
        }
      }
    |}
  ];

  module MultiMiniByTxHashConfig = [%graphql
    {|
      subscription RequestsMiniByTxHashCon($tx_hash:bytea!) {
        requests(where: {tx_hash: {_eq: $tx_hash}}) {
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
          reportsAggregate: reports_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          sufficientValidatorCount: sufficient_validator_count @bsDecoder(fn: "GraphQLParser.int64")
          requestedValidatorsAgregate: requested_validators_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          result @bsDecoder(fn: "optionBuffer")
        }
      }
    |}
  ];

  let toExternal =
      (
        {
          id,
          requester,
          oracleScript,
          transaction: {txHash, blockHeight, timestamp},
          reportsAggregate,
          sufficientValidatorCount,
          requestedValidatorsAgregate,
          result,
        },
      ) => {
    id,
    requester,
    oracleScriptID: oracleScript.id,
    oracleScriptName: oracleScript.name,
    txHash,
    blockHeight,
    timestamp,
    reportsCount:
      reportsAggregate.aggregate->Belt_Option.map(({count}) => count)->Belt_Option.getExn,
    sufficientValidatorCount,
    requestedValidatorsCount:
      requestedValidatorsAgregate.aggregate
      ->Belt_Option.map(({count}) => count)
      ->Belt_Option.getExn,
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
               oracleScriptID: y##oracle_script.id,
               oracleScriptName: y##oracle_script.name,
               txHash: y##transaction.txHash,
               blockHeight: y##transaction.blockHeight,
               timestamp: y##transaction.timestamp,
               reportsCount:
                 y##reportsAggregate.aggregate
                 ->Belt_Option.map(({count}) => count)
                 ->Belt_Option.getExn,
               sufficientValidatorCount: y##sufficientValidatorCount,
               requestedValidatorsCount:
                 y##requestedValidatorsAgregate.aggregate
                 ->Belt_Option.map(({count}) => count)
                 ->Belt_Option.getExn,
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
               oracleScriptID: y##oracleScript.id,
               oracleScriptName: y##oracleScript.name,
               txHash: y##transaction.txHash,
               blockHeight: y##transaction.blockHeight,
               timestamp: y##transaction.timestamp,
               reportsCount:
                 y##reportsAggregate.aggregate
                 ->Belt_Option.map(({count}) => count)
                 ->Belt_Option.getExn,
               sufficientValidatorCount: y##sufficientValidatorCount,
               requestedValidatorsCount:
                 y##requestedValidatorsAgregate.aggregate
                 ->Belt_Option.map(({count}) => count)
                 ->Belt_Option.getExn,
               result: y##result,
             }
           )
       );
  };
};

module RequestCountByDataSourceConfig = [%graphql
  {|
    subscription RequestsMiniCountByDataSource($id: bigint!) {
      raw_data_requests_aggregate(where: {data_source_id: $id}) {
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
      requests_aggregate(where: {oracle_script_id: $id}) {
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

type oracle_script_code_internal_t = {schema: string};

type oracle_script_internal_t = {
  oracleScriptID: ID.OracleScript.t,
  name: string,
  oracleScriptCode: oracle_script_code_internal_t,
};

type raw_data_request_t = {
  externalID: int,
  dataSource: data_source_internal_t,
  calldata: JsBuffer.t,
};

type requested_validator_internal_t = {validator: ValidatorSub.Mini.t};

type resolve_status_t =
  | Pending
  | Success
  | Failure
  | Unknown;

type t = {
  id: ID.Request.t,
  oracleScript: oracle_script_internal_t,
  calldata: JsBuffer.t,
  requestedValidators: array(requested_validator_internal_t),
  sufficientValidatorCount: int,
  expirationHeight: int,
  resolveStatus: resolve_status_t,
  requester: Address.t,
  transaction: TxSub.Mini.t,
  rawDataRequests: array(raw_data_request_t),
  reports: array(report_t),
  result: option(JsBuffer.t),
};

let optionBuffer = Belt_Option.map(_, GraphQLParser.buffer);
let parseResolveStatus = x =>
  switch (x) {
  | "Pending" => Pending
  | "Success" => Success
  | "Failure" => Failure
  | _ => Unknown
  };

module SingleRequestConfig = [%graphql
  {|
    subscription Request($id: bigint!) {
      requests_by_pk(id: $id) @bsRecord {
        id @bsDecoder(fn: "ID.Request.fromJson")
        oracleScript: oracle_script @bsRecord {
          oracleScriptID:id @bsDecoder(fn: "ID.OracleScript.fromJson")
          name
          oracleScriptCode: oracle_script_code  @bsRecord {
            schema @bsDecoder(fn: "Belt_Option.getExn")
          }
        }
        calldata @bsDecoder(fn: "GraphQLParser.buffer")
        requestedValidators: requested_validators @bsRecord {
          validator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
          }
        }
        sufficientValidatorCount: sufficient_validator_count @bsDecoder(fn: "GraphQLParser.int64")
        expirationHeight: expiration_height @bsDecoder(fn: "GraphQLParser.int64")
        resolveStatus: resolve_status  @bsDecoder(fn: "parseResolveStatus")
        requester @bsDecoder(fn: "Address.fromBech32")
        transaction @bsRecord {
          txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
          blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
          timestamp @bsDecoder(fn: "GraphQLParser.time")
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
            timestamp @bsDecoder(fn: "GraphQLParser.time")
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
        result @bsDecoder(fn: "optionBuffer")
      }
    }
  |}
];

module MultiRequestConfig = [%graphql
  {|
    subscription Requests($limit: Int!, $offset: Int!) {
      requests(limit: $limit, offset: $offset) @bsRecord {
        id @bsDecoder(fn: "ID.Request.fromJson")
        oracleScript: oracle_script @bsRecord {
          oracleScriptID:id @bsDecoder(fn: "ID.OracleScript.fromJson")
          name
          oracleScriptCode: oracle_script_code  @bsRecord {
            schema @bsDecoder(fn: "Belt_Option.getExn")
          }
        }
        calldata @bsDecoder(fn: "GraphQLParser.buffer")
        requestedValidators: requested_validators @bsRecord {
          validator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
          }
        }
        sufficientValidatorCount: sufficient_validator_count @bsDecoder(fn: "GraphQLParser.int64")
        expirationHeight: expiration_height @bsDecoder(fn: "GraphQLParser.int64")
        resolveStatus: resolve_status  @bsDecoder(fn: "parseResolveStatus")
        requester @bsDecoder(fn: "Address.fromBech32")
        transaction @bsRecord {
          txHash: tx_hash @bsDecoder(fn: "GraphQLParser.hash")
          blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
          timestamp @bsDecoder(fn: "GraphQLParser.time")
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
            timestamp @bsDecoder(fn: "GraphQLParser.time")
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
        result @bsDecoder(fn: "optionBuffer")
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
