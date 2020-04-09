module Request = {
  type t = {
    id: ID.Request.t,
    oracleScriptID: ID.OracleScript.t,
    oracleScriptName: string,
    calldata: JsBuffer.t,
    requestedValidatorCount: int,
    sufficientValidatorCount: int,
    expiration: int,
    prepareGas: int,
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
    calldata: JsBuffer.t,
  };
};
