type t = {
  id: ID.OracleScript.t,
  owner: Address.t,
  name: string,
  description: string,
  schema: string,
  sourceCodeURL: string,
  timestamp: MomentRe.Moment.t,
  relatedDataSources: list(ID.DataSource.t),
  request: int,
  responseTime: int,
};

type related_data_source_t = {dataSourceID: ID.DataSource.t};
type block_t = {timestamp: MomentRe.Moment.t};
type transaction_t = {block: block_t};

type internal_t = {
  id: ID.OracleScript.t,
  owner: Address.t,
  name: string,
  description: string,
  schema: string,
  sourceCodeURL: string,
  transaction: option(transaction_t),
  // related_data_sources: array(related_data_source_t),
};

let toExternal = ({id, owner, name, description, schema, sourceCodeURL, transaction}) => {
  id,
  owner,
  name,
  description,
  schema,
  sourceCodeURL,
  timestamp:
    switch (transaction) {
    | Some({block}) => block.timestamp
    // TODO: Please revisit again.
    | _ => MomentRe.momentNow()
    },
  relatedDataSources: [],
  // TODO: These will be removed after the data adding to schema
  request: Js.Math.random_int(300, 200000),
  responseTime: Js.Math.random_int(100, 999),
  //   related_data_sources->Belt.Array.map(x => x.dataSourceID)->Belt.List.fromArray,
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
  let ID.OracleScript.ID(id_) = id;
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=SingleConfig.makeVariables(~id=id_, ()),
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
