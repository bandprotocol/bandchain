type data_source_t = {
  dataSourceID: ID.DataSource.t,
  dataSourceName: string,
};
type related_data_sources = {dataSource: data_source_t};
type block_t = {timestamp: MomentRe.Moment.t};
type transaction_t = {block: block_t};
type request_stat_t = {count: int};
type response_last_1_day_t = {
  id: ID.OracleScript.t,
  responseTime: float,
  resolveStatus: string,
};

type internal_t = {
  id: ID.OracleScript.t,
  owner: Address.t,
  name: string,
  description: string,
  schema: string,
  sourceCodeURL: string,
  transaction: option(transaction_t),
  relatedDataSources: array(related_data_sources),
  requestStat: option(request_stat_t),
  responsesLast1Day: array(response_last_1_day_t),
};

type t = {
  id: ID.OracleScript.t,
  owner: Address.t,
  name: string,
  description: string,
  schema: string,
  sourceCodeURL: string,
  timestamp: option(MomentRe.Moment.t),
  relatedDataSources: list(data_source_t),
  requestCount: int,
  responseTime: option(float),
};

let toExternal =
    (
      {
        id,
        owner,
        name,
        description,
        schema,
        sourceCodeURL,
        transaction: txOpt,
        relatedDataSources,
        requestStat: requestStatOpt,
        responsesLast1Day,
      },
    ) => {
  id,
  owner,
  name,
  description,
  schema,
  sourceCodeURL,
  timestamp: {
    let%Opt tx = txOpt;
    Some(tx.block.timestamp);
  },
  relatedDataSources:
    relatedDataSources->Belt.Array.map(({dataSource}) => dataSource)->Belt.List.fromArray,
  // Note: requestCount can't be nullable value.
  requestCount: requestStatOpt->Belt.Option.map(({count}) => count)->Belt.Option.getExn,
  responseTime:
    responsesLast1Day->Belt.Array.map(({responseTime}) => responseTime)->Belt.Array.get(0),
};
module MultiConfig = [%graphql
  {|
  subscription OracleScripts($limit: Int!, $offset: Int!, $searchTerm: String!) {
    oracle_scripts(limit: $limit, offset: $offset,where: {name: {_ilike: $searchTerm}}, order_by: {request_stat: {count: desc}, transaction: {block: {timestamp: desc}}, id: desc})  {
      id
      owner
      name
      description
      schema
      sourceCodeURL: source_code_url
      transaction  {
        block  {
          timestamp
        }
      }
      relatedDataSources: related_data_source_oracle_scripts  {
        dataSource: data_source  {
          dataSourceID: id
          dataSourceName: name
        }
      }
      requestStat: request_stat  {
        count
      }
    }
  }
|}
];

module SingleConfig = [%graphql
  {|
  subscription OracleScript($id: Int!) {
    oracle_scripts_by_pk(id: $id) @bsRecord {
      id @bsDecoder(fn: "ID.OracleScript.fromInt")
      owner @bsDecoder(fn: "Address.fromBech32")
      name
      description
      schema
      sourceCodeURL: source_code_url
      transaction @bsRecord {
        block @bsRecord {
          timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
        }
      }
      relatedDataSources: related_data_source_oracle_scripts @bsRecord {
        dataSource: data_source @bsRecord {
          dataSourceID: id  @bsDecoder(fn: "ID.DataSource.fromInt")
          dataSourceName: name
        }
      }
      requestStat: request_stat @bsRecord {
        count
      }
      responsesLast1Day: response_last_1_day @bsRecord {
        id @bsDecoder(fn: "ID.OracleScript.fromIntExn")
        responseTime: response_time @bsDecoder(fn: "GraphQLParser.floatWithDefault")
        resolveStatus: resolve_status @bsDecoder(fn: "GraphQLParser.jsonToStringExn")
      }
    }
  },
|}
];

module OracleScriptsCountConfig = [%graphql
  {|
  subscription OracleScriptsCount($searchTerm: String!) {
    oracle_scripts_aggregate(where: {name: {_ilike: $searchTerm}}){
      aggregate{
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

module OracleScriptsStatConfig = [%graphql
  {|
  subscription OracleScriptsStatConfig {
    oracle_script_statistic_last_1_day(where: {resolve_status: {_eq: "Success"}}) {
      id,
      resolveStatus: resolve_status
      responseTime: response_time
    }
  }
|}
];

let get = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=SingleConfig.makeVariables(~id=id |> ID.OracleScript.toInt, ()),
    );
  let%Sub x = result;
  switch (x##oracle_scripts_by_pk) {
  | Some(data) => Sub.resolve(data |> toExternal)
  | None => NoData
  };
};

let getList = (~page, ~pageSize, ~searchTerm, ()) => {
  let offset = (page - 1) * pageSize;
  let keyword = {j|%$searchTerm%|j};
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ~searchTerm=keyword, ()),
    );
  let (oracleScriptStat, _) = ApolloHooks.useSubscription(OracleScriptsStatConfig.definition);

  let%Sub x = oracleScriptStat;

  let oracleScriptStatList =
    x##oracle_script_statistic_last_1_day
    ->Belt.Array.map(each =>
        {
          id: each##id |> ID.OracleScript.fromIntExn,
          resolveStatus: each##resolveStatus |> GraphQLParser.jsonToStringExn,
          responseTime: each##responseTime |> GraphQLParser.floatWithDefault,
        }
      );

  let%Sub y = result;
  let oracleScriptList =
    y##oracle_scripts
    ->Belt.Array.map(each =>
        {
          id: each##id |> ID.OracleScript.fromInt,
          owner: each##owner |> Address.fromBech32,
          name: each##name,
          description: each##description,
          schema: each##schema,
          sourceCodeURL: each##sourceCodeURL,
          timestamp: {
            let%Opt tx = each##transaction;
            Some(tx##block##timestamp |> GraphQLParser.timestamp);
          },
          relatedDataSources:
            each##relatedDataSources
            ->Belt.Array.map(ds =>
                {
                  dataSourceID: ds##dataSource##dataSourceID |> ID.DataSource.fromInt,
                  dataSourceName: ds##dataSource##dataSourceName,
                }
              )
            ->Belt.List.fromArray,
          // Note: requestCount can't be nullable value.
          requestCount: each##requestStat->Belt.Option.map(c => c##count)->Belt.Option.getExn,
          responseTime: {
            oracleScriptStatList
            ->Belt.Array.keep(({id}) => id |> ID.OracleScript.toInt == each##id)
            ->Belt_Array.map(({responseTime}) => responseTime)
            ->Belt.Array.get(0);
          },
        }
      );
  Sub.resolve(oracleScriptList);
};

let count = (~searchTerm, ()) => {
  let keyword = {j|%$searchTerm%|j};
  let (result, _) =
    ApolloHooks.useSubscription(
      OracleScriptsCountConfig.definition,
      ~variables=OracleScriptsCountConfig.makeVariables(~searchTerm=keyword, ()),
    );
  result
  |> Sub.map(_, x =>
       x##oracle_scripts_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
     );
};
