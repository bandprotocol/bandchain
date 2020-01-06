module Validator = {
  type t = {
    operatorAddress: string,
    consensusPubkey: string,
    moniker: string,
    identity: string,
    website: string,
    details: string,
  };

  let decodeValidator = json =>
    JsonUtils.Decode.{
      operatorAddress: json |> field("operator_address", string),
      consensusPubkey: json |> field("consensus_pubkey", string),
      moniker: json |> at(["description", "moniker"], string),
      identity: json |> at(["description", "identity"], string),
      website: json |> at(["description", "website"], string),
      details: json |> at(["description", "details"], string),
    };

  let decodeValidators = json =>
    JsonUtils.Decode.(json |> field("result", list(decodeValidator)));
};

let get = (~pollInterval=?, ()) => {
  let json = Axios.use({j|staking/validators|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Validator.decodeValidators);
};
