module RawDataReport = {
  type t = {
    externalDataID: int,
    data: JsBuffer.t,
  };

  let decode = json =>
    JsonUtils.Decode.{
      externalDataID: json |> field("externalDataID", intstr),
      data: json |> field("data", string) |> JsBuffer.fromBase64,
    };
};

module Report = {
  type t = {
    reporter: Address.t,
    txHash: Hash.t,
    reportedAtHeight: int,
    reportedAtTime: MomentRe.Moment.t,
    values: list(RawDataReport.t),
  };

  let decode = json =>
    JsonUtils.Decode.{
      reporter: json |> field("reporter", string) |> Address.fromBech32,
      txHash: json |> at(["tx", "hash"], string) |> Hash.fromHex,
      reportedAtHeight: json |> at(["tx", "height"], intstr),
      reportedAtTime: json |> at(["tx", "timestamp"], moment),
      values: json |> field("value", list(RawDataReport.decode)),
    };
};

module RawDataRequest = {
  type t = {
    externalID: int,
    dataSourceID: int,
    calldata: JsBuffer.t,
  };

  let decode = json =>
    JsonUtils.Decode.{
      externalID: json |> field("externalID", intstr),
      dataSourceID: json |> at(["detail", "dataSourceID"], intstr),
      calldata: json |> at(["detail", "calldata"], string) |> JsBuffer.fromBase64,
    };
};

module Request = {
  type resolve_status_t =
    | Open
    | Success
    | Failure
    | Unknown;

  type t = {
    id: int,
    oracleScriptID: int,
    oracleScriptName: string,
    calldata: JsBuffer.t,
    requestedValidators: list(Address.t),
    sufficientValidatorCount: int,
    expirationHeight: int,
    resolveStatus: resolve_status_t,
    requester: Address.t,
    txHash: Hash.t,
    requestedAtHeight: int,
    requestedAtTime: MomentRe.Moment.t,
    rawDataRequests: list(RawDataRequest.t),
    reports: list(Report.t),
    result: option(JsBuffer.t),
  };

  let getResolveStatus =
    fun
    | 0 => Open
    | 1 => Success
    | 2 => Failure
    | _ => Unknown;

  let decodeResult = json =>
    JsonUtils.Decode.{
      id: json |> field("id", intstr),
      oracleScriptID: json |> field("oracleScriptID", intstr),
      oracleScriptName: "Mean Crypto Price",
      calldata: json |> field("calldata", string) |> JsBuffer.fromBase64,
      requestedValidators:
        json
        |> field(
             "requestedValidators",
             list(validator => validator |> string |> Address.fromBech32),
           ),
      sufficientValidatorCount: json |> field("sufficientValidatorCount", intstr),
      expirationHeight: json |> field("expirationHeight", intstr),
      resolveStatus: Unknown, // TODO , fix this Mock
      requester: json |> at(["requester"], string) |> Address.fromBech32,
      txHash: json |> at(["requestTx", "hash"], string) |> Hash.fromHex,
      requestedAtHeight: json |> at(["requestTx", "height"], intstr),
      requestedAtTime: json |> at(["requestTx", "timestamp"], moment),
      rawDataRequests: json |> field("rawDataRequests", list(RawDataRequest.decode)),
      reports: json |> field("reports", list(Report.decode)),
      result:
        json
        |> optional(at(["result", "data"], string))
        |> Belt.Option.map(_, JsBuffer.fromBase64),
    };

  let decode = json => JsonUtils.Decode.(json |> field("result", decodeResult));

  let decodeList = json =>
    JsonUtils.Decode.(
      json
      |> optional(field("result", list(decodeResult)))
      |> Belt.Option.getWithDefault(_, [])
    );
};

let get = id => {
  let json = AxiosHooks.use({j|zoracle/request/$id|j});
  json |> Belt.Option.map(_, Request.decode);
};

let getList = (~page=1, ~limit=10, ()) => {
  let json = AxiosHooks.use({j|zoracle/requests?page=$page&limit=$limit|j});
  json |> Belt.Option.map(_, Request.decodeList);
};
