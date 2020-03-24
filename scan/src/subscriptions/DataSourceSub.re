type t = {
  id: ID.DataSource.t,
  owner: Address.t,
  name: string,
  description: string,
  executable: JsBuffer.t,
  timestamp: MomentRe.Moment.t,
};

module MultiConfig = [%graphql
  {|
  subscription DataSources {
    data_sources @bsRecord {
      id @bsDecoder(fn: "ID.DataSource.fromJson")
      owner @bsDecoder(fn: "Address.fromBech32")
      name
      description
      executable @bsDecoder(fn: "GraphQLParser.buffer")
      timestamp: last_updated @bsDecoder(fn: "GraphQLParser.time")
    }
  }
|}
];

module SingleConfig = [%graphql
  {|
  subscription DataSource($id: bigint!) {
    data_sources_by_pk(id: $id) @bsRecord {
      id @bsDecoder(fn: "ID.DataSource.fromJson")
      owner @bsDecoder(fn: "Address.fromBech32")
      name
      description
      executable @bsDecoder(fn: "GraphQLParser.buffer")
      timestamp: last_updated @bsDecoder(fn: "GraphQLParser.time")
    }
  },
|}
];

let get = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=SingleConfig.makeVariables(~id=id |> ID.DataSource.toJson, ()),
    );
  result |> Sub.map(_, x => x##data_sources_by_pk);
};

let getList = () => {
  let (result, _) = ApolloHooks.useSubscription(MultiConfig.definition);
  result |> Sub.map(_, x => x##data_sources);
};
