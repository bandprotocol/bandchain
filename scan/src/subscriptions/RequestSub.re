open ValidatorSub.Mini;
open TxSub.Mini;

module Mini = {
  type oracle_script_internal_t = {
    scriptID: ID.OracleScript.t,
    name: string,
    schema: string,
  };

  type aggregate_internal_t = {count: int};

  type aggregate_wrapper_intenal_t = {aggregate: option(aggregate_internal_t)};

  type request_internal = {
    id: ID.Request.t,
    sender: option(Address.t),
    clientID: string,
    requestTime: option(MomentRe.Moment.t),
    resolveTime: option(MomentRe.Moment.t),
    calldata: JsBuffer.t,
    oracleScript: oracle_script_internal_t,
    transactionOpt: option(TxSub.Mini.t),
    reportsAggregate: aggregate_wrapper_intenal_t,
    minCount: int,
    resolveStatus: RequestStatus.t,
    requestedValidatorsAggregate: aggregate_wrapper_intenal_t,
    result: option(JsBuffer.t),
  };

  type t = {
    id: ID.Request.t,
    sender: option(Address.t),
    clientID: string,
    requestTime: option(MomentRe.Moment.t),
    resolveTime: option(MomentRe.Moment.t),
    calldata: JsBuffer.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    txHash: option(Hash.t),
    txTimestamp: option(MomentRe.Moment.t),
    blockHeight: option(ID.Block.t),
    reportsCount: int,
    minCount: int,
    askCount: int,
    resolveStatus: RequestStatus.t,
    result: option(JsBuffer.t),
  };

  module MultiMiniByDataSourceConfig = [%graphql
    {|
      subscription RequestsMiniByDataSource($id: Int!, $limit: Int!, $offset: Int!) {
        requests(
          where: {raw_requests: {data_source_id: {_eq: $id}}}
          limit: $limit
          offset: $offset
          order_by: {id: desc}
        ) @bsRecord {
          id @bsDecoder(fn: "ID.Request.fromInt")
          clientID: client_id
          requestTime: request_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
          resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
          sender @bsDecoder(fn: "GraphQLParser.addressOpt")
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
          oracleScript: oracle_script @bsRecord {
            scriptID: id @bsDecoder(fn: "ID.OracleScript.fromInt")
            name
            schema
          }
          transactionOpt: transaction @bsRecord {
            hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
            block @bsRecord {
              timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
            }
            gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
          }
          reportsAggregate: reports_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          resolveStatus: resolve_status  @bsDecoder(fn: "RequestStatus.fromJsonString")
          minCount: min_count
          requestedValidatorsAggregate: val_requests_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          result @bsDecoder(fn: "GraphQLParser.optionBuffer")
        }
      }
    |}
  ];

