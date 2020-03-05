module DataSource = {
  type revision_t = {
    name: string,
    timestamp: MomentRe.Moment.t,
    block: int,
    txHash: Hash.t,
  };

  type t = {
    owner: Address.t,
    name: string,
    fee: list(TxHook.Coin.t),
    executable: JsBuffer.t,
    requests: list(RequestHook.Request.t),
    revisions: list(revision_t),
  };

  let decode = json =>
    JsonUtils.Decode.{
      owner: json |> field("owner", string) |> Address.fromBech32,
      name: json |> field("name", string),
      fee: json |> field("fee", list(TxHook.Coin.decodeCoin)),
      executable: json |> field("executable", string) |> JsBuffer.fromBase64,
      requests: [],
      revisions: [],
    };
};

let getDataSource = dataSourceID => {
  let json = AxiosHooks.use({j|zoracle/data_source/$dataSourceID|j});
  json |> Belt.Option.map(_, DataSource.decode);
  // TODO: Add requests that use this data source
  // TODO: Add revision txs that create and change this data source
};
