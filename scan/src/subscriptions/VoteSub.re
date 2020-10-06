type block_t = {timestamp: MomentRe.Moment.t};
type validator_t = {
  moniker: string,
  operatorAddress: Address.t,
  identity: string,
};
type account_t = {
  address: Address.t,
  validator: option(validator_t),
};
type transaction_t = {
  hash: Hash.t,
  block: block_t,
};

type internal_t = {
  account: account_t,
  transaction: transaction_t,
};

type t = {
  voter: Address.t,
  txHash: Hash.t,
  timestamp: MomentRe.Moment.t,
  validator: option(validator_t),
};

let toExternal = ({account: {address, validator}, transaction: {hash, block}}) => {
  voter: address,
  txHash: hash,
  timestamp: block.timestamp,
  validator,
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
  validatorID: int,
  valPower: float,
  valVote: option(vote_t),
  delVotes: vote_t => float,
  proposalID: ID.Proposal.t,
};

type internal_vote_t = {
  validatorID: int,
  answer: vote_t,
  power: Coin.t,
  proposalID: ID.Proposal.t,
};

type result_val_t = {
  validatorID: int,
  validatorPower: float,
  validatorAns: option(vote_t),
  proposalID: ID.Proposal.t,
};

type vote_stat_t = {
  proposalID: ID.Proposal.t,
  totalYes: float,
  totalYesPercent: float,
  totalNo: float,
  totalNoPercent: float,
  totalNoWithVeto: float,
  totalNoWithVetoPercent: float,
  totalAbstain: float,
  totalAbstainPercent: float,
  total: float,
};

let getAnswer = json => {
  exception NoChoice(string);
  let answer = json |> GraphQLParser.jsonToStringExn;
  switch (answer) {
  | "Yes" => Yes
  | "No" => No
  | "NoWithVeto" => NoWithVeto
  | "Abstain" => Abstain
  | _ => raise(NoChoice("There is no choice"))
  };
};

let getValVote = (valVotes, answer) =>
  valVotes->Belt_Array.reduce(0., (a, {validatorPower, validatorAns}) =>
    switch (validatorAns) {
    | Some(vote) => vote == answer ? a +. validatorPower : a
    | None => a
    }
  );

let getDelVote = (delVotes, answer_) =>
  delVotes->Belt_Array.reduce(0., (a, {power, answer}) =>
    answer == answer_ ? a +. (power |> Coin.getBandAmountFromCoin) : a
  );

