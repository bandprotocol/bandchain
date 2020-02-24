module Report = {
  type t = {
    reporter: Address.t,
    txHash: Hash.t,
    reportedAtHeight: int,
    reportedAtTime: MomentRe.Moment.t,
    values: option(array((string, Js.Json.t))),
  };

  let decode = json =>
    JsonUtils.Decode.{
      reporter: json |> at(["reporter"], string) |> Address.fromBech32,
      txHash: json |> at(["txhash"], string) |> Hash.fromHex,
      reportedAtHeight: json |> at(["reportedAtHeight"], intstr),
      reportedAtTime: json |> at(["reportedAtTime"], moment),
      values:
        json |> optional(at(["value"], dict(x => x))) |> Belt_Option.map(_, Js.Dict.entries),
    };

  let decodeReports = json => JsonUtils.Decode.(json |> field("reports", list(decode)));
};

module Request = {
  type t = {
    info: ScriptHook.ScriptInfo.t,
    codeHash: Hash.t,
    params: array((string, Js.Json.t)),
    targetBlock: int,
    requester: Address.t,
    txHash: Hash.t,
    requestedAtHeight: int,
    requestedAtTime: MomentRe.Moment.t,
    reports: list(Report.t),
    result: array((string, Js.Json.t)),
  };

  let decodeResult = json =>
    JsonUtils.Decode.{
      info: json |> field("scriptInfo", ScriptHook.ScriptInfo.decode),
      codeHash: json |> at(["codeHash"], string) |> Hash.fromHex,
      params: json |> at(["params"], dict(x => x)) |> Js.Dict.entries,
      targetBlock: json |> at(["targetBlock"], intstr),
      requester: json |> at(["requester"], string) |> Address.fromHex,
      txHash: json |> at(["txhash"], string) |> Hash.fromHex,
      requestedAtHeight: json |> at(["requestedAtHeight"], intstr),
      requestedAtTime: json |> at(["requestedAtTime"], moment),
      reports: json |> Report.decodeReports,
      result: json |> at(["result"], dict(x => x)) |> Js.Dict.entries,
    };

  let decode = json => JsonUtils.Decode.(json |> field("result", decodeResult));
};

let getRequest = (reqID, ~pollInterval=?, ()) => {
  let json = AxiosHooks.use({j|zoracle/request/$reqID|j}, ~pollInterval?, ());
  json |> Belt.Option.map(_, Request.decode);
};
