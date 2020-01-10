module Block = {
  type t = {
    hash: Hash.t,
    height: int,
    timestamp: MomentRe.Moment.t,
    proposer: Address.t,
  };

  let decodeBlock = json =>
    JsonUtils.Decode.{
      hash: json |> at(["block_meta", "block_id", "hash"], string) |> Hash.fromHex,
      height: json |> at(["block_meta", "header", "height"], intstr),
      timestamp: json |> at(["block_meta", "header", "time"], moment),
      proposer:
        json |> at(["block_meta", "header", "proposer_address"], string) |> Address.fromHex,
    };

  let decodeBlocks = json => JsonUtils.Decode.(json |> list(decodeBlock));
};

let latest = (~page=1, ~limit=10, ~pollInterval=?, ()) => {
  let json = Axios.use({j|d3n/blocks/latest?page=$page&limit=$limit|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Block.decodeBlocks);
};

let at_height = (height, ~pollInterval=?, ()) => {
  let json = Axios.use({j|blocks/$height|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Block.decodeBlock);
};
