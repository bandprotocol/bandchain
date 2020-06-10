// TODO: Replace this mock by the real
type request_t =
  | Request(RequestSub.t)
  | RequestMini(RequestSub.Mini.t);

type request_packet_t = {
  clientID: string,
  oracleScriptID: int,
  calldata: JsBuffer.t,
  askCount: int,
  minCount: int,
};

type response_packet_t = {
  clientID: string,
  requestID: int,
  ansCount: int,
  requestTime: int,
  resolveTimne: int,
  resolveStatus: int,
  result: JsBuffer.t,
};

let int8ToHex: int => string = [%bs.raw
  {|
  function int32ToHex(x) {
    return x.toString(16).padStart(2,0)
  }
|}
];

let int32ToHex: int => string = [%bs.raw
  {|
  function int32ToHex(x) {
    return x.toString(16).padStart(8,0)
  }
|}
];

let int64ToHex: int => string = [%bs.raw
  {|
  function int32ToHex(x) {
    return x.toString(16).padStart(16,0)
  }
|}
];

let encodeString = (s: string) => {
  let buf = JsBuffer.fromUTF8(s);
  JsBuffer.concat([|buf->JsBuffer.byteLength->int32ToHex->JsBuffer.fromHex, buf|]);
};

let encodeBuffer = (buf: JsBuffer.t) =>
  JsBuffer.concat([|buf->JsBuffer.byteLength->int32ToHex->JsBuffer.fromHex, buf|]);

let encodeRequest = (req: request_packet_t) => {
  JsBuffer.concat([|
    req.clientID->encodeString,
    req.oracleScriptID->int64ToHex->JsBuffer.fromHex,
    req.calldata->encodeBuffer,
    req.askCount->int64ToHex->JsBuffer.fromHex,
    req.minCount->int64ToHex->JsBuffer.fromHex,
  |]);
};

let encodeResponse = (res: response_packet_t) => {
  JsBuffer.concat([|
    res.clientID->encodeString,
    res.requestID->int64ToHex->JsBuffer.fromHex,
    res.ansCount->int64ToHex->JsBuffer.fromHex,
    res.requestTime->int64ToHex->JsBuffer.fromHex,
    res.resolveTimne->int64ToHex->JsBuffer.fromHex,
    res.resolveStatus->int8ToHex->JsBuffer.fromHex,
    res.result->encodeBuffer,
  |]);
};

let resolveStatusToInt = (rs: RequestSub.resolve_status_t) =>
  switch (rs) {
  | Pending => 0
  | Success => 1
  | Failure => 2
  | Unknown => 3
  };

let toPackets = (request: request_t) => {
  switch (request) {
  | Request({
      clientID,
      oracleScript: {oracleScriptID: ID(oracleScriptID)},
      calldata,
      requestedValidators,
      minCount,
      id: ID(requestID),
      reports,
      requestTime,
      resolveTime,
      resolveStatus,
      result,
    }) => (
      {
        clientID,
        oracleScriptID,
        calldata,
        askCount: requestedValidators |> Belt_Array.length,
        minCount,
      },
      {
        clientID,
        requestID,
        ansCount: reports |> Belt_Array.length,
        requestTime: requestTime |> Belt_Option.getExn |> MomentRe.Moment.toUnix,
        resolveTimne: resolveTime |> Belt_Option.getExn |> MomentRe.Moment.toUnix,
        resolveStatus: resolveStatus |> resolveStatusToInt,
        result: result |> Belt_Option.getExn,
      },
    )
  | RequestMini({
      clientID,
      oracleScriptID: ID(oracleScriptID),
      calldata,
      askCount,
      minCount,
      id: ID(requestID),
      reportsCount: ansCount,
      requestTime,
      resolveTime,
      resolveStatus,
      result,
    }) => (
      {clientID, oracleScriptID, calldata, askCount, minCount},
      {
        clientID,
        requestID,
        ansCount,
        requestTime: requestTime |> Belt_Option.getExn |> MomentRe.Moment.toUnix,
        resolveTimne: resolveTime |> Belt_Option.getExn |> MomentRe.Moment.toUnix,
        resolveStatus: resolveStatus |> resolveStatusToInt,
        result: result |> Belt_Option.getExn,
      },
    )
  };
};

let createProof = (request: request_t) => {
  let (req, res) = request->toPackets;
  JsBuffer.concat([|req->encodeRequest, res->encodeResponse|]);
};
