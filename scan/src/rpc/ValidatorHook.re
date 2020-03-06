module Validator = {
  type status_t = {
    uptime: float,
    avgResponseTime: int,
  };
  type request_t = {
    completedRequestCount: int,
    missedRequestCount: int,
  };
  type proposed_block_t = {
    height: int,
    timestamp: MomentRe.Moment.t,
    blockHash: string,
    txn: int,
  };
  type delegator_t = {
    delegator: string,
    sharePercentage: float,
    amount: int,
  };
  type external_t = {
    externalID: int,
    externalValue: int,
  };
  type report_t = {
    requestID: int,
    txHash: string,
    oracleScriptID: int,
    oracleScriptName: string,
    dataSourceIDList: list(int),
    externalList: list(external_t),
  };

  type t = {
    isActive: bool,
    operatorAddress: Address.t,
    consensusPubkey: PubKey.t,
    rewardDestinationAddress: string,
    votingPower: float,
    moniker: string,
    identity: string,
    website: string,
    details: string,
    tokens: float,
    commission: float,
    bondedHeight: int,
    uptime: float,
    completedRequestCount: int,
    missedRequestCount: int,
    indexNodeStatus: status_t,
    indexRequestCountResponse: request_t,
    proposedBlockList: list(proposed_block_t),
    delegatorList: list(delegator_t),
    reportList: list(report_t),
  };

  let decodeValidator = json =>
    JsonUtils.Decode.{
      isActive: true,
      operatorAddress: json |> field("operator_address", string) |> Address.fromBech32,
      consensusPubkey: json |> field("consensus_pubkey", string) |> PubKey.fromBech32,
      rewardDestinationAddress: "band17ljds2gj3kds234lkg",
      votingPower: 45.34,
      moniker: json |> at(["description", "moniker"], string),
      identity: json |> at(["description", "identity"], string),
      website: json |> at(["description", "website"], string),
      details: json |> at(["description", "details"], string),
      tokens: json |> at(["tokens"], uamount),
      commission:
        json |> at(["commission", "commission_rates", "rate"], JsonUtils.Decode.floatstr),
      bondedHeight: 1,
      // TODO: hard code
      uptime: 100.0,
      completedRequestCount: 23459,
      missedRequestCount: 100,
      indexNodeStatus: {
        uptime: 100.00,
        avgResponseTime: 2,
      },
      indexRequestCountResponse: {
        completedRequestCount: 23459,
        missedRequestCount: 100,
      },
      proposedBlockList: [
        {
          height: 10,
          timestamp: MomentRe.momentWithUnix(1583465551),
          blockHash: "bandvaloperwklefk234sdhf2jsadhfkalshdfk13e42",
          txn: 10,
        },
        {
          height: 11,
          timestamp: MomentRe.momentWithUnix(1583465599),
          blockHash: "bandvaloperwklefkasdadjsadhfkalshdfk13e42",
          txn: 12,
        },
      ],
      delegatorList: [
        {
          delegator: "bandvaloperwklefk234sdhf2jsadhfkalshdfk13e42",
          sharePercentage: 12.0,
          amount: 12,
        },
        {
          delegator: "bandvaloperw123123123312f2jsadhfkalshdfk13e42",
          sharePercentage: 88.0,
          amount: 88,
        },
      ],
      reportList: [
        {
          requestID: 10,
          txHash: "29NMNMSD3312SF21DF3DS",
          oracleScriptID: 213,
          oracleScriptName: "Mean Crypto Price",
          dataSourceIDList: [1, 2],
          externalList: [{externalID: 1, externalValue: 10}, {externalID: 2, externalValue: 2}],
        },
        {
          requestID: 13,
          txHash: "AKJS123FNK213SL3DF",
          oracleScriptID: 112,
          oracleScriptName: "US powerball",
          dataSourceIDList: [3, 4],
          externalList: [
            {externalID: 3, externalValue: 110},
            {externalID: 4, externalValue: 32},
          ],
        },
      ],
    };

  let decodeValidators = json =>
    JsonUtils.Decode.(json |> field("result", list(decodeValidator)));
};

module GlobalInfo = {
  type t = {
    allBondedAmount: int,
    totalSupply: int,
    inflationRate: float,
    avgBlockTime: float,
  };
};

module ValidatorStatus = {
  type t = {activeValidatorCount: int};

  let decodeValidators = json =>
    JsonUtils.Decode.{
      activeValidatorCount: json |> field("result", list(_ => ())) |> Belt_List.length,
    };
};

module ValidatorCount = {
  type t = {validatorCount: int};
  let decodeValidators = json =>
    JsonUtils.Decode.{
      validatorCount: json |> field("result", list(_ => ())) |> Belt_List.length,
    };
};

let getValidators = (~limit=10, ~page=1, ~status="bonded", ()) => {
  let json = AxiosHooks.use({j|staking/validators?limit=$limit&page=$page&status=$status|j});
  json |> Belt.Option.map(_, Validator.decodeValidators);
};

let getGlobalInfo = _ => {
  GlobalInfo.{
    allBondedAmount: 5353500,
    totalSupply: 10849023,
    inflationRate: 12.45,
    avgBlockTime: 2.59,
  };
};

let getValidatorStatus = (~status="bonded", ()) => {
  let json = AxiosHooks.use({j|staking/validators?status=$status|j});
  Js.Console.log(json |> Belt.Option.map(_, ValidatorStatus.decodeValidators));
  json |> Belt.Option.map(_, ValidatorStatus.decodeValidators);
};

let getValidatorCount = _ => {
  let json = AxiosHooks.use({j|staking/validators|j});
  Js.Console.log(json |> Belt.Option.map(_, ValidatorCount.decodeValidators));
  json |> Belt.Option.map(_, ValidatorCount.decodeValidators);
};
