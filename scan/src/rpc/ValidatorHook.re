module Validator = {
  type t = {
    operatorAddress: Address.t,
    consensusPubkey: PubKey.t,
    moniker: string,
    identity: string,
    website: string,
    details: string,
    tokens: float,
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
    };

  let decodeValidators = json =>
    JsonUtils.Decode.(json |> field("result", list(decodeValidator)));
};

let get = (~pollInterval=?, ()) => {
  let json = AxiosHooks.use({j|staking/validators|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Validator.decodeValidators);
};
