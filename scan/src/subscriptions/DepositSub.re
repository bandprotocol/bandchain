type t = {
  depositor: Address.t,
  amount: list(Coin.t),
  txHashOpt: option(Hash.t),
};

type account_t = {address: Address.t};
type transaction_t = {hash: Hash.t};

type internal_t = {
  account: account_t,
  amount: list(Coin.t),
  transactionOpt: option(transaction_t),
};

let toExternal = ({account, amount, transactionOpt}) => {
  depositor: account.address,
  amount,
  txHashOpt: transactionOpt->Belt.Option.map(({hash}) => hash),
};

module MultiConfig = [%graphql
  {|
    subscription Deposits($limit: Int!, $offset: Int!, $proposal_id: Int!) {
      deposits(limit: $limit, offset: $offset,  where: {proposal_id: {_eq: $proposal_id}}, order_by: {depositor_id: asc}) @bsRecord {
        account @bsRecord {
          address @bsDecoder(fn:"Address.fromBech32")
        }
        amount @bsDecoder(fn: "GraphQLParser.coins")
        transactionOpt: transaction @bsRecord {
          hash @bsDecoder(fn: "GraphQLParser.hash")
        }
      }
    }
|}
];

module DepositCountConfig = [%graphql
  {|
    subscription DepositCount($proposal_id: Int!) {
      deposits_aggregate(where: {proposal_id: {_eq: $proposal_id}}) {
        aggregate {
          count
        }
      }
    }
  |}
];

let getList = (proposalID, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=
        MultiConfig.makeVariables(
          ~proposal_id=proposalID |> ID.Proposal.toInt,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Sub.map(_, x => x##deposits->Belt.Array.map(toExternal));
};

let count = proposalID => {
  let (result, _) =
    ApolloHooks.useSubscription(
      DepositCountConfig.definition,
      ~variables=
        DepositCountConfig.makeVariables(~proposal_id=proposalID |> ID.Proposal.toInt, ()),
    );
  result
  |> Sub.map(_, x =>
       x##deposits_aggregate##aggregate
       |> Belt_Option.getExn
       |> (y => y##count)
       |> Belt.Option.getExn
     );
};
