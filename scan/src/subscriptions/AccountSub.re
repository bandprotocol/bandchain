type t = {balance: list(Coin.t)};

module SingleConfig = [%graphql
  {|
  subscription Account($address: String!) {
    accounts_by_pk(address: $address) @bsRecord {
      balance @bsDecoder(fn: "GraphQLParser.coins")
    }
  }
  |}
];

let get = address => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=SingleConfig.makeVariables(~address=address |> Address.toBech32, ()),
    );
  result |> Sub.map(_, x => x##accounts_by_pk |> Belt_Option.getWithDefault(_, {balance: []}));
};
