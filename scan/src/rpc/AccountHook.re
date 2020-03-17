let decodeBalances = json =>
  JsonUtils.Decode.(json |> field("result", list(json => json |> field("amount", uamount))));
let decodeStakes = json =>
  JsonUtils.Decode.(
    json |> at(["result", "total"], list(json => json |> field("amount", uamount)))
  );

let getBalanceStake = address => {
  let addressStr = address |> Address.toBech32;
  let json = AxiosHooks.use({j|distribution/delegators/$addressStr/rewards|j});

  let balances = json |> Belt.Option.map(_, decodeStakes);
  switch (balances) {
  | None => None
  | Some(x) => Some(x |> Belt_List.reduce(_, 0.0, (x, y) => x +. y))
  };
};

let decodeDelegation = json =>
  JsonUtils.Decode.(
    json |> field("delegator_address", string) |> Address.fromBech32,
    json |> field("shares", floatstr),
    10,
  );
let decodeDelegations = json =>
  JsonUtils.Decode.(json |> field("result", list(decodeDelegation)));

let getBalance = address => {
  let addressStr = address |> Address.toBech32;
  let json = AxiosHooks.use({j|bank/balances/$addressStr|j});

  let balances = json |> Belt.Option.map(_, decodeBalances);

  switch (balances) {
  | None => None
  | Some(x) => Some(x |> Belt_List.reduce(_, 0.0, (x, y) => x +. y))
  };
};

let getDelegations = address => {
  let addressStr = address |> Address.toBech32;
  let json = AxiosHooks.use({j|staking/delegators/$addressStr/delegations|j});
  json |> Belt.Option.map(_, decodeDelegations);
};
