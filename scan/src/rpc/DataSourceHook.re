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

  let decodeResult = json =>
    JsonUtils.Decode.{
      owner: json |> field("owner", string) |> Address.fromBech32,
      name: json |> field("name", string),
      fee: json |> field("fee", list(TxHook.Coin.decodeCoin)),
      executable: json |> field("executable", string) |> JsBuffer.fromBase64,
      requests: [
        {
          info: {
            name: "Oracle script1",
            codeHash:
              "6E1EAE347E7F2E27DFE6F21328DF7EB6A599D4F0ED73D54B356C77646FBEC33D" |> Hash.fromHex,
            params: [],
            dataSources: [],
            result: [],
            creator: "band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9" |> Address.fromBech32,
          },
          codeHash:
            "6E1EAE347E7F2E27DFE6F21328DF7EB6A599D4F0ED73D54B356C77646FBEC33D" |> Hash.fromHex,
          params: [|("symbol", "ETH" |> Json.Encode.string)|],
          targetBlock: 50000,
          requester: "band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9" |> Address.fromBech32,
          txHash:
            "D3C77B93B10169E9D3C5ACA9A4A049CED40D7BE231E5D1A79FFAE7498952A032" |> Hash.fromHex,
          requestedAtHeight: 45000,
          requestedAtTime: MomentRe.momentWithUnix(1583465551),
          reports: [],
          result: [|("price_in_usd", 201.89 |> Json.Encode.float)|],
        },
        {
          info: {
            name: "Oracle script5",
            codeHash:
              "6E1EAE347E7F2E27DFE6F21328DF7EB6A599D4F0ED73D54B356C77646FBEC33D" |> Hash.fromHex,
            params: [],
            dataSources: [],
            result: [],
            creator: "band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9" |> Address.fromBech32,
          },
          codeHash:
            "6E1EAE347E7F2E27DFE6F21328DF7EB6A599D4F0ED73D54B356C77646FBEC33D" |> Hash.fromHex,
          params: [|("symbol", "BTC" |> Json.Encode.string)|],
          targetBlock: 50000,
          requester: "band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9" |> Address.fromBech32,
          txHash:
            "D3C77B93B10169E9D3C5ACA9A4A049CED40D7BE231E5D1A79FFAE7498952A032" |> Hash.fromHex,
          requestedAtHeight: 45000,
          requestedAtTime: MomentRe.momentWithUnix(1583465551),
          reports: [],
          result: [|("price_in_usd", 9065.89 |> Json.Encode.float)|],
        },
      ],
      revisions: [
        {
          name: "Coingecko script v2",
          timestamp: MomentRe.momentWithUnix(1583465551),
          block: 472395,
          txHash:
            "6E1EAE347E7F2E27DFE6F21328DF7EB6A599D4F0ED73D54B356C77646FBEC33D" |> Hash.fromHex,
        },
        {
          name: "Coingecko script",
          timestamp: MomentRe.momentWithUnix(1583465050),
          block: 472295,
          txHash:
            "D3C77B93B10169E9D3C5ACA9A4A049CED40D7BE231E5D1A79FFAE7498952A032" |> Hash.fromHex,
        },
      ],
    };

  let decode = json => JsonUtils.Decode.(json |> field("result", decodeResult));
};

let getDataSource = dataSourceID => {
  let json = AxiosHooks.use({j|zoracle/data_source/$dataSourceID|j});
  json |> Belt.Option.map(_, DataSource.decode);
  // TODO: Add requests that use this data source
  // TODO: Add revision txs that create and change this data source
};
