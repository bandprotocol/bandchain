module Validator = {
  type t = {
    operatorAddress: Address.t,
    consensusPubkey: PubKey.t,
    moniker: string,
    identity: string,
    website: string,
    details: string,
    tokens: float,
    commission: float,
    uptime: float,
    repoertRate: float,
  };

  let decodeValidator = json =>
    JsonUtils.Decode.{
      operatorAddress: json |> field("operator_address", string) |> Address.fromBech32,
      consensusPubkey: json |> field("consensus_pubkey", string) |> PubKey.fromBech32,
      moniker: json |> at(["description", "moniker"], string),
      identity: json |> at(["description", "identity"], string),
      website: json |> at(["description", "website"], string),
      details: json |> at(["description", "details"], string),
      tokens: json |> at(["tokens"], uamount),
      commission:
        json |> at(["commission", "commission_rates", "rate"], JsonUtils.Decode.floatstr),
      // TODO: hard code
      uptime: 100.0,
      repoertRate: 100.0,
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

module ValidatorIndexInfo = {
  type t = {
    isActive: bool,
    operatorAddress: string,
    rewardDestinationAddress: string,
    votingPower: float,
    commission: float,
    bondedHeight: int,
    website: string,
    details: string,
  };
};

module ValidatorIndexNodeSatus = {
  type t = {
    uptime: float,
    avgResponseTime: int,
  };
};

module ValidatorIndexRequestResponse = {
  type t = {
    completedRequestCount: int,
    missedRequestCount: int,
  };
};

module ValidatorIndexProposedBlocks = {
  type t = {proposedBlockCount: int};
};

module ValidatorIndexDelegators = {
  type t = {delegatorCount: int};
};

module ValidatorIndexReports = {
  type t = {reportCount: int};
};

module ValidatorIndexProposedBlock = {
  type t = {
    height: int,
    timestamp: MomentRe.Moment.t,
    blockHash: string,
    txn: int,
  };
};

module ValidatorIndexDelegator = {
  type t = {
    delegator: string,
    sharePercentage: float,
    amount: int,
  };
};

module External = {
  type t = {
    externalID: int,
    externalValue: int,
  };
};

module ValidatorIndexReport = {
  type t = {
    requestID: int,
    txHash: string,
    oracleScriptID: int,
    oracleScriptName: string,
    dataSourceIDList: list(int),
    externalList: list(External.t),
  };
};

let getValidator = (~limit=10, ~page=1, ~status="bonded", ()) => {
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

let getValidatorIndexInfo = _ => {
  ValidatorIndexInfo.{
    isActive: true,
    operatorAddress: "bandvaloperwklefk234sdhf2jsadhfkalshdfk13e42",
    rewardDestinationAddress: "band17ljds2gj3kds234lkg",
    votingPower: 45.34,
    commission: 3.00,
    bondedHeight: 1,
    website: "https://coingeko.node",
    details: "We are the leading staking service provider for blockchain projects.",
  };
};

let getValidatorNodeStatus = _ => {
  ValidatorIndexNodeSatus.{uptime: 100.00, avgResponseTime: 2};
};

let getValidatorIndexRequestResponse = _ => {
  ValidatorIndexRequestResponse.{completedRequestCount: 23459, missedRequestCount: 100};
};

let getValidatorProposedBlocks = _ => {
  ValidatorIndexProposedBlocks.{proposedBlockCount: 2390};
};

let getValidatorIndexDelegators = _ => {
  ValidatorIndexDelegators.{delegatorCount: 4};
};

let getValidatorIndexReports = _ => {
  ValidatorIndexReports.{reportCount: 2};
};

let getValidatorIndexProposedBlockList = _ => {
  [
    ValidatorIndexProposedBlock.{
      height: 10,
      timestamp: MomentRe.momentWithUnix(1583465551),
      blockHash: "bandvaloperwklefk234sdhf2jsadhfkalshdfk13e42",
      txn: 10,
    },
    ValidatorIndexProposedBlock.{
      height: 11,
      timestamp: MomentRe.momentWithUnix(1583465599),
      blockHash: "bandvaloperwklefkasdadjsadhfkalshdfk13e42",
      txn: 12,
    },
  ];
};

let getValidatorIndexDelegatorList = _ => {
  [
    ValidatorIndexDelegator.{
      delegator: "bandvaloperwklefk234sdhf2jsadhfkalshdfk13e42",
      sharePercentage: 12.0,
      amount: 12,
    },
    ValidatorIndexDelegator.{
      delegator: "bandvaloperw123123123312f2jsadhfkalshdfk13e42",
      sharePercentage: 88.0,
      amount: 88,
    },
  ];
};

let getValidatorIndexReportList = _ => {
  [
    ValidatorIndexReport.{
      requestID: 10,
      txHash: "29NMNMSD3312SF21DF3DS",
      oracleScriptID: 213,
      oracleScriptName: "Mean Crypto Price",
      dataSourceIDList: [1, 2],
      externalList:
        External.[{externalID: 1, externalValue: 10}, {externalID: 2, externalValue: 2}],
    },
    ValidatorIndexReport.{
      requestID: 13,
      txHash: "AKJS123FNK213SL3DF",
      oracleScriptID: 112,
      oracleScriptName: "US powerball",
      dataSourceIDList: [3, 4],
      externalList:
        External.[{externalID: 3, externalValue: 110}, {externalID: 4, externalValue: 32}],
    },
  ];
};