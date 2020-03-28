module Account = {
  type delegator_t = {
    validatorAddress: string,
    balance: float,
  };

  type reward_t = {
    validatorAddress: string,
    reward: float,
  };

  type delegation_t = {
    validatorAddress: string,
    balance: float,
    reward: float,
  };

  type t = {
    balance: float,
    balanceStake: float,
    reward: float,
    delegations: list(delegation_t),
  };

  let filterUband = coins =>
    coins |> Belt_List.reduce(_, 0.0, (_, (x, y)) => +. (compare("uband", y) == 0 ? x : 0.));

  let decodeCoin = json =>
    JsonUtils.Decode.(json |> field("amount", uamount), json |> field("denom", string));

  let decodeReward = json =>
    JsonUtils.Decode.{
      validatorAddress: json |> field("validator_address", string),
      reward: json |> field("reward", list(decodeCoin)) |> filterUband,
    };
  let decodeRewards = json =>
    JsonUtils.Decode.(json |> at(["result", "rewards"], list(decodeReward)));

  let decodeDelegator = json =>
    JsonUtils.Decode.{
      validatorAddress: json |> field("validator_address", string),
      balance: json |> field("balance", uamount),
    };
  let decodeDelegators = json =>
    JsonUtils.Decode.(json |> field("result", list(decodeDelegator)));

  let decodeDelegations = (rewardsJson, delegationsJson) => {
    let rewards = rewardsJson |> decodeRewards;
    let delegators = delegationsJson |> decodeDelegators;

    let%IterList reward = rewards;
    let%IterList delegator = delegators;

    if (compare(reward.validatorAddress, delegator.validatorAddress) == 0) {
      [
        {
          validatorAddress: reward.validatorAddress,
          balance: delegator.balance,
          reward: reward.reward,
        },
      ];
    } else {
      [];
    };
  };
  let decode = (balancesJson, balanceStakejson, rewardsJson, delegatorsJson) =>
    JsonUtils.Decode.{
      balance: balancesJson |> field("result", list(decodeCoin)) |> filterUband,
      balanceStake:
        balanceStakejson
        |> field("result", list(json => json |> field("balance", uamount)))
        |> Belt_List.reduce(_, 0.0, (+.)),
      reward: rewardsJson |> at(["result", "total"], list(decodeCoin)) |> filterUband,
      delegations: decodeDelegations(rewardsJson, delegatorsJson),
    };
};
let get = address => {
  let addressStr = address |> Address.toBech32;

  let balancesJsonOpt = AxiosHooks.use({j|bank/balances/$addressStr|j});
  let balanceStakeJsonOpt = AxiosHooks.use({j|staking/delegators/$addressStr/delegations|j});
  let rewardsJsonOpt = AxiosHooks.use({j|distribution/delegators/$addressStr/rewards|j});
  let delegationsJsonOpt = AxiosHooks.use({j|staking/delegators/$addressStr/delegations|j});

  let%Opt balancesJson = balancesJsonOpt;
  let%Opt balanceStakeJson = balanceStakeJsonOpt;
  let%Opt rewardsJson = rewardsJsonOpt;
  let%Opt delegationsJson = delegationsJsonOpt;
  Some(Account.decode(balancesJson, balanceStakeJson, rewardsJson, delegationsJson));
};
