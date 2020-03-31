type t = {
  id: ID.OracleScript.t,
  owner: Address.t,
  name: string,
  description: string,
  code: JsBuffer.t,
  timestamp: MomentRe.Moment.t,
};

module MultiConfig = [%graphql
  {|
  subscription OracleScripts($limit: Int!, $offset: Int!) {
    oracle_scripts(limit: $limit, offset: $offset) @bsRecord {
      id @bsDecoder(fn: "ID.OracleScript.fromJson")
      owner @bsDecoder(fn: "Address.fromBech32")
      name
      description
      code @bsDecoder(fn: "GraphQLParser.buffer")
      timestamp: last_updated @bsDecoder(fn: "GraphQLParser.time")
    }
  }
|}
];

module SingleConfig = [%graphql
  {|
  subscription OracleScript($id: bigint!) {
    oracle_scripts_by_pk(id: $id) @bsRecord {
      id @bsDecoder(fn: "ID.OracleScript.fromJson")
      owner @bsDecoder(fn: "Address.fromBech32")
      name
      description
      code @bsDecoder(fn: "GraphQLParser.buffer")
      timestamp: last_updated @bsDecoder(fn: "GraphQLParser.time")
    }
  },
|}
];

module OracleScriptsCountConfig = [%graphql
  {|
  subscription OracleScriptsCount {
    oracle_scripts_aggregate{
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
      ~variables=SingleConfig.makeVariables(~id=id |> ID.OracleScript.toJson, ()),
    );
  let%Sub x = result;
  switch (x##oracle_scripts_by_pk) {
  | Some(data) => Sub.resolve(data)
  | None => NoData
  };
};

let getList = (~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ()),
    );
  result |> Sub.map(_, x => x##oracle_scripts);
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(OracleScriptsCountConfig.definition);
  result
  |> Sub.map(_, x =>
       x##oracle_scripts_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count)
     );
};