module MultiConfig = [%graphql
  {|
    subscription Votes($limit: Int!, $offset: Int!, $proposal_id: Int!, $answer: voteoption!) {
      votes(limit: $limit, offset: $offset, where: {proposal_id: {_eq: $proposal_id}, answer: {_eq: $answer}}, order_by: {transaction: {block_height: desc}}) @bsRecord {
        account @bsRecord {
          address @bsDecoder(fn:"Address.fromBech32")
          validator @bsRecord {
            moniker
            operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
            identity
          }
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
      validator_vote_proposals_view(where: {proposal_id: {_eq: $proposal_id}}) @bsRecord {
        validatorID: id @bsDecoder(fn: "Belt_Option.getExn")
        answer @bsDecoder(fn: "getAnswer")
        power: tokens @bsDecoder(fn: "GraphQLParser.coinExn")
        proposalID: proposal_id @bsDecoder(fn: "ID.Proposal.fromIntExn")
      }
    }
  |}
];

module DelegatorVoteByProposalIDConfig = [%graphql
  {|
    subscription DelegatorVoteByProposalID($proposal_id: Int!) {
      non_validator_vote_proposals_view(where: {proposal_id: {_eq: $proposal_id}}) @bsRecord {
        validatorID: validator_id @bsDecoder(fn: "Belt_Option.getExn")
        answer @bsDecoder(fn: "getAnswer")
        power: tokens @bsDecoder(fn: "GraphQLParser.coinExn")
        proposalID: proposal_id @bsDecoder(fn: "ID.Proposal.fromIntExn")
      }
    }
  |}
];

module ValidatorVotesConfig = [%graphql
  {|
    subscription ValidatorVoteByProposalID {
      validator_vote_proposals_view @bsRecord {
        validatorID: id @bsDecoder(fn: "Belt_Option.getExn")
        answer @bsDecoder(fn: "getAnswer")
        power: tokens @bsDecoder(fn: "GraphQLParser.coinExn")
        proposalID: proposal_id @bsDecoder(fn: "ID.Proposal.fromIntExn")
      }
    }
  |}
];

module DelegatorVotesConfig = [%graphql
  {|
    subscription DelegatorVoteByProposalID {
      non_validator_vote_proposals_view @bsRecord {
        validatorID: validator_id @bsDecoder(fn: "Belt_Option.getExn")
        answer @bsDecoder(fn: "getAnswer")
        power: tokens @bsDecoder(fn: "GraphQLParser.coinExn")
        proposalID: proposal_id @bsDecoder(fn: "ID.Proposal.fromIntExn")
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

let parse = (valVotes, delVotes) => {
  let valVotesMap =
    valVotes->Belt_Array.reduce(
      Belt_MapString.empty, (acc, {validatorID, answer, power, proposalID}) => {
      acc->Belt_MapString.set(
        (validatorID |> string_of_int) ++ (proposalID |> ID.Proposal.toString),
        {
          validatorID,
          valPower: power |> Coin.getBandAmountFromCoin,
          valVote: Some(answer),
          delVotes: _ => 0.,
          proposalID,
        },
      )
    });
  let valVotesWithDelVotesMap =
    delVotes
    ->Belt_Array.reduce(valVotesMap, (acc, {validatorID, answer, power, proposalID}) => {
        acc->Belt_MapString.update(
          (validatorID |> string_of_int) ++ (proposalID |> ID.Proposal.toString),
          v => {
            let entry =
              v->Belt_Option.getWithDefault({
                validatorID,
                valPower: 0.,
                valVote: None,
                delVotes: _ => 0.,
                proposalID,
              });
            Some({
              ...entry,
              delVotes: ans =>
                ans == answer ? power |> Coin.getBandAmountFromCoin : entry.delVotes(ans),
            });
          },
        )
      })
    ->Belt_MapString.valuesToArray;
  let parsedData =
    valVotesWithDelVotesMap
    ->Belt_Array.reduce(
        Belt_MapString.empty, (acc, {validatorID, valPower, valVote, delVotes, proposalID}) => {
        acc->Belt_MapString.set(
          (validatorID |> string_of_int) ++ (proposalID |> ID.Proposal.toString),
          {
            validatorID,
            validatorPower:
              switch (valVote) {
              | Some(_) =>
                valPower
                -. delVotes(Yes)
                -. delVotes(No)
                -. delVotes(NoWithVeto)
                -. delVotes(Abstain)
              | None => 0.
              },
            validatorAns: valVote,
            proposalID,
          },
        )
      })
    ->Belt_MapString.valuesToArray;
  parsedData;
};

let getVoteStatByProposalID = proposalID => {
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
      DelegatorVoteByProposalIDConfig.definition,
      ~variables=
        DelegatorVoteByProposalIDConfig.makeVariables(
          ~proposal_id=proposalID |> ID.Proposal.toInt,
          (),
        ),
    );

  let%Sub valVotes = validatorVotes;
  let%Sub delVotes = delegatorVotes;

  let parsedData =
    parse(valVotes##validator_vote_proposals_view, delVotes##non_validator_vote_proposals_view);

  let delegatorData = delVotes##non_validator_vote_proposals_view;

  let validatorPower =
    parsedData->Belt_Array.reduce(0., (a, {validatorPower}) => a +. validatorPower);

  let delegatorPower =
    delVotes##non_validator_vote_proposals_view
    ->Belt_Array.reduce(0., (a, {power}) => a +. (power |> Coin.getBandAmountFromCoin));

  let totalPower = validatorPower +. delegatorPower;

  let totalYesPower = parsedData->getValVote(Yes) +. delegatorData->getDelVote(Yes);
  let totalNoPower = parsedData->getValVote(No) +. delegatorData->getDelVote(No);
  let totalNoWithVetoPower =
    parsedData->getValVote(NoWithVeto) +. delegatorData->getDelVote(NoWithVeto);
  let totalAbstainPower = parsedData->getValVote(Abstain) +. delegatorData->getDelVote(Abstain);

  Sub.resolve({
    proposalID,
    totalYes: totalYesPower,
    totalYesPercent: totalPower == 0. ? 0. : totalYesPower /. totalPower *. 100.,
    totalNo: totalNoPower,
    totalNoPercent: totalPower == 0. ? 0. : totalNoPower /. totalPower *. 100.,
    totalNoWithVeto: totalNoWithVetoPower,
    totalNoWithVetoPercent: totalPower == 0. ? 0. : totalNoWithVetoPower /. totalPower *. 100.,
    totalAbstain: totalAbstainPower,
    totalAbstainPercent: totalPower == 0. ? 0. : totalAbstainPower /. totalPower *. 100.,
    total: totalPower,
  });
};

let getVoteStats = () => {
  let (validatorVotes, _) = ApolloHooks.useSubscription(ValidatorVotesConfig.definition);
  let (delegatorVotes, _) = ApolloHooks.useSubscription(DelegatorVotesConfig.definition);

  let%Sub valVotes = validatorVotes;
  let%Sub delVotes = delegatorVotes;

  //TODO: Used too many mapping, revisit later.
  let parsedData =
    parse(valVotes##validator_vote_proposals_view, delVotes##non_validator_vote_proposals_view);

  let parsedMap =
    parsedData
    ->Belt_Array.reduce(Belt_MapInt.empty, (acc, {validatorID, proposalID}) =>
        acc->Belt_MapInt.set(
          proposalID |> ID.Proposal.toInt,
          {
            proposalID,
            validatorID,
            validatorPower:
              parsedData->Belt_Array.reduce(0., (a, {validatorPower, proposalID: delProposalID}) => {
                proposalID == delProposalID ? a +. validatorPower : a
              }),
            validatorAns: None,
          },
        )
      )
    ->Belt_MapInt.valuesToArray;

  let voteMap =
    parsedMap->Belt_Array.reduce(Belt_MapInt.empty, (acc, {proposalID, validatorPower}) =>
      acc->Belt_MapInt.set(
        proposalID |> ID.Proposal.toInt,
        {
          validatorPower
          +. delVotes##non_validator_vote_proposals_view
             ->Belt_Array.reduce(0., (a, {power, proposalID: delProposalID}) => {
                 proposalID == delProposalID ? a +. (power |> Coin.getBandAmountFromCoin) : a
               });
        },
      )
    );

  Sub.resolve(voteMap);
};
