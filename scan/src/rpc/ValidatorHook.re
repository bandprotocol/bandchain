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

module ActiveValidator = {
  type t = {
    address: string,
    pubKey: string,
    proposerPriority: int,
    votingPower: int,
  };

  let decodeActiveValidator = json =>
    JsonUtils.Decode.{
      address: json |> field("address", string),
      pubKey: json |> field("pubKey", string),
      proposerPriority: json |> field("proposer_priority", intstr),
      votingPower: json |> field("voting_power", intstr),
    };

  let decodeActiveValidators = json =>
    JsonUtils.Decode.(
      json |> field("result", field("validators", list(decodeActiveValidator)))
    );
};

let getAll = (~pollInterval=?, ()) => {
  let json = Axios.use({j|staking/validators|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Validator.decodeValidators);
};

let getActive = (~pollInterval=?, ()) => {
  let json = Axios.use({j|validatorsets/latest|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, ActiveValidator.decodeActiveValidators);
};
