module Delegators = {
  type delegator_t = {
    validatorAddress: string,
    balance: float,
  };

  type reward_t = {
    validatorAddress: string,
    reward: float,
  };

  type merge_t = {
    validatorAddress: string,
    balance: float,
    reward: float,
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

let decodeDelegation = (delegationJson, rewardJson) =>
  JsonUtils.Decode.(
    delegationJson |> field("validator", string),
    rewardJson |> field("amount", uamount),
  );

let decodeDelegations = (delegationJsons, rewardJsons) =>
  JsonUtils.Decode.(
    delegationJsons |> field("result"),
    rewardJsons |> at(["result", "rewards"]),
  );

let decodeBalances = json =>
  JsonUtils.Decode.(json |> field("result", list(json => json |> field("amount", uamount))));

let decodeBalanceStakes = json =>
  JsonUtils.Decode.(json |> field("result", list(json => json |> field("balance", uamount))));

let decodeReward = json =>
  JsonUtils.Decode.(
    json |> at(["result", "total"], list(json => json |> field("amount", uamount)))
  );

let getBalanceStake = address => {
  let addressStr = address |> Address.toBech32;
  let json = AxiosHooks.use({j|staking/delegators/$addressStr/delegations|j});

  let balances = json |> Belt.Option.map(_, decodeBalanceStakes);
  switch (balances) {
  | None => None
  | Some(x) => Some(x |> Belt_List.reduce(_, 0.0, (+.)))
  };
};

let getBalance = address => {
  let addressStr = address |> Address.toBech32;
  let json = AxiosHooks.use({j|bank/balances/$addressStr|j});

  let balances = json |> Belt.Option.map(_, decodeBalances);

  switch (balances) {
  | None => None
  | Some(x) => Some(x |> Belt_List.reduce(_, 0.0, (+.)))
  };
};

let getDelegations = address => {
  let addressStr = address |> Address.toBech32;

  let delegationsJson = AxiosHooks.use({j|staking/delegators/$addressStr/delegations|j});
  let rewardsJson = AxiosHooks.use({j|distribution/delegators/$addressStr/rewards|j});

  let%Opt pairReward = rewardsJson |> Belt_Option.map(_, Delegators.decodeRewards);
  let%Opt pairDelegation = delegationsJson |> Belt_Option.map(_, Delegators.decodeDelegators);

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

let getReward = address => {
  let addressStr = address |> Address.toBech32;
  let json = AxiosHooks.use({j|distribution/delegators/$addressStr/rewards|j});

  let%Opt rewards = json |> Belt.Option.map(_, decodeReward);

  Some(rewards |> Belt_List.reduce(_, 0.0, (+.)));
};
