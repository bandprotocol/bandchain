open ValidatorSub.Mini;
open TxSub.Mini;

type resolve_status_t =
  | Pending
  | Success
  | Failure
  | Expired
  | Unknown;

let parseResolveStatus = json => {
  let status = json |> Js.Json.decodeString |> Belt_Option.getExn;
  switch (status) {
  | "Open" => Pending
  | "Success" => Success
  | "Failure" => Failure
  | "Expired" => Expired
  | _ => Unknown
  };
};

type oracle_script_internal_t = {
  scriptID: ID.OracleScript.t,
  name: string,
  schema: string,
};

type aggregate_internal_t = {count: int};

type aggregate_wrapper_intenal_t = {aggregate: option(aggregate_internal_t)};

type request_internal = {
  id: ID.Request.t,
  sender: Address.t,
  clientID: string,
  requestTime: option(MomentRe.Moment.t),
  resolveTime: option(MomentRe.Moment.t),
  calldata: JsBuffer.t,
  oracleScript: oracle_script_internal_t,
  transactionOpt: option(TxSub.Mini.t),
  reportsAggregate: aggregate_wrapper_intenal_t,
  minCount: int,
  resolveStatus: resolve_status_t,
  requestedValidatorsAggregate: aggregate_wrapper_intenal_t,
  result: option(JsBuffer.t),
};

type t = {
  id: ID.Request.t,
  sender: Address.t,
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
      query RequestsMiniByDataSource($id: Int!, $limit: Int!, $offset: Int!) {
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
          sender @bsDecoder(fn: "GraphQLParser.addressExn")
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
          resolveStatus: resolve_status  @bsDecoder(fn: "parseResolveStatus")
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
  txHash: transactionOpt->Belt.Option.map(({hash}) => hash)->Belt.Option.getExn,
  txTimestamp: transactionOpt->Belt.Option.map(({block}) => block.timestamp)->Belt.Option.getExn,
  blockHeight:
    transactionOpt->Belt.Option.map(({blockHeight}) => blockHeight)->Belt.Option.getExn,
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

let getListByDataSource = (id, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useQuery(
      MultiMiniByDataSourceConfig.definition,
      ~pollInterval=500,
      ~variables=
        MultiMiniByDataSourceConfig.makeVariables(
          ~id=id |> ID.DataSource.toInt,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Query.map(_, x => x##requests->Belt_Array.map(toExternal));
};
