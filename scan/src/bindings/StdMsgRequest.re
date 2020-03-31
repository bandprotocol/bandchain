type amount_t = {
  amount: string,
  denom: string,
};

type send_value_t = {
  oracleScriptID: string,
  calldata: string,
  requestedValidatorCount: string,
  sufficientValidatorCount: string,
  expiration: string,
  prepareGas: string,
  executeGas: string,
  sender: string,
};

type msg_t = {
  [@bs.as "type"]
  type_: string,
  value: send_value_t,
};

type fee_t = {
  amount: array(amount_t),
  gas: string,
};

type t = {
  msgs: array(msg_t),
  chain_id: string,
  fee: fee_t,
  memo: string,
  account_number: string,
  sequence: string,
};

let create =
    (
      oracleScriptID,
      calldata,
      requestedValidatorCount,
      sufficientValidatorCount,
      sender,
      feeAmount,
      accountNumber,
      sequence,
    ) => {
  {
    msgs: [|
      {
        type_: "zoracle/Request",
        value: {
          oracleScriptID: oracleScriptID |> string_of_int,
          calldata: calldata |> JsBuffer.toBase64,
          requestedValidatorCount: requestedValidatorCount |> string_of_int,
          sufficientValidatorCount: sufficientValidatorCount |> string_of_int,
          expiration: 20 |> string_of_int,
          prepareGas: 20000 |> string_of_int,
          executeGas: 150000 |> string_of_int,
          sender: sender |> Address.toBech32,
        },
      },
    |],
    chain_id: "bandchain",
    fee: {
      amount: [|{amount: feeAmount, denom: "uband"}|],
      gas: 3000000 |> string_of_int,
    },
    memo: "",
    account_number: accountNumber,
    sequence,
  };
};
