type t = {
  voter: Address.t,
  txHash: Hash.t,
  timestamp: MomentRe.Moment.t,
};

type block_t = {timestamp: MomentRe.Moment.t};
type account_t = {address: Address.t};
type transaction_t = {
  hash: Hash.t,
  block: block_t,
};

type internal_t = {
  account: account_t,
  transaction: transaction_t,
};

let toExternal = ({account, transaction: {hash, block}}) => {
  voter: account.address,
  txHash: hash,
  timestamp: block.timestamp,
};

type vote_t =
  | Yes
  | No
  | NoWithVeto
  | Abstain;

let toString = (~withSpace=false) =>
  fun
  | Yes => "Yes"
  | No => "No"
  | NoWithVeto => withSpace ? "No With Veto" : "NoWithVeto"
  | Abstain => "Abstain";

module MultiConfig = [%graphql
  {|
    subscription Votes($limit: Int!, $offset: Int!, $proposal_id: Int!, $answer: voteoption!) {
      votes(limit: $limit, offset: $offset, where: {proposal_id: {_eq: $proposal_id}, answer: {_eq: $answer}}, order_by: {transaction: {block_height: desc}}) @bsRecord {
        account @bsRecord {
          address @bsDecoder(fn:"Address.fromBech32")
        }
        transaction @bsRecord {
          hash @bsDecoder(fn: "GraphQLParser.hash")
          block @bsRecord {
            timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
          }
        }
      }
    }
|}
];

module VoteCountConfig = [%graphql
  {|
    subscription DepositCount($proposal_id: Int!, $answer: voteoption!) {
      votes_aggregate(where: {proposal_id: {_eq: $proposal_id}, answer: {_eq: $answer}}) {
        aggregate {
          count
        }
      }
    }
  |}
];

let getList = (proposalID, answer, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=
        MultiConfig.makeVariables(
          ~proposal_id=proposalID |> ID.Proposal.toInt,
          ~answer=answer |> toString |> Js.Json.string,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Sub.map(_, x => x##votes->Belt.Array.map(toExternal));
};

let count = (proposalID, answer) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      VoteCountConfig.definition,
      ~variables=
        VoteCountConfig.makeVariables(
          ~proposal_id=proposalID |> ID.Proposal.toInt,
          ~answer=answer |> toString |> Js.Json.string,
          (),
        ),
    );

  result
  |> Sub.map(_, x =>
       x##votes_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count) |> Belt.Option.getExn
     );
};
