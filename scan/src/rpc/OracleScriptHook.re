module OracleScript = {
  type t = {
    id: int,
    owner: Address.t,
    name: string,
    code: JsBuffer.t,
  };

  let decodeResult = json =>
    JsonUtils.Decode.{
      id: json |> field("id", intstr),
      owner: json |> field("owner", string) |> Address.fromBech32,
      name: json |> field("name", string),
      code: json |> field("code", string) |> JsBuffer.fromBase64,
    };

  let decode = json => JsonUtils.Decode.(json |> field("result", decodeResult));

  let decodeList = json => JsonUtils.Decode.(json |> field("result", list(decodeResult)));
};

let get = oracleScriptID => {
  let json = AxiosHooks.use({j|zoracle/oracle_script/$oracleScriptID|j});
  json |> Belt.Option.map(_, OracleScript.decode);
};

let getScriptList = (~page=1, ~limit=10, ()) => {
  let json = AxiosHooks.use({j|zoracle/oracle_scripts?page=$page&limit=$limit|j});
  json |> Belt.Option.map(_, OracleScript.decodeList);
};
