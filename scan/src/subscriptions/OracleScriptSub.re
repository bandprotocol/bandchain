type data_source_t = {
  dataSourceID: ID.DataSource.t,
  dataSourceName: string,
};
type related_data_sources = {dataSource: data_source_t};
type block_t = {timestamp: MomentRe.Moment.t};
type transaction_t = {block: block_t};
type request_stat_t = {count: int};
type response_last_1_day_t = {
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

let toExternal = ({id, owner, name, description, schema, sourceCodeURL, transaction: txOpt}) => {
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
  relatedDataSources: [],
  // Note: requestCount can't be nullable value.
  requestCount: 0,
  responseTime: None,
};

module MultiConfig = [%graphql
  {|
  subscription OracleScripts($limit: Int!, $offset: Int!, $searchTerm: String!) {
    oracle_scripts(limit: $limit, offset: $offset,where: {name: {_ilike: $searchTerm}}, order_by: {transaction: {block: {timestamp: desc}}, id: desc}) @bsRecord {
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
