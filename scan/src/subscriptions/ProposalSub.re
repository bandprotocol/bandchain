type proposal_status_t =
  | Deposit
  | Voting
  | Passed
  | Rejected
  | Failed;

let parseProposalStatus = json => {
  exception NotFound(string);
  let status = json |> Js.Json.decodeString |> Belt_Option.getExn;
  switch (status) {
  | "DepositPeriod" => Deposit
  | "VotingPeriod" => Voting
  | "Passed" => Passed
  | "Rejected" => Rejected
  | "Failed" => Failed
  | _ => raise(NotFound("The proposal status is not existing"))
  };
};

type account_t = {address: Address.t};

type deposit_t = {amount: list(Coin.t)};

type internal_t = {
  id: ID.Proposal.t,
  title: string,
  status: proposal_status_t,
  description: string,
  submitTime: MomentRe.Moment.t,
  depositEndTime: MomentRe.Moment.t,
  votingStartTime: MomentRe.Moment.t,
  votingEndTime: MomentRe.Moment.t,
  account: option(account_t),
  proposalType: string,
  totalDeposit: list(Coin.t),
};

type t = {
  id: ID.Proposal.t,
  name: string,
  status: proposal_status_t,
  description: string,
  submitTime: MomentRe.Moment.t,
  depositEndTime: MomentRe.Moment.t,
  votingStartTime: MomentRe.Moment.t,
  votingEndTime: MomentRe.Moment.t,
  proposerAddress: option(Address.t),
  proposalType: string,
  totalDeposit: list(Coin.t),
};

let toExternal =
    (
      {
        id,
        title,
        status,
        description,
        submitTime,
        depositEndTime,
        votingStartTime,
        votingEndTime,
        account,
        proposalType,
        totalDeposit,
      },
    ) => {
  id,
  name: title,
  status,
  description,
  submitTime,
  depositEndTime,
  votingStartTime,
  votingEndTime,
  proposerAddress: account->Belt.Option.map(({address}) => address),
  proposalType,
  totalDeposit,
};

module MultiConfig = [%graphql
  {|
  subscription Proposals($limit: Int!, $offset: Int!) {
    proposals(limit: $limit, offset: $offset, order_by: {id: desc}, where: {status: {_neq: "Inactive"}}) @bsRecord {
      id @bsDecoder(fn: "ID.Proposal.fromInt")
      title
      status @bsDecoder(fn: "parseProposalStatus")
      description
      submitTime: submit_time @bsDecoder(fn: "GraphQLParser.timestamp")
      depositEndTime: deposit_end_time @bsDecoder(fn: "GraphQLParser.timestamp")
      votingStartTime: voting_time @bsDecoder(fn: "GraphQLParser.timestamp")
      votingEndTime: voting_end_time @bsDecoder(fn: "GraphQLParser.timestamp")
      proposalType: type
      account @bsRecord {
        address @bsDecoder(fn: "Address.fromBech32")
      }
      totalDeposit: total_deposit @bsDecoder(fn: "GraphQLParser.coins")
    }
  }
|}
];

module SingleConfig = [%graphql
  {|
  subscription Proposal($id: Int!) {
    proposals_by_pk(id: $id) @bsRecord {
      id @bsDecoder(fn: "ID.Proposal.fromInt")
      title
      status @bsDecoder(fn: "parseProposalStatus")
      description
      submitTime: submit_time @bsDecoder(fn: "GraphQLParser.timestamp")
      depositEndTime: deposit_end_time @bsDecoder(fn: "GraphQLParser.timestamp")
      votingStartTime: voting_time @bsDecoder(fn: "GraphQLParser.timestamp")
      votingEndTime: voting_end_time @bsDecoder(fn: "GraphQLParser.timestamp")
      proposalType: type
      account @bsRecord {
          address @bsDecoder(fn: "Address.fromBech32")
      }
      totalDeposit: total_deposit @bsDecoder(fn: "GraphQLParser.coins")
    }
  }
|}
];

module ProposalsCountConfig = [%graphql
  {|
  subscription ProposalsCount {
    proposals_aggregate{
      aggregate{
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

let getList = (~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ()),
    );
  result |> Sub.map(_, internal => internal##proposals->Belt_Array.map(toExternal));
};

let get = id => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=SingleConfig.makeVariables(~id=id |> ID.Proposal.toInt, ()),
    );

  let%Sub x = result;
  switch (x##proposals_by_pk) {
  | Some(data) => Sub.resolve(data |> toExternal)
  | None => NoData
  };
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(ProposalsCountConfig.definition);
  result
  |> Sub.map(_, x => x##proposals_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};
