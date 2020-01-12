module Script = {
  type field = {
    name: string,
    dataType: string,
  };

  let decodeField = json =>
    JsonUtils.Decode.{
      name: json |> field("name", string),
      dataType: json |> field("type", string),
    };

  type t = {
    name: string,
    params: list(field),
    dataSources: list(field),
    creator: Address.t,
    txHash: Hash.t,
    createdAtHeight: int,
    createdAtTime: MomentRe.Moment.t,
  };

  let decodeResultScript = json =>
    JsonUtils.Decode.{
      name: json |> at(["info", "name"], string),
      params: json |> at(["info", "params"], list(decodeField)),
      dataSources: json |> at(["info", "dataSources"], list(decodeField)),
      creator: json |> at(["info", "creator"], string) |> Address.fromBech32,
      txHash: json |> field("txhash", string) |> Hash.fromHex,
      createdAtHeight: json |> field("createdAtHeight", intstr),
      createdAtTime: json |> field("createdAtTime", moment),
    };

  let decodeScript = json => JsonUtils.Decode.(json |> field("result", decodeResultScript));
};

let getInfo = codeHash => {
  let codeHashHex = codeHash->Hash.toHex;
  let json = Axios.use({j|zoracle/script/$codeHashHex|j}, ());
  json |> Belt.Option.map(_, Script.decodeScript);
};
