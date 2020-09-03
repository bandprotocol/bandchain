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

type answer_vote_t = {
  id: int,
  valPower: int,
  // valVote: option(vote_t),
  // delVotes: vote_t => int,
};

type total_vote_t = {
  validatorID: option(int),
  answer: vote_t,
  power: Coin.t,
};

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

module ValidatorVoteByProposalIDConfig = [%graphql
  {|
    subscription ValidatorVoteByProposalID($proposal_id: Int!) {
      validator_vote_proposals_view(where: {proposal_id: {_eq: $proposal_id}}) {
        amount
        answer
        id
        proposal_id
      }
    }
  |}
];

module DeligatorVoteByProposalIDConfig = [%graphql
  {|
    subscription DeligatorVoteByProposalID($proposal_id: Int!) {
      non_validator_vote_proposals_view(where: {proposal_id: {_eq: $proposal_id}}) {
        amount
        answer
        validator_id
        proposal_id
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

let getValidatorVoteByProposalID = proposalID => {
  let (validatorVotes, _) =
    ApolloHooks.useSubscription(
      ValidatorVoteByProposalIDConfig.definition,
      ~variables=
        ValidatorVoteByProposalIDConfig.makeVariables(
          ~proposal_id=proposalID |> ID.Proposal.toInt,
          (),
        ),
    );
  let (delegatorVotes, _) =
    ApolloHooks.useSubscription(
      DeligatorVoteByProposalIDConfig.definition,
      ~variables=
        DeligatorVoteByProposalIDConfig.makeVariables(
          ~proposal_id=proposalID |> ID.Proposal.toInt,
          (),
        ),
    );

  let%Sub x = validatorVotes;
  let%Sub y = delegatorVotes;

  let valVotes =
    x##validator_vote_proposals_view
    ->Belt.Array.map(each =>
        {
          validatorID: each##id,
          answer:
            switch (each##answer |> ) {
            | "Yes" => Yes
            | "No" => No
            | "NoWithVeto" => NoWithVeto
            | "Abstain" => Abstain
            },
          power: each##amount |> GraphQLParser.coinExn,
        }
      );

  let m =
    valVotes##validator_vote_proposals_view
    ->Belt_Array.reduce(
        Belt_MapInt.empty,
        (acc, x) => {
          let (id, power, choice) = x;
          acc->Belt_MapInt.set(
            id,
            {id, valPower: power, valVote: Some(choice), delVotes: _ => 0},
          );
        },
      );
  let n =
    delVotes##non_validator_vote_proposals_view
    ->Belt_Array.reduce(
        m,
        (acc, x) => {
          let (id, power, choice) = x;
          acc->Belt_MapInt.update(
            id,
            v => {
              let entry =
                v->Belt_Option.getWithDefault({id, valPower: 0, valVote: None, delVotes: _ => 0});
              let delVotes = ch => ch == choice ? power : entry.delVotes(ch);
              // Js.Console.log3(id, delVotes, power);
              Some({...entry, delVotes: ch => ch == choice ? power : entry.delVotes(ch)});
            },
          );
        },
      )
    ->Belt_MapInt.valuesToArray;

  Sub.resolve({id: 1, valPower: 20});
};