  module MultiMiniByOracleScriptConfig = [%graphql
    {|
      subscription RequestsMiniByOracleScript($id: Int!, $limit: Int!, $offset: Int!) {
        requests(
          where: {oracle_script_id: {_eq: $id}}
          limit: $limit
          offset: $offset
          order_by: {id: desc}
        ) @bsRecord {
          id @bsDecoder(fn: "ID.Request.fromInt")
          clientID: client_id
          requestTime: request_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
          resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
          sender @bsDecoder(fn: "GraphQLParser.addressOpt")
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
          oracleScript: oracle_script @bsRecord {
            scriptID: id @bsDecoder(fn: "ID.OracleScript.fromInt")
            name
            schema
          }
          transactionOpt: transaction @bsRecord {
            hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
            block @bsRecord {
              timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
            }
            gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
          }
          reportsAggregate: reports_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          resolveStatus: resolve_status  @bsDecoder(fn: "RequestStatus.fromJsonString")
          minCount: min_count
          requestedValidatorsAggregate: val_requests_aggregate @bsRecord {
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
        requests(where: {transaction: {hash: {_eq: $tx_hash}}}) @bsRecord {
          id @bsDecoder(fn: "ID.Request.fromInt")
          clientID: client_id
          requestTime: request_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
          resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
          sender @bsDecoder(fn: "GraphQLParser.addressOpt")
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
          oracleScript: oracle_script @bsRecord {
            scriptID: id @bsDecoder(fn: "ID.OracleScript.fromInt")
            name
            schema
          }
          transactionOpt: transaction @bsRecord {
            hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
            block @bsRecord {
              timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
            }
            gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
          }
          reportsAggregate: reports_aggregate @bsRecord {
            aggregate @bsRecord {
              count @bsDecoder(fn: "Belt_Option.getExn")
            }
          }
          resolveStatus: resolve_status  @bsDecoder(fn: "RequestStatus.fromJsonString")
          minCount: min_count
          requestedValidatorsAggregate: val_requests_aggregate @bsRecord {
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
          sender,
          clientID,
          requestTime,
          resolveTime,
          calldata,
          oracleScript,
          transactionOpt,
          reportsAggregate,
          minCount,
          resolveStatus,
          requestedValidatorsAggregate,
          result,
        },
      ) => {
    id,
    sender,
    clientID,
    requestTime,
    resolveTime,
    calldata,
    oracleScriptID: oracleScript.scriptID,
    oracleScriptName: oracleScript.name,
    txHash: transactionOpt->Belt.Option.map(({hash}) => hash),
    txTimestamp: transactionOpt->Belt.Option.map(({block}) => block.timestamp),
    blockHeight: transactionOpt->Belt.Option.map(({blockHeight}) => blockHeight),
    reportsCount:
      reportsAggregate.aggregate->Belt_Option.map(({count}) => count)->Belt.Option.getExn,
    minCount,
    askCount:
      requestedValidatorsAggregate.aggregate
      ->Belt_Option.map(({count}) => count)
      ->Belt_Option.getExn,
    resolveStatus,
    result,
  };

  let getListByDataSource = (id, ~page, ~pageSize, ()) => {
    let offset = (page - 1) * pageSize;
    let (result, _) =
      ApolloHooks.useSubscription(
        MultiMiniByDataSourceConfig.definition,
        ~variables=
          MultiMiniByDataSourceConfig.makeVariables(
            ~id=id |> ID.DataSource.toInt,
            ~limit=pageSize,
            ~offset,
            (),
          ),
      );
    result |> Sub.map(_, x => x##requests->Belt_Array.map(toExternal));
  };

  let getListByOracleScript = (id, ~page, ~pageSize, ()) => {
    let offset = (page - 1) * pageSize;
    let (result, _) =
      ApolloHooks.useSubscription(
        MultiMiniByOracleScriptConfig.definition,
        ~variables=
          MultiMiniByOracleScriptConfig.makeVariables(
            ~id=id |> ID.OracleScript.toInt,
            ~limit=pageSize,
            ~offset,
            (),
          ),
      );
    result |> Sub.map(_, x => x##requests->Belt_Array.map(toExternal));
  };

  let getListByTxHash = (txHash: Hash.t) => {
    let (result, _) =
      ApolloHooks.useSubscription(
        MultiMiniByTxHashConfig.definition,
        ~variables=
          MultiMiniByTxHashConfig.makeVariables(
            ~tx_hash=txHash |> Hash.toHex |> (x => "\\x" ++ x) |> Js.Json.string,
            (),
          ),
      );
    result |> Sub.map(_, x => x##requests->Belt_Array.map(toExternal));
  };
};

module RequestCountByDataSourceConfig = [%graphql
  {|
    subscription RequestsMiniCountByDataSource($id: Int!) {
      raw_requests_aggregate(where: {data_source_id: {_eq: $id}}) {
        aggregate {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

module RequestCountByOracleScriptConfig = [%graphql
  {|
    subscription RequestsCountMiniByOracleScript($id: Int!) {
      requests_aggregate(where: {oracle_script_id: {_eq: $id}}) {
        aggregate {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

type report_detail_t = {
  externalID: string,
  exitCode: string,
  data: JsBuffer.t,
};

type report_t = {
  transactionOpt: option(TxSub.Mini.t),
  reportDetails: array(report_detail_t),
  reportValidator: ValidatorSub.Mini.t,
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
  externalID: string,
  dataSource: data_source_internal_t,
  calldata: JsBuffer.t,
};

type requested_validator_internal_t = {validator: ValidatorSub.Mini.t};

type internal_t = {
  id: ID.Request.t,
  clientID: string,
  requestTime: option(MomentRe.Moment.t),
  resolveTime: option(MomentRe.Moment.t),
  oracleScript: oracle_script_internal_t,
  calldata: JsBuffer.t,
  requestedValidators: array(requested_validator_internal_t),
  minCount: int,
  resolveStatus: RequestStatus.t,
  sender: option(Address.t),
  transactionOpt: option(TxSub.Mini.t),
  rawDataRequests: array(raw_data_request_t),
  reports: array(report_t),
  result: option(JsBuffer.t),
};

type t = {
  id: ID.Request.t,
  clientID: string,
  requestTime: option(MomentRe.Moment.t),
  resolveTime: option(MomentRe.Moment.t),
  oracleScript: oracle_script_internal_t,
  calldata: JsBuffer.t,
  requestedValidators: array(requested_validator_internal_t),
  minCount: int,
  resolveStatus: RequestStatus.t,
  requester: option(Address.t),
  transactionOpt: option(TxSub.Mini.t),
  rawDataRequests: array(raw_data_request_t),
  reports: array(report_t),
  result: option(JsBuffer.t),
};

let toExternal =
    (
      {
        id,
        clientID,
        requestTime,
        resolveTime,
        oracleScript,
        calldata,
        requestedValidators,
        minCount,
        resolveStatus,
        sender,
        transactionOpt,
        rawDataRequests,
        reports,
        result,
      },
    ) => {
  id,
  clientID,
  requestTime,
  resolveTime,
  oracleScript,
  calldata,
  requestedValidators,
  minCount,
  resolveStatus,
  requester: sender,
  transactionOpt,
  rawDataRequests,
  reports,
  result,
};

module SingleRequestConfig = [%graphql
  {|
    subscription Request($id: Int!) {
      requests_by_pk(id: $id) @bsRecord {
        id @bsDecoder(fn: "ID.Request.fromInt")
        clientID: client_id
        requestTime: request_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
        resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
        oracleScript: oracle_script @bsRecord {
          oracleScriptID:id @bsDecoder(fn: "ID.OracleScript.fromInt")
          name
          schema
        }
        calldata @bsDecoder(fn: "GraphQLParser.buffer")
        reports(order_by: {validator_id: asc}) @bsRecord {
          transactionOpt: transaction @bsRecord {
            hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
            block @bsRecord {
              timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
            }
            gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
          }
          reportDetails: raw_reports(order_by: {external_id: asc}) @bsRecord {
            externalID: external_id @bsDecoder (fn: "GraphQLParser.string")
            exitCode: exit_code @bsDecoder (fn: "GraphQLParser.string")
            data @bsDecoder(fn: "GraphQLParser.buffer")
          }
          reportValidator: validator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
            identity
          }
        }
        requestedValidators: val_requests(order_by: {validator_id: asc}) @bsRecord {
          validator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
            identity
          }
        }
        minCount: min_count
        resolveStatus: resolve_status  @bsDecoder(fn: "RequestStatus.fromJsonString")
        sender @bsDecoder(fn: "GraphQLParser.addressOpt")
        transactionOpt: transaction @bsRecord {
          hash @bsDecoder(fn: "GraphQLParser.hash")
          blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
          block @bsRecord {
            timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
          }
          gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
        }
        rawDataRequests: raw_requests(order_by: {external_id: asc}) @bsRecord {
          externalID: external_id @bsDecoder (fn: "GraphQLParser.string")
          dataSource: data_source @bsRecord {
            dataSourceID: id @bsDecoder(fn: "ID.DataSource.fromInt")
            name
          }
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
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
        id @bsDecoder(fn: "ID.Request.fromInt")
        clientID: client_id
        requestTime: request_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
        resolveTime: resolve_time @bsDecoder(fn: "GraphQLParser.fromUnixSecondOpt")
        oracleScript: oracle_script @bsRecord {
          oracleScriptID:id @bsDecoder(fn: "ID.OracleScript.fromInt")
          name
          schema
        }
        calldata @bsDecoder(fn: "GraphQLParser.buffer")
        reports @bsRecord {
          transactionOpt: transaction @bsRecord {
            hash @bsDecoder(fn: "GraphQLParser.hash")
            blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
            block @bsRecord {
              timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
            }
            gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
          }
          reportDetails: raw_reports(order_by: {external_id: asc}) @bsRecord {
            externalID: external_id @bsDecoder (fn: "GraphQLParser.string")
            exitCode: exit_code @bsDecoder (fn: "GraphQLParser.string")
            data @bsDecoder(fn: "GraphQLParser.buffer")
          }
          reportValidator: validator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
            identity
          }
        }
        requestedValidators: val_requests @bsRecord {
          validator @bsRecord {
            consensusAddress: consensus_address
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            moniker
            identity
          }
        }
        minCount: min_count
        resolveStatus: resolve_status  @bsDecoder(fn: "RequestStatus.fromJsonString")
        sender @bsDecoder(fn: "GraphQLParser.addressOpt")
        transactionOpt: transaction @bsRecord {
          hash @bsDecoder(fn: "GraphQLParser.hash")
          blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
          block @bsRecord {
            timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
          }
          gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
        }
        rawDataRequests: raw_requests(order_by: {external_id: asc}) @bsRecord {
          externalID: external_id @bsDecoder (fn: "GraphQLParser.string")
          dataSource: data_source @bsRecord {
            dataSourceID: id @bsDecoder(fn: "ID.DataSource.fromInt")
            name
          }
          calldata @bsDecoder(fn: "GraphQLParser.buffer")
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
      ~variables=SingleRequestConfig.makeVariables(~id=id |> ID.Request.toInt, ()),
    );
  switch (result) {
  | ApolloHooks.Subscription.Data(data) =>
    switch (data##requests_by_pk) {
    | Some(x) => ApolloHooks.Subscription.Data(x |> toExternal)
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
  result |> Sub.map(_, x => x##requests->Belt.Array.map(toExternal));
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
        RequestCountByOracleScriptConfig.makeVariables(~id=id |> ID.OracleScript.toInt, ()),
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
      ~variables=RequestCountByDataSourceConfig.makeVariables(~id=id |> ID.DataSource.toInt, ()),
    );
  result
  |> Sub.map(_, x => {
       {
         let%Opt aggregate = x##raw_requests_aggregate##aggregate;
         Some(aggregate##count);
       }
       ->Belt_Option.getExn
     });
};
