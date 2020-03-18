module OracleScript = {
  type revision_t = {
    name: string,
    timestamp: MomentRe.Moment.t,
    height: int,
    txHash: Hash.t,
  };

  type t = {
    id: int,
    owner: Address.t,
    name: string,
    description: string,
    relatedDataSource: list(ID.DataSource.t),
    code: JsBuffer.t,
    requests: list(RequestHook.Request.t),
    revisions: list(revision_t),
  };

  let decodeResult = json =>
    JsonUtils.Decode.{
      id: json |> field("id", intstr),
      owner: json |> field("owner", string) |> Address.fromBech32,
      name: json |> field("name", string),
      description: json |> field("description", string),
      relatedDataSource: [ID.DataSource.ID(1), ID.DataSource.ID(2), ID.DataSource.ID(3)],
      code: json |> field("code", string) |> JsBuffer.fromBase64,
      requests: [
        {
          id: 1,
          oracleScriptID: 1,
          oracleScriptName: "name",
          calldata: "AAAAAAAAV0M=" |> JsBuffer.fromBase64,
          requestedValidators: [
            "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
            "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
          ],
          sufficientValidatorCount: 2,
          expirationHeight: 3000,
          resolveStatus: Success,
          requester: "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
          txHash:
            "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
          requestedAtHeight: 40000,
          requestedAtTime: MomentRe.momentNow(),
          rawDataRequests: [],
          reports: [],
          result: Some("AAAAAAAAV0M=" |> JsBuffer.fromBase64),
        },
      ],
      revisions: [
        {
          name: "Coingecko script v2",
          timestamp: MomentRe.momentWithUnix(1583465551),
          height: 472395,
          txHash:
            "6E1EAE347E7F2E27DFE6F21328DF7EB6A599D4F0ED73D54B356C77646FBEC33D" |> Hash.fromHex,
        },
        {
          name: "Coingecko script",
          timestamp: MomentRe.momentWithUnix(1583465050),
          height: 472295,
          txHash:
            "D3C77B93B10169E9D3C5ACA9A4A049CED40D7BE231E5D1A79FFAE7498952A032" |> Hash.fromHex,
        },
      ],
    };

  let decode = json => JsonUtils.Decode.(json |> field("result", decodeResult));

  let decodeList = json =>
    JsonUtils.Decode.(
      json
      |> optional(field("result", list(decodeResult)))
      |> Belt.Option.getWithDefault(_, [])
    );
};

let get = oracleScriptID => {
  let json = AxiosHooks.use({j|zoracle/oracle_script/$oracleScriptID|j});
  json |> Belt.Option.map(_, OracleScript.decode);
  // TODO: Add requests that use this oracle script.
  // TODO: Add revision txs that create and change this oracle script.
};

let getList = (~page=1, ~limit=10, ()) => {
  let json = AxiosHooks.use({j|zoracle/oracle_scripts?page=$page&limit=$limit|j});
  json |> Belt.Option.map(_, OracleScript.decodeList);
};
