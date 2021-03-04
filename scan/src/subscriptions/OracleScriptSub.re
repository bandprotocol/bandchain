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
  requestCount: requestStatOpt->Belt.Option.mapWithDefault(0, ({count}) => count),
};

module MultiConfig = [%graphql
  {|
  subscription OracleScripts($limit: Int!, $offset: Int!, $searchTerm: String!) {
    oracle_scripts(limit: $limit, offset: $offset,where: {name: {_ilike: $searchTerm}}, order_by: {request_stat: {count: desc}, transaction: {block: {timestamp: desc}}, id: desc}) @bsRecord {
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

module MultiOracleScriptStatLast1DayConfig = [%graphql
  {|
  subscription MultiOracleScriptStatLast1DayConfig {
    oracle_script_statistic_last_1_day(where: {resolve_status: {_eq: "Success"}}) @bsRecord {
      id @bsDecoder(fn: "ID.OracleScript.fromIntExn")
      responseTime: response_time @bsDecoder(fn: "GraphQLParser.floatWithDefault")
    }
  }
  |}
];

module SingleOracleScriptStatLast1DayConfig = [%graphql
  {|
  subscription SingleOracleScriptStatLast1DayConfig($id: Int!) {
    oracle_script_statistic_last_1_day(where: {resolve_status: {_eq: "Success"}, id: {_eq: $id}}) @bsRecord {
      id @bsDecoder(fn: "ID.OracleScript.fromIntExn")
      responseTime: response_time @bsDecoder(fn: "GraphQLParser.floatWithDefault")
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
  result |> Sub.map(_, internal => internal##oracle_scripts->Belt.Array.map(toExternal));
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

let getResponseTimeList = () => {
  let (result, _) = ApolloHooks.useSubscription(MultiOracleScriptStatLast1DayConfig.definition);
  result |> Sub.map(_, internal => internal##oracle_script_statistic_last_1_day);
};

let getResponseTime = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleOracleScriptStatLast1DayConfig.definition,
      ~variables=SingleConfig.makeVariables(~id=id |> ID.OracleScript.toInt, ()),
    );
  result
  |> Sub.map(_, internal => internal##oracle_script_statistic_last_1_day |> Belt.Array.get(_, 0));
};
