module Mini = {
  type t = {
    consensusAddress: string,
    operatorAddress: Address.t,
    moniker: string,
  };
};

type internal_t = {
  operatorAddress: Address.t,
  consensusAddress: Address.t,
  moniker: string,
  identity: string,
  website: string,
  tokens: Coin.t,
  commissionRate: float,
  consensusPubKey: PubKey.t,
  bondedHeight: int,
  jailed: bool,
  details: string,
};

type t = {
  avgResponseTime: int,
  isActive: bool,
  operatorAddress: Address.t,
  consensusAddress: Address.t,
  consensusPubKey: PubKey.t,
  votingPower: float,
  moniker: string,
  identity: string,
  website: string,
  details: string,
  tokens: Coin.t,
  commission: float,
  bondedHeight: int,
  completedRequestCount: int,
  missedRequestCount: int,
  uptime: option(float),
};

let toExternal =
    (
      {
        operatorAddress,
        consensusAddress,
        moniker,
        identity,
        website,
        tokens,
        commissionRate,
        consensusPubKey,
        bondedHeight,
        jailed,
        details,
      }: internal_t,
    ) => {
  isActive: !jailed,
  operatorAddress,
  consensusAddress,
  consensusPubKey,
  votingPower: tokens.amount,
  moniker,
  identity,
  website,
  details,
  tokens,
  commission: commissionRate *. 100.,
  bondedHeight,
  // TODO: remove hardcoded when somewhere use it
  avgResponseTime: 2,
  completedRequestCount: 23459,
  missedRequestCount: 20,
  uptime: None,
};

type validator_vote_t = {
  consensusAddress: Address.t,
  count: int,
  voted: bool,
};

module SingleConfig = [%graphql
  {|
      subscription Validator($operator_address: String!) {
        validators_by_pk(operator_address: $operator_address) @bsRecord {
          operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
          consensusAddress: consensus_address @bsDecoder(fn: "Address.fromHex")
          moniker
          identity
          website
          tokens @bsDecoder(fn: "GraphQLParser.coin")
          commissionRate: commission_rate @bsDecoder(fn: "float_of_string")
          consensusPubKey: consensus_pubkey @bsDecoder(fn: "PubKey.fromBech32")
          bondedHeight: bonded_height @bsDecoder(fn: "GraphQLParser.int64")
          jailed
          details
        }
      }
  |}
];

module MultiConfig = [%graphql
  {|
      subscription Validators($limit: Int!, $offset: Int!, $jailed: Boolean!) {
        validators(limit: $limit, offset: $offset, where: {jailed: {_eq: $jailed}}, order_by: {tokens: desc}) @bsRecord {
          operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
          consensusAddress: consensus_address @bsDecoder(fn: "Address.fromHex")
          moniker
          identity
          website
          tokens @bsDecoder(fn: "GraphQLParser.coin")
          commissionRate: commission_rate @bsDecoder(fn: "float_of_string")
          consensusPubKey: consensus_pubkey @bsDecoder(fn: "PubKey.fromBech32")
          bondedHeight: bonded_height @bsDecoder(fn: "GraphQLParser.int64")
          jailed
          details
        }
      }
  |}
];

module TotalBondedAmountConfig = [%graphql
  {|
  subscription TotalBondedAmount{
    validators_aggregate{
      aggregate{
        sum{
          tokens @bsDecoder(fn: "GraphQLParser.coinWithDefault")
        }
      }
    }
  }
  |}
];

module ValidatorCountConfig = [%graphql
  {|
    subscription ValidatorCount {
      validators_aggregate{
        aggregate{
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

module ValidatorCountByJailedConfig = [%graphql
  {|
    subscription ValidatorCountByJailed($jailed: Boolean!) {
      validators_aggregate(where: {jailed: {_eq: $jailed}}) {
        aggregate{
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
      }
    }
  |}
];

module SingleLast250VotedConfig = [%graphql
  {|
  subscription ValidatorLast25Voted($consensusAddress: String!) {
    validator_last_250_votes(where: {consensus_address: {_eq: $consensusAddress}}) {
      count
      voted
    }
  }
|}
];

module MultiLast250VotedConfig = [%graphql
  {|
  subscription ValidatorsLast25Voted {
    validator_last_250_votes {
      consensus_address
      count
      voted
    }
  }
|}
];

let get = operator_address => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=
        SingleConfig.makeVariables(
          ~operator_address=operator_address |> Address.toOperatorBech32,
          (),
        ),
    );
  let%Sub x = result;
  switch (x##validators_by_pk) {
  | Some(data) => Sub.resolve(data |> toExternal)
  | None => NoData
  };
};

let getList = (~page, ~pageSize, ~isActive, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ~jailed=!isActive, ()),
    );
  result |> Sub.map(_, x => x##validators->Belt_Array.map(toExternal));
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(ValidatorCountConfig.definition);
  result
  |> Sub.map(_, x => x##validators_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};

let countByActive = isActive => {
  let (result, _) =
    ApolloHooks.useSubscription(
      ValidatorCountByJailedConfig.definition,
      ~variables=ValidatorCountByJailedConfig.makeVariables(~jailed=!isActive, ()),
    );
  result
  |> Sub.map(_, x => x##validators_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};

let getTotalBondedAmount = () => {
  let (result, _) = ApolloHooks.useSubscription(TotalBondedAmountConfig.definition);
  result
  |> Sub.map(_, a =>
       ((a##validators_aggregate##aggregate |> Belt_Option.getExn)##sum |> Belt_Option.getExn)##tokens
     );
};

let getUptime = consensusAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleLast250VotedConfig.definition,
      ~variables=
        SingleLast250VotedConfig.makeVariables(
          ~consensusAddress=consensusAddress |> Address.toHex(~upper=true),
          (),
        ),
    );

  let%Sub x = result;
  let validatorVotes = x##validator_last_250_votes;
  let signedBlock =
    validatorVotes
    ->Belt.Array.keep(each => each##voted == Some(true))
    ->Belt.Array.get(0)
    ->Belt.Option.flatMap(each => each##count)
    ->Belt.Option.mapWithDefault(0, GraphQLParser.int64)
    |> float_of_int;

  let missedBlock =
    validatorVotes
    ->Belt.Array.keep(each => each##voted == Some(false))
    ->Belt.Array.get(0)
    ->Belt.Option.flatMap(each => each##count)
    ->Belt.Option.mapWithDefault(0, GraphQLParser.int64)
    |> float_of_int;

  if (signedBlock == 0. && missedBlock == 0.) {
    Sub.resolve(None);
  } else {
    let uptime = signedBlock /. (signedBlock +. missedBlock) *. 100.;
    Sub.resolve(Some(uptime));
  };
};

// For computing uptime on Validator home page
let getListVotesBlock = () => {
  let (result, _) = ApolloHooks.useSubscription(MultiLast250VotedConfig.definition);

  let%Sub x = result;
  let validatorVotes =
    x##validator_last_250_votes
    ->Belt.Array.map(each =>
        {
          consensusAddress: each##consensus_address->Belt.Option.getExn->Address.fromHex,
          count: each##count->Belt.Option.getExn->GraphQLParser.int64,
          voted: each##voted->Belt.Option.getExn,
        }
      );

  Sub.resolve(validatorVotes);
};
