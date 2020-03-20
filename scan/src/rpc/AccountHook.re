module Account = {
  type delegator_t = {
    validatorAddress: string,
    balance: float,
  };

  type reward_t = {
    validatorAddress: string,
    reward: float,
  };

  type delelegation_t = {
    validatorAddress: string,
    balance: float,
    reward: float,
  };

  type t = {
    balance: float,
    balanceStake: float,
    reward: float,
  };

  let filterUband = coins =>
    coins |> Belt_List.reduce(_, 0.0, (_, (x, y)) => +. (compare("uband", y) == 0 ? x : 0.));

  let decodeCoin = json =>
    JsonUtils.Decode.(json |> field("amount", uamount), json |> field("denom", string));

  let decode = (balancesJson, balanceStakejson, rewardsJson) =>
    JsonUtils.Decode.{
      balance: balancesJson |> field("result", list(decodeCoin)) |> filterUband,
      balanceStake:
        balanceStakejson
        |> field("result", list(json => json |> field("balance", uamount)))
        |> Belt_List.reduce(_, 0.0, (+.)),
      reward: rewardsJson |> at(["result", "total"], list(decodeCoin)) |> filterUband,
    };

  let decodeDelegator = json =>
    JsonUtils.Decode.(
      json |> field("validator_address", string),
      json |> field("balance", uamount),
    );
  let decodeDelegators = json =>
    JsonUtils.Decode.(json |> field("result", list(decodeDelegator)));

  let decodeRewardUband = json =>
    JsonUtils.Decode.(json |> field("amount", uamount), json |> field("denom", string));
  let decodeReward = json =>
    JsonUtils.Decode.(
      json |> field("validator_address", string),
      json
      |> field("reward", list(decodeRewardUband))
      |> Belt_List.reduce(_, 0.0, (_, (x, y)) => +. (compare("uband", y) == 0 ? x : 0.)),
    );
  let decodeRewards = json =>
    JsonUtils.Decode.(json |> at(["result", "rewards"], list(decodeReward)));
};
let get = address => {
  let addressStr = address |> Address.toBech32;

  let balancesJsonOpt = AxiosHooks.use({j|bank/balances/$addressStr|j});
  let balanceStakeJsonOpt = AxiosHooks.use({j|staking/delegators/$addressStr/delegations|j});
  let rewardsJsonOpt = AxiosHooks.use({j|distribution/delegators/$addressStr/rewards|j});

  let%Opt balancesJson = balancesJsonOpt;
  let%Opt balanceStakeJson = balanceStakeJsonOpt;
  let%Opt rewardsJson = rewardsJsonOpt;

  Some(Account.decode(balancesJson, balanceStakeJson, rewardsJson));
};

let getDelegations = address => {
  let addressStr = address |> Address.toBech32;

  let delegationsJson = AxiosHooks.use({j|staking/delegators/$addressStr/delegations|j});
  let rewardsJson = AxiosHooks.use({j|distribution/delegators/$addressStr/rewards|j});

  let%Opt pairReward = rewardsJson |> Belt_Option.map(_, Account.decodeRewards);
  let%Opt pairDelegation = delegationsJson |> Belt_Option.map(_, Account.decodeDelegators);

  Some(
    {
      let%IterList (rewardAddress, rewardValue) = pairReward;
      let%IterList (delegationAddress, balance) = pairDelegation;

      if (compare(rewardAddress, delegationAddress) == 0) {
        [(rewardAddress, balance, rewardValue)];
      } else {
        [];
      };
    },
  );
};
