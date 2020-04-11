module Request = {
  type t = {
    id: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    calldata: JsBuffer.t,
    requestedValidatorCount: int,
    sufficientValidatorCount: int,
    expiration: int,
    executeGas: int,
    sender: Address.t,
  };
};

module Response = {
  type status_t =
    | Success
    | Fail;

  type t = {
    requestID: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    status: status_t,
    result: option(JsBuffer.t),
  };
};

type packet_direction_t =
  | Incoming
  | Outgoing;

type packet_t =
  | Unknown
  | Request(Request.t)
  | Response(Response.t);

type t = {
  direction: packet_direction_t,
  chainID: string,
  chennel: string,
  port: string,
  blockHeight: ID.Block.t,
  packet: packet_t,
};

// TODO: replace this mock when wireup
let getMockList = () => [|
  {
    direction: Incoming,
    chainID: "wenchang testnet v0",
    chennel: "htjvlvazyj",
    port: "bibc1",
    blockHeight: ID.Block.ID(9999),
    packet:
      Request({
        id: ID.Request.ID(888),
        oracleScriptID: ID.OracleScript.ID(7777),
        oracleScriptName: "Mock Oracle Script",
        calldata: "aa" |> JsBuffer.fromHex,
        requestedValidatorCount: 4,
        sufficientValidatorCount: 3,
        expiration: 6666,
        executeGas: 1000000,
        sender: "band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs" |> Address.fromBech32,
      }),
  },
  {
    direction: Outgoing,
    chainID: "Gawa testnet 01",
    chennel: "tjlvazhyjv",
    port: "oracle",
    blockHeight: ID.Block.ID(10999),
    packet:
      Response({
        requestID: ID.Request.ID(10888),
        oracleScriptID: ID.OracleScript.ID(10777),
        oracleScriptName: "Mock Oracle Script",
        status: Response.Success,
        result: Some("aa" |> JsBuffer.fromHex),
      }),
  },
  {
    direction: Outgoing,
    chainID: "Mumu network 01",
    chennel: "azjlvvthyj",
    port: "mumian_port",
    blockHeight: ID.Block.ID(20999),
    packet:
      Response({
        requestID: ID.Request.ID(20888),
        oracleScriptID: ID.OracleScript.ID(20777),
        oracleScriptName: "Mock Oracle Script",
        status: Response.Fail,
        result: None,
      }),
  },
|];
