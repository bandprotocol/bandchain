module Mini = {
  type t = {
    consensusAddress: string,
    operatorAddress: Address.t,
    moniker: string,
    identity: string,
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
  commissionMaxChange: float,
  commissionMaxRate: float,
  consensusPubKey: PubKey.t,
  jailed: bool,
  oracleStatus: bool,
  details: string,
};

type t = {
  rank: int,
  avgResponseTime: int,
  isActive: bool,
  oracleStatus: bool,
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
  commissionMaxChange: float,
  commissionMaxRate: float,
  completedRequestCount: int,
  missedRequestCount: int,
  uptime: option(float),
  oracleReports: int,
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
        commissionMaxChange,
        commissionMaxRate,
        jailed,
        oracleStatus,
        details,
      }: internal_t,
      rank,
    ) => {
  rank,
  isActive: !jailed,
  oracleStatus,
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
  commissionMaxChange: commissionMaxChange *. 100.,
  commissionMaxRate: commissionMaxRate *. 100.,
  // TODO: remove hardcoded when somewhere use it
  avgResponseTime: 2,
  completedRequestCount: 23459,
  missedRequestCount: 20,
  uptime: None,
  oracleReports: 3000,
};

type validator_voted_status_t =
  | Missed
  | Signed
  | Proposed;

type validator_single_uptime_t = {
  blockHeight: ID.Block.t,
  status: validator_voted_status_t,
};

type validator_single_uptime_status_t = {
  validatorVotes: array(validator_single_uptime_t),
  proposedCount: int,
  missedCount: int,
  signedCount: int,
};

type validator_vote_t = {
  consensusAddress: Address.t,
  count: int,
  voted: bool,
};

type historical_oracle_statuses_t = {
  status: bool,
  timestamp: MomentRe.Moment.t,
};

type historical_oracle_statuses_count_t = {
  oracleStatusReports: array(historical_oracle_statuses_t),
  uptimeCount: int,
  downtimeCount: int,
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
          commissionMaxChange: commission_max_change @bsDecoder(fn: "float_of_string")
          commissionMaxRate: commission_max_rate @bsDecoder(fn: "float_of_string")
          consensusPubKey: consensus_pubkey @bsDecoder(fn: "PubKey.fromBech32")
          jailed
          details
          oracleStatus: status
        }
      }
  |}
];

