module ReportersByValidatorAddressConfig = [%graphql
  {|
subscription HistoricalBondedToken($operator_address: String!, $limit: Int!, $offset: Int!) {
    reporters(where: {validator: {operator_address: {_eq: $operator_address}}}, offset: $offset, limit: $limit, order_by: {reporter_id: asc}) {
      account {
        address
      }
    }
  }
|}
];

module ReportersCountConfig = [%graphql
  {|
    subscription DelegatorCount($operator_address: String!) {
      reporters_aggregate(where: {validator: {operator_address: {_eq: $operator_address}}}) {
        aggregate {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

let getList = (~operatorAddress, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;

  let (result, _) =
    ApolloHooks.useSubscription(
      ReportersByValidatorAddressConfig.definition,
      ~variables=
        ReportersByValidatorAddressConfig.makeVariables(
          ~operator_address=operatorAddress |> Address.toOperatorBech32,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );

  result
  |> Sub.map(_, x =>
       x##reporters->Belt_Array.map(each => each##account##address |> Address.fromBech32)
     );
};

let count = operatorAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      ReportersCountConfig.definition,
      ~variables=
        ReportersCountConfig.makeVariables(
          ~operator_address=operatorAddress |> Address.toOperatorBech32,
          (),
        ),
    );
  result
  |> Sub.map(_, x => x##reporters_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};
