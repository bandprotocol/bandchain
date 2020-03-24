type revision_t = {
  name: string,
  timestamp: MomentRe.Moment.t,
  height: int,
  txHash: Hash.t,
};

type t = {
  id: int,
  owner: Address.t,
  name: string,
  executable: JsBuffer.t,
};

module Config = [%graphql
  {|
  subscription DataSources {
    data_sources @bsRecord {
      id @bsDecoder(fn: "GraphQLParser.int64")
      owner @bsDecoder(fn: "Address.fromBech32")
      name
      executable @bsDecoder(fn: "GraphQLParser.bytea")
    }
  }
|}
];

let use = () => {
  let (result, _) = ApolloHooks.useSubscription(Config.definition);
  result |> Sub.map(_, x => x##data_sources);
};