module MultiConfig = [%graphql
  {|
      subscription Validators($jailed: Boolean!) {
        validators(where: {jailed: {_eq: $jailed}}, order_by: {tokens: desc, moniker: asc}) @bsRecord {
          operatorAddress: operator_address @bsDecoder(fn: "Address.fromBech32")
          consensusAddress: consensus_address @bsDecoder(fn: "Address.fromHex")
          moniker
          identity
          website
          tokens @bsDecoder(fn: "GraphQLParser.coin")
          commissionRate: commission_rate @bsDecoder(fn: "float_of_string")
          commissionMaxChange: commission_max_change @bsDecoder(fn: "float_of_string")
          commissionMaxRate: commission_max_rate @bsDecoder(fn: "float_of_string")
          consensusPubKey: consensus_pubkey @bsDecoder(fn: "PubKey.fromBech32")
          jailed
          details
          oracleStatus: status
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

module SingleLast100VotedConfig = [%graphql
  {|
  subscription SingleLast100Voted($consensusAddress: String!) {
    validator_votes(limit: 100, where: {validator: {consensus_address: {_eq: $consensusAddress}}}, order_by: {block_height: desc}) {
    block_height
    consensus_address
    voted
      block {
        proposer
      }
    }
  }
|}
];

module HistoricalOracleStatusesConfig = [%graphql
  {|
  subscription HistoricalOracleStatuses($operatorAddress: String!, $greater: timestamp) {
    historical_oracle_statuses(where: {operator_address: {_eq: $operatorAddress}, timestamp: {_gte: $greater}}) {
      operator_address
      status
      timestamp
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
  | Some(data) => Sub.resolve(data->toExternal(0)) // 0 is arbitrary rank.
  | None => NoData
  };
};

let getList = (~isActive, ()) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=MultiConfig.makeVariables(~jailed=!isActive, ()),
    );
  result
  |> Sub.map(_, x =>
       x##validators->Belt_Array.mapWithIndex((idx, each) => toExternal(each, idx + 1))
     );
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
          ~consensusAddress=consensusAddress |> Address.toHex,
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

let getBlockUptimeByValidator = consensusAddress => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleLast100VotedConfig.definition,
      ~variables=
        SingleLast100VotedConfig.makeVariables(
          ~consensusAddress=consensusAddress |> Address.toHex,
          (),
        ),
    );
  let%Sub x = result;
  let validatorVotes =
    x##validator_votes
    ->Belt.Array.map(each =>
        {
          blockHeight: each##block_height |> ID.Block.fromInt,
          status:
            switch (each##voted, each##block##proposer == consensusAddress->Address.toHex) {
            | (false, _) => Missed
            | (true, false) => Signed
            | (true, true) => Proposed
            },
        }
      );
  Sub.resolve({
    validatorVotes,
    proposedCount:
      validatorVotes->Belt.Array.keep(({status}) => status == Proposed)->Belt.Array.size,
    signedCount:
      validatorVotes->Belt.Array.keep(({status}) => status == Signed)->Belt.Array.size,
    missedCount:
      validatorVotes->Belt.Array.keep(({status}) => status == Missed)->Belt.Array.size,
  });
};

let getHistoricalOracleStatus = (operatorAddress, greater, oracleStatus) => {
  let (result, _) =
    ApolloHooks.useSubscription(
      HistoricalOracleStatusesConfig.definition,
      ~variables=
        HistoricalOracleStatusesConfig.makeVariables(
          ~operatorAddress=operatorAddress |> Address.toOperatorBech32,
          ~greater=greater |> MomentRe.Moment.format(Config.timestampUseFormat) |> Js.Json.string,
          (),
        ),
    );
  let%Sub x = result;
  let oracleStatusReports =
    x##historical_oracle_statuses->Belt.Array.size > 0
      ? x##historical_oracle_statuses
        ->Belt.Array.map(each =>
            {status: each##status, timestamp: each##timestamp |> GraphQLParser.timestamp}
          )
        ->Belt.List.fromArray
      : [{timestamp: greater |> MomentRe.Moment.startOf(`day), status: oracleStatus}];

  let normalizedDate =
    oracleStatusReports->Belt_List.map(({timestamp, status}) =>
      if (status) {
        {
          timestamp:
            timestamp
            |> MomentRe.Moment.startOf(`day)
            |> MomentRe.Moment.add(~duration=MomentRe.duration(1., `days)),
          status: true,
        };
      } else {
        {timestamp: timestamp |> MomentRe.Moment.startOf(`day), status: false};
      }
    );

  let optimizedDate =
    normalizedDate
    ->Belt_List.add({
        timestamp: greater |> MomentRe.Moment.startOf(`day),
        status: !normalizedDate->Belt_List.headExn.status,
      })
    ->Belt_List.zip(
        normalizedDate->Belt_List.concat([
          {
            timestamp:
              MomentRe.momentNow()
              |> MomentRe.Moment.defaultUtc
              |> MomentRe.Moment.startOf(`day)
              |> MomentRe.Moment.add(~duration=MomentRe.duration(1., `days)),
            status: false,
          },
        ]),
      );
  let response =
    {
      let%IterList ({timestamp: st, status}, {timestamp: en, _}) = optimizedDate;
      Belt_List.makeBy(en->MomentRe.diff(st, `days)->int_of_float, idx =>
        {
          timestamp:
            st |> MomentRe.Moment.add(~duration=MomentRe.duration(idx |> float_of_int, `days)),
          status,
        }
      );
    }
    ->Belt.List.toArray
    ->Belt.Array.sliceToEnd(1);

  Sub.resolve({
    oracleStatusReports: response,
    uptimeCount: response->Belt.Array.keep(({status}) => status)->Belt.Array.size,
    downtimeCount: response->Belt.Array.keep(({status}) => !status)->Belt.Array.size,
  });
};
