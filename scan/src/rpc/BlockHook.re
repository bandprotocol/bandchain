module Block = {
  type t = {
    hash: string,
    height: int,
    timestamp: MomentRe.Moment.t,
  };

  let decode_block = json =>
    JsonUtils.Decode.{
      hash: json |> at(["block_meta", "block_id", "hash"], string),
      height: json |> at(["block_meta", "header", "height"], intstr),
      timestamp: json |> at(["block_meta", "header", "time"], moment),
    };
};

let latest = (~pollInterval=?, ()) => {
  let json = Axios.use({j|blocks/latest|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Block.decode_block);
};

let at_height = (height, ~pollInterval=?, ()) => {
  let json = Axios.use({j|blocks/$height|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Block.decode_block);
};
