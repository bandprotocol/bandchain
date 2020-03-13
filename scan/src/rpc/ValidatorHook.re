type status_t =
  | Bonded
  | Unbonded
  | Unbonding;

module Validator = {
  type node_status_t = {
    uptime: float,
    avgResponseTime: int,
  };
  type delegator_t = {
    delegator: string,
    sharePercentage: float,
    amount: int,
  };
  type block_t = {
    hash: Hash.t,
    height: int,
    timestamp: MomentRe.Moment.t,
    proposer: Address.t,
    numTxs: int,
    totalTxs: int,
  };
  type report_t = {
    requestID: int,
    hash: Hash.t,
    oracleScriptID: int,
    oracleScriptName: string,
    dataSourceIDList: list(int),
    externalIDs: list(int),
    externalValues: list(int),
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
    nodeStatus: node_status_t,
    proposedBlocks: list(block_t),
    delegators: list(delegator_t),
    reports: list(report_t),
  };

  let decodeResult = json =>
    JsonUtils.Decode.{
      isActive: true,
      operatorAddress: json |> field("operator_address", string) |> Address.fromBech32,
      consensusPubkey: json |> field("consensus_pubkey", string) |> PubKey.fromBech32,
      rewardDestinationAddress: "band17ljds2gj3kds234lkg",
      votingPower: 25.,
      moniker: json |> at(["description", "moniker"], string),
      identity: json |> at(["description", "identity"], string),
      website: json |> at(["description", "website"], string),
      details: json |> at(["description", "details"], string),
      tokens: json |> at(["tokens"], uamount),
      commission:
        (json |> at(["commission", "commission_rates", "rate"], JsonUtils.Decode.floatstr))
        *. 100.,
      bondedHeight: 1,
      // TODO: mock for now
      uptime: 100.0,
      completedRequestCount: 23459,
      missedRequestCount: 20,
      nodeStatus: {
        uptime: 100.00,
        avgResponseTime: 2,
      },
      proposedBlocks: [
        {
          hash: Hash.fromHex("6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"),
          height: 10,
          timestamp: MomentRe.momentWithUnix(1583465551),
          proposer:
            Address.fromHex("15b6a7d60a1ba577524250833761c7afd6405d5a739791dca68feb3f11448506"),
          numTxs: 10,
          totalTxs: 100,
        },
        {
          hash: Hash.fromHex("4cbd10df46c72786dc91b78e44bea2e0718271495d46b63f59ff3cf2e6e86f96"),
          height: 10,
          timestamp: MomentRe.momentWithUnix(1583465251),
          proposer:
            Address.fromHex("414a1428edf228cd7233502534e462270a723e1a0f31a6fe591b108883dcbf74"),
          numTxs: 32,
          totalTxs: 100,
        },
      ],
      delegators: [
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
      reports: [
        {
          requestID: 10,
          hash:
            Hash.fromBase64("1be50992b1a00d9ea17acaafcf11615e4f37cbc50160e70be4992854b57264e8"),
          oracleScriptID: 213,
          oracleScriptName: "Mean Crypto Price",
          dataSourceIDList: [1, 2],
          externalIDs: [1, 2],
          externalValues: [1, 2],
        },
        {
          requestID: 13,
          hash:
            Hash.fromBase64("a469056894d91a4987ef2a07f16fc7bbae04f8166b43aa168e71e48b71cac651"),
          oracleScriptID: 112,
          oracleScriptName: "US powerball",
          dataSourceIDList: [3, 4],
          externalIDs: [3, 110],
          externalValues: [4, 32],
        },
      ],
    };

  let decodeList = json => JsonUtils.Decode.(json |> field("result", list(decodeResult)));

  let decode = json => JsonUtils.Decode.(json |> field("result", decodeResult));
};

module GlobalInfo = {
  type t = {
    allBondedAmount: int,
    totalSupply: int,
    inflationRate: float,
    avgBlockTime: float,
  };
};

let get = address => {
  let addressStr = address |> Address.toBech32;
  let json = AxiosHooks.use({j|staking/validator/$addressStr|j});
  json |> Belt.Option.map(_, Validator.decode);
};

let getList = (~limit=10, ~page=1, ~status="bonded", ()) => {
  let json = AxiosHooks.use({j|staking/validators?limit=$limit&page=$page&status=$status|j});
  json |> Belt.Option.map(_, Validator.decodeList);
};

// TODO: mock for now
let getGlobalInfo = _ => {
  GlobalInfo.{
    allBondedAmount: 5353500,
    totalSupply: 10849023,
    inflationRate: 12.45,
    avgBlockTime: 2.59,
  };
};

let toString =
  fun
  | Bonded => "bonded"
  | Unbonded => "unbonded"
  | Unbonding => "unbonding";

let getValidatorCount = (~status=Bonded, ()) => {
  let statusStr = status |> toString;
  let json = AxiosHooks.use({j|staking/validators?status=$statusStr|j});
  Belt_Option.mapWithDefault(json, [], JsonUtils.Decode.(field("result", list(_ => ()))))
  |> Belt_List.length;
};
