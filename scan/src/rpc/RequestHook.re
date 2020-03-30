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
    id: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
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
      id: json |> field("id", ID.Request.fromJson),
      oracleScriptID: json |> field("oracleScriptID", ID.OracleScript.fromJson),
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
      resolveStatus: json |> field("resolveStatus", int) |> getResolveStatus,
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

// Mock
let mockReports: list(Report.t) = [
  {
    reporter: "21304A6071c15A0d18f3101a621b70673b1a6eeA" |> Address.fromHex,
    txHash: "983d54647a55f9c7965700ea7170a437c1f93e3997486023f725736b37975fd6" |> Hash.fromHex,
    reportedAtHeight: 15,
    reportedAtTime: 1585560116000. |> MomentRe.momentWithTimestampMS,
    values: [
      {externalDataID: 2, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 3, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 1, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 4, data: "6c756c75" |> JsBuffer.fromHex},
    ],
  },
  {
    reporter: "f21304A6071b70673c15A0d183101a621b1a6eeA" |> Address.fromHex,
    txHash: "bf7f29332c129628b2fd9e344ac1dfa8704001356966289201e22a0408382820" |> Hash.fromHex,
    reportedAtHeight: 15,
    reportedAtTime: 1585560116000. |> MomentRe.momentWithTimestampMS,
    values: [
      {externalDataID: 2, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 4, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 1, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 3, data: "6c756c75" |> JsBuffer.fromHex},
    ],
  },
  {
    reporter: "0760979e23e829dd41ab346b0f1112ae5aa911c3" |> Address.fromHex,
    txHash: "001afa2a6d4c817f952430f241fb6cef93903d5c3bb532818292bfe436e57693" |> Hash.fromHex,
    reportedAtHeight: 15,
    reportedAtTime: 1585560116000. |> MomentRe.momentWithTimestampMS,
    values: [
      {externalDataID: 2, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 4, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 1, data: "6c756c75" |> JsBuffer.fromHex},
      {externalDataID: 3, data: "6c756c75" |> JsBuffer.fromHex},
    ],
  },
];

let get: 'a => option(Request.t) =
  id => {
    // let json = AxiosHooks.use({j|zoracle/request/$id|j});
    // json |> Belt.Option.map(_, Request.decode);
    let random = Js.Math.random_int(0, 4);
    Some({
      id,
      oracleScriptID: 999,
      oracleScriptName: "Mock Crypto price script",
      calldata: "626974636f696e" |> JsBuffer.fromHex,
      requestedValidators: [
        "band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun" |> Address.fromBech32,
        "band1j9vk75jjty02elhwqqjehaspfslaem8p0utr4q" |> Address.fromBech32,
        "band1cg26m90y3wk50p9dn8pema8zmaa22plxa0hnrg" |> Address.fromBech32,
      ],
      sufficientValidatorCount: 3,
      expirationHeight: 35235,
      resolveStatus: Open,
      requester: "ceAd18f31A021b70671521304A6071b1a6e301a6" |> Address.fromHex,
      txHash: "486023f72a90ea7796570575f6b3797755f9cd683d546170a437c1f93e399743" |> Hash.fromHex,
      requestedAtHeight: 10000,
      requestedAtTime: 1585560116000. |> MomentRe.momentWithTimestampMS,
      rawDataRequests: [
        {externalID: 3, dataSourceID: 1, calldata: "68656c6c6f" |> JsBuffer.fromHex},
        {externalID: 4, dataSourceID: 353, calldata: "6d756d75" |> JsBuffer.fromHex},
        {externalID: 1, dataSourceID: 5, calldata: "70686f746f" |> JsBuffer.fromHex},
        {externalID: 2, dataSourceID: 6, calldata: "636c6f7564" |> JsBuffer.fromHex},
      ],
      reports: mockReports->Belt_List.drop(random)->Belt_Option.getWithDefault([]),
      result: Some("6c756c75" |> JsBuffer.fromHex),
    });
  };

let getList = (~page=1, ~limit=10, ()) => {
  let json = AxiosHooks.use({j|zoracle/requests?page=$page&limit=$limit|j});
  json |> Belt.Option.map(_, Request.decodeList);
};
