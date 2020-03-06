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

  let decodeReports = json => JsonUtils.Decode.(json |> field("reports", list(decode)));
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
  type t = {
    requestID: int,
    oracleScriptID: int,
    calldata: JsBuffer.t,
    requestedValidators: list(Address.t),
    sufficientValidatorCount: int,
    expirationHeight: int,
    isResolved: bool,
    requester: Address.t,
    txHash: Hash.t,
    requestedAtHeight: int,
    requestedAtTime: MomentRe.Moment.t,
    rawDataRequests: list(RawDataRequest.t),
    reports: list(Report.t),
    result: JsBuffer.t,
  };

  let decodeResult = json =>
    JsonUtils.Decode.{
      requestID: json |> field("id", intstr),
      oracleScriptID: json |> field("oracleScriptID", intstr),
      calldata: json |> field("calldata", string) |> JsBuffer.fromBase64,
      requestedValidators:
        json
        |> field(
             "requestedValidators",
             list(validator => validator |> string |> Address.fromBech32),
           ),
      sufficientValidatorCount: json |> field("sufficientValidatorCount", intstr),
      expirationHeight: json |> field("expirationHeight", intstr),
      isResolved: json |> field("isResolved", bool),
      requester: json |> at(["requester"], string) |> Address.fromHex,
      txHash: json |> at(["requestTx", "hash"], string) |> Hash.fromHex,
      requestedAtHeight: json |> at(["requestTx", "height"], intstr),
      requestedAtTime: json |> at(["requestTx", "timestamp"], moment),
      rawDataRequests: json |> field("rawDataRequests", list(RawDataRequest.decode)),
      reports: json |> Report.decodeReports,
      result: json |> field("result", string) |> JsBuffer.fromBase64,
    };

  let decode = json => JsonUtils.Decode.(json |> field("result", decodeResult));
};

let getRequest = reqID => {
  let json = AxiosHooks.use({j|zoracle/request/$reqID|j});
  json |> Belt.Option.map(_, Request.decode);
};

// TODO: mock for now
let getRequestList = (~page=1, ~limit=10, ()) => {
  Request.[
    {
      requestID: 1,
      oracleScriptID: 1,
      calldata: "AAAAAAAAV0M=" |> JsBuffer.fromBase64,
      requestedValidators: [
        "bandvaloper13zmknvkq2sj920spz90g4r9zjan8g58423y76e" |> Address.fromBech32,
        "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
      ],
      sufficientValidatorCount: 2,
      expirationHeight: 3000,
      isResolved: true,
      requester: "bandvaloper1fwffdxysc5a0hu0falsq4lyneucj05cwryzfp0" |> Address.fromBech32,
      txHash: "AC006D7136B0041DA4568A4CA5B7C1F8E8E0B4A74F11213B99EC4956CC8A247C" |> Hash.fromHex,
      requestedAtHeight: 40000,
      requestedAtTime: MomentRe.momentNow(),
      rawDataRequests: [],
      reports: [],
      result: "AAAAAAAAV0M=" |> JsBuffer.fromBase64,
    },
  ];
};
