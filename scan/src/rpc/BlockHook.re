module Block = {
  type t = {
    hash: Hash.t,
    height: int,
    timestamp: MomentRe.Moment.t,
    proposer: Address.t,
    numTxs: int,
    totalTxs: int,
  };

  let decodeBlock = json =>
    JsonUtils.Decode.{
      hash: json |> at(["block_meta", "block_id", "hash"], string) |> Hash.fromHex,
      height: json |> at(["block_meta", "header", "height"], intstr),
      timestamp: json |> at(["block_meta", "header", "time"], moment),
      proposer:
        json |> at(["block_meta", "header", "proposer_address"], string) |> Address.fromHex,
      numTxs: json |> at(["block_meta", "header", "num_txs"], intstr),
      totalTxs: json |> at(["block_meta", "header", "total_txs"], intstr),
    };

  let decodeBlocks = json => JsonUtils.Decode.(json |> list(decodeBlock));

  let getProposerMoniker = (block: t, validators: list(ValidatorHook.Validator.t)) =>
    validators
    ->Belt_List.keepMap(validator =>
        validator.consensusPubkey |> PubKey.toAddress == block.proposer
          ? Some(validator.moniker) : None
      )
    ->Belt_List.get(0)
    ->Belt_Option.getWithDefault("Unknown");
};

let latest = (~page=1, ~limit=10, ~pollInterval=?, ()) => {
  let json = AxiosHooks.use({j|d3n/blocks/latest?page=$page&limit=$limit|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Block.decodeBlocks);
};

let atHeight = (height, ~pollInterval=?, ()) => {
  let json = AxiosHooks.use({j|blocks/$height|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Block.decodeBlock);
};
