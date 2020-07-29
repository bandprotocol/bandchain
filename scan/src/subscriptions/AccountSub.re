type t = {
  balance: list(Coin.t),
  commission: list(Coin.t),
};

type validator_t = {commission: list(Coin.t)};

type internal_t = {
  balance: list(Coin.t),
  validator: option(validator_t),
};

let toExternal = ({balance, validator}) => {
  {
    balance,
    commission:
      switch (validator) {
      | Some(validator') => validator'.commission
      | None => []
      },
  };
};

module SingleConfig = [%graphql
  {|
  subscription Account($address: String!) {
    accounts_by_pk(address: $address) @bsRecord {
      balance @bsDecoder(fn: "GraphQLParser.coins")
      validator @bsRecord {
        commission: accumulated_commission @bsDecoder(fn: "GraphQLParser.coins")
      }
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
  result
  |> Sub.map(_, x =>
       x##accounts_by_pk
       |> Belt_Option.mapWithDefault(_, {balance: [], commission: []}, toExternal)
     );
};
